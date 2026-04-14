# M001: Article viewing for daily news reading

**Gathered:** 2026-04-14
**Status:** Ready for planning

## Project Description

A command-line RSS reader that completes the existing feed management tool with article viewing capability, enabling users to read full news articles in the terminal or open them in a browser.

## Why This Milestone

The existing tool can manage feeds and list article titles, but cannot actually display article content. This makes it incomplete for daily news reading. This milestone adds the missing piece: fetching and displaying full article content.

## User-Visible Outcome

### When this milestone is complete, the user can:

- Run `article view [id]` to fetch and display full article content in markdown format
- Run `article view [id] --open` to open the article in their default browser
- See clear error messages when article fetch fails (timeout, 404, paywall, parse error)

### Entry point / environment

- Entry point: CLI command `rss-cli article view [id]`
- Environment: Local terminal
- Live dependencies involved: Remote article URLs (HTTP/HTTPS)

## Completion Class

- Contract complete means: All commands work as specified, error handling tested
- Integration complete means: Article view command works with real RSS feeds and real article URLs
- Operational complete means: None (no daemon, no long-running processes)

## Final Integrated Acceptance

To call this milestone complete, we must prove:

- `article view [id]` successfully fetches, converts, and displays a real article from a known RSS feed
- `article view [id] --open` opens the article in the browser
- Error messages are shown for network failures, HTTP errors, and parse failures

## Architectural Decisions

### HTML to Markdown library

**Decision:** Use `github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown`

**Rationale:** 
- High-performance (Rust-based with Go bindings via CGO)
- Preserves article structure (headings, lists, code blocks, links)
- Well-maintained and actively developed
- Simple API: `Convert(html string) (string, error)`

**Alternatives Considered:**
- Hand-rolled HTML stripper — loses formatting, more code to maintain
- Other Go libraries — less maintained, fewer features

### Browser opening approach

**Decision:** Use `github.com/toqueteos/webbrowser` library

**Rationale:**
- Cross-platform (Windows, macOS, Linux)
- Simple API: `webbrowser.Open(url string) error`
- Handles OS-specific browser detection
- Well-established library with good adoption

**Alternatives Considered:**
- Manual `exec.Command` with OS detection — more code, error-prone
- `github.com/JasonLovesDoggo/gopen` — similar, but webbrowser has more adoption

### Content fetching strategy

**Decision:** Fetch full article content on-demand when viewing, not during feed update

**Rationale:**
- Avoids database storage bloat
- Always shows current content (no stale cache)
- Respects publisher bandwidth (only fetch what user actually reads)
- Simpler implementation (no caching logic)

**Alternatives Considered:**
- Cache during feed update — faster viewing, but storage intensive
- Hybrid (cache if available, fetch if not) — more complex, marginal benefit

### Article display format

**Decision:** Display full article at once, let user pipe to `less` if needed

**Rationale:**
- Simplest implementation
- User has full control (scroll, search, pipe)
- No custom pagination logic needed

**Alternatives Considered:**
- Built-in pagination — more code, reinvents `less`
- Auto-pipe through `less` — takes control away from user

## Error Handling Strategy

**Fetch failures:**
- Network timeout (30s) — error: "Failed to fetch article: connection timed out"
- HTTP 404/410 — error: "Article not found (404)"
- HTTP 403/paywall — error: "Access denied (403) — article may be paywalled"
- Invalid HTML — error: "Failed to parse article content"
- Empty content — error: "Article content is empty"

**Browser open failures:**
- No default browser — error: "No default browser found"
- Browser launch fails — error: "Failed to open browser: [error]"

All errors use `printer.Error()` for consistent JSON/text output formatting.

## Risks and Unknowns

- **HTML structure varies widely** — different news sites use different HTML structures; the converter must handle common patterns
- **Paywalls and JavaScript-rendered content** — some articles won't be accessible via simple HTTP GET
- **CGO dependency** — html-to-markdown uses CGO to call Rust FFI; may complicate cross-compilation

## Existing Codebase / Prior Art

- `pkg/rss/fetcher.go` — existing HTTP client with timeout and user-agent; can reuse patterns
- `cmd/rss-cli/article_cmd.go` — existing article commands; new `view` command follows same pattern
- `pkg/ui/output.go` — error and output formatting; use `printer.Error()` and `printer.Output()`

## Relevant Requirements

- R001 — Article view command (M001/S01)
- R002 — Browser open flag (M001/S02)
- R003 — HTML to Markdown conversion (M001/S01)
- R004 — Error handling (M001/S03)
- R005 — On-demand fetching (M001/S01)

## Scope

### In Scope

- `article view [id]` command implementation
- HTML to Markdown conversion
- `--open` flag for browser integration
- Error handling for fetch and conversion failures
- JSON and text output modes

### Out of Scope / Non-Goals

- OPML export (explicitly excluded)
- Feed categorization/folders
- Search across articles
- Dashboard/summary views
- Auto-refresh daemon
- Caching full content in database

## Technical Constraints

- Must work with existing database schema (no migrations)
- Must follow existing command patterns (Cobra, `RunE`, `printer.Error()`)
- Must support both JSON and text output modes

## Integration Points

- Existing `article` command group in `cmd/rss-cli/article_cmd.go`
- Existing database layer for article lookup by ID
- External article URLs (HTTP/HTTPS)

## Testing Requirements

- Unit test: HTML to Markdown conversion with sample HTML
- Integration test: `article view` with a known working RSS feed
- Manual test: Verify 3-5 real news sites render correctly

## Acceptance Criteria

**S01 (Article view command):**
- `article view [id]` fetches HTML from the article's stored link URL
- HTML is converted to Markdown with structure preserved
- Markdown is displayed to stdout
- Works with JSON and text output modes

**S02 (Browser open flag):**
- `article view [id] --open` opens the article URL in the default browser
- Command exits successfully after opening browser
- Works cross-platform (macOS, Linux, Windows)

**S03 (Error handling):**
- Network timeout shows clear error message
- HTTP 404 shows "Article not found" error
- HTTP 403 shows "Access denied — may be paywalled" error
- Parse failures show clear error message
- All errors use `printer.Error()` for consistent formatting

## Open Questions

- None — all decisions made during discussion
