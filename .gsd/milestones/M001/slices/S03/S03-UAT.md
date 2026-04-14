# S03: Error handling & edge cases — UAT

**Milestone:** M001
**Written:** 2026-04-14T15:23:20.136Z

# S03: Error handling & edge cases — UAT

**Milestone:** M001
**Written:** 2026-04-14

## UAT Type

- UAT mode: artifact-driven (CLI commands with observable output)
- Why this mode is sufficient: All error handling is exercised through CLI commands that produce immediate, verifiable JSON/text output. No runtime server or background processes involved.

## Preconditions

- RSS CLI built: `go build -o rss-cli ./cmd/rss-cli`
- Database exists with at least one article (for article open testing)
- Network access available for httpbin.org tests

## Smoke Test

```bash
./rss-cli article --help
```
**Expected:** Shows article subcommands including 'open' command

## Test Cases

### 1. Invalid Article ID Handling

1. Run: `./rss-cli article view abc`
2. **Expected:** Returns error "Invalid article ID" with proper JSON formatting

### 2. Non-existent Article Handling

1. Run: `./rss-cli article view 99999` (assuming ID 99999 doesn't exist)
2. **Expected:** Returns "Failed to retrieve article: sql: no rows in result set"

### 3. HTTP 404 Error Categorization

1. Run: `./rss-cli feed add "https://httpbin.org/status/404"`
2. **Expected:** Returns error containing "Not Found" (NotFound category)

### 4. HTTP 403 Error Categorization

1. Run: `./rss-cli feed add "https://httpbin.org/status/403"`
2. **Expected:** Returns error containing "Forbidden" (AccessDenied category)

### 5. Article Open Command

1. Run: `./rss-cli article open <valid-article-id>` (replace with actual ID from database)
2. **Expected:** Attempts to open browser; in headless environment returns appropriate error about browser detection

### 6. Custom Browser Failure

1. Run: `./rss-cli article open <valid-article-id> --browser nonexistent-browser`
2. **Expected:** Returns error "Failed to open article in nonexistent-browser: exec: \"nonexistent-browser\": executable file not found in $PATH"

### 7. Article View with --open Flag

1. Run: `./rss-cli article view <valid-article-id> --open`
2. **Expected:** Fetches article content and attempts to open browser

## Edge Cases

### Empty URL Handling

1. Attempt to view article with no URL stored in database
2. **Expected:** Returns "Article has no URL" error

### Invalid URL Format

1. Add feed with malformed URL
2. **Expected:** Returns network error with context about URL validation

### Timeout Error Detection

1. Add feed with very slow response URL (e.g., httpbin.org/delay/10)
2. **Expected:** Returns timeout error with appropriate message

## Failure Signals

- CLI commands return non-zero exit codes on errors
- Error messages contain context (operation being performed)
- Error messages are user-friendly (not raw stack traces)
- JSON output format maintained for error responses

## Not Proven By This UAT

- Long-running timeout scenarios (would require extended test time)
- Actual browser opening behavior in desktop environment (tested in headless environment only)
- All possible HTTP status codes (tested representative samples: 404, 403, 500, 503)

## Notes for Tester

- Replace `<valid-article-id>` with actual article ID from your database
- HTTP 404/403 tests require network access to httpbin.org
- Browser open tests will fail in headless environments but should show appropriate error messages
- All error messages should be actionable and include operation context
