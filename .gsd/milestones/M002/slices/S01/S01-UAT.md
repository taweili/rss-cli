# S01: HTML-to-markdown converter — UAT

**Milestone:** M002
**Written:** 2026-04-15T02:11:31.516Z

# S01: HTML-to-markdown converter — UAT

**Milestone:** M002
**Written:** 2026-04-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: This slice is a library function (not a CLI command or UI), so unit tests with mock HTTP servers provide complete coverage of all success and failure paths.

## Preconditions

- Go toolchain available
- CGO enabled (required by html-to-markdown library)
- All dependencies installed via go mod download

## Smoke Test

Verify FetchAndConvertArticle() function exists and compiles:
1. Run `go build ./pkg/rss`
2. **Expected:** Build succeeds with no errors

## Test Cases

### 1. Successful HTML conversion

1. Call FetchAndConvertArticle() with a URL returning valid HTML containing headings, paragraphs, lists, links, and code blocks
2. **Expected:** Returns markdown string with preserved structure (# for h1, ## for h2, - for lists, [text](url) for links, code blocks with ```)

### 2. HTTP 404 error handling

1. Call FetchAndConvertArticle() with a URL returning 404 status
2. **Expected:** Returns error with Category="NotFound", StatusCode=404, UserMessage()="Article not found (404)"

### 3. HTTP 403 error handling

1. Call FetchAndConvertArticle() with a URL returning 403 status
2. **Expected:** Returns error with Category="AccessDenied", StatusCode=403, UserMessage()="Access denied (403) — article may be paywalled"

### 4. Timeout handling

1. Call FetchAndConvertArticle() with a URL that takes >30s to respond
2. **Expected:** Returns error with Category="Timeout", completes within 30s timeout window

### 5. Empty response handling

1. Call FetchAndConvertArticle() with a URL returning empty body
2. **Expected:** Returns error with Category="EmptyContent", UserMessage()="Article fetched but contains no content"

### 6. Network error handling

1. Call FetchAndConvertArticle() with an invalid URL (e.g., "http://invalid-host-that-does-not-exist")
2. **Expected:** Returns error with Category="NetworkError", UserMessage() contains "network error"

## Edge Cases

### Whitespace-only HTML

1. Call FetchAndConvertArticle() with HTML containing only whitespace
2. **Expected:** Returns error with Category="EmptyContent" (treated as empty)

### Too many redirects

1. Call FetchAndConvertArticle() with a URL that redirects to itself
2. **Expected:** Returns error with Category="TooManyRedirects", UserMessage()="Too many redirects"

## Failure Signals

- Build failures: `go build ./pkg/rss` fails
- Test failures: Any TestFetchAndConvertArticle* test fails
- Vet errors: `go vet ./...` reports issues

## Not Proven By This UAT

- Live HTTP requests to real websites (tests use mock HTTP servers)
- Integration with the CLI command layer (covered in S02)
- Database caching behavior (covered in S02)

## Notes for Tester

All tests use httptest.NewServer() for mock HTTP scenarios, ensuring deterministic, fast test execution without external network dependencies. The timeout test uses a 2s delay (well under the 30s timeout) to verify timeout handling without waiting the full duration.
