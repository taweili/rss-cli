---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Implement article fetch command

Create cmd/rss-cli/article_fetch_cmd.go implementing the article fetch [id] command following the articleViewCmd pattern: parse article ID, fetch from database, call rss.FetchAndConvertArticle(url), update database with result, output markdown via ui.Printer with --json/--text mode support.

## Inputs

- `cmd/rss-cli/article_cmd.go`
- `pkg/rss/converter.go`
- `pkg/database/article.go`

## Expected Output

- `cmd/rss-cli/article_fetch_cmd.go`

## Verification

go build -o rss-cli ./cmd/rss-cli
