package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

// ToDoItem is something
type ToDoItem struct {
	Title string
	Body  string
}

var items []ToDoItem

var c = cache.New(5*time.Minute, 10*time.Minute)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	log.Println("Healthcheck call")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		strs := make([]string, len(items))
		for i, v := range items {
			strs[i] = v.Title
		}
		// io.WriteString(w, strings.Join(strs, ", "))
		d := c.Items()
		for _, v := range d {
			io.WriteString(w, v.Object.(string))
		}

	case "POST":
		var toDoItem ToDoItem
		json.NewDecoder(r.Body).Decode(&toDoItem)
		items = append(items, toDoItem)
		c.Add(toDoItem.Title, toDoItem.Body, cache.NoExpiration)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/list", listHandler).Methods("POST", "GET")
	r.HandleFunc("/health", healthCheckHandler)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
