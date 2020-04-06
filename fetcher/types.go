package fetcher

import "time"

type (
	Id       int
	Interval float64

	Task struct {
		Url      string `json:"url"`
		Interval `json:"interval"`
	}

	Fetcher struct {
		Task
		Id `json:"id"`
	}

	Result struct {
		Response  string
		Duration  time.Duration
		CreatedAt time.Time
	}
)
