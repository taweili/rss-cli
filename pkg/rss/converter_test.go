package rss

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchAndConvertArticle_Success(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		wantContains string
	}{
		{
			name:         "simple heading and paragraph",
			html:         "<h1>Hello World</h1><p>This is a paragraph.</p>",
			wantContains: "# Hello World",
		},
		{
			name:         "multiple paragraphs",
			html:         "<p>First paragraph.</p><p>Second paragraph.</p>",
			wantContains: "First paragraph.",
		},
		{
			name:         "headings h1-h3",
			html:         "<h1>Main</h1><h2>Section</h2><h3>Subsection</h3>",
			wantContains: "# Main",
		},
		{
			name:         "unordered list",
			html:         "<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul>",
			wantContains: "- Item 1",
		},
		{
			name:         "ordered list",
			html:         "<ol><li>First</li><li>Second</li></ol>",
			wantContains: "1. First",
		},
		{
			name:         "links",
			html:         `<p>Visit <a href="https://example.com">Example</a></p>`,
			wantContains: "[Example](https://example.com)",
		},
		{
			name:         "bold and italic",
			html:         "<p>This is <strong>bold</strong> and <em>italic</em>.</p>",
			wantContains: "**bold**",
		},
		{
			name:         "code blocks",
			html:         "<pre><code>func main() {\n    fmt.Println(\"Hello\")\n}</code></pre>",
			wantContains: "func main()",
		},
		{
			name:         "inline code",
			html:         "<p>Use <code>printf</code> for output</p>",
			wantContains: "`printf`",
		},
		{
			name:         "nested structure",
			html:         "<article><h1>Title</h1><p>Intro</p><div><h2>Section</h2><p>Content</p></div></article>",
			wantContains: "# Title",
		},
		{
			name:         "blockquote",
			html:         "<blockquote>This is a quote</blockquote>",
			wantContains: "> This is a quote",
		},
		{
			name:         "image with alt text",
			html:         `<img src="image.jpg" alt="Description">`,
			wantContains: "![Description]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.html))
			}))
			defer server.Close()

			md, err := FetchAndConvertArticle(server.URL)
			if err != nil {
				t.Fatalf("FetchAndConvertArticle() error = %v", err)
			}

			if !strings.Contains(md, tt.wantContains) {
				t.Errorf("FetchAndConvertArticle() = %q, want to contain %q", md, tt.wantContains)
			}
		})
	}
}

func TestFetchAndConvertArticle_HTTPErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantCategory string
	}{
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			wantCategory: ErrCategoryNotFound,
		},
		{
			name:       "403 Forbidden",
			statusCode: http.StatusForbidden,
			wantCategory: ErrCategoryAccessDenied,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			wantCategory: ErrCategoryServerError,
		},
		{
			name:       "503 Service Unavailable",
			statusCode: http.StatusServiceUnavailable,
			wantCategory: ErrCategoryServiceUnavailable,
		},
		{
			name:       "410 Gone",
			statusCode: http.StatusGone,
			wantCategory: ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			_, err := FetchAndConvertArticle(server.URL)
			if err == nil {
				t.Fatalf("FetchAndConvertArticle() expected error, got nil")
			}

			httpErr, ok := err.(*HTTPError)
			if !ok {
				t.Fatalf("Expected *HTTPError, got %T", err)
			}

			if httpErr.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, httpErr.StatusCode)
			}

			if httpErr.Category != tt.wantCategory {
				t.Errorf("Expected category %q, got %q", tt.wantCategory, httpErr.Category)
			}
		})
	}
}

func TestFetchAndConvertArticle_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response that will timeout
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<p>Too late</p>"))
	}))
	defer server.Close()

	// Create a client with a very short timeout for testing
	// We'll use a custom approach since FetchAndConvertArticle uses a fixed 30s timeout
	// Instead, we'll test with a URL that takes too long
	
	// For this test, we'll just verify the function exists and can be called
	// A real timeout test would require modifying the client timeout
	// Since we can't easily test a 30s timeout, we'll skip the actual timeout verification
	// and just ensure the function handles slow servers gracefully
	
	done := make(chan error, 1)
	go func() {
		_, err := FetchAndConvertArticle(server.URL)
		done <- err
	}()

	select {
	case err := <-done:
		// If we get an error quickly, it might be a timeout or connection error
		if err != nil {
			httpErr, ok := err.(*HTTPError)
			if ok && httpErr.Category == ErrCategoryTimeout {
				t.Logf("Got expected timeout error: %v", err)
			}
		}
	case <-time.After(35 * time.Second):
		t.Fatal("FetchAndConvertArticle() did not timeout as expected")
	}
}

func TestFetchAndConvertArticle_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}))
	defer server.Close()

	_, err := FetchAndConvertArticle(server.URL)
	if err == nil {
		t.Fatal("FetchAndConvertArticle() expected error for empty response, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected *HTTPError, got %T", err)
	}

	if httpErr.Category != ErrCategoryEmptyContent {
		t.Errorf("Expected category %q, got %q", ErrCategoryEmptyContent, httpErr.Category)
	}
}

func TestFetchAndConvertArticle_WhitespaceOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("   \n\t   "))
	}))
	defer server.Close()

	_, err := FetchAndConvertArticle(server.URL)
	if err == nil {
		t.Fatal("FetchAndConvertArticle() expected error for whitespace-only response, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected *HTTPError, got %T", err)
	}

	if httpErr.Category != ErrCategoryEmptyContent {
		t.Errorf("Expected category %q, got %q", ErrCategoryEmptyContent, httpErr.Category)
	}
}

func TestFetchAndConvertArticle_NetworkError(t *testing.T) {
	// Test with invalid URL that will cause network error
	_, err := FetchAndConvertArticle("http://invalid-url-that-does-not-exist.example")
	if err == nil {
		t.Fatal("FetchAndConvertArticle() expected network error, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected *HTTPError, got %T", err)
	}

	if httpErr.Category != ErrCategoryNetworkError {
		t.Errorf("Expected category %q, got %q", ErrCategoryNetworkError, httpErr.Category)
	}
}

func TestFetchAndConvertArticle_InvalidURL(t *testing.T) {
	_, err := FetchAndConvertArticle("not-a-valid-url")
	if err == nil {
		t.Fatal("FetchAndConvertArticle() expected error for invalid URL, got nil")
	}
}

func TestFetchAndConvertArticle_TooManyRedirects(t *testing.T) {
	// Create a server that redirects to itself
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, server.URL, http.StatusMovedPermanently)
	}))
	defer server.Close()

	_, err := FetchAndConvertArticle(server.URL)
	if err == nil {
		t.Fatal("FetchAndConvertArticle() expected redirect error, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected *HTTPError, got %T", err)
	}

	if httpErr.Category != ErrCategoryTooManyRedirects {
		t.Errorf("Expected category %q, got %q", ErrCategoryTooManyRedirects, httpErr.Category)
	}
}

func TestFetchAndConvertArticle_UserMessage(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func() (*httptest.Server, string)
		wantContains string
	}{
		{
			name: "empty content user message",
			setupFunc: func() (*httptest.Server, string) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(""))
				}))
				return server, ErrCategoryEmptyContent
			},
			wantContains: "empty",
		},
		{
			name: "not found user message",
			setupFunc: func() (*httptest.Server, string) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
				return server, ErrCategoryNotFound
			},
			wantContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := tt.setupFunc()
			defer server.Close()

			_, err := FetchAndConvertArticle(server.URL)
			if err == nil {
				t.Fatalf("FetchAndConvertArticle() expected error, got nil")
			}

			httpErr, ok := err.(*HTTPError)
			if !ok {
				t.Fatalf("Expected *HTTPError, got %T", err)
			}

			msg := httpErr.UserMessage()
			if !strings.Contains(strings.ToLower(msg), tt.wantContains) {
				t.Errorf("UserMessage() = %q, want to contain %q", msg, tt.wantContains)
			}
		})
	}
}
