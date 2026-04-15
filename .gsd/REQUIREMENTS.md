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

### R007 — User can run `article fetch [id]` to fetch full article HTML from source URL, convert to markdown, and display in terminal
- Class: core-capability
- Status: validated
- Description: User can run `article fetch [id]` to fetch full article HTML from source URL, convert to markdown, and display in terminal
- Why it matters: Provides a dedicated command for fetching full article content with explicit caching behavior
- Source: user
- Primary owning slice: M002/S02
- Supporting slices: M002/S01
- Validation: Implemented in S02: article fetch [id] command fetches HTML from source URL, converts to Markdown using rss.FetchAndConvertArticle(), displays in terminal, and caches result to database using UpdateArticleContent(). Supports both --json and --text output modes. Verified via build tests, command registration (./rss-cli article fetch --help), and manual testing with real articles from The Register (IDs 19531, 19532).
- Notes: New command separate from `article view` for clarity

### R008 — Fetch full article HTML from the article's source URL and convert it to readable markdown
- Class: core-capability
- Status: validated
- Description: Fetch full article HTML from the article's source URL and convert it to readable markdown
- Why it matters: Users get the complete article content, not just RSS feed snippets
- Source: user
- Primary owning slice: M002/S01
- Supporting slices: none
- Validation: Implemented in S01: FetchAndConvertArticle(url) fetches HTML from source URL using existing HTTP client (30s timeout, 10 redirect limit), converts to markdown using htmltomarkdown.Convert(), validates non-empty output. Verified via 12+ converter tests covering various HTML structures (headings, paragraphs, lists, links, code blocks, blockquotes, images) and manual fetch of articles 19531 and 19532 returning full markdown content with preserved structure.
- Notes: Uses html-to-markdown library for conversion

### R009 — Network errors, HTTP errors (404, 403), timeouts, and conversion failures show clear, actionable error messages
- Class: failure-visibility
- Status: validated
- Description: Network errors, HTTP errors (404, 403), timeouts, and conversion failures show clear, actionable error messages
- Why it matters: User needs to know why fetching failed and whether to retry
- Source: user
- Primary owning slice: M002/S01
- Supporting slices: none
- Validation: Implemented in S01: HTTPError type extended with ErrCategoryConversionFailure category and UserMessage() method for all failure modes (404, 403, 500, 503, 410, timeout, empty content, network errors, too many redirects, conversion failures). Verified via 25+ error categorization tests in error_test.go and converter_test.go, plus manual verification showing proper error messages for failed feeds.
- Notes: Reuse existing HTTPError pattern with categories

### R010 — After first successful fetch, article content is cached in the database content field for subsequent views
- Class: constraint
- Status: validated
- Description: After first successful fetch, article content is cached in the database content field for subsequent views
- Why it matters: Avoids repeated network requests, improves performance for re-reading
- Source: user
- Primary owning slice: M002/S02
- Supporting slices: none
- Validation: Implemented in S02: After fetching article content via rss.FetchAndConvertArticle(), the markdown result is cached to the database content field using the UpdateArticleContent() method added in T01. Uses parameterized UPDATE query. Verified via unit tests in article_test.go (TestUpdateArticleContent, TestUpdateArticleContent_NonExistentArticle) and integration with the article fetch command showing content persisted after first fetch.
- Notes: Updates existing content field, supersedes R005 no-cache constraint

### R011 — Converted markdown is displayed as-is in the terminal, preserving formatting
- Class: quality-attribute
- Status: validated
- Description: Converted markdown is displayed as-is in the terminal, preserving formatting
- Why it matters: Users see the full formatted article with headings, lists, links, etc.
- Source: user
- Primary owning slice: M002/S02
- Supporting slices: none
- Validation: Implemented in S02: Markdown output is displayed as-is via ui.Printer with --text mode support, preserving headings, lists, links, code blocks, and formatting. Verified via integration with existing converter tests that confirm structured markdown output from HTML conversion, and manual testing showing full formatted article content in terminal.
- Notes: No plain-text stripping, raw markdown output

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
| R007 | core-capability | validated | M002/S02 | M002/S01 | Implemented in S02: article fetch [id] command fetches HTML from source URL, converts to Markdown using rss.FetchAndConvertArticle(), displays in terminal, and caches result to database using UpdateArticleContent(). Supports both --json and --text output modes. Verified via build tests, command registration (./rss-cli article fetch --help), and manual testing with real articles from The Register (IDs 19531, 19532). |
| R008 | core-capability | validated | M002/S01 | none | Implemented in S01: FetchAndConvertArticle(url) fetches HTML from source URL using existing HTTP client (30s timeout, 10 redirect limit), converts to markdown using htmltomarkdown.Convert(), validates non-empty output. Verified via 12+ converter tests covering various HTML structures (headings, paragraphs, lists, links, code blocks, blockquotes, images) and manual fetch of articles 19531 and 19532 returning full markdown content with preserved structure. |
| R009 | failure-visibility | validated | M002/S01 | none | Implemented in S01: HTTPError type extended with ErrCategoryConversionFailure category and UserMessage() method for all failure modes (404, 403, 500, 503, 410, timeout, empty content, network errors, too many redirects, conversion failures). Verified via 25+ error categorization tests in error_test.go and converter_test.go, plus manual verification showing proper error messages for failed feeds. |
| R010 | constraint | validated | M002/S02 | none | Implemented in S02: After fetching article content via rss.FetchAndConvertArticle(), the markdown result is cached to the database content field using the UpdateArticleContent() method added in T01. Uses parameterized UPDATE query. Verified via unit tests in article_test.go (TestUpdateArticleContent, TestUpdateArticleContent_NonExistentArticle) and integration with the article fetch command showing content persisted after first fetch. |
| R011 | quality-attribute | validated | M002/S02 | none | Implemented in S02: Markdown output is displayed as-is via ui.Printer with --text mode support, preserving headings, lists, links, code blocks, and formatting. Verified via integration with existing converter tests that confirm structured markdown output from HTML conversion, and manual testing showing full formatted article content in terminal. |

## Coverage Summary

- Active requirements: 0
- Mapped to slices: 0
- Validated: 10 (R001, R002, R003, R004, R005, R007, R008, R009, R010, R011)
- Unmapped active requirements: 0
