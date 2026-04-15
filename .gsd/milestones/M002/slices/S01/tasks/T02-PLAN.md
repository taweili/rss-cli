---
estimated_steps: 8
estimated_files: 1
skills_used: []
---

# T02: Add unit tests for FetchAndConvertArticle()

Add comprehensive unit tests for FetchAndConvertArticle() in pkg/rss/converter_test.go. Tests should cover:
1. Successful conversion with various HTML structures (headings, paragraphs, lists, links, code blocks)
2. HTTP error scenarios (404, 403, 500) using mock HTTP server
3. Timeout scenario
4. Empty HTML response
5. HTML that converts to empty markdown
6. Network errors

Use table-driven tests following the existing pattern in converter_test.go and fetcher_test.go. Use httptest.NewServer() for mock HTTP scenarios.

## Inputs

- ``pkg/rss/converter.go` — Function to test`
- ``pkg/rss/fetcher_test.go` — Reference for test patterns with mock HTTP server`

## Expected Output

- ``pkg/rss/converter_test.go` — Updated with comprehensive unit tests for FetchAndConvertArticle()`

## Verification

go test ./pkg/rss -run TestFetchAndConvertArticle -v

## Observability Impact

None — test coverage only, no runtime changes
