---
estimated_steps: 8
estimated_files: 2
skills_used: []
---

# T03: Write comprehensive error handling tests

Write comprehensive unit tests for error handling scenarios.

Steps:
1. Add TestHTTPError_Categorization table-driven test for all status codes
2. Add TestFetchArticleContent_EmptyResponse test
3. Add TestDetectBrowser tests for platform detection
4. Add TestOpenBrowser_Failure tests for invalid URLs
5. Add integration test with mock HTTP server for 404, 403, 410, timeout scenarios
6. Ensure all tests verify error messages match expected format

## Inputs

- ``pkg/rss/fetcher.go``
- ``cmd/rss-cli/article_cmd.go``

## Expected Output

- ``pkg/rss/error_test.go``
- ``cmd/rss-cli/article_open_test.go``

## Verification

go test -v ./... -run 'TestHTTP|TestFetch|TestDetect|TestOpen'

## Observability Impact

none
