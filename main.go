package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/havaker/TWljaGHFgi1TYWxh/middleware"
)

type Task struct {
	Url      string
	Interval int64
}

const MB = 1000 * 1000

func test(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	for {
		var t Task
		err := decoder.Decode(&t)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("%s\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if t.Interval == 0 || t.Url == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("%v\n", t)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.BodyLimit(MB))
	r.Post("/", test)
	http.ListenAndServe(":8080", r)
}
