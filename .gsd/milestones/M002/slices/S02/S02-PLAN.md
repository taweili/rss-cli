# S02: Article fetch command

**Goal:** Implement `article fetch [id]` CLI command that fetches full article HTML from source URL, converts to markdown using FetchAndConvertArticle(), displays markdown in terminal, and caches result to database content field
**Demo:** article fetch [id] displays markdown in terminal and caches to database

## Must-Haves

- Not provided.

## Proof Level

- This slice proves: Not provided.

## Integration Closure

Not provided.

## Verification

- Not provided.

## Tasks

- [x] **T01: Add UpdateArticleContent() database method** `est:15m`
  Add UpdateArticleContent(id int, content string) error method to pkg/database/article.go to update the content field of an existing article. This enables caching fetched markdown after successful fetch.
  - Files: `pkg/database/article.go`
  - Verify: go test ./pkg/database -run TestUpdateArticleContent -v

- [x] **T02: Implement article fetch command** `est:45m`
  Create cmd/rss-cli/article_fetch_cmd.go implementing the article fetch [id] command following the articleViewCmd pattern: parse article ID, fetch from database, call rss.FetchAndConvertArticle(url), update database with result, output markdown via ui.Printer with --json/--text mode support.
  - Files: `cmd/rss-cli/article_fetch_cmd.go`, `pkg/rss/converter.go`
  - Verify: go build -o rss-cli ./cmd/rss-cli

- [x] **T03: Register article fetch command** `est:10m`
  Register articleFetchCmd subcommand in article_cmd.go init() function and add the command to the articleCmd.AddCommand() call.
  - Files: `cmd/rss-cli/article_cmd.go`
  - Verify: go build -o rss-cli ./cmd/rss-cli && ./rss-cli article --help | grep -q fetch

## Files Likely Touched

- pkg/database/article.go
- cmd/rss-cli/article_fetch_cmd.go
- pkg/rss/converter.go
- cmd/rss-cli/article_cmd.go
