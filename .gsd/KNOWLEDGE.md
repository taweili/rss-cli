# Knowledge Base

This file captures cross-cutting lessons, patterns, and gotchas that emerged during milestone execution. Use it to avoid reinventing solutions and to understand why certain patterns exist.

## Error Handling Pattern (M001/S03)

**Context:** Network operations, HTTP requests, and external API calls can fail in many ways.

**Pattern:**
1. Define error category constants (NotFound, AccessDenied, Timeout, ParseFailure, etc.)
2. Add `Category` field to error types
3. Implement `UserMessage()` method for user-friendly output
4. Use `errors.As()` to preserve error categorization when wrapping
5. Map HTTP status codes to categories in helper functions

**Example:**
```go
const (
    ErrCategoryNotFound     = "NotFound"
    ErrCategoryAccessDenied = "AccessDenied"
    ErrCategoryTimeout      = "Timeout"
)

type HTTPError struct {
    Category string
    StatusCode int
    Err error
}

func (h *HTTPError) UserMessage() string {
    switch h.Category {
    case ErrCategoryNotFound:
        return "Article not found (404)"
    case ErrCategoryAccessDenied:
        return "Access denied (403) — article may be paywalled"
    default:
        return fmt.Sprintf("Request failed (%d)", h.StatusCode)
    }
}
```

**Why it matters:** Users get actionable messages instead of raw HTTP status codes. Tests can verify error categorization independently from message formatting.

## Browser Opening Pattern (M001/S02, M001/S03)

**Context:** Opening URLs in the default browser across platforms.

**Pattern:**
1. Use `github.com/toqueteos/webbrowser` as primary opener
2. Implement platform-specific fallbacks (xdg-open, open, start)
3. Detach browser processes (nil stdout/stderr) to prevent blocking
4. Support `--browser` flag for custom browser selection

**Example:**
```go
func openBrowser(url string) error {
    // Try webbrowser library first
    if err := webbrowser.Open(url); err == nil {
        return nil
    }
    // Fallback to platform-specific commands
    // ...
}
```

**Why it matters:** CLI doesn't hang waiting for browser to close. Works reliably across macOS, Linux, and Windows.

## HTML to Markdown Conversion (M001/S01)

**Context:** Converting article HTML to readable markdown while preserving structure.

**Pattern:**
1. Use `github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown`
2. Fetch HTML with proper user-agent and timeout
3. Validate content is not empty before conversion
4. Handle conversion errors gracefully

**Note:** Library uses CGO to call Rust FFI. Requires CGO-enabled build environment.

## On-Demand Fetching Strategy (M001/S01)

**Context:** When to fetch full article content.

**Decision:** Fetch on-demand when viewing, not during feed update.

**Rationale:**
- Avoids database storage bloat
- Always shows current content (no stale cache)
- Respects publisher bandwidth
- Simpler implementation

**Trade-off:** Slightly slower article viewing, but acceptable for reading workflow.

## Testing Pattern for Error Handling (M001/S03)

**Context:** Comprehensive test coverage for error scenarios.

**Pattern:**
1. Table-driven tests for HTTP status code categorization
2. Mock HTTP server for integration tests
3. Test both error type and UserMessage() output
4. Manual testing against real services (httpbin.org)

**Example:**
```go
func TestHTTPErrorCategorization(t *testing.T) {
    tests := []struct {
        name     string
        statusCode int
        wantCategory string
    }{
        {"404", 404, ErrCategoryNotFound},
        {"403", 403, ErrCategoryAccessDenied},
        {"503", 503, ErrCategoryServiceUnavailable},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

## CLI Command Pattern (M001)

**Context:** Consistent CLI command structure across the application.

**Pattern:**
1. Use `RunE` for error-returning command execution
2. Use `cobra.ExactArgs(n)` for required positional arguments
3. Access flags via `cmd.Flags().GetString()`, `cmd.Flags().GetBool()`
4. Use `printer.Error()` for consistent error JSON output
5. Defer `db.Close()` immediately after successful connection
6. Support both `--json` and `--text` output modes

## Database Query Pattern (M001/S01)

**Context:** Consistent database access patterns.

**Pattern:**
- Single result: `db.QueryRow()` with `Scan()`
- Multiple results: `db.Query()` with `rows.Next()` loop
- Always defer `rows.Close()`
- Use parameterized queries to prevent SQL injection
- Use transactions for multi-table operations

## Verification Evidence

All patterns above were verified through:
- Unit tests (25+ tests in error_test.go, article_open_test.go, converter_test.go)
- Integration tests with mock HTTP server
- Manual testing against real RSS feeds and httpbin.org
- Build verification (`go build -o rss-cli ./cmd/rss-cli`)
