package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type Entry struct {
	AuthorID string
	Author   string
	TitleID  string
	Title    string
	InfoURL  string
	ZipURL   string
}

func findEntires(siteURL string) ([]Entry, error) {
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func main() {
	listURL = "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := findEntires(listURL)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Printf(entry.Title, entry.ZipURL)
	}
}
