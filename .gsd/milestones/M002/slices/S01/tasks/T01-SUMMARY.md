---
id: T01
parent: S01
milestone: M002
key_files:
  - pkg/rss/converter.go
  - pkg/rss/fetcher.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-04-15T02:04:18.889Z
blocker_discovered: false
---

# T01: Implement FetchAndConvertArticle() HTML-to-markdown converter function

**Implement FetchAndConvertArticle() HTML-to-markdown converter function**

## What Happened

Implemented the FetchAndConvertArticle() function in pkg/rss/converter.go that:
1. Fetches HTML from a URL using the same HTTP client pattern as FetchArticleContent (30s timeout, 10 redirect limit, User-Agent header)
2. Validates the HTML response is not empty
3. Converts HTML to markdown using htmltomarkdown.Convert()
4. Validates the converted markdown is not empty
5. Returns categorized errors (HTTPError with appropriate categories) for all failure modes

Added a new error category ErrCategoryConversionFailure to fetcher.go for conversion failures, along with a corresponding user-friendly error message in the UserMessage() method.

The implementation reuses the existing HTTPError type and error categorization constants from fetcher.go, maintaining consistency with the existing codebase patterns.

## Verification

Built the package successfully with `go build ./pkg/rss`. Verified converter.go exists and contains the FetchAndConvertArticle function. All existing tests in pkg/rss pass (27 tests, 1 skipped).

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `go build ./pkg/rss && test -f pkg/rss/converter.go && grep -q 'func FetchAndConvertArticle' pkg/rss/converter.go` | 0 | ✅ pass | 500ms |
| 2 | `go test -v ./pkg/rss/...` | 0 | ✅ pass | 1515ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `pkg/rss/converter.go`
- `pkg/rss/fetcher.go`
