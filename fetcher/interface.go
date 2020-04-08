package fetcher

import (
	"errors"
	"sync"
)

// Fetcher holds info about state and associated workers
type Fetcher struct {
	workers   map[ID]*worker
	db        map[ID]*state
	urlToID   map[string]ID
	available ID
	mutex     *sync.Mutex
}

// NewFetcher returns new Fetcher object
func NewFetcher() *Fetcher {
	return &Fetcher{
		workers:   make(map[ID]*worker),
		db:        make(map[ID]*state),
		urlToID:   make(map[string]ID),
		available: 1,
		mutex:     &sync.Mutex{},
	}
}

// ErrNotFound is an error returned when resource could not be found
var ErrNotFound = errors.New("Resource not found")

// Save starts or updates worker, which processes given TaskRequest
func (f *Fetcher) Save(t TaskRequest) ID {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	id, found := f.urlToID[t.URL]
	if !found {
		id = f.available
		f.available++

		w := &worker{
			finish: make(chan struct{}),
			update: make(chan Interval),
		}

		s := &state{
			TaskInfo: TaskInfo{TaskRequest: t, ID: id},
			results:  []Result{},
			mutex:    f.mutex,
		}

		f.workers[id] = w
		f.db[id] = s
		f.urlToID[t.URL] = id

		go w.run(s.Interval, s)
	} else {
		f.workers[id].updateInterval(t.Interval)
		f.db[id].Interval = t.Interval
	}

	return id
}

// Remove deletes worker with given id, and state associated with it
func (f *Fetcher) Remove(id ID) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	w, found := f.workers[id]
	if !found {
		return ErrNotFound
	}

	w.stop()

	delete(f.workers, id)
	delete(f.urlToID, f.db[id].URL)
	delete(f.db, id)

	return nil
}

// List returns list of all workers
func (f *Fetcher) List() []TaskInfo {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	t := []TaskInfo{}
	for _, s := range f.db {
		t = append(t, s.TaskInfo)
	}

	return t
}

// History returns data downloaded by worker identified by id
// returned slice is read only
func (f *Fetcher) History(id ID) ([]Result, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	s, found := f.db[id]
	if !found {
		return nil, ErrNotFound
	}

	return s.results, nil
}

// Shutdown stops all workers associated with f
func (f *Fetcher) Shutdown() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	for _, w := range f.workers {
		w.stop()
	}
}
