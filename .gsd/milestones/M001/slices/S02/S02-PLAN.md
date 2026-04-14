# S02: Browser open flag

**Goal:** Add --open flag to article view command that opens the article URL in the default browser
**Demo:** 

## Must-Haves

- Not provided.

## Proof Level

- This slice proves: Not provided.

## Integration Closure

Not provided.

## Verification

- Not provided.

## Tasks

- [x] **T01: Implement article view command with --open flag** `est:45m`
  Implement the article view command with --open flag support. This task creates the articleViewCmd subcommand that retrieves an article by ID from the database, displays its content, and optionally opens the article URL in the default browser when --open flag is provided.
  - Files: `cmd/rss-cli/article_cmd.go`, `pkg/database/article.go`
  - Verify: go build -o rss-cli ./cmd/rss-cli && ./rss-cli article view --help | grep -q 'open'

## Files Likely Touched

- cmd/rss-cli/article_cmd.go
- pkg/database/article.go
