---
id: S01
parent: M002
milestone: M002
provides:
  - (none)
requires:
  []
affects:
  []
key_files:
  - ["pkg/rss/converter.go", "pkg/rss/converter_test.go", "pkg/rss/fetcher.go", "pkg/database/article.go"]
key_decisions:
  - (none)
patterns_established:
  - ["Reuse existing HTTP client pattern (30s timeout, 10 redirect limit, User-Agent) for new fetch operations", "Add new error categories to existing HTTPError type rather than creating new error types", "Use table-driven tests with mock HTTP servers for comprehensive coverage", "Validate both HTTP response and converted output for empty content"]
observability_surfaces:
  - none
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-04-15T02:11:31.516Z
blocker_discovered: false
---

# S01: HTML-to-markdown converter

**Implemented FetchAndConvertArticle(url) function that fetches HTML from URLs, converts to markdown using html-to-markdown library, and returns categorized errors for all failure modes**

## What Happened

Implemented the FetchAndConvertArticle() function in pkg/rss/converter.go that fetches HTML from a URL using the existing HTTP client pattern (30s timeout, 10 redirect limit, User-Agent header), validates the response is not empty, converts HTML to markdown using htmltomarkdown.Convert(), validates the converted markdown is not empty, and returns categorized errors for all failure modes.

Added a new error category ErrCategoryConversionFailure to fetcher.go for conversion failures, along with a corresponding user-friendly error message in the UserMessage() method.

Created comprehensive unit tests in pkg/rss/converter_test.go covering:
- 12 successful conversion scenarios with various HTML structures (headings, paragraphs, lists, links, code blocks, blockquotes, images, nested structures)
- 5 HTTP error scenarios (404, 403, 500, 503, 410)
- Timeout handling
- Empty and whitespace-only response handling
- Network error handling
- Too many redirects handling
- User message verification

During testing, fixed two issues: removed a duplicate GetArticleByID() method in pkg/database/article.go and added the missing ErrCategoryConversionFailure constant to fetcher.go.

All 22 sub-tests pass, full test suite passes (4 packages), and go vet reports no issues.

## Verification

All tests pass: go test ./... (4 packages). All 9 test functions with 22 sub-tests in TestFetchAndConvertArticle pass. go vet ./... reports no issues. go build ./pkg/rss succeeds.

## Requirements Advanced

- R009 — Error categorization implemented with HTTPError type returning Category and UserMessage() for all failure modes (404, 403, 500, 503, 410, timeout, empty content, network errors, too many redirects). Verified via dedicated test cases.

## Requirements Validated

- R008 — FetchAndConvertArticle() implemented in converter.go, fetches HTML from source URL and converts to markdown using htmltomarkdown.Convert(). Verified via 22 unit tests covering success scenarios with various HTML structures.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

None.

## Follow-ups

None.

## Files Created/Modified

None.
