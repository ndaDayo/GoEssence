package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Entry struct {
	AuthorID string
	Author   string
	TitleID  string
	Title    string
	SiteURL  string
	ZipURL   string
}

func main() {
	listURL := "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := findEntires(listURL)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Printf(entry.Title, entry.ZipURL)
	}
}

func findEntires(siteURL string) ([]Entry, error) {
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return nil, err
	}
	pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)
	doc.Find("ol li a").Each(func(n int, elem *goquery.Selection) {
		token := pat.FindStringSubmatch(elem.AttrOr("href", ""))
		if len(token) != 3 {
			return
		}

		pageURL := fmt.Sprintf("https://www.aozora.gr.jp/cards/%s/card%s.html",
			token[1], token[2])
		author, zipURL := findAuthorAndZIP(pageURL)
		println(zipURL)
	})

	return nil, nil
}

func findAuthorAndZIP(siteURL string) (string, string) {
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
	return author, zipURL
}
