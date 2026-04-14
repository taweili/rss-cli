---
id: M001
title: "Article viewing for daily news reading"
status: complete
completed_at: 2026-04-14T15:29:16.166Z
key_decisions:
  - Used github.com/kreuzberg-dev/html-to-markdown for HTML-to-Markdown conversion (D001)
  - Used github.com/toqueteos/webbrowser for cross-platform browser opening (D002)
  - Fetch article content on-demand rather than caching in database (D003)
  - Display full article at once, let user pipe to less if needed
  - Used errors.As() pattern to preserve redirect error categorization when wrapping errors
  - Detached browser processes (nil stdout/stderr) to prevent CLI from blocking
key_files:
  - cmd/rss-cli/article_cmd.go
  - cmd/rss-cli/article_open_test.go
  - pkg/database/article.go
  - pkg/rss/fetcher.go
  - pkg/rss/fetcher_test.go
  - pkg/rss/error_test.go
  - pkg/ui/output.go
  - go.mod
lessons_learned:
  - html-to-markdown library with CGO dependency works well but requires CGO-enabled build environment
  - Error categorization with UserMessage() method provides excellent user experience while maintaining structured error handling
  - Detaching browser processes is essential to prevent CLI from blocking while browser is open
  - Table-driven tests provide comprehensive coverage for HTTP status code categorization
  - Platform-specific browser detection benefits from fallback chain: webbrowser library → platform-specific commands
---

# M001: Article viewing for daily news reading

**Implemented complete article viewing capability with HTML-to-Markdown conversion, browser integration, and comprehensive error handling for RSS CLI.**

## What Happened

This milestone implemented the missing piece for daily RSS reading: the ability to fetch and display full article content in the terminal.

Three slices delivered the complete capability:

**S01 (Article view command)** implemented the core functionality: database layer with GetArticleByID(), RSS fetcher with FetchArticleContent() reusing HTTP client patterns (30s timeout, user-agent), CLI command with html-to-markdown conversion, and comprehensive error handling. The command supports both JSON and text output modes, displaying full article content fetched on-demand from source URLs.

**S02 (Browser open flag)** added the --open flag to article view command, using github.com/toqueteos/webbrowser library for cross-platform browser opening. The implementation follows established patterns with RunE, cobra.ExactArgs, and consistent JSON output.

**S03 (Error handling & edge cases)** enhanced error categorization with 10 distinct error categories (NotFound, AccessDenied, Timeout, ParseFailure, EmptyContent, TooManyRedirects, ServiceUnavailable, ServerError, NetworkError), added UserMessage() method for user-friendly output, implemented article open command with cross-platform browser detection and --browser flag, and created 25+ comprehensive tests across error_test.go and article_open_test.go. Manual testing against httpbin.org confirmed proper error categorization for 404, 403, and 500 errors.

All tests pass (go test ./...), build succeeds, and manual testing confirmed the commands work with real articles from Slashdot in both output modes.

## Success Criteria Results

**S01 Success Criteria:**
- ✅ `article view [id]` fetches HTML from the article's stored link URL - Implemented via FetchArticleContent() in pkg/rss/fetcher.go, verified via manual testing with real articles
- ✅ HTML is converted to Markdown with structure preserved - Using github.com/kreuzberg-dev/html-to-markdown library, verified via unit tests in converter_test.go
- ✅ Markdown is displayed to stdout - Implemented in cmd/rss-cli/article_cmd.go with both JSON and text output modes
- ✅ Works with JSON and text output modes - Verified via manual testing: `./rss-cli article view 17494 --json` and `./rss-cli article view 17494 --text`

**S02 Success Criteria:**
- ✅ `article view [id] --open` opens the article URL in the default browser - Implemented using github.com/toqueteos/webbrowser library, verified via manual testing
- ✅ Command exits successfully after opening browser - Browser processes detached to prevent blocking, command returns immediately
- ✅ Works cross-platform (macOS, Linux, Windows) - webbrowser library handles OS-specific detection, fallbacks implemented for reliability

**S03 Success Criteria:**
- ✅ Network timeout shows clear error message - isTimeoutError() helper detects timeout, UserMessage() returns "Connection timed out"
- ✅ HTTP 404 shows "Article not found" error - httpErr() categorizes 404/410 as NotFound, verified via httpbin.org testing
- ✅ HTTP 403 shows "Access denied — may be paywalled" error - 403 categorized as AccessDenied with paywall message
- ✅ Parse failures show clear error message - ParseFailure category with descriptive message
- ✅ All errors use printer.Error() for consistent formatting - All error paths use printer.Error() for JSON/text output consistency

## Definition of Done Results

**Definition of Done:**
- ✅ All slices are complete - S01, S02, S03 all marked complete with all tasks done (9/9 tasks)
- ✅ All slice summaries exist - S01-SUMMARY.md, S02-SUMMARY.md, S03-SUMMARY.md all present
- ✅ Cross-slice integration points work correctly - Article view command integrates database lookup, HTTP fetching, HTML conversion, and browser opening seamlessly
- ✅ All tests pass - go test ./... passes with 25+ tests across error_test.go, article_open_test.go, converter_test.go
- ✅ Build succeeds - go build -o rss-cli ./cmd/rss-cli completes without errors
- ✅ Manual testing confirmed - Real article viewing, browser opening, and error handling all verified

## Requirement Outcomes

**Requirement Status Transitions:**

- R001: active → validated — S01 implemented article view command with HTML-to-Markdown conversion, on-demand fetching, and comprehensive error handling. Verified via unit tests and manual testing with real articles.

- R002: active → validated — S02 added --open flag using webbrowser library for cross-platform browser opening. Verified via manual testing and unit tests.

- R003: active → validated — S01 used html-to-markdown library that preserves headings, lists, links, and formatting. Verified via unit tests with sample HTML and manual testing.

- R004: active → validated — S03 implemented comprehensive error categorization with 10 error categories and user-friendly messages. Verified via 25+ tests and manual testing against httpbin.org.

- R005: active → validated — S01 fetches content from source URL on-demand without database caching. Verified via code review showing no database storage of article content.

## Deviations

None.

## Follow-ups

None.
