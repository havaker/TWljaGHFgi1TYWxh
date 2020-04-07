package fetcher

import (
	"errors"
	"sync"
)

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

	id, found := urlToID[t.URL]
	if !found {
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
		db[id].Interval = t.Interval
	}

	return id
}

// Remove deletes worker with given id, and state associated with it
func Remove(id ID) error {
	mutex.Lock()
	defer mutex.Unlock()

	w, found := workers[id]
	if !found {
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
// returned slice is read only
func History(id ID) ([]Result, error) {
	mutex.Lock()
	defer mutex.Unlock()

	s, found := db[id]
	if !found {
		return nil, ErrNotFound
	}

	return s.results, nil
}
