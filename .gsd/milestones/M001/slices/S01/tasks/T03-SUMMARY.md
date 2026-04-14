---
id: T03
parent: S01
milestone: M001
key_files:
  - cmd/rss-cli/article_cmd.go
  - pkg/ui/output.go
  - go.mod
key_decisions:
  - (none)
duration: 
verification_result: untested
completed_at: 2026-04-14T14:28:16.081Z
blocker_discovered: false
---

# T03: Implemented article view command with HTML-to-Markdown conversion.

**Implemented article view command with HTML-to-Markdown conversion.**

## What Happened

Implemented articleViewCmd in cmd/rss-cli/article_cmd.go with html-to-markdown conversion. Added html-to-markdown dependency via go get. Command looks up article by ID, fetches HTML, converts to Markdown, and outputs result in JSON or text mode. Extended pkg/ui/output.go to handle article view output format. Added comprehensive error handling for HTTP errors (404, 403, timeout) with distinct messages. Manually tested with real articles from Slashdot feed.

## Verification

go build -o rss-cli ./cmd/rss-cli succeeded. Manual testing: ./rss-cli article view 17494 displays content in both JSON and text modes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| — | No verification commands discovered | — | — | — |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `cmd/rss-cli/article_cmd.go`
- `pkg/ui/output.go`
- `go.mod`
