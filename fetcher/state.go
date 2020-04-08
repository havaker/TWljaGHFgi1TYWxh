package fetcher

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type state struct {
	TaskInfo
	results []Result
	mutex   *sync.Mutex
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

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.results = append(s.results, res)
}
