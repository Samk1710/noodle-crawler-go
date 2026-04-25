package worker

import (
	"github.com/Samk1710/noodle-crawler-go/internal/fetcher"
	"github.com/Samk1710/noodle-crawler-go/internal/parser"
	"github.com/Samk1710/noodle-crawler-go/pkg/models"
)

func StartWorker(id int, jobs <-chan models.Job, results chan<- models.Result) {
	for job := range jobs {
		body, status, err := fetcher.Fetch(job.URL)

		if err != nil {
			results <- models.Result{
				URL:   job.URL,
				Err:   err,
				Depth: job.Depth,
			}
			continue
		}

		links, err := parser.ExtractLinks(body)

		results <- models.Result{
			URL:        job.URL,
			StatusCode: status,
			Links:      links,
			Err:        err,
			Depth:      job.Depth,
		}
	}
}
