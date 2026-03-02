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
