---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T01: Run automated test suite and build verification

Run the full test suite and build verification to confirm all automated tests pass and the CLI builds successfully.

## Inputs

- `pkg/rss/converter.go`
- `pkg/rss/fetcher.go`
- `pkg/database/article.go`
- `cmd/rss-cli/article_fetch_cmd.go`

## Expected Output

- `Verification evidence captured in task summary`

## Verification

go test ./... && go build -o rss-cli ./cmd/rss-cli
