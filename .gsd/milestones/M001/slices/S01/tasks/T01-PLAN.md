---
estimated_steps: 6
estimated_files: 1
skills_used: []
---

# T01: Add database method to get article by ID

Add GetArticleByID() method to pkg/database/article.go to retrieve a single article by ID.

Steps:
1. Read existing article.go to understand query patterns
2. Add GetArticleByID(id int) (*Article, error) function
3. Follow existing pattern with timestamp parsing
4. Test with go build

## Inputs

- `pkg/database/article.go`

## Expected Output

- `pkg/database/article.go with GetArticleByID method`

## Verification

go build ./... && go test ./pkg/database/...
