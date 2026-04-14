# S01 — Research

**Date:** 2026-04-14

## Summary

Slice S01 implements the `article view [id]` command that fetches full article content from the stored link URL, converts HTML to Markdown, and displays it in the terminal. The slice owns requirements R001 (view article content), R003 (HTML to Markdown conversion), and R005 (on-demand fetching).

The implementation requires:
1. Adding a new `article view [id]` subcommand following existing Cobra patterns
2. Fetching HTML from the article's stored `Link` field using HTTP client patterns from `pkg/rss/fetcher.go`
3. Converting HTML to Markdown using `github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown`
4. Displaying the result using the existing `ui.Printer` for both JSON and text output modes

The existing codebase provides all necessary patterns: database access via `pkg/database/article.go`, HTTP fetching with timeout/user-agent in `pkg/rss/fetcher.go`, and output formatting in `pkg/ui/output.go`.

## Recommendation

Implement the article view command as a new subcommand `articleViewCmd` in `cmd/rss-cli/article_cmd.go`. The command should:

1. Accept a single positional argument (article ID)
2. Look up the article by ID in the database to get the `Link` URL
3. Fetch HTML from the URL with 30s timeout, user-agent header, and redirect limit (reuse patterns from `fetcher.go`)
4. Convert HTML to Markdown using `htmltomarkdown.Convert()`
5. Output the result using `printer.Output()` with the article metadata and markdown content

The implementation is straightforward—no novel architecture needed. Follow existing command patterns exactly: use `RunE`, check errors, call `printer.Error()` on failures, and defer `db.Close()`.

## Implementation Landscape

### Key Files

- `cmd/rss-cli/article_cmd.go` — Add `articleViewCmd` following the same pattern as `articleListCmd` and `articleMarkCmd`
- `pkg/rss/fetcher.go` — Reuse HTTP client pattern (timeout, user-agent, redirect limit) for fetching article HTML; may extract a new `FetchArticleContent(url string) (string, error)` function
- `pkg/ui/output.go` — May need to extend `outputArticle()` to handle displaying fetched markdown content, or create new output method for view command
- `pkg/database/article.go` — May need `GetArticleByID(id int) (*Article, error)` if not already present

### Build Order

1. **First**: Add `GetArticleByID()` to `pkg/database/article.go` — this unblocks everything else
2. **Second**: Add HTML fetching function to `pkg/rss/fetcher.go` — can be tested independently
3. **Third**: Add `article view` command to `cmd/rss-cli/article_cmd.go` — wires everything together
4. **Fourth**: Add tests for HTML→Markdown conversion and command integration

### Verification Approach

```bash
# Build the CLI
go build -o rss-cli ./cmd/rss-cli

# Run tests
go test ./...

# Test with real data (manual)
./rss-cli feed add https://example.com/feed.xml
./rss-cli feed update-all
./rss-cli article list --unread
./rss-cli article view [id]
./rss-cli article view [id] --json
```

## Don't Hand-Roll

| Problem | Existing Solution | Why Use It |
|---------|------------------|------------|
| HTML to Markdown conversion | `github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown` | High-performance (Rust-based via CGO), preserves structure (headings, lists, code blocks, links), simple API: `Convert(html string) (string, error)` |
| HTTP fetching with timeout | Existing pattern in `pkg/rss/fetcher.go` | Already has 30s timeout, user-agent header, redirect limit—reuse exactly |
| CLI command structure | Existing Cobra commands in `article_cmd.go` | Follows established patterns for flags, error handling, output formatting |

## Constraints

- **CGO dependency**: html-to-markdown uses CGO to call Rust FFI; requires CGO-enabled Go build (already satisfied since go-sqlite3 also requires CGO)
- **Database schema**: Cannot modify schema for this milestone; must work with existing `articles` table structure (already has `link` field)
- **Output format**: Must support both JSON and text modes via existing `ui.Printer`

## Common Pitfalls

- **Empty content handling**: Some articles may have empty body after conversion; must check and return clear error
- **HTTP error codes**: Distinguish 404 (not found), 403 (paywall), timeout, and parse errors with distinct messages per error handling strategy in M001-CONTEXT.md
- **HTML encoding**: Ensure fetched HTML is properly decoded (UTF-8) before conversion

## Open Risks

- **HTML structure varies**: Different news sites use different HTML structures; the converter handles common patterns but edge cases may produce poor output
- **JavaScript-rendered content**: Articles that require JS to render won't be accessible via simple HTTP GET (out of scope for this milestone)

## Skills Discovered

| Technology | Skill | Status |
|------------|-------|--------|
| Go | jeffallan/claude-skills@golang-pro | available (9.4K installs) |
| Go | affaan-m/everything-claude-code@golang-patterns | available (5.1K installs) |
| Go | affaan-m/everything-claude-code@golang-testing | available (4K installs) |

## Sources

- html-to-markdown Go API documentation (https://github.com/kreuzberg-dev/html-to-markdown/blob/main/docs/reference/api-go.md)
- webbrowser package (https://pkg.go.dev/github.com/toqueteos/webbrowser)
- Existing codebase patterns in `cmd/rss-cli/`, `pkg/rss/`, `pkg/database/`, `pkg/ui/`
