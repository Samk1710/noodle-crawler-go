package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Samk1710/noodle-crawler-go/internal/scheduler"
	"github.com/Samk1710/noodle-crawler-go/internal/worker"
	"github.com/Samk1710/noodle-crawler-go/pkg/models"
	"github.com/Samk1710/noodle-crawler-go/pkg/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: crawler <url>")
		return
	}

	startURL := os.Args[1]

	jobs := make(chan models.Job, 100)
	results := make(chan models.Result, 100)

	visited := scheduler.NewVisited()

	var wg sync.WaitGroup

	workerCount := 5
	maxDepth := 2

	// start workers
	for i := 0; i < workerCount; i++ {
		go worker.StartWorker(i, jobs, results)
	}

	// aggregation storage
	var allResults []models.Result
	var mu sync.Mutex

	// result processor (aggregator)
	go func() {
		for result := range results {

			mu.Lock()
			allResults = append(allResults, result)
			mu.Unlock()

			if result.Err != nil {
				fmt.Println("Error:", result.URL, result.Err)
				wg.Done()
				continue
			}

			fmt.Println("Visited:", result.URL)

			if result.Depth >= maxDepth {
				wg.Done()
				continue
			}

			for _, link := range result.Links {

				normalized, err := utils.NormalizeURL(result.URL, link)
				if err != nil || normalized == "" {
					continue
				}

				if !utils.IsSameDomain(startURL, normalized) {
					continue
				}

				if !visited.Add(normalized) {
					continue
				}

				wg.Add(1)
				jobs <- models.Job{
					URL:   normalized,
					Depth: result.Depth + 1,
				}
			}

			wg.Done()
		}
	}()

	// seed
	visited.Add(startURL)
	wg.Add(1)
	jobs <- models.Job{URL: startURL, Depth: 0}

	// wait for crawl completion
	wg.Wait()

	close(jobs)
	close(results)

	// output JSON
	output, _ := json.MarshalIndent(allResults, "", "  ")
	fmt.Println(string(output))
}
