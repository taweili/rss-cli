package rss

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// Error categories for RSS feed fetching
const (
	ErrCategoryNone          = ""
	ErrCategoryNotFound      = "NotFound"
	ErrCategoryAccessDenied  = "AccessDenied"
	ErrCategoryTimeout       = "Timeout"
	ErrCategoryParseFailure  = "ParseFailure"
	ErrCategoryEmptyContent  = "EmptyContent"
	ErrCategoryTooManyRedirects = "TooManyRedirects"
	ErrCategoryServiceUnavailable = "ServiceUnavailable"
	ErrCategoryServerError   = "ServerError"
	ErrCategoryNetworkError  = "NetworkError"
)

// FetchAndParseFeed fetches an RSS/Atom feed from the given URL and parses it
func FetchAndParseFeed(url string) (*gofeed.Feed, error) {
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
		return nil, err
	}

	req.Header.Set("User-Agent", "rss-cli/1.0")

	resp, err := client.Do(req)
	if err != nil {
		// Check if this is already a categorized HTTP error (e.g., too many redirects)
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}
		// Check for timeout errors
		if isTimeoutError(err) {
			return nil, &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryTimeout,
				Message:    fmt.Sprintf("request timed out: %v", err),
			}
		}
		// Network errors (DNS, connection refused, etc.)
		return nil, &HTTPError{
			StatusCode: 0,
			Category:   ErrCategoryNetworkError,
			Message:    fmt.Sprintf("network error: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, httpErr(resp.StatusCode)
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	if err != nil {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Category:   ErrCategoryParseFailure,
			Message:    fmt.Sprintf("failed to parse feed: %v", err),
		}
	}

	// Validate feed content
	if err := ValidateFeedContent(feed); err != nil {
		return nil, err
	}

	return feed, nil
}

// ValidateFeedContent checks if a parsed feed has valid content
func ValidateFeedContent(feed *gofeed.Feed) error {
	if feed == nil {
		return &HTTPError{
			StatusCode: 0,
			Category:   ErrCategoryEmptyContent,
			Message:    "feed is nil",
		}
	}

	if strings.TrimSpace(feed.Title) == "" {
		return &HTTPError{
			StatusCode: 0,
			Category:   ErrCategoryEmptyContent,
			Message:    "feed has no title",
		}
	}

	return nil
}

// isTimeoutError checks if an error is a timeout error
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	var netErr interface{ Timeout() bool }
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	// Check for timeout in error message as fallback
	errStr := err.Error()
	return strings.Contains(errStr, "timeout") || strings.Contains(errStr, "i/o timeout")
}

func httpErr(statusCode int) error {
	category := categorizeStatusCode(statusCode)
	return &HTTPError{
		StatusCode: statusCode,
		Category:   category,
		Message:    http.StatusText(statusCode),
	}
}

// categorizeStatusCode maps HTTP status codes to error categories
func categorizeStatusCode(statusCode int) string {
	switch statusCode {
	case http.StatusNotFound, http.StatusGone:
		return ErrCategoryNotFound
	case http.StatusForbidden:
		return ErrCategoryAccessDenied
	case http.StatusServiceUnavailable:
		return ErrCategoryServiceUnavailable
	default:
		if statusCode >= 500 {
			return ErrCategoryServerError
		}
		return ErrCategoryNone
	}
}

type HTTPError struct {
	StatusCode int
	Category   string
	Message    string
}

func (e *HTTPError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return http.StatusText(e.StatusCode)
}

// UserMessage returns a user-friendly error message
func (e *HTTPError) UserMessage() string {
	switch e.Category {
	case ErrCategoryNotFound:
		return "Feed not found. The URL may be incorrect or the feed has been removed."
	case ErrCategoryAccessDenied:
		return "Access denied. The feed may require authentication or has blocked access."
	case ErrCategoryTimeout:
		return "Request timed out. The server took too long to respond."
	case ErrCategoryParseFailure:
		return "Failed to parse feed. The feed format may be invalid or corrupted."
	case ErrCategoryEmptyContent:
		return "Feed is empty or has no valid content."
	case ErrCategoryTooManyRedirects:
		return "Too many redirects. The feed URL may be in a redirect loop."
	case ErrCategoryServiceUnavailable:
		return "Service temporarily unavailable. Please try again later."
	case ErrCategoryServerError:
		return "Server error. The feed server encountered an internal error."
	case ErrCategoryNetworkError:
		return "Network error. Unable to connect to the feed server."
	default:
		if e.StatusCode > 0 {
			return fmt.Sprintf("HTTP error %d: %s", e.StatusCode, e.Message)
		}
		return e.Message
	}
}
