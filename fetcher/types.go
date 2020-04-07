package fetcher

type (
	// ID is resource identifier
	ID int
	// Interval is type used to store time duration, in seconds
	Interval float64

	// Task stores information used to create worker
	Task struct {
		URL      string `json:"url"`
		Interval `json:"interval"`
	}

	// Fetcher stores information used to describe worker
	Fetcher struct {
		Task
		ID `json:"id"`
	}

	// Result stores data downloaded by worker
	Result struct {
		Response  *string `json:"response"`
		Duration  float64 `json:"duration"`
		CreatedAt float64 `json:"created_at"`
	}
)
