---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Implement article view command with --open flag

Implement the article view command with --open flag support. This task creates the articleViewCmd subcommand that retrieves an article by ID from the database, displays its content, and optionally opens the article URL in the default browser when --open flag is provided.

## Inputs

- `pkg/database/article.go - provides GetArticleByID method`
- `go.mod - already has github.com/toqueteos/webbrowser dependency`

## Expected Output

- `cmd/rss-cli/article_cmd.go`

## Verification

go build -o rss-cli ./cmd/rss-cli && ./rss-cli article view --help | grep -q 'open'
