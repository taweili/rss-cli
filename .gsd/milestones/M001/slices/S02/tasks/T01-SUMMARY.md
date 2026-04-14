---
id: T01
parent: S02
milestone: M001
key_files:
  - cmd/rss-cli/article_cmd.go
  - pkg/database/article.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-04-14T14:39:13.157Z
blocker_discovered: false
---

# T01: Implement article view command with --open flag

**Implement article view command with --open flag**

## What Happened

Implemented the article view command with --open flag support.

Changes made:
1. Added `GetArticleByID(id int) (*Article, error)` method to `pkg/database/article.go` to retrieve a single article by its ID from the database.

2. Added `articleViewCmd` subcommand to `cmd/rss-cli/article_cmd.go`:
   - Takes an article ID as a required positional argument
   - Retrieves the article from the database using GetArticleByID
   - Supports `--open` flag that opens the article URL in the default browser using the webbrowser library
   - Outputs the article as JSON (consistent with other commands in the CLI)

3. Registered the new command and flag in the init() function.

Verification:
- Built the CLI successfully with `go build -o rss-cli ./cmd/rss-cli`
- Verified the --open flag appears in help output
- The command follows the existing patterns in the codebase for error handling and JSON output

## Verification

Built CLI and verified --open flag is present in article view --help output. Command: go build -o rss-cli ./cmd/rss-cli && ./rss-cli article view --help | grep -q 'open'

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `go build -o rss-cli ./cmd/rss-cli && ./rss-cli article view --help | grep -q 'open'` | 0 | ✅ pass | 5000ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `cmd/rss-cli/article_cmd.go`
- `pkg/database/article.go`
