package tests

import (
	c "LinksService/controllers"
	f "LinksService/functions"
	"bytes"
	"database/sql"
	"errors"
	"github.com/betable/sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateLink(t *testing.T) {
	ip, url := "1", "2"
	link1 := f.CreateLink(ip, url)
	link2 := f.CreateLink(ip, url)

	if link1 != link2 {
		t.Errorf("Not equal!")
	}

}

func TestValidateUserShortURL(t *testing.T) {
	str := "fghjrt"

	ret := f.ValidateUserShortURL(str)

	if ret == false {
		t.Errorf(str + " is invalid. Why?")
	}

	str = "123sfg"

	ret = f.ValidateUserShortURL(str)

	if ret == false {
		t.Errorf(str + " is invalid. Why?")
	}

	str = "абвгд"

	ret = f.ValidateUserShortURL(str)

	if ret == true {
		t.Errorf(str + " is valid. Why?")
	}

	str = "абвгдабвгдабвгдабвгдабвгдабвгд" //30

	ret = f.ValidateUserShortURL(str)

	if ret == true {
		t.Errorf(str + " is valid. Why?")
	}

	str = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" //30

	ret = f.ValidateUserShortURL(str)

	if ret == true {
		t.Errorf(str + " is valid. Why?")
	}

	str = "aaa" //30

	ret = f.ValidateUserShortURL(str)

	if ret == true {
		t.Errorf(str + " is valid. Why?")
	}
}

func TestValidateLink(t *testing.T) {
	link := "http://google.com"
	res := f.ValidateLink(link)

	if res == false {
		t.Errorf("http://google.com isnt valid!")
	}

	link = "1.com"
	res = f.ValidateLink(link)

	if res == true {
		t.Errorf("1.com is valid!")
	}

	link = "google.ru"
	res = f.ValidateLink(link)

	if res == true {
		t.Errorf("google.ru is valid!")
	}
}

func TestSetShortLink_Valid(t *testing.T) {
	req, err := http.NewRequest("POST", "/setShort", bytes.NewBuffer([]byte(`{"URI":"http://google.com", "Short":"tgdfvc"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-type", "application/json")

	bh := c.BaseHandler{
		Db: nil,
		Select: func(shortURL string, db *sql.DB) (string, error) {
			return "", nil
		},
		Insert: func(fullURL, shortURL string, db *sql.DB) error {
			return nil
		},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bh.SetShortLink)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"URL":"` + req.Host + `/tgdfvc"}`
	res := rr.Body.String()
	res = res[:len(res)-1]
	if res != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSetShortLink_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/setShort", nil)
	if err != nil {
		t.Fatal(err)
	}

	bh := c.BaseHandler{
		Db: nil,
		Select: func(shortURL string, db *sql.DB) (string, error) {
			return "", nil
		},
		Insert: func(fullURL, shortURL string, db *sql.DB) error {
			return nil
		},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bh.SetShortLink)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRedirectHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/fgdfhfg", nil)
	if err != nil {
		t.Fatal(err)
	}

	bh := c.BaseHandler{
		Db: nil,
		Select: func(shortURL string, db *sql.DB) (string, error) {
			return "hello", nil
		},
		Insert: func(fullURL, shortURL string, db *sql.DB) error {
			return nil
		},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bh.Index)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	bh = c.BaseHandler{
		Db: nil,
		Select: func(shortURL string, db *sql.DB) (string, error) {
			return "", nil
		},
		Insert: func(fullURL, shortURL string, db *sql.DB) error {
			return nil
		},
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(bh.Index)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	bh = c.BaseHandler{
		Db: nil,
		Select: func(shortURL string, db *sql.DB) (string, error) {
			return "", errors.New("Mock")
		},
		Insert: func(fullURL, shortURL string, db *sql.DB) error {
			return nil
		},
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(bh.Index)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadGateway {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestInsertURLs(t *testing.T) {
	db, mock, err := sqlmock.New()
	mock.ExpectExec("INSERT INTO links(.*, .*) VALUES (.*, .*)").
		WithArgs("http://google.com", "fgnbrt").
		WillReturnResult(sqlmock.NewResult(0, 1))
	err = f.InsertURLs("http://google.com", "fgnbrt", db)
	if err != nil {
		t.Errorf("Insert method failed!")
	}
}

func TestSelectShortURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT userLink FROM links WHERE shortLink = (.*)").
		WithArgs("hi").
		WillReturnRows(sqlmock.NewRows([]string{"userLink"}).AddRow("http://google.com"))

	_, err = f.SelectShortURL("hi", db)
	if err != nil {
		t.Errorf("Database null connection passed!")
	}
}
