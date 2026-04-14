# S03 Research: Error Handling & Edge Cases

## Summary

The current codebase has a solid foundation for error handling with `printer.Error()` in `pkg/ui/output.go` providing consistent JSON/text output formatting. The RSS fetcher in `pkg/rss/fetcher.go` already implements a 30-second timeout and custom `HTTPError` type for HTTP status codes. Database operations return errors that are properly wrapped with context in CLI commands.

However, several critical error scenarios are NOT yet covered:
1. **No article open command exists** - browser opening functionality needs to be implemented from scratch
2. **HTTP error differentiation** - the current `httpErr()` returns generic status text without distinguishing 404, 403, 410, etc.
3. **Parse failure messages** - gofeed parse errors are returned raw without user-friendly messaging
4. **Empty content detection** - no validation that fetched articles have actual content
5. **Default browser detection** - no mechanism to detect or handle missing browser

## Recommendation

Implement error handling in three layers:

1. **RSS Fetcher Layer** (`pkg/rss/fetcher.go`): Enhance `HTTPError` to categorize errors (not-found, access-denied, timeout, parse-failure) with user-friendly messages matching the architecture spec.

2. **Article Command Layer** (`cmd/rss-cli/article_cmd.go`): Add new `open` subcommand with browser detection using `exec.Command` and platform-specific fallbacks. Wrap all fetcher errors with `printer.Error()`.

3. **Validation Layer**: Add content validation after parsing - check for empty feeds, empty article content, and provide appropriate error messages.

This approach keeps error logic close to the source while ensuring all CLI output flows through `printer.Error()` for consistency.

## Implementation Landscape

### Key Files

| File | Purpose | Changes Needed |
|------|---------|----------------|
| `pkg/rss/fetcher.go` | HTTP fetching & parsing | Enhance error types, add content validation |
| `cmd/rss-cli/article_cmd.go` | Article commands | Add `open` subcommand with browser logic |
| `pkg/ui/output.go` | Output formatting | `printer.Error()` already exists - no changes |

### Build Order

1. **First**: Enhance `pkg/rss/fetcher.go` error types
   - Add error categorization (NotFound, AccessDenied, Timeout, ParseFailure, EmptyContent)
   - Add user-friendly error messages per architecture spec
   - Add content validation helper

2. **Second**: Add `cmd/rss-cli/article_cmd.go` open command
   - Implement browser detection (xdg-open, open, start)
   - Wire up error handling with `printer.Error()`

3. **Third**: Integration testing
   - Test network timeout simulation
   - Test HTTP error codes (404, 403, 410)
   - Test invalid XML/HTML parsing
   - Test browser open on current platform

### Verification Approach

```bash
# Unit tests for error types
go test -v ./pkg/rss -run TestHTTPError

# Integration test with mock server
go test -v ./pkg/rss -run TestFetchAndParseFeed

# Manual CLI verification
go run ./cmd/rss-cli article open <id>  # With various failure modes
```

## Common Pitfalls

1. **Browser command platform differences**: Linux uses `xdg-open`, macOS uses `open`, Windows uses `start`. Must detect platform or try fallbacks.

2. **Timeout errors wrapped incorrectly**: Go's `http.Client` timeout errors don't implement `net.Error` interface consistently. Check for `context.DeadlineExceeded` or use `errors.As()` with `net.Error`.

3. **Parse errors from gofeed**: These can be XML parsing errors, encoding errors, or feed format errors. Wrap with generic "Failed to parse" message rather than exposing library internals.

4. **Empty content vs empty feed**: Distinguish between "feed parsed but has no items" vs "article content field is empty" - different error messages.

5. **Error count tracking**: The database has `IncrementErrorCount()` on feeds but it's not being called in the fetch flow. Should be wired into the update commands.

## Error Scenarios Matrix

| Scenario | Current State | S03 Requirement |
|----------|--------------|-----------------|
| Network timeout (30s) | ✅ Timeout configured, raw error | ❌ Needs user-friendly message |
| HTTP 404/410 | ✅ Status captured, generic text | ❌ Needs "Article not found (404)" |
| HTTP 403/paywall | ✅ Status captured, generic text | ❌ Needs "Access denied (403) — article may be paywalled" |
| Invalid HTML/XML | ❌ Raw gofeed error | ❌ Needs "Failed to parse article content" |
| Empty content | ❌ No validation | ❌ Needs "Article content is empty" |
| Browser open - no default | ❌ No open command | ❌ Needs "No default browser found" |
| Browser open - launch failure | ❌ No open command | ❌ Needs "Failed to open browser: [error]" |
