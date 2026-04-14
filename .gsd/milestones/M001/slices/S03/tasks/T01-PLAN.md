---
estimated_steps: 9
estimated_files: 1
skills_used: []
---

# T01: Enhance RSS fetcher error types with categorization

Enhance the RSS fetcher error handling to provide categorized, user-friendly error messages.

Steps:
1. Create error category constants in pkg/rss/fetcher.go: ErrNotFound, ErrAccessDenied, ErrTimeout, ErrParseFailure, ErrEmptyContent, ErrTooManyRedirects
2. Enhance HTTPError type to include Category field and UserMessage() method
3. Update httpErr() to categorize status codes: 404→NotFound, 410→NotFound, 403→AccessDenied, 503→ServiceUnavailable, 5xx→ServerError
4. Add timeout error detection wrapper that returns Timeout category
5. Add ValidateFeedContent() helper to check for empty feeds/items
6. Update FetchArticleContent() to validate non-empty response
7. Ensure all error messages match architecture spec format

## Inputs

- ``pkg/rss/fetcher.go``

## Expected Output

- ``pkg/rss/fetcher.go``

## Verification

go test -v ./pkg/rss -run TestHTTPError

## Observability Impact

none
