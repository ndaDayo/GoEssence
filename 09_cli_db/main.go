package main

import (
	"archive/zip"
	"bytes"
	"cli_db/domain"
	"cli_db/repository"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/encoding/japanese"
)

func main() {
	db, err := repository.SetupDB("database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	listURL := "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := findEntires(listURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("found %s entries", len(entries))

	for _, entry := range entries {
		log.Printf("adding %+v", entry)
		content, err := extractText(entry.ZipURL)
		if err != nil {
			log.Println(err)
			continue
		}

		err = repository.AddEntry(db, &entry, content)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func findEntires(siteURL string) ([]domain.Entry, error) {
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return nil, err
	}
	pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)
	entries := []domain.Entry{}
	doc.Find("ol li a").Each(func(n int, elem *goquery.Selection) {
		token := pat.FindStringSubmatch(elem.AttrOr("href", ""))
		if len(token) != 3 {
			return
		}

		title := elem.Text()
		pageURL := fmt.Sprintf("https://www.aozora.gr.jp/cards/%s/card%s.html",
			token[1], token[2])
		author, zipURL := findAuthorAndZIP(pageURL)
		if zipURL != "" {
			entries = append(entries, domain.Entry{
				AuthorID: token[1],
				Author:   author,
				TitleID:  token[2],
				Title:    title,
				SiteURL:  siteURL,
				ZipURL:   zipURL,
			})
		}
	})

	return entries, nil
}

func findAuthorAndZIP(siteURL string) (string, string) {
	log.Println("query", siteURL)
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return "", ""
	}

	author := doc.Find("table[summary=作家データ] tr:nth-child(2) td:nth-child(2)").First().Text()
	zipURL := ""
	doc.Find("table.download a").Each(func(n int, elem *goquery.Selection) {
		href := elem.AttrOr("href", "")
		if strings.HasSuffix(href, ".zip") {
			zipURL = href
		}
	})

	if zipURL == "" {
		return author, ""
	}

	if strings.HasPrefix(zipURL, "http://") || strings.HasPrefix(zipURL, "http://") {
		return author, zipURL
	}

	u, err := url.Parse(siteURL)

	if err != nil {
		return author, ""
	}

	u.Path = path.Join(path.Dir(u.Path), zipURL)
	return author, zipURL
}

func extractText(zipURL string) (string, error) {
	resp, err := http.Get(zipURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))

	for _, file := range r.File {
		if path.Ext(file.Name) == ".txt" {
			f, err := file.Open()
			if err != nil {
				return "", err
			}
			b, err := ioutil.ReadAll(f)
			f.Close()
			if err != nil {
				return "", err
			}
			b, err = japanese.ShiftJIS.NewDecoder().Bytes(b)
			if err != nil {
				return "", err
			}

			return string(b), nil
		}
	}

	return "", errors.New("contents not found")
}
