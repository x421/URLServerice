package main

import (
	c "LinksService/controllers"
	f "LinksService/functions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

/*
Валидацию ВСЕГО пользовательского ввода по полной!
Отдельный файл под константы и глоб перем и ошибки?
*/
func main() {
	/*
		если джсон кривой
	*/
	db := f.GetSQLConnection("mysql", "root", "root", "localhost", "3306", "db", nil)
	bh := c.BaseHandler{
		Db:     db,
		Select: f.SelectShortURL,
		Insert: f.InsertURLs,
	}

	http.HandleFunc("/setShort", bh.SetShortLink)
	/*
		пользователь подает бред
	*/
	http.HandleFunc("/", bh.Index)
	http.ListenAndServe(":80", nil)
}
