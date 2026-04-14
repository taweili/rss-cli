---
estimated_steps: 7
estimated_files: 1
skills_used: []
---

# T02: Add HTML fetching function

Add FetchArticleContent() function to pkg/rss/fetcher.go for fetching raw HTML from article URLs.

Steps:
1. Read existing fetcher.go to understand HTTP client pattern
2. Add FetchArticleContent(url string) (string, error) function
3. Reuse 30s timeout, user-agent, redirect limit
4. Add io import for ReadAll
5. Test with go build

## Inputs

- `pkg/rss/fetcher.go`

## Expected Output

- `pkg/rss/fetcher.go with FetchArticleContent function`

## Verification

go build ./... && go test ./pkg/rss/...
