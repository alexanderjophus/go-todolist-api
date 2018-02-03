package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func newHelloServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/list", listHandler)
	return r
}

func TestMyHandler(t *testing.T) {
	var jsonStr = []byte(`{"title":"title of item","body":"this is a todo item"}`)
	req, _ := http.NewRequest("POST", "/list", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()
	newHelloServer().ServeHTTP(res, req)

	req, _ = http.NewRequest("GET", "/list", nil)
	res = httptest.NewRecorder()
	newHelloServer().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("expected 200 response code, got", res.Code)
	}

	todoList := []TodoItem{}
	json.Unmarshal(res.Body.Bytes(), &todoList)
	todoItem := todoList[0]

	if todoItem.TITLE != "title of item" {
		t.Error("expected title 'title of item', got", todoItem.TITLE)
	}

	if todoItem.BODY != "this is a todo item" {
		t.Error("expected body 'this is a todo item', got", todoItem.BODY)
	}

	if todoItem.COMPLETED != false {
		t.Error("expected item to not be completed, got", todoItem.COMPLETED)
	}
}
