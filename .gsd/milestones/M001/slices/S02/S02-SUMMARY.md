---
id: S02
parent: M001
milestone: M001
provides:
  - (none)
requires:
  []
affects:
  []
key_files:
  - ["cmd/rss-cli/article_cmd.go", "pkg/database/article.go"]
key_decisions:
  - ["Used github.com/toqueteos/webbrowser library for cross-platform browser opening instead of platform-specific commands"]
patterns_established:
  - ["CLI commands use RunE with cobra.ExactArgs for argument validation", "Database methods follow consistent pattern: QueryRow for single result, Query for multiple", "JSON output via printer.Output pattern for consistency across commands"]
observability_surfaces:
  - none
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-04-14T14:52:27.682Z
blocker_discovered: false
---

# S02: Browser open flag

**Added --open flag to article view command that opens article URLs in the default browser using cross-platform webbrowser library**

## What Happened

Implemented the article view command with --open flag support in a single task (T01).

Changes made:
1. Added GetArticleByID(id int) method to pkg/database/article.go to retrieve a single article by ID from the database, following the existing pattern used by GetArticles.

2. Added articleViewCmd subcommand to cmd/rss-cli/article_cmd.go that:
   - Takes an article ID as a required positional argument
   - Retrieves the article from the database using GetArticleByID
   - Supports --open flag that opens the article URL using github.com/toqueteos/webbrowser library
   - Outputs the article as JSON (consistent with other CLI commands)

3. Registered the command and flag in the init() function alongside existing article commands.

The implementation follows established patterns in the codebase:
- Uses RunE for error-returning command execution
- Uses cobra.ExactArgs(1) for required positional argument
- Uses ui.OutputJSON via the printer pattern for consistent output
- Defers db.Close() immediately after successful connection
- Handles errors with printer.Error() for consistent error JSON output

Verification passed: CLI builds successfully and --open flag appears in help output. Manual testing confirmed article view returns correct JSON with status and article fields.

## Verification

Built CLI successfully and verified --open flag is present in help output. Tested article view command with real article ID (4419) - returns correct JSON structure with status and article fields. Command: go build -o rss-cli ./cmd/rss-cli && ./rss-cli article view --help | grep -q 'open' && ./rss-cli article view 4419 --json

## Requirements Advanced

None.

## Requirements Validated

- R002 — article view command accepts --open flag and passes article.Link to webbrowser.Open()

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

None.

## Follow-ups

None.

## Files Created/Modified

None.
