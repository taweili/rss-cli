---
id: T02
parent: S02
milestone: M002
key_files:
  - cmd/rss-cli/article_fetch_cmd.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-04-15T02:17:24.648Z
blocker_discovered: false
---

# T02: Implement article fetch [id] CLI command with markdown conversion and caching

**Implement article fetch [id] CLI command with markdown conversion and caching**

## What Happened

Implemented the `article fetch [id]` CLI command following the existing articleViewCmd pattern. The command:

1. Parses the article ID from command-line arguments
2. Fetches the article from the SQLite database using GetArticleByID()
3. Validates that the article has a URL
4. Calls rss.FetchAndConvertArticle(url) to fetch HTML from the source URL and convert it to markdown
5. Caches the markdown result to the database using UpdateArticleContent()
6. Outputs the result via ui.Printer with support for both --json and --text modes

The implementation reuses existing patterns from article_cmd.go and integrates with the rss.FetchAndConvertArticle() function from pkg/rss/converter.go and UpdateArticleContent() from pkg/database/article.go (added in T01).

Verified:
- Build succeeds: `go build -o rss-cli ./cmd/rss-cli`
- Command is properly registered in the CLI help
- All existing tests pass: `go test ./...`

## Verification

Build verification: `go build -o rss-cli ./cmd/rss-cli` completed successfully with no errors. Command registration verified via `./rss-cli article --help` showing 'fetch' command. All package tests pass via `go test ./...`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `go build -o rss-cli ./cmd/rss-cli` | 0 | ✅ pass | 2500ms |
| 2 | `./rss-cli article --help` | 0 | ✅ pass | 50ms |
| 3 | `go test ./...` | 0 | ✅ pass | 1533ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `cmd/rss-cli/article_fetch_cmd.go`
