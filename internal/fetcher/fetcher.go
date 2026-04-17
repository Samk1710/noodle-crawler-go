package fetcher

import (
	"io"
	"net/http"
	"time"
)

func Fetch(url string) (string, int, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}

	return string(bodyBytes), resp.StatusCode, nil
}