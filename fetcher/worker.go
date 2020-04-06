package fetcher

import "time"

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
			go s.fetch()
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
