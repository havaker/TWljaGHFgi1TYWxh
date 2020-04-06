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

func (s *state) downloadUrl(c http.Client) *string {
	resp, err := c.Get(s.Url)
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
	data := s.downloadUrl(client)
	end := time.Now()

	res.Response = data
	res.CreatedAt = float64(end.UnixNano()) / 1e9
	res.Duration = end.Sub(start).Seconds()

	mutex.Lock()
	defer mutex.Unlock()
	s.results = append(s.results, res)
}

var (
	workers      = make(map[Id]*worker)
	db           = make(map[Id]*state)
	urlToId      = make(map[string]Id)
	available Id = 1
	mutex        = &sync.Mutex{}
)

var NotFoundErr = errors.New("Resource not found")

func Save(t Task) Id {
	mutex.Lock()
	defer mutex.Unlock()

	id, ok := urlToId[t.Url]
	if !ok {
		id = available
		available++

		w := &worker{
			finish: make(chan struct{}),
			update: make(chan Interval),
		}

		s := &state{
			Fetcher: Fetcher{Task: t, Id: id},
			results: []Result{},
		}

		workers[id] = w
		db[id] = s
		urlToId[t.Url] = id

		go w.run(s.Interval, s)
	} else {
		workers[id].updateInterval(t.Interval)
	}

	return id
}

func Remove(id Id) error {
	mutex.Lock()
	defer mutex.Unlock()

	w, ok := workers[id]
	if !ok {
		return NotFoundErr
	}

	w.stop()

	delete(workers, id)
	delete(urlToId, db[id].Url)
	delete(db, id)

	return nil
}

func List() []Fetcher {
	mutex.Lock()
	defer mutex.Unlock()

	f := []Fetcher{}
	for _, s := range db {
		f = append(f, s.Fetcher)
	}

	return f
}

func History(id Id) ([]Result, error) {
	mutex.Lock()
	defer mutex.Unlock()

	s, ok := db[id]
	if !ok {
		return nil, NotFoundErr
	}

	// TODO copy?
	return s.results, nil
}
