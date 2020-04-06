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
	Fetcher
}

func (w *worker) run() {
	ticker := time.NewTicker(time.Second * time.Duration(w.Interval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("tick\n")
		case <-w.finish:
			return
		case interval := <-w.update:
			ticker.Stop()
			w.Interval = interval
			ticker = time.NewTicker(time.Second * time.Duration(w.Interval))
		}
	}
}

func (w *worker) stop() {
	close(w.finish)
}

func (w *worker) updateInterval(i Interval) {
	w.update <- i
}

var (
	workers      = make(map[Id]*worker)
	results      = make(map[Id]([]Result))
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
			finish:  make(chan struct{}),
			update:  make(chan Interval),
			Fetcher: Fetcher{Task: t, Id: id},
		}

		workers[id] = w
		results[id] = []Result{}
		urlToId[t.Url] = id

		go w.run()
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
	delete(results, id)
	delete(urlToId, w.Url)

	return nil
}

func List() []Fetcher {
	return nil
}

func History(id Id) ([]Result, error) {
	mutex.Lock()
	defer mutex.Unlock()
	return nil, nil
}
