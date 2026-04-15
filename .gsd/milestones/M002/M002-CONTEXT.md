---
depends_on: [M001]
---

# M002: Article Fetch with Caching

**Gathered:** 2026-04-15
**Status:** Ready for planning

## Project Description

A new `article fetch [id]` command that fetches full article HTML from the source URL, converts it to markdown, displays it in the terminal, and caches the result in the database for subsequent views.

## Why This Milestone

M001 implemented on-demand fetching (R005) but did not cache the results. This milestone adds a dedicated fetch command with explicit caching behavior, improving performance for re-reading and reducing repeated network requests.

## User-Visible Outcome

### When this milestone is complete, the user can:

- Run `article fetch [id]` to fetch full article content and see markdown output in terminal
- Run `article fetch [id]` again to see cached content (no network request)
- See clear error messages when fetching fails (404, 403/paywall, timeout, conversion failure)

### Entry point / environment

- Entry point: `article fetch [id]` CLI command
- Environment: Local terminal
- Live dependencies involved: Remote article URLs (HTTP/HTTPS)

## Completion Class

- Contract complete means: Unit tests pass, converter works with mock HTML
- Integration complete means: `article fetch [id]` works against real RSS feeds
- Operational complete means: Caching works correctly (first fetch = network, second = cached)

## Final Integrated Acceptance

To call this milestone complete, we must prove:

- `article fetch [id]` fetches HTML from article URL, converts to markdown, displays output
- Database content field is updated after first fetch
- Subsequent `article fetch [id]` calls return cached content (verified by network monitoring or timestamps)
- Error messages are clear for 404, 403, timeout scenarios

## Architectural Decisions

### Command Interface

**Decision:** New `article fetch [id]` command, separate from `article view`

**Rationale:** Keeps mental model clean — `view` shows cached data, `fetch` retrieves fresh full content. User explicitly chose this over a `--full` flag.

**Alternatives Considered:**
- `article view [id] --full` flag — rejected for clarity
- Always fetch full content on view — rejected to preserve existing behavior

### Caching Strategy

**Decision:** Cache full article markdown in the existing `content` field after first fetch

**Rationale:** The full article is more valuable than the RSS snippet. Reusing the existing field avoids schema migration. Supersedes M001 decision D003 (no caching).

**Alternatives Considered:**
- Add new `full_content` column — rejected to avoid schema migration
- Never cache (ephemeral only) — rejected per user request

### Converter Implementation

**Decision:** Single `FetchAndConvertArticle(url string) (string, error)` function in `pkg/rss/converter.go`

**Rationale:** Matches existing pattern in `fetcher.go`. Simple, focused function that does one thing.

**Alternatives Considered:**
- Separate fetch and convert functions — rejected for simplicity
- Add conversion as method on existing fetcher — rejected for separation of concerns

## Error Handling Strategy

Reuse the existing `HTTPError` pattern from M001/S03:
- Categories: `NotFound`, `AccessDenied`, `Timeout`, `NetworkError`, `ConversionFailure`, `EmptyContent`
- `UserMessage()` returns actionable messages like "Article not found (404)" or "Failed to convert article content"
- Empty content after conversion → `ErrCategoryEmptyContent` with message "Article appears to be empty or contains no readable text"

## Risks and Unknowns

- Some RSS feeds may have relative URLs or no URL at all in the `Link` field — need to handle gracefully
- Some articles may use non-HTML formats (JavaScript-rendered, paywalls) — conversion may fail

## Existing Codebase / Prior Art

- `pkg/rss/fetcher.go` — existing `FetchArticleContent(url)` function, HTTP error handling pattern
- `pkg/rss/converter_test.go` — tests for html-to-markdown library (no converter implementation yet)
- `cmd/rss-cli/article_cmd.go` — existing `article view` and `article open` commands
- `pkg/database/article.go` — database operations, `content` field stores RSS snippets

## Relevant Requirements

- R007 — New `article fetch [id]` command (this milestone)
- R008 — Fetch and convert full article HTML (this milestone)
- R009 — Clear error messages (this milestone, extends M001/S03 pattern)
- R010 — Cache on first fetch (this milestone, supersedes R005)
- R011 — Direct markdown output (this milestone)

## Scope

### In Scope

- Create `pkg/rss/converter.go` with `FetchAndConvertArticle()` function
- Create `cmd/rss-cli/article_fetch_cmd.go` with `article fetch` command
- Update database `content` field on first fetch
- Add unit tests for converter
- Manual verification against real feeds

### Out of Scope / Non-Goals

- Modifying `article view` command behavior
- Adding new database columns
- Stripping markdown formatting for plain-text output
- Fetching articles during feed update (still on-demand only)

## Technical Constraints

- Must use existing `htmltomarkdown.Convert()` function (CGO dependency)
- Must reuse existing `HTTPError` pattern for consistency
- Database schema cannot change (reuse `content` field)

## Integration Points

- `pkg/rss/fetcher.go` — reuse HTTP client and error handling
- `pkg/database/article.go` — update article content caching
- `cmd/rss-cli/article_cmd.go` — add new subcommand

## Testing Requirements

- Unit tests for `FetchAndConvertArticle()` with mock HTML
- Integration tests with httptest server for HTTP error scenarios
- Manual testing against 2-3 real RSS feeds (e.g., Slashdot, Hacker News, tech blogs)

## Acceptance Criteria

### S01 (Converter)
- `FetchAndConvertArticle(url)` fetches HTML and returns markdown
- Returns categorized errors for 404, 403, timeout, conversion failure
- Empty HTML or unconvertible content returns `ErrCategoryEmptyContent` or `ErrCategoryConversionFailure`

### S02 (Fetch Command)
- `article fetch [id]` displays markdown in terminal
- Database `content` field is updated after successful fetch
- Subsequent fetches can return cached content (implementation decision: always fresh or check cache first)

### S03 (Tests)
- All unit tests pass
- Manual verification against real feeds succeeds
- Build succeeds (`go build -o rss-cli ./cmd/rss-cli`)

## Open Questions

- Should `article fetch` always fetch fresh, or check cache first and skip network if already cached?
  - **Decision:** Always fetch fresh for `article fetch` command (explicit fetch = fresh). The caching is for subsequent `article view` calls.
