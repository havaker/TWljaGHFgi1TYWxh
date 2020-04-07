package fetcher

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type state struct {
	Fetcher
	results []Result
}

func (s *state) downloadURL(c http.Client) *string {
	resp, err := c.Get(s.URL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	str := string(data)
	return &str
}

func (s *state) fetch() {
	res := Result{}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	data := s.downloadURL(client)
	end := time.Now()

	res.Response = data
	res.CreatedAt = float64(end.UnixNano()) / 1e9
	res.Duration = end.Sub(start).Seconds()

	mutex.Lock()
	defer mutex.Unlock()
	s.results = append(s.results, res)
}

var (
	workers      = make(map[ID]*worker)
	db           = make(map[ID]*state)
	urlToID      = make(map[string]ID)
	available ID = 1
	mutex        = &sync.Mutex{}
)

// ErrNotFound is an error returned when resource could not be found
var ErrNotFound = errors.New("Resource not found")

// Save starts or updates worker, which processes given Task
func Save(t Task) ID {
	mutex.Lock()
	defer mutex.Unlock()

	id, ok := urlToID[t.URL]
	if !ok {
		id = available
		available++

		w := &worker{
			finish: make(chan struct{}),
			update: make(chan Interval),
		}

		s := &state{
			Fetcher: Fetcher{Task: t, ID: id},
			results: []Result{},
		}

		workers[id] = w
		db[id] = s
		urlToID[t.URL] = id

		go w.run(s.Interval, s)
	} else {
		workers[id].updateInterval(t.Interval)
	}

	return id
}

// Remove deletes worker with given id, and state associated with it
func Remove(id ID) error {
	mutex.Lock()
	defer mutex.Unlock()

	w, ok := workers[id]
	if !ok {
		return ErrNotFound
	}

	w.stop()

	delete(workers, id)
	delete(urlToID, db[id].URL)
	delete(db, id)

	return nil
}

// List returns list of all workers
func List() []Fetcher {
	mutex.Lock()
	defer mutex.Unlock()

	f := []Fetcher{}
	for _, s := range db {
		f = append(f, s.Fetcher)
	}

	return f
}

// History returns data downloaded by worker identified by id
func History(id ID) ([]Result, error) {
	mutex.Lock()
	defer mutex.Unlock()

	s, ok := db[id]
	if !ok {
		return nil, ErrNotFound
	}

	// TODO copy?
	return s.results, nil
}
