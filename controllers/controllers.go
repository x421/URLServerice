package controllers

import (
	f "LinksService/functions"
	s "LinksService/structs"
	"encoding/json"
	"fmt"
	"net/http"
)

type BaseHandler s.BaseHandler

func (bh *BaseHandler) SetShortLink(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost { // кроме пост есть еще что то
		http.Error(writer, "POST only", http.StatusBadRequest)
		return
	}

	var ClientData s.URL
	err := json.NewDecoder(request.Body).Decode(&ClientData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	res := f.ValidateLink(ClientData.URI)
	if res == false {
		http.Error(writer, "URL is not valid!", http.StatusBadRequest)
		return
	}

	if ClientData.Short == "" {
		ClientData.Short = f.CreateLink(request.RemoteAddr, ClientData.URI)
	}

	// отображать схему таблицы не самая лучшая затея+
	db := bh.Db

	_, err = db.Exec("INSERT INTO links(`userLink`, `shortLink`) VALUES (?, ?)", ClientData.URI, ClientData.Short)

	// не факт конечно что эти ошибки, но эти вероятны в 99 случаях
	if err != nil {
		http.Error(writer, "Same short link already ib DB", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	var ret s.ReturnURL
	ret.URL = request.Host +"/"+ ClientData.Short
	json.NewEncoder(writer).Encode(ret)
}

func (bh *BaseHandler) RedirectHandler(writer http.ResponseWriter, request *http.Request, path string) {
	db := bh.Db

	rows, err := db.Query("SELECT userLink FROM links WHERE shortLink = ?", path)

	if err != nil {
		http.Error(writer, "DB select error", http.StatusBadRequest)
		return
	}

	defer rows.Close()
	var link string

	if rows.NextResultSet() == true {
		rows.Next()
		rows.Scan(&link)
	} else {
		http.Error(writer, "Not Found", 404)
		return
	}

	http.Redirect(writer, request, link, 301)
}

func (bh *BaseHandler) Index(writer http.ResponseWriter, request *http.Request){
	if request.URL.Path != "/" {
		bh.RedirectHandler(writer, request, string([]rune(request.URL.Path)[1:]))
		return
	}
	fmt.Fprintf(writer, "Hello")
}