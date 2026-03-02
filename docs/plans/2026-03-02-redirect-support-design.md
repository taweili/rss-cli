# Redirect Support Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add explicit redirect handling with a 10-hop limit to the RSS feed fetcher to prevent infinite redirect loops while allowing legitimate redirects.

**Architecture:** Extend the existing `http.Client` configuration in `pkg/rss/fetcher.go` with a `CheckRedirect` function that limits redirects to 10 hops. This provides safety without changing the existing behavior for normal cases.

**Tech Stack:** Go standard library (net/http), gofeed for RSS parsing

---

### Task 1: Write failing test for redirect limit

**Files:**
- Create: `pkg/rss/fetcher_test.go`

**Step 1: Create test file with redirect limit test**

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/rss -run TestFetchAndParseFeed_ExcessiveRedirects -v`

Expected: FAIL - test will fail because current implementation doesn't have redirect limit

**Step 3: Commit test file**

```bash
git add pkg/rss/fetcher_test.go
git commit -m "test: add test for redirect limit enforcement"
```

---

### Task 2: Implement redirect limit in fetcher

**Files:**
- Modify: `pkg/rss/fetcher.go:13-15`

**Step 1: Add CheckRedirect to http.Client**

```go
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
```

**Step 2: Add fmt import**

Add to imports at top of file:
```go
import (
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)
```

**Step 3: Run tests to verify implementation**

Run: `go test ./pkg/rss -run TestFetchAndParseFeed_ExcessiveRedirects -v`

Expected: PASS - test should now pass with redirect limit working

**Step 4: Run all tests to ensure no regressions**

Run: `go test ./...`

Expected: All tests pass

**Step 5: Commit implementation**

```bash
git add pkg/rss/fetcher.go
git commit -m "feat: add redirect limit (10 hops) to feed fetcher

Prevents infinite redirect loops while allowing legitimate redirects.
Returns clear error when limit is exceeded."
```

---

### Task 3: Add test for successful redirect following

**Files:**
- Modify: `pkg/rss/fetcher_test.go`

**Step 1: Add test for successful redirect**

```go
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
```

**Step 2: Run new test**

Run: `go test ./pkg/rss -run TestFetchAndParseFeed_SuccessfulRedirect -v`

Expected: PASS

**Step 3: Commit additional test**

```bash
git add pkg/rss/fetcher_test.go
git commit -m "test: add test for successful redirect following

Verifies that legitimate redirects (301) work correctly
and the final feed content is parsed."
```

---

## Summary

This plan implements redirect support by:
1. Adding a 10-hop limit to prevent infinite redirect loops
2. Testing the limit enforcement
3. Testing that legitimate redirects still work
4. Following TDD principles throughout

The implementation is minimal and focused, adding only what's needed for safety without changing normal behavior.
