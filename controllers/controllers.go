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

func SendAnswer(writer http.ResponseWriter, str string, num int, result string) {
	var ret s.ReturnURL
	ret.ErrStr = str
	ret.ErrCode = num
	ret.URL = result
	json.NewEncoder(writer).Encode(ret)
}

func (bh *BaseHandler) SetShortLink(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost { // кроме пост есть еще что то
		SendAnswer(writer, "POST only", http.StatusBadRequest, "")
		return
	}

	var ClientData s.URL

	err := json.NewDecoder(request.Body).Decode(&ClientData)
	if err != nil {
		SendAnswer(writer, "JSON decode error", http.StatusInternalServerError, "")
		return
	}

	res1 := f.ValidateLink(ClientData.URI)
	res2 := f.ValidateUserShortURL(ClientData.Short)
	if (res1 || res2) == false {
		SendAnswer(writer, "URL is not valid!", http.StatusBadRequest, "")
		return
	}

	if ClientData.Short == "" {
		ClientData.Short = f.CreateLink(request.RemoteAddr, ClientData.URI)
	}

	err = bh.Insert(ClientData.URI, ClientData.Short, bh.Db)

	if err != nil {
		SendAnswer(writer, "Same short link already ib DB", http.StatusBadRequest, "")
		return
	}
	SendAnswer(writer, "", http.StatusOK, request.Host+"/"+ClientData.Short)
}

func (bh *BaseHandler) RedirectHandler(writer http.ResponseWriter, request *http.Request, path string) {
	writer.Header().Set("Content-Type", "application/json")
	link, err := bh.Select(path, bh.Db)

	if err != nil {
		SendAnswer(writer, "DB select error", http.StatusInternalServerError, "")
		return
	}

	if link != "" {
		http.Redirect(writer, request, link, http.StatusFound)
	} else {
		SendAnswer(writer, "Not Found", http.StatusNotFound, "")
	}
}

func (bh *BaseHandler) Index(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		bh.RedirectHandler(writer, request, string([]rune(request.URL.Path)[1:]))
		return
	}
	fmt.Fprintf(writer, "Hello")
}
