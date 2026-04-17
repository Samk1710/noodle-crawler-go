package models

type Page struct {
	URL        string
	StatusCode int
	Body       string
	Links      []string
}