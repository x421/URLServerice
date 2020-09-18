package controllers

import (
	f "LinksService/functions"
	s "LinksService/structs"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type BaseHandler struct {
	Select s.SelectShortURL
	Insert s.InsertURLs
	Db     *sql.DB
}

func (bh *BaseHandler) SetShortLink(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost { // кроме пост есть еще что то
		http.Error(writer, "POST only", http.StatusBadRequest)
		return
	}

	var ClientData s.URL
	err := json.NewDecoder(request.Body).Decode(&ClientData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
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

	err = bh.Insert(ClientData.URI, ClientData.Short, bh.Db)

	// не факт конечно что эти ошибки, но эти вероятны в 99 случаях
	if err != nil {
		http.Error(writer, "Same short link already ib DB", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	var ret s.ReturnURL
	ret.URL = request.Host + "/" + ClientData.Short
	json.NewEncoder(writer).Encode(ret)
}

func (bh *BaseHandler) RedirectHandler(writer http.ResponseWriter, request *http.Request, path string) {
	link, err := bh.Select(path, bh.Db)

	if err != nil {
		http.Error(writer, "DB select error", http.StatusBadGateway)
		return
	}

	if link != "" {
		http.Redirect(writer, request, link, http.StatusFound)
	} else {
		http.Error(writer, "Not Found", http.StatusNotFound)
	}
}

func (bh *BaseHandler) Index(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		bh.RedirectHandler(writer, request, string([]rune(request.URL.Path)[1:]))
		return
	}
	fmt.Fprintf(writer, "Hello")
}
