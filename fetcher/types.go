package fetcher

type (
	// ID is resource identifier
	ID int
	// Interval is type used to store time duration, in seconds
	Interval float64

	// TaskRequest stores information used to create worker
	TaskRequest struct {
		URL      string `json:"url"`
		Interval `json:"interval"`
	}

	// TaskInfo stores information used to describe worker
	TaskInfo struct {
		TaskRequest
		ID `json:"id"`
	}

	// Result stores data downloaded by worker
	Result struct {
		Response  *string `json:"response"`
		Duration  float64 `json:"duration"`
		CreatedAt float64 `json:"created_at"`
	}
)
