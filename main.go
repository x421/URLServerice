package main

import (
	c "LinksService/controllers"
	f "LinksService/functions"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
)

func main() {
	str := os.Getenv("User") + ":" + os.Getenv("Pass") + "@(" + os.Getenv("Ip") + ":" + os.Getenv("PortDB") + ")/db"
	conn, err := sql.Open("mysql", str)
	if err != nil {
		panic("DB connect error")
	}
	bh := c.BaseHandler{
		Db:     conn,
		Select: f.SelectShortURL,
		Insert: f.InsertURLs,
	}

	http.HandleFunc("/setShort", bh.SetShortLink)
	http.HandleFunc("/", bh.Index)
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
