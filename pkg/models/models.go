package models

type Page struct {
	URL        string
	StatusCode int
	Body       string
	Links      []string
}

type Job struct {
	URL   string
	Depth int
}

type Result struct {
	URL        string
	StatusCode int
	Links      []string
	Err        error
}