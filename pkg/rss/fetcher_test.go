package rss

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func TestHTTPError_UserMessage(t *testing.T) {
	tests := []struct {
		name       string
		httpError  *HTTPError
		wantContains string
	}{
		{
			name: "NotFound 404",
			httpError: &HTTPError{
				StatusCode: http.StatusNotFound,
				Category:   ErrCategoryNotFound,
			},
			wantContains: "not found",
		},
		{
			name: "NotFound 410",
			httpError: &HTTPError{
				StatusCode: http.StatusGone,
				Category:   ErrCategoryNotFound,
			},
			wantContains: "removed",
		},
		{
			name: "AccessDenied 403",
			httpError: &HTTPError{
				StatusCode: http.StatusForbidden,
				Category:   ErrCategoryAccessDenied,
			},
			wantContains: "Access denied",
		},
		{
			name: "Timeout",
			httpError: &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryTimeout,
			},
			wantContains: "timed out",
		},
		{
			name: "ParseFailure",
			httpError: &HTTPError{
				StatusCode: 200,
				Category:   ErrCategoryParseFailure,
			},
			wantContains: "parse",
		},
		{
			name: "EmptyContent",
			httpError: &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryEmptyContent,
			},
			wantContains: "empty",
		},
		{
			name: "TooManyRedirects",
			httpError: &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryTooManyRedirects,
			},
			wantContains: "redirect",
		},
		{
			name: "ServiceUnavailable",
			httpError: &HTTPError{
				StatusCode: http.StatusServiceUnavailable,
				Category:   ErrCategoryServiceUnavailable,
			},
			wantContains: "unavailable",
		},
		{
			name: "ServerError",
			httpError: &HTTPError{
				StatusCode: http.StatusInternalServerError,
				Category:   ErrCategoryServerError,
			},
			wantContains: "Server error",
		},
		{
			name: "NetworkError",
			httpError: &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryNetworkError,
			},
			wantContains: "Network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.httpError.UserMessage()
			if !strings.Contains(strings.ToLower(msg), strings.ToLower(tt.wantContains)) {
				t.Errorf("UserMessage() = %q, want to contain %q", msg, tt.wantContains)
			}
		})
	}
}

func TestCategorizeStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantCategory string
	}{
		{"404 Not Found", http.StatusNotFound, ErrCategoryNotFound},
		{"410 Gone", http.StatusGone, ErrCategoryNotFound},
		{"403 Forbidden", http.StatusForbidden, ErrCategoryAccessDenied},
		{"503 Service Unavailable", http.StatusServiceUnavailable, ErrCategoryServiceUnavailable},
		{"500 Internal Server Error", http.StatusInternalServerError, ErrCategoryServerError},
		{"502 Bad Gateway", http.StatusBadGateway, ErrCategoryServerError},
		{"200 OK", http.StatusOK, ErrCategoryNone},
		{"301 Moved Permanently", http.StatusMovedPermanently, ErrCategoryNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := categorizeStatusCode(tt.statusCode)
			if got != tt.wantCategory {
				t.Errorf("categorizeStatusCode(%d) = %q, want %q", tt.statusCode, got, tt.wantCategory)
			}
		})
	}
}

func TestHTTPError_Error(t *testing.T) {
	tests := []struct {
		name      string
		httpError *HTTPError
		wantContains string
	}{
		{
			name: "with message",
			httpError: &HTTPError{
				StatusCode: http.StatusNotFound,
				Message:    "Custom error message",
			},
			wantContains: "Custom error message",
		},
		{
			name: "without message uses status text",
			httpError: &HTTPError{
				StatusCode: http.StatusNotFound,
			},
			wantContains: "Not Found",
		},
		{
			name: "zero status code",
			httpError: &HTTPError{
				StatusCode: 0,
				Message:    "Network error",
			},
			wantContains: "Network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.httpError.Error()
			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("Error() = %q, want to contain %q", got, tt.wantContains)
			}
		})
	}
}

func TestValidateFeedContent(t *testing.T) {
	tests := []struct {
		name    string
		feed    *gofeed.Feed
		wantErr bool
		errCategory string
	}{
		{
			name:    "nil feed",
			feed:    nil,
			wantErr: true,
			errCategory: ErrCategoryEmptyContent,
		},
		{
			name: "empty title",
			feed: &gofeed.Feed{
				Title: "",
			},
			wantErr: true,
			errCategory: ErrCategoryEmptyContent,
		},
		{
			name: "whitespace title",
			feed: &gofeed.Feed{
				Title: "   ",
			},
			wantErr: true,
			errCategory: ErrCategoryEmptyContent,
		},
		{
			name: "valid title",
			feed: &gofeed.Feed{
				Title: "Valid Feed Title",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFeedContent(tt.feed)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFeedContent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				if httpErr, ok := err.(*HTTPError); ok {
					if httpErr.Category != tt.errCategory {
						t.Errorf("ValidateFeedContent() error category = %q, want %q", httpErr.Category, tt.errCategory)
					}
				} else {
					t.Errorf("Expected *HTTPError, got %T", err)
				}
			}
		})
	}
}

func TestFetchAndParseFeed_NetworkError(t *testing.T) {
	// Test with invalid URL that will cause network error
	_, err := FetchAndParseFeed("http://invalid-url-that-does-not-exist.example")
	if err == nil {
		t.Fatal("Expected network error, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected *HTTPError, got %T", err)
	}

	if httpErr.Category != ErrCategoryNetworkError {
		t.Errorf("Expected category %q, got %q", ErrCategoryNetworkError, httpErr.Category)
	}
}

func TestFetchAndParseFeed_HTTPErrorNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := FetchAndParseFeed(server.URL)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected *HTTPError, got %T", err)
	}

	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, httpErr.StatusCode)
	}

	if httpErr.Category != ErrCategoryNotFound {
		t.Errorf("Expected category %q, got %q", ErrCategoryNotFound, httpErr.Category)
	}

	// Verify user message is helpful
	msg := httpErr.UserMessage()
	if !strings.Contains(strings.ToLower(msg), "not found") {
		t.Errorf("UserMessage() = %q, want to contain 'not found'", msg)
	}
}

func TestFetchAndParseFeed_HTTPErrorForbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	_, err := FetchAndParseFeed(server.URL)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected *HTTPError, got %T", err)
	}

	if httpErr.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, httpErr.StatusCode)
	}

	if httpErr.Category != ErrCategoryAccessDenied {
		t.Errorf("Expected category %q, got %q", ErrCategoryAccessDenied, httpErr.Category)
	}
}

func TestFetchAndParseFeed_HTTPErrorGone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGone)
	}))
	defer server.Close()

	_, err := FetchAndParseFeed(server.URL)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected *HTTPError, got %T", err)
	}

	if httpErr.StatusCode != http.StatusGone {
		t.Errorf("Expected status code %d, got %d", http.StatusGone, httpErr.StatusCode)
	}

	if httpErr.Category != ErrCategoryNotFound {
		t.Errorf("Expected category %q, got %q", ErrCategoryNotFound, httpErr.Category)
	}
}

func TestFetchAndParseFeed_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := FetchAndParseFeed(server.URL)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected *HTTPError, got %T", err)
	}

	if httpErr.Category != ErrCategoryServerError {
		t.Errorf("Expected category %q, got %q", ErrCategoryServerError, httpErr.Category)
	}
}

func TestFetchAndParseFeed_Timeout(t *testing.T) {
	// Create a server that delays response longer than timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Longer than we'd want to wait in a test
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a custom client with very short timeout for testing
	// Note: The actual FetchAndParseFeed uses 30s timeout, so we can't easily test
	// timeout without modifying the function. This test documents the behavior.
	t.Skip("Timeout test requires injectable client - documented behavior")
}

func TestIsTimeoutError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		want    bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "timeout in message",
			err:  &mockError{msg: "i/o timeout"},
			want: true,
		},
		{
			name: "timeout word in message",
			err:  &mockError{msg: "request timeout"},
			want: true,
		},
		{
			name: "no timeout",
			err:  &mockError{msg: "connection refused"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTimeoutError(tt.err)
			if got != tt.want {
				t.Errorf("isTimeoutError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchAndParseFeed_ExcessiveRedirects(t *testing.T) {
	// Create a server that redirects infinitely
	redirectCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectCount++
		if redirectCount > 15 {
			t.Fatal("Server redirected more than expected, limit not working")
		}
		http.Redirect(w, r, "/redirect-loop", http.StatusFound)
	}))
	defer server.Close()

	_, err := FetchAndParseFeed(server.URL + "/redirect-loop")
	if err == nil {
		t.Fatal("Expected error for excessive redirects, got nil")
	}

	// Verify it's categorized as TooManyRedirects
	if httpErr, ok := err.(*HTTPError); ok {
		if httpErr.Category != ErrCategoryTooManyRedirects {
			t.Errorf("Expected category %q, got %q", ErrCategoryTooManyRedirects, httpErr.Category)
		}
	}

	// Should have stopped after 10 redirects
	if redirectCount > 10 {
		t.Fatalf("Expected at most 10 redirects, got %d", redirectCount)
	}
}

func TestFetchAndParseFeed_SuccessfulRedirect(t *testing.T) {
	// Create a server that redirects once then serves valid RSS
	redirects := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/old-feed" {
			redirects++
			http.Redirect(w, r, "/new-feed", http.StatusMovedPermanently)
			return
		}
		// Serve minimal valid RSS
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<rss version="2.0">
<channel>
<title>Test Feed</title>
<item><title>Test Item</title></item>
</channel>
</rss>`))
	}))
	defer server.Close()

	feed, err := FetchAndParseFeed(server.URL + "/old-feed")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if redirects != 1 {
		t.Fatalf("Expected exactly 1 redirect, got %d", redirects)
	}

	if feed.Title != "Test Feed" {
		t.Fatalf("Expected feed title 'Test Feed', got '%s'", feed.Title)
	}
}

// mockError implements error interface for testing
type mockError struct {
	msg string
}

func (m *mockError) Error() string {
	return m.msg
}
