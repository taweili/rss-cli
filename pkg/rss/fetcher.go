package rss

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

// FetchAndParseFeed fetches an RSS/Atom feed from the given URL and parses it
func FetchAndParseFeed(url string) (*gofeed.Feed, error) {
	// Create an http client with a reasonable timeout and redirect limit
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	// Add User Agent header to avoid being blocked by servers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "rss-cli/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, httpErr(resp.StatusCode)
	}

	fp := gofeed.NewParser()
	return fp.Parse(resp.Body)
}

func httpErr(statusCode int) error {
	return &HTTPError{StatusCode: statusCode}
}

type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	return http.StatusText(e.StatusCode)
}
