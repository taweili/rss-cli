# S03: Error handling & edge cases

**Goal:** Implement comprehensive error handling for network timeouts, HTTP errors (404, 403, 410), parse failures, empty content detection, and browser open failures with clear user-friendly messages
**Demo:** 

## Must-Haves

- Not provided.

## Proof Level

- This slice proves: Not provided.

## Integration Closure

Not provided.

## Verification

- Not provided.

## Tasks

- [x] **T01: Enhance RSS fetcher error types with categorization** `est:1h`
  Enhance the RSS fetcher error handling to provide categorized, user-friendly error messages.

Steps:
1. Create error category constants in pkg/rss/fetcher.go: ErrNotFound, ErrAccessDenied, ErrTimeout, ErrParseFailure, ErrEmptyContent, ErrTooManyRedirects
2. Enhance HTTPError type to include Category field and UserMessage() method
3. Update httpErr() to categorize status codes: 404→NotFound, 410→NotFound, 403→AccessDenied, 503→ServiceUnavailable, 5xx→ServerError
4. Add timeout error detection wrapper that returns Timeout category
5. Add ValidateFeedContent() helper to check for empty feeds/items
6. Update FetchArticleContent() to validate non-empty response
7. Ensure all error messages match architecture spec format
  - Files: ``pkg/rss/fetcher.go``
  - Verify: go test -v ./pkg/rss -run TestHTTPError

- [x] **T02: Add article open command with browser detection** `est:1h`
  Implement the article open command that opens article URLs in the default browser with proper error handling.

Steps:
1. Add articleOpenCmd to cmd/rss-cli/article_cmd.go
2. Implement detectBrowser() function that checks platform-specific commands: Linux (xdg-open), macOS (open), Windows (start)
3. Implement openBrowser(url) function using exec.Command with fallback chain
4. Wire up error handling: no browser found, launch failure
5. Add --browser flag to article view command as alternative
6. Use printer.Error() for all error output
7. Handle edge cases: empty URL, invalid URL format
  - Files: ``cmd/rss-cli/article_cmd.go``
  - Verify: go build -o rss-cli ./cmd/rss-cli && ./rss-cli article open --help

- [ ] **T03: Write comprehensive error handling tests** `est:1h`
  Write comprehensive unit tests for error handling scenarios.

Steps:
1. Add TestHTTPError_Categorization table-driven test for all status codes
2. Add TestFetchArticleContent_EmptyResponse test
3. Add TestDetectBrowser tests for platform detection
4. Add TestOpenBrowser_Failure tests for invalid URLs
5. Add integration test with mock HTTP server for 404, 403, 410, timeout scenarios
6. Ensure all tests verify error messages match expected format
  - Files: ``pkg/rss/error_test.go``, ``cmd/rss-cli/article_open_test.go``
  - Verify: go test -v ./... -run 'TestHTTP|TestFetch|TestDetect|TestOpen'

- [ ] **T04: Manual verification and integration testing** `est:45m`
  Perform manual verification of all error scenarios and edge cases.

Steps:
1. Test article view with invalid ID (database error)
2. Test article view with non-existent article (empty result)
3. Test article view with 404 URL (simulate with httpbin.org)
4. Test article view with 403 URL (simulate with httpbin.org)
5. Test article open command with valid article
6. Test article open command with empty URL
7. Verify all error messages are clear and actionable
8. Run full test suite: go test ./...
  - Verify: go test ./... && echo 'All tests passed'

## Files Likely Touched

- `pkg/rss/fetcher.go`
- `cmd/rss-cli/article_cmd.go`
- `pkg/rss/error_test.go`
- `cmd/rss-cli/article_open_test.go`
