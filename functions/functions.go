package functions

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetSQLConnection(driver, user, pass, server, port, db string, writer http.ResponseWriter) *sql.DB {
	conn, err := sql.Open(driver, user+":"+pass+"@("+server+":"+port+")/"+db)
	if err != nil {
		http.Error(writer, "DB connect error", http.StatusBadRequest)
	}
	return conn
}

func ValidateLink(link string) bool {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}

	u, err := url.Parse(link)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func CreateLink(ip, url string) string {
	str := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	str += ip
	str += url
	h := md5.New()
	h.Write([]byte(str))

	return (hex.EncodeToString(h.Sum(nil)))[3:8]
}

func InsertURLs(fullURL, shortURL string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO links(`userLink`, `shortLink`) VALUES (?, ?)", fullURL, shortURL)
	return err
}
func SelectShortURL(shortURL string, db *sql.DB) (string, error) {
	row := db.QueryRow("SELECT userLink FROM links WHERE shortLink = ?", shortURL)
	link := ""
	err := row.Scan(&link)

	return link, err
}
