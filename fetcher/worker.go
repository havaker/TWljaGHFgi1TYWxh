package fetcher

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type worker struct {
	finish chan struct{}
	update chan Interval
}

func (w *worker) run(initial Interval, s *state) {
	ticker := time.NewTicker(time.Second * time.Duration(initial))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("tick\n")
		case <-w.finish:
			return
		case interval := <-w.update:
			ticker.Stop()
			ticker = time.NewTicker(time.Second * time.Duration(interval))
		}
	}
}

func (w *worker) stop() {
	close(w.finish)
}

func (w *worker) updateInterval(i Interval) {
	w.update <- i
}

type state struct {
	Fetcher
	results []Result
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
	return nil, nil
}
