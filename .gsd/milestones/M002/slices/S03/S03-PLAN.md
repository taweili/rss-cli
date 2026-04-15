# S03: Tests and verification

**Goal:** Verify all tests pass and manually verify article fetch command works against real RSS feeds
**Demo:** All tests pass, manual verification against real feeds succeeds

## Must-Haves

- Not provided.

## Proof Level

- This slice proves: Not provided.

## Integration Closure

Not provided.

## Verification

- Not provided.

## Tasks

- [x] **T01: Run automated test suite and build verification** `est:15m`
  Run the full test suite and build verification to confirm all automated tests pass and the CLI builds successfully.
  - Files: `pkg/rss/converter_test.go`, `pkg/rss/error_test.go`, `pkg/database/article_test.go`, `cmd/rss-cli/article_fetch_cmd.go`
  - Verify: go test ./... && go build -o rss-cli ./cmd/rss-cli

- [x] **T02: Manual verification against real RSS feeds** `est:30m`
  Manually test the article fetch command against 2-3 real RSS feeds to prove end-to-end functionality: fetching HTML, converting to markdown, caching to database, and displaying formatted output.
  - Verify: ./rss-cli feed update-all && ./rss-cli article list --limit 5 && ./rss-cli article fetch [id]

## Files Likely Touched

- pkg/rss/converter_test.go
- pkg/rss/error_test.go
- pkg/database/article_test.go
- cmd/rss-cli/article_fetch_cmd.go
