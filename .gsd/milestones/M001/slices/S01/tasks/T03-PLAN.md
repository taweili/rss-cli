---
estimated_steps: 8
estimated_files: 3
skills_used: []
---

# T03: Implement article view command

Implement article view command in cmd/rss-cli/article_cmd.go.

Steps:
1. Add html-to-markdown dependency
2. Add articleViewCmd following existing command patterns
3. Implement: lookup article, fetch HTML, convert to markdown, output result
4. Add error handling for HTTP errors (404, 403, timeout)
5. Update pkg/ui/output.go for article view output format
6. Test manually with real articles

## Inputs

- `cmd/rss-cli/article_cmd.go`
- `pkg/ui/output.go`

## Expected Output

- `cmd/rss-cli/article_cmd.go with articleViewCmd`
- `pkg/ui/output.go updated`

## Verification

go build -o rss-cli ./cmd/rss-cli && ./rss-cli article view [id]
