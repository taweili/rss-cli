---
id: S01
parent: M001
milestone: M001
provides:
  - Working article view command that fetches, converts, and displays full article content in terminal.
requires:
  []
affects:
  []
key_files:
  - ["pkg/database/article.go", "pkg/rss/fetcher.go", "cmd/rss-cli/article_cmd.go", "pkg/ui/output.go", "pkg/rss/converter_test.go"]
key_decisions:
  - ["Used html-to-markdown library for conversion instead of custom HTML stripping", "Fetch article content on-demand rather than caching in database", "Display full article at once, let user pipe to less if needed"]
patterns_established:
  - ["Reuse existing HTTP client patterns for article fetching", "Use html-to-markdown library for structure-preserving conversion", "Distinguish HTTP error codes (404, 403) with specific error messages", "Support both JSON and text output modes consistently"]
observability_surfaces:
  - none
drill_down_paths:
  - ["milestones/M001/slices/S01/tasks/T01-SUMMARY.md", "milestones/M001/slices/S01/tasks/T02-SUMMARY.md", "milestones/M001/slices/S01/tasks/T03-SUMMARY.md", "milestones/M001/slices/S01/tasks/T04-SUMMARY.md"]
duration: ""
verification_result: passed
completed_at: 2026-04-14T14:28:41.516Z
blocker_discovered: false
---

# S01: Article view command

**Implemented article view command with HTML-to-Markdown conversion, on-demand fetching, and comprehensive error handling.**

## What Happened

Implemented the article view command across four tasks:

1. **Database layer** (T01): Added GetArticleByID() method following existing query patterns.

2. **RSS fetcher layer** (T02): Added FetchArticleContent() function reusing HTTP client pattern (30s timeout, user-agent, redirect limit).

3. **CLI command layer** (T03): Added articleViewCmd with html-to-markdown conversion, comprehensive error handling for HTTP errors (404, 403, timeout), and support for both JSON and text output modes.

4. **Tests** (T04): Added unit tests for HTML-to-Markdown conversion and HTTP error handling.

All tests pass. Manual testing confirmed the command works with real articles from Slashdot in both output modes.

## Verification

All tests pass (go test ./...). Manual testing confirmed:
- ./rss-cli article view 17494 displays article content in JSON mode
- ./rss-cli article view 17494 --text displays article in text mode
- ./rss-cli article view 99999 shows 'Failed to retrieve article' error
- ./rss-cli article view invalid shows 'Invalid article ID' error

## Requirements Advanced

- R001 — Implemented article view command that fetches and displays full article content
- R003 — Used html-to-markdown library that preserves headings, lists, links, and formatting
- R005 — Fetches content from source URL on-demand without database caching

## Requirements Validated

None.

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

- `pkg/database/article.go` — Added GetArticleByID() method to retrieve single article by ID
- `pkg/rss/fetcher.go` — Added FetchArticleContent() function and io import for fetching raw HTML
- `cmd/rss-cli/article_cmd.go` — Added articleViewCmd with HTML-to-Markdown conversion and error handling
- `pkg/ui/output.go` — Added outputText handling for article view content with metadata
- `pkg/rss/converter_test.go` — Added unit tests for HTML-to-Markdown conversion and HTTP error handling
