package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	if t.Interval == 0 || t.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := fetcher.Save(t)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"id\": %d}\n", id)
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := fetcher.Remove(fetcher.ID(id)); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func listHandler(w http.ResponseWriter, req *http.Request) {
	fetchers := fetcher.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fetchers)
}

func historyHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	history, err := fetcher.History(fetcher.ID(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

const mb = 1000 * 1000

func main() {
	r := chi.NewRouter()
	r.Use(middleware.BodyLimit(mb, http.MethodPost))
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
