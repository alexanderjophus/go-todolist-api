package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/etherlabsio/healthcheck"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/satori/go.uuid"
)

// TodoItem represents the items in the todolist
type TodoItem struct {
	ID        uuid.UUID `json:"id"`
	TITLE     string    `json:"title"`
	BODY      string    `json:"body"`
	COMPLETED bool      `json:"completed"`
}

var c = cache.New(5*time.Minute, 10*time.Minute)

func listHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		items := make([]TodoItem, c.ItemCount())
		i := 0
		for _, v := range c.Items() {
			items[i] = v.Object.(TodoItem)
			i++
		}
		val, _ := json.Marshal(items)
		io.WriteString(w, string(val))

	case "POST":
		var todoItem TodoItem
		json.NewDecoder(r.Body).Decode(&todoItem)
		ID, _ := uuid.NewV4()
		todoItem.ID = ID
		c.Add(ID.String(), todoItem, cache.NoExpiration)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/list", listHandler).Methods("POST", "GET")
	r.Handle("/health", healthcheck.Handler(
		healthcheck.WithTimeout(2*time.Second),
		healthcheck.WithChecker(
			"database", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					// need to find way of actually testing databases health
					return nil
				},
			)),
	))

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Println("Server running")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
