package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/havaker/TWljaGHFgi1TYWxh/fetcher"
	"github.com/havaker/TWljaGHFgi1TYWxh/middleware"
)

func createHandler(w http.ResponseWriter, req *http.Request) {
	var t fetcher.Task
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)

	if err != nil {
		log.Printf("%s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if t.Interval == 0 || t.Url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := fetcher.Save(t)

	fmt.Fprintf(w, "{\"id\": %d}\n", id)
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
}

func listHandler(w http.ResponseWriter, req *http.Request) {
}

func historyHandler(w http.ResponseWriter, req *http.Request) {
}

const MB = 1000 * 1000

func main() {
	r := chi.NewRouter()
	r.Use(middleware.BodyLimit(MB, http.MethodPost))
	r.Route("/api/fetcher", func(r chi.Router) {
		r.Post("/", createHandler)
		r.Get("/", listHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", deleteHandler)
			r.Get("/history", historyHandler)
		})
	})

	http.ListenAndServe(":8080", r)
}
