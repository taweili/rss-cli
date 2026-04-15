---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Add UpdateArticleContent() database method

Add UpdateArticleContent(id int, content string) error method to pkg/database/article.go to update the content field of an existing article. This enables caching fetched markdown after successful fetch.

## Inputs

- `pkg/database/article.go`

## Expected Output

- `pkg/database/article.go`

## Verification

go test ./pkg/database -run TestUpdateArticleContent -v
