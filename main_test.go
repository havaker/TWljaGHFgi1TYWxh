package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/havaker/TWljaGHFgi1TYWxh/fetcher"
)

func TestCreateHandler(t *testing.T) {
	fetch = fetcher.NewFetcher() // create new mock fetcher state

	str := []byte(`{"url":"https://httpbin.org/range/15","interval":60}`)

	req, err := http.NewRequest("POST", "api/fetcher", bytes.NewBuffer(str))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(createHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: %d", rec.Code)
	}

	expected := "{\"id\": 1}\n"
	actual := rec.Body.String()

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}

	fetch.Shutdown()
}

func TestListHandler(t *testing.T) {
	fetch = fetcher.NewFetcher() // create new mock fetcher state
	taskRequest := fetcher.TaskRequest{
		URL:      "https://httpbin.org/range/15",
		Interval: 60,
	}
	fetch.Save(taskRequest)

	req, err := http.NewRequest("GET", "api/fetcher", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(listHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: %d", rec.Code)
	}

	expected := []fetcher.TaskInfo{{
		TaskRequest: taskRequest,
		ID:          1,
	}}

	var actual []fetcher.TaskInfo
	err = json.Unmarshal([]byte(rec.Body.String()), &actual)
	if err != nil {
		t.Errorf("handler returned wrong json: %v", rec.Body.String())
	}

	if len(actual) != len(expected) || actual[0] != expected[0] {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}

	fetch.Shutdown()
}
