package rss

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
