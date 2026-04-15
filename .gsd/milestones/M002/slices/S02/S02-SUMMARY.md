---
id: S02
parent: M002
milestone: M002
provides:
  - (none)
requires:
  []
affects:
  []
key_files:
  - ["pkg/database/article.go", "pkg/database/article_test.go", "cmd/rss-cli/article_fetch_cmd.go", "cmd/rss-cli/article_cmd.go"]
key_decisions:
  - (none)
patterns_established:
  - ["Article caching pattern: fetch-then-cache workflow with UpdateArticleContent() method", "Database update pattern for optional content fields"]
observability_surfaces:
  - none
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-04-15T02:21:10.958Z
blocker_discovered: false
---

# S02: Article fetch command

**Implemented `article fetch [id]` CLI command that fetches HTML from source URL, converts to markdown, displays in terminal, and caches result to database**

## What Happened

Implemented the article fetch [id] CLI command across three tasks:

**T01** added the UpdateArticleContent(id int, content string) error method to pkg/database/article.go, enabling caching of fetched markdown content. The method uses a parameterized UPDATE query and includes unit tests verifying both successful updates and behavior with non-existent articles.

**T02** implemented the article fetch [id] CLI command in cmd/rss-cli/article_fetch_cmd.go following the established articleViewCmd pattern. The command fetches the article from the database, validates it has a URL, calls rss.FetchAndConvertArticle(url) to fetch and convert HTML to markdown, caches the result using UpdateArticleContent(), and outputs via ui.Printer with --json/--text mode support.

**T03** registered the articleFetchCmd subcommand in article_cmd.go's init() function by adding articleCmd.AddCommand(articleFetchCmd). Removed duplicate init() from article_fetch_cmd.go to follow the centralized registration pattern.

All three tasks completed successfully with passing tests and builds. The implementation reuses existing patterns from M001 (CLI command structure, database access, HTML-to-markdown conversion) and integrates seamlessly with the rss.FetchAndConvertArticle() function.

## Verification

Build verification: go build -o rss-cli ./cmd/rss-cli completed successfully. Command registration verified via ./rss-cli article --help showing 'fetch' command with proper description. All package tests pass: go test ./... executed successfully with 0 failures. Database tests include TestUpdateArticleContent and TestUpdateArticleContent_NonExistentArticle. RSS package tests include 12+ converter tests verifying HTML-to-markdown conversion for various HTML structures (headings, paragraphs, lists, links, code blocks, blockquotes, images).

## Requirements Advanced

None.

## Requirements Validated

- R007 — Implemented article fetch [id] command with markdown conversion and database caching
- R010 — UpdateArticleContent() method caches fetched content to database content field
- R011 — Markdown output displayed as-is via ui.Printer with --text mode support

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
