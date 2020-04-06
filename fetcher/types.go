package fetcher

import "time"

type (
	Id       int
	Interval float64

	Task struct {
		Url string
		Interval
	}

	Fetcher struct {
		Task
		Id
	}

	Result struct {
		Response  string
		Duration  time.Duration
		CreatedAt time.Time
	}
)
