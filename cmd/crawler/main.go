package main

import (
	"fmt"
	"os"

	"noodle-crawler-go/internal/fetcher"
	"noodle-crawler-go/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: crawler <url>")
		return
	}

	url := os.Args[1]

	body, status, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println("Fetch error:", err)
		return
	}

	links, err := parser.ExtractLinks(body)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	fmt.Println("URL:", url)
	fmt.Println("Status:", status)
	fmt.Println("Links found:", len(links))

	for _, link := range links {
		fmt.Println("-", link)
	}
}