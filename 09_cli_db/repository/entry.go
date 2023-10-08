package repository

import (
	"cli_db/domain"
	"database/sql"
	"log"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func SetupDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	for _, query := range queries() {
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	return db, nil

}

func queries() []string {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS authors(author_id TEXT, author TEXT, PRIMARY KEY(author_id))`,
		`CREATE TABLE IF NOT EXISTS contents(author_id TEXT, title_id TEXT, title TEXT, content TEXT, PRIMARY KEY(author_id, title_id))`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS contents_fts USING fts4(words)`,
	}

	return queries
}

func AddEntry(db *sql.DB, entry *domain.Entry, content string) error {
	_, err := db.Exec(`
        REPLACE INTO authors(author_id, author) values(?, ?)
    `,
		entry.AuthorID,
		entry.Author,
	)

	if err != nil {
		return err
	}
	res, err := db.Exec(`
        REPLACE INTO contents(author_id, title_id, title, conten) values(?, ?, ?, ?)
    `,
		entry.AuthorID,
		entry.TitleID,
		entry.Title,
		content,
	)

	if err != nil {
		return err
	}
	docID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return err
	}
	seg := t.Wakati(content)
	_, err = db.Exec(
		`REPLACE INTO contents_fts(docid, words) values (?, ?)`,
		docID,
		strings.Join(seg, " "),
	)
	if err != nil {
		return err
	}

	return nil
}
