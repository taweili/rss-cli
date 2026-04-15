---
estimated_steps: 7
estimated_files: 2
skills_used: []
---

# T01: Implement FetchAndConvertArticle() converter function

Create the core converter implementation in pkg/rss/converter.go. The function FetchAndConvertArticle(url string) (string, error) should:
1. Fetch HTML using the same HTTP client pattern as FetchArticleContent() (30s timeout, 10 redirect limit, User-Agent header)
2. Validate the HTML response is not empty
3. Convert HTML to markdown using htmltomarkdown.Convert()
4. Validate the converted markdown is not empty
5. Return categorized errors (HTTPError with appropriate categories) for all failure modes

Reuse the existing HTTPError type and error categorization constants from fetcher.go. Add a new error category ErrCategoryConversionFailure if the HTML converts to empty markdown.

## Inputs

- ``pkg/rss/fetcher.go` — Reference for HTTP client pattern and HTTPError type`
- ``pkg/rss/converter_test.go` — Reference for htmltomarkdown.Convert() usage`

## Expected Output

- ``pkg/rss/converter.go` — New file with FetchAndConvertArticle() function implementation`

## Verification

go build ./pkg/rss && test -f pkg/rss/converter.go && grep -q 'func FetchAndConvertArticle' pkg/rss/converter.go

## Observability Impact

None — internal library function, no new runtime boundaries or async flows
