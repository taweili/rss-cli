# M002/S01 — Research

**Date:** 2026-04-15

## Summary

Slice S01 requires implementing an HTML-to-markdown converter function in `pkg/rss/converter.go`. The converter will fetch full article HTML from a URL and convert it to markdown using the existing `github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown` library already in the dependency tree.

The implementation follows the established pattern from `pkg/rss/fetcher.go`: a single `FetchAndConvertArticle(url string) (string, error)` function that handles HTTP fetching, error categorization, and markdown conversion. Error handling reuses the existing `HTTPError` pattern with categories like `NotFound`, `AccessDenied`, `Timeout`, `ConversionFailure`, and `EmptyContent`.

## Recommendation

Implement `FetchAndConvertArticle()` as a standalone function in `pkg/rss/converter.go` that:
1. Fetches HTML using the same HTTP client pattern as `FetchArticleContent()`
2. Validates the HTML is not empty
3. Converts HTML to markdown using `htmltomarkdown.Convert()`
4. Validates the converted markdown is not empty
5. Returns categorized errors for all failure modes

This approach keeps the converter focused and testable, matches the existing codebase patterns, and provides clear error messages for CLI output.

## Implementation Landscape

### Key Files

- `pkg/rss/converter.go` — **NEW FILE**: Create `FetchAndConvertArticle(url string) (string, error)` function with HTTP fetching, validation, and markdown conversion
- `pkg/rss/converter_test.go` — **UPDATE**: Add unit tests for `FetchAndConvertArticle()` with mock HTTP server scenarios (404, 403, timeout, empty content, conversion failure)
- `pkg/rss/fetcher.go` — **REFERENCE ONLY**: Reuse HTTP client pattern, error categorization, and `HTTPError` type structure
- `cmd/rss-cli/article_fetch_cmd.go` — **NEW FILE** (S02): Will call `FetchAndConvertArticle()` from the CLI command

### Build Order

1. **Create `pkg/rss/converter.go`** — Core conversion logic, unblocks all downstream work
2. **Add unit tests in `pkg/rss/converter_test.go`** — Validate conversion works with mock HTML and error scenarios
3. **Manual verification** — Test against real RSS feed articles

The converter is the foundation — S02 (CLI command) and S03 (integration tests) depend on it working correctly.

### Verification Approach

```bash
# Unit tests for converter
go test ./pkg/rss -run TestFetchAndConvertArticle -v

# Build verification
go build -o rss-cli ./cmd/rss-cli

# Manual testing (after S02 CLI command is implemented)
./rss-cli article fetch [id]
```

## Don't Hand-Roll

| Problem | Existing Solution | Why Use It |
|---------|------------------|------------|
| HTML to markdown conversion | `github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown` | Already in go.mod, tested, preserves formatting (headings, lists, links, code blocks) |
| HTTP client with timeout/redirects | Pattern from `pkg/rss/fetcher.go` | Reuse 30s timeout, 10-redirect limit, User-Agent header |
| Error categorization | `HTTPError` type in `fetcher.go` | Consistent error messages across CLI commands |

## Constraints

- Must use existing `htmltomarkdown.Convert()` function (CGO dependency via html-to-markdown Rust library)
- Must reuse existing `HTTPError` pattern for consistency with M001/S03
- Cannot change database schema (S02 will update existing `content` field)
- Function signature must be `FetchAndConvertArticle(url string) (string, error)` to match existing patterns

## Common Pitfalls

- **Empty HTML after fetch** — Validate response body is not empty before conversion, return `ErrCategoryEmptyContent`
- **HTML parses but produces empty markdown** — Some pages may have only scripts/styles; validate markdown output is not empty, return `ErrCategoryConversionFailure` or `ErrCategoryEmptyContent`
- **Relative URLs in article content** — The converter preserves links as-is; relative URLs in the original HTML remain relative (acceptable for terminal display)
- **JavaScript-rendered content** — Articles that require JS execution won't render; conversion will produce empty/partial markdown (return `ErrCategoryEmptyContent`)

## Open Risks

- Some RSS feeds may have non-HTML content types (JSON, plain text) — conversion may fail or produce garbage
- Paywalled articles may return HTML with login forms instead of content — will convert but produce unexpected markdown

## Skills Discovered

| Technology | Skill | Status |
|------------|-------|--------|
| Go | `golang-pro` | available in `<available_skills>` |
| Go testing | `golang-testing` | available in `<available_skills>` |

## Sources

- `pkg/rss/fetcher.go` — Existing HTTP fetching and error handling pattern
- `pkg/rss/converter_test.go` — Existing tests showing `htmltomarkdown.Convert()` usage
- [html-to-markdown Go API docs](https://github.com/kreuzberg-dev/html-to-markdown/blob/main/docs/reference/api-go.md) — `Convert()` function signature and examples
- `go.mod` — Confirms `github.com/kreuzberg-dev/html-to-markdown/packages/go/v2 v2.30.0` dependency
