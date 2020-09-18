package functions

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

func ValidateUserShortURL(link string) bool {
	ret, err := regexp.MatchString("[a-zA-Z0-9-]", link)
	if err != nil {
		return false
	}

	return ret && (len(link) < 25 && len(link) > 4)
}

func ValidateLink(link string) bool {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}

	u, err := url.Parse(link)
	if err != nil || u.Scheme == "" || u.Host == "" || len(link) > 100 {
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
	// с этой ошибкой сложно stub
	row.Scan(&link)

	return link, nil
}
