package fetcher

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
		Response  *string `json:"response"`
		Duration  float64 `json:"duration"`
		CreatedAt float64 `json:"created_at"`
	}
)
