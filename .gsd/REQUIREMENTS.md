# Requirements

This file is the explicit capability and coverage contract for the project.

## Validated

### R001 — User can run `article view [id]` to fetch and display full article content in the terminal
- Class: core-capability
- Status: validated
- Description: User can run `article view [id]` to fetch and display full article content in the terminal
- Why it matters: Enables actual news reading without leaving the terminal
- Source: user
- Primary owning slice: M001/S01
- Supporting slices: none
- Validation: Implemented in S01: article view [id] command fetches HTML from article URL, converts to Markdown using html-to-markdown library, displays in terminal. Verified via unit tests (converter_test.go) and manual testing with real articles from Slashdot. Supports both JSON and text output modes.
- Notes: Fetches from stored link URL on-demand, converts HTML to Markdown

### R002 — User can run `article view [id] --open` to open the article URL in their default browser
- Class: core-capability
- Status: validated
- Description: User can run `article view [id] --open` to open the article URL in their default browser
- Why it matters: Fallback for articles that don't render well in terminal, or when user prefers browser
- Source: user
- Primary owning slice: M001/S02
- Supporting slices: none
- Validation: Implemented in S02: article view [id] --open flag opens article URL in default browser using github.com/toqueteos/webbrowser library. Verified via manual testing and unit tests (article_open_test.go). Cross-platform support confirmed.
- Notes: Uses cross-platform browser opening library

### R003 — Conversion preserves headings, paragraphs, lists, links, code blocks, and basic formatting
- Class: quality-attribute
- Status: validated
- Description: Conversion preserves headings, paragraphs, lists, links, code blocks, and basic formatting
- Why it matters: Articles must remain readable and structured, not just plain text
- Source: user
- Primary owning slice: M001/S01
- Supporting slices: none
- Validation: Implemented in S01: html-to-markdown library preserves headings, paragraphs, lists, links, code blocks, and basic formatting. Verified via unit tests in converter_test.go with sample HTML containing various elements. Manual testing confirmed readable, structured output.
- Notes: Uses github.com/kreuzberg-dev/html-to-markdown library

### R004 — Network errors, HTTP errors (404, 403), and parse failures show clear, actionable error messages
- Class: failure-visibility
- Status: validated
- Description: Network errors, HTTP errors (404, 403), and parse failures show clear, actionable error messages
- Why it matters: User needs to know why an article didn't load and whether to retry
- Source: user
- Primary owning slice: M001/S03
- Supporting slices: none
- Validation: Implemented in S03: Comprehensive error categorization with 10 error categories (NotFound, AccessDenied, Timeout, ParseFailure, EmptyContent, TooManyRedirects, ServiceUnavailable, ServerError, NetworkError). UserMessage() method provides user-friendly messages. Verified via 25+ tests in error_test.go and manual testing against httpbin.org for 404/403/500 errors.
- Notes: Timeout, 404, 403/paywall, parse errors all have distinct messages

### R005 — Full article content is fetched from the source URL when viewing, not cached in the database
- Class: constraint
- Status: validated
- Description: Full article content is fetched from the source URL when viewing, not cached in the database
- Why it matters: Avoids storage bloat, always shows current content, respects publisher bandwidth
- Source: user
- Primary owning slice: M001/S01
- Supporting slices: none
- Validation: Implemented in S01: Article content fetched on-demand from source URL when viewing via FetchArticleContent() in pkg/rss/fetcher.go. No database caching - content not stored in database. Verified via code review and manual testing showing fresh content fetched each time.
- Notes: No database schema changes needed

## Out of Scope

### R006 — Export feeds to OPML file format
- Class: anti-feature
- Status: out-of-scope
- Description: Export feeds to OPML file format
- Why it matters: Explicitly excluded to prevent scope creep
- Source: user
- Primary owning slice: none
- Supporting slices: none
- Validation: n/a
- Notes: User confirmed not needed for this milestone

## Traceability

| ID | Class | Status | Primary owner | Supporting | Proof |
|---|---|---|---|---|---|
| R001 | core-capability | validated | M001/S01 | none | Implemented in S01: article view [id] command fetches HTML from article URL, converts to Markdown using html-to-markdown library, displays in terminal. Verified via unit tests (converter_test.go) and manual testing with real articles from Slashdot. Supports both JSON and text output modes. |
| R002 | core-capability | validated | M001/S02 | none | Implemented in S02: article view [id] --open flag opens article URL in default browser using github.com/toqueteos/webbrowser library. Verified via manual testing and unit tests (article_open_test.go). Cross-platform support confirmed. |
| R003 | quality-attribute | validated | M001/S01 | none | Implemented in S01: html-to-markdown library preserves headings, paragraphs, lists, links, code blocks, and basic formatting. Verified via unit tests in converter_test.go with sample HTML containing various elements. Manual testing confirmed readable, structured output. |
| R004 | failure-visibility | validated | M001/S03 | none | Implemented in S03: Comprehensive error categorization with 10 error categories (NotFound, AccessDenied, Timeout, ParseFailure, EmptyContent, TooManyRedirects, ServiceUnavailable, ServerError, NetworkError). UserMessage() method provides user-friendly messages. Verified via 25+ tests in error_test.go and manual testing against httpbin.org for 404/403/500 errors. |
| R005 | constraint | validated | M001/S01 | none | Implemented in S01: Article content fetched on-demand from source URL when viewing via FetchArticleContent() in pkg/rss/fetcher.go. No database caching - content not stored in database. Verified via code review and manual testing showing fresh content fetched each time. |
| R006 | anti-feature | out-of-scope | none | none | n/a |

## Coverage Summary

- Active requirements: 0
- Mapped to slices: 0
- Validated: 5 (R001, R002, R003, R004, R005)
- Unmapped active requirements: 0
