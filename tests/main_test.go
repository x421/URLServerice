package tests

import (
	c "LinksService/controllers"
	f "LinksService/functions"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTest(t *testing.T) {
	ip, url:="1", "2"
	link1:=f.CreateLink(ip, url)
	link2:=f.CreateLink(ip, url)

	if link1 != link2 {
		t.Errorf("Not equal!")
	}

}

func TestValidateLink(t *testing.T) {
	link:="http://google.com"
	res:=f.ValidateLink(link)

	if res == false {
		t.Errorf("http://google.com isnt valid!")
	}

	link="htp://google.com"
	res=f.ValidateLink(link)

	if res == true {
		t.Errorf("HTP://google.com is valid!")
	}

	link="http://google."
	res=f.ValidateLink(link)

	if res == true {
		t.Errorf("HTP://google. is valid!")
	}
}

func TestShortLink(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/setShort", bytes.NewBuffer([]byte(`{"URI":"http://google.com", "Short":"tgdfvc"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-type", "application/json")
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.SetShortLink)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"URL": "`+req.Host+`/tgdfvc"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}