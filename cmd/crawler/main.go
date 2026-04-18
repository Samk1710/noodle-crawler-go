package main

import (
	"fmt"
	"os"

	"noodle-crawler-go/internal/worker"
	"noodle-crawler-go/pkg/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: crawler <url>")
		return
	}

	startURL := os.Args[1]

	jobs := make(chan models.Job, 100)
	results := make(chan models.Result, 100)

	visited := make(map[string]bool)

	workerCount := 5

	// start workers
	for i := 0; i < workerCount; i++ {
		go worker.StartWorker(i, jobs, results)
	}

	// seed first job
	jobs <- models.Job{URL: startURL, Depth: 0}
	visited[startURL] = true

	// maxDepth := 2

	for {
		result := <-results

		if result.Err != nil {
			fmt.Println("Error:", result.URL, result.Err)
			continue
		}

		fmt.Println("Visited:", result.URL)

		// enqueue new links
		for _, link := range result.Links {

			if visited[link] {
				continue
			}

			visited[link] = true

			jobs <- models.Job{
				URL:   link,
				Depth: 1, // simplified for now
			}
		}

		// stop condition (temporary)
		if len(visited) > 20 {
			break
		}
	}
}