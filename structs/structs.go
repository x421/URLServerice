package structs

import "database/sql"

type URL struct {
	URI   string
	Short string
}

type ReturnURL struct {
	URL     string
	ErrStr  string
	ErrCode int
}

type InsertURLs func(fullURL, shortURL string, db *sql.DB) error
type SelectShortURL func(shortURL string, db *sql.DB) (string, error)
