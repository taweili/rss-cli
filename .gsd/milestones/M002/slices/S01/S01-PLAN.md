# S01: HTML-to-markdown converter

**Goal:** Implement FetchAndConvertArticle(url string) (string, error) function that fetches HTML from a URL, converts it to markdown using html-to-markdown library, and returns categorized errors for all failure modes
**Demo:** FetchAndConvertArticle(url) fetches HTML from URL, converts to markdown, returns error with category for failures

## Must-Haves

- Not provided.

## Proof Level

- This slice proves: Not provided.

## Integration Closure

Not provided.

## Verification

- Not provided.

## Tasks

- [x] **T01: Implement FetchAndConvertArticle() converter function** `est:45m`
  Create the core converter implementation in pkg/rss/converter.go. The function FetchAndConvertArticle(url string) (string, error) should:
1. Fetch HTML using the same HTTP client pattern as FetchArticleContent() (30s timeout, 10 redirect limit, User-Agent header)
2. Validate the HTML response is not empty
3. Convert HTML to markdown using htmltomarkdown.Convert()
4. Validate the converted markdown is not empty
5. Return categorized errors (HTTPError with appropriate categories) for all failure modes

Reuse the existing HTTPError type and error categorization constants from fetcher.go. Add a new error category ErrCategoryConversionFailure if the HTML converts to empty markdown.
  - Files: `pkg/rss/converter.go`, `pkg/rss/fetcher.go`
  - Verify: go build ./pkg/rss && test -f pkg/rss/converter.go && grep -q 'func FetchAndConvertArticle' pkg/rss/converter.go

- [x] **T02: Add unit tests for FetchAndConvertArticle()** `est:45m`
  Add comprehensive unit tests for FetchAndConvertArticle() in pkg/rss/converter_test.go. Tests should cover:
1. Successful conversion with various HTML structures (headings, paragraphs, lists, links, code blocks)
2. HTTP error scenarios (404, 403, 500) using mock HTTP server
3. Timeout scenario
4. Empty HTML response
5. HTML that converts to empty markdown
6. Network errors

Use table-driven tests following the existing pattern in converter_test.go and fetcher_test.go. Use httptest.NewServer() for mock HTTP scenarios.
  - Files: `pkg/rss/converter_test.go`
  - Verify: go test ./pkg/rss -run TestFetchAndConvertArticle -v

## Files Likely Touched

- pkg/rss/converter.go
- pkg/rss/fetcher.go
- pkg/rss/converter_test.go
