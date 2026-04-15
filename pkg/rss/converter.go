package rss

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown"
)

// FetchAndConvertArticle fetches HTML from a URL and converts it to markdown
func FetchAndConvertArticle(url string) (string, error) {
	// Fetch HTML content using the same pattern as FetchArticleContent
	html, err := fetchHTML(url)
	if err != nil {
		return "", err
	}

	// Validate HTML response is not empty
	if strings.TrimSpace(html) == "" {
		return "", &HTTPError{
			StatusCode: 0,
			Category:   ErrCategoryEmptyContent,
			Message:    "HTML content is empty",
		}
	}

	// Convert HTML to markdown
	md, err := htmltomarkdown.Convert(html)
	if err != nil {
		return "", &HTTPError{
			StatusCode: 0,
			Category:   ErrCategoryConversionFailure,
			Message:    fmt.Sprintf("failed to convert HTML to markdown: %v", err),
		}
	}

	// Validate converted markdown is not empty
	if strings.TrimSpace(md) == "" {
		return "", &HTTPError{
			StatusCode: 0,
			Category:   ErrCategoryConversionFailure,
			Message:    "converted markdown is empty",
		}
	}

	return md, nil
}

// fetchHTML fetches HTML content from a URL with timeout and redirect limits
func fetchHTML(url string) (string, error) {
	// Create an http client with a reasonable timeout and redirect limit
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return &HTTPError{
					StatusCode: 0,
					Category:   ErrCategoryTooManyRedirects,
					Message:    "stopped after 10 redirects",
				}
			}
			return nil
		},
	}

	// Add User Agent header to avoid being blocked by servers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "rss-cli/1.0")

	resp, err := client.Do(req)
	if err != nil {
		// Check if this is already a categorized HTTP error (e.g., too many redirects)
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			return "", httpErr
		}
		// Check for timeout errors
		if isTimeoutError(err) {
			return "", &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryTimeout,
				Message:    fmt.Sprintf("request timed out: %v", err),
			}
		}
		// Network errors (DNS, connection refused, etc.)
		return "", &HTTPError{
			StatusCode: 0,
			Category:   ErrCategoryNetworkError,
			Message:    fmt.Sprintf("network error: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", httpErr(resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
