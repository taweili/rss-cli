---
estimated_steps: 6
estimated_files: 1
skills_used: []
---

# T04: Add unit tests

Add unit tests for HTML-to-Markdown conversion and HTTP error handling.

Steps:
1. Create pkg/rss/converter_test.go
2. Add table-driven tests for HTML-to-Markdown conversion
3. Add tests for FetchArticleContent error cases
4. Run tests and verify all pass

## Inputs

- `pkg/rss/fetcher.go`

## Expected Output

- `pkg/rss/converter_test.go`

## Verification

go test ./pkg/rss/... -v
