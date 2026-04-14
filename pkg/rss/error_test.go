package rss

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHTTPError_Categorization tests all HTTP status code categorizations
func TestHTTPError_Categorization(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		wantCategory string
		wantMessage  string
	}{
		{
			name:         "404 Not Found",
			statusCode:   http.StatusNotFound,
			wantCategory: ErrCategoryNotFound,
			wantMessage:  "Not Found",
		},
		{
			name:         "410 Gone",
			statusCode:   http.StatusGone,
			wantCategory: ErrCategoryNotFound,
			wantMessage:  "Gone",
		},
		{
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			wantCategory: ErrCategoryAccessDenied,
			wantMessage:  "Forbidden",
		},
		{
			name:         "503 Service Unavailable",
			statusCode:   http.StatusServiceUnavailable,
			wantCategory: ErrCategoryServiceUnavailable,
			wantMessage:  "Service Unavailable",
		},
		{
			name:         "500 Internal Server Error",
			statusCode:   http.StatusInternalServerError,
			wantCategory: ErrCategoryServerError,
			wantMessage:  "Internal Server Error",
		},
		{
			name:         "502 Bad Gateway",
			statusCode:   http.StatusBadGateway,
			wantCategory: ErrCategoryServerError,
			wantMessage:  "Bad Gateway",
		},
		{
			name:         "504 Gateway Timeout",
			statusCode:   http.StatusGatewayTimeout,
			wantCategory: ErrCategoryServerError,
			wantMessage:  "Gateway Timeout",
		},
		{
			name:         "200 OK - no category",
			statusCode:   http.StatusOK,
			wantCategory: ErrCategoryNone,
			wantMessage:  "OK",
		},
		{
			name:         "301 Moved Permanently - no category",
			statusCode:   http.StatusMovedPermanently,
			wantCategory: ErrCategoryNone,
			wantMessage:  "Moved Permanently",
		},
		{
			name:         "400 Bad Request - no specific category",
			statusCode:   http.StatusBadRequest,
			wantCategory: ErrCategoryNone,
			wantMessage:  "Bad Request",
		},
		{
			name:         "401 Unauthorized - no specific category",
			statusCode:   http.StatusUnauthorized,
			wantCategory: ErrCategoryNone,
			wantMessage:  "Unauthorized",
		},
		{
			name:         "408 Request Timeout - no specific category",
			statusCode:   http.StatusRequestTimeout,
			wantCategory: ErrCategoryNone,
			wantMessage:  "Request Timeout",
		},
		{
			name:         "429 Too Many Requests - no specific category",
			statusCode:   http.StatusTooManyRequests,
			wantCategory: ErrCategoryNone,
			wantMessage:  "Too Many Requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := httpErr(tt.statusCode)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			httpErr, ok := err.(*HTTPError)
			if !ok {
				t.Fatalf("Expected *HTTPError, got %T", err)
			}

			if httpErr.Category != tt.wantCategory {
				t.Errorf("Category = %q, want %q", httpErr.Category, tt.wantCategory)
			}

			if httpErr.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", httpErr.StatusCode, tt.statusCode)
			}

			if !strings.Contains(httpErr.Message, tt.wantMessage) {
				t.Errorf("Message = %q, want to contain %q", httpErr.Message, tt.wantMessage)
			}
		})
	}
}

// TestFetchArticleContent_EmptyResponse tests empty response handling
func TestFetchArticleContent_EmptyResponse(t *testing.T) {
	tests := []struct {
		name        string
		responseBody string
		contentType string
		wantErr     bool
		errCategory string
	}{
		{
			name:        "completely empty response",
			responseBody: "",
			contentType: "application/rss+xml",
			wantErr:     true,
			errCategory: ErrCategoryParseFailure,
		},
		{
			name:        "whitespace only",
			responseBody: "   \n\t  ",
			contentType: "application/rss+xml",
			wantErr:     true,
			errCategory: ErrCategoryParseFailure,
		},
		{
			name:        "invalid XML",
			responseBody: "<invalid>",
			contentType: "application/rss+xml",
			wantErr:     true,
			errCategory: ErrCategoryParseFailure,
		},
		{
			name:        "HTML instead of RSS",
			responseBody: "<!DOCTYPE html><html><body>Not a feed</body></html>",
			contentType: "text/html",
			wantErr:     true,
			errCategory: ErrCategoryParseFailure,
		},
		{
			name:        "valid RSS with empty title",
			responseBody: `<?xml version="1.0"?>
<rss version="2.0">
<channel>
<title></title>
</channel>
</rss>`,
			contentType: "application/rss+xml",
			wantErr:     true,
			errCategory: ErrCategoryEmptyContent,
		},
		{
			name:        "valid RSS with whitespace title",
			responseBody: `<?xml version="1.0"?>
<rss version="2.0">
<channel>
<title>   </title>
</channel>
</rss>`,
			contentType: "application/rss+xml",
			wantErr:     true,
			errCategory: ErrCategoryEmptyContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.contentType != "" {
					w.Header().Set("Content-Type", tt.contentType)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			_, err := FetchAndParseFeed(server.URL)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAndParseFeed() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				if httpErr, ok := err.(*HTTPError); ok {
					if httpErr.Category != tt.errCategory {
						t.Errorf("Error category = %q, want %q", httpErr.Category, tt.errCategory)
					}
				} else {
					t.Errorf("Expected *HTTPError, got %T", err)
				}
			}
		})
	}
}

// TestFetchAndParseFeed_Integration tests integration scenarios with mock server
func TestFetchAndParseFeed_Integration(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		errCategory string
		statusCode  int
	}{
		{
			name: "404 Not Found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr:     true,
			errCategory: ErrCategoryNotFound,
			statusCode:  http.StatusNotFound,
		},
		{
			name: "403 Forbidden",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			},
			wantErr:     true,
			errCategory: ErrCategoryAccessDenied,
			statusCode:  http.StatusForbidden,
		},
		{
			name: "410 Gone",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusGone)
			},
			wantErr:     true,
			errCategory: ErrCategoryNotFound,
			statusCode:  http.StatusGone,
		},
		{
			name: "500 Internal Server Error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr:     true,
			errCategory: ErrCategoryServerError,
			statusCode:  http.StatusInternalServerError,
		},
		{
			name: "503 Service Unavailable",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
			},
			wantErr:     true,
			errCategory: ErrCategoryServiceUnavailable,
			statusCode:  http.StatusServiceUnavailable,
		},
		{
			name: "Successful RSS feed",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/rss+xml")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0"?>
<rss version="2.0">
<channel>
<title>Test Feed</title>
<item><title>Test Item</title><link>http://example.com/item</link></item>
</channel>
</rss>`))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			feed, err := FetchAndParseFeed(server.URL)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAndParseFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				if httpErr, ok := err.(*HTTPError); ok {
					if httpErr.Category != tt.errCategory {
						t.Errorf("Error category = %q, want %q", httpErr.Category, tt.errCategory)
					}
					if tt.statusCode > 0 && httpErr.StatusCode != tt.statusCode {
						t.Errorf("StatusCode = %d, want %d", httpErr.StatusCode, tt.statusCode)
					}
				} else {
					t.Errorf("Expected *HTTPError, got %T", err)
				}
			}

			if !tt.wantErr && feed == nil {
				t.Error("Expected valid feed, got nil")
			}
		})
	}
}

// TestHTTPError_UserMessageFormat verifies user-friendly message format
func TestHTTPError_UserMessageFormat(t *testing.T) {
	tests := []struct {
		name         string
		httpError    *HTTPError
		wantContains []string
	}{
		{
			name: "NotFound message format",
			httpError: &HTTPError{
				StatusCode: http.StatusNotFound,
				Category:   ErrCategoryNotFound,
			},
			wantContains: []string{"not found", "URL"},
		},
		{
			name: "AccessDenied message format",
			httpError: &HTTPError{
				StatusCode: http.StatusForbidden,
				Category:   ErrCategoryAccessDenied,
			},
			wantContains: []string{"access denied", "authentication"},
		},
		{
			name: "Timeout message format",
			httpError: &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryTimeout,
			},
			wantContains: []string{"timed out", "server"},
		},
		{
			name: "ParseFailure message format",
			httpError: &HTTPError{
				StatusCode: 200,
				Category:   ErrCategoryParseFailure,
			},
			wantContains: []string{"parse", "format"},
		},
		{
			name: "EmptyContent message format",
			httpError: &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryEmptyContent,
			},
			wantContains: []string{"empty", "content"},
		},
		{
			name: "TooManyRedirects message format",
			httpError: &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryTooManyRedirects,
			},
			wantContains: []string{"redirect", "loop"},
		},
		{
			name: "ServiceUnavailable message format",
			httpError: &HTTPError{
				StatusCode: http.StatusServiceUnavailable,
				Category:   ErrCategoryServiceUnavailable,
			},
			wantContains: []string{"unavailable", "try again"},
		},
		{
			name: "ServerError message format",
			httpError: &HTTPError{
				StatusCode: http.StatusInternalServerError,
				Category:   ErrCategoryServerError,
			},
			wantContains: []string{"server", "error"},
		},
		{
			name: "NetworkError message format",
			httpError: &HTTPError{
				StatusCode: 0,
				Category:   ErrCategoryNetworkError,
			},
			wantContains: []string{"network", "connect"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.httpError.UserMessage()
			msgLower := strings.ToLower(msg)

			for _, want := range tt.wantContains {
				if !strings.Contains(msgLower, strings.ToLower(want)) {
					t.Errorf("UserMessage() = %q, want to contain %q", msg, want)
				}
			}
		})
	}
}

// TestFetchAndParseFeed_InvalidURL tests invalid URL handling
func TestFetchAndParseFeed_InvalidURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		// Empty URL returns a network error (unsupported protocol scheme)
		{"empty URL", "", true},
		{"invalid URL", "not-a-valid-url", true},
		{"URL without scheme", "example.com/feed", true},
		{"malformed URL", "http://[invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FetchAndParseFeed(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAndParseFeed(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}
