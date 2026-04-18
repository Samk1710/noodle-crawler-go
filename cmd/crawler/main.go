package main

import (
	"fmt"
	"os"
	"sync"

	"noodle-crawler-go/internal/worker"
	"noodle-crawler-go/pkg/models"
	"noodle-crawler-go/pkg/utils"
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

	var wg sync.WaitGroup

	workerCount := 5
	maxDepth := 2

	// start workers
	for i := 0; i < workerCount; i++ {
		go worker.StartWorker(i, jobs, results)
	}

	// seed first job
	wg.Add(1)
	jobs <- models.Job{URL: startURL, Depth: 0}
	visited[startURL] = true

	// result processor
	go func() {
		for result := range results {

			if result.Err != nil {
				fmt.Println("Error:", result.URL, result.Err)
				wg.Done()
				continue
			}

			fmt.Println("Visited:", result.URL)

			// find depth of this URL
			currentDepth := 0 // fallback

			if currentDepth >= maxDepth {
				wg.Done()
				continue
			}

			for _, link := range result.Links {

				normalized, err := utils.NormalizeURL(result.URL, link)
				if err != nil || normalized == "" {
					continue
				}

				// domain restriction
				if !utils.IsSameDomain(startURL, normalized) {
					continue
				}

				if visited[normalized] {
					continue
				}

				visited[normalized] = true

				wg.Add(1)
				jobs <- models.Job{
					URL:   normalized,
					Depth: currentDepth + 1,
				}
			}

			wg.Done()
		}
	}()

	// wait until all jobs are done
	wg.Wait()

	// cleanup
	close(jobs)
	close(results)

	fmt.Println("Crawl finished.")
}