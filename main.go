package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	chimid "github.com/go-chi/chi/middleware"
	"github.com/havaker/TWljaGHFgi1TYWxh/fetcher"
	"github.com/havaker/TWljaGHFgi1TYWxh/middleware"
)

var fetch *fetcher.Fetcher

func createHandler(w http.ResponseWriter, req *http.Request) {
	var t fetcher.TaskRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)

	if err != nil {
		log.Printf("createHandler: %s\n", err.Error())
		if err == middleware.ErrBodyLimitExceeded {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	if t.Interval <= 0 || t.URL == "" {
		log.Printf("createHandler: invalid object\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := fetch.Save(t)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"id\": %d}\n", id)
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		log.Printf("deleteHandler: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := fetch.Remove(fetcher.ID(id)); err != nil {
		log.Printf("deleteHandler: %s\n", err.Error())
		w.WriteHeader(http.StatusNotFound)
	}
}

func listHandler(w http.ResponseWriter, req *http.Request) {
	tasks := fetch.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func historyHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		log.Printf("historyHandler: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	history, err := fetch.History(fetcher.ID(id))
	if err != nil {
		log.Printf("historyHandler: %s\n", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func main() {
	fetch = fetcher.NewFetcher()

	r := chi.NewRouter()

	MB := 1000 * 1000
	r.Use(middleware.BodyLimit(MB, http.MethodPost))
	r.Use(chimid.Logger)

	r.Route("/api/fetcher", func(r chi.Router) {
		r.Post("/", createHandler)
		r.Get("/", listHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", deleteHandler)
			r.Get("/history", historyHandler)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
