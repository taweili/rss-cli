# S02: Article Fetch Command — Research

**Date:** 2026-04-15

## Summary

S02 implements the `article fetch [id]` CLI command that fetches full article HTML from the source URL, converts it to markdown using the `FetchAndConvertArticle()` function from S01, displays the markdown output in the terminal, and caches the result in the database `content` field for subsequent views.

The command follows the existing CLI patterns established in M001: uses `RunE` for error handling, `cobra.ExactArgs(1)` for the article ID argument, supports both `--json` and `--text` output modes via the `ui.Printer`, and defers `db.Close()` immediately after connection. The implementation requires creating a new command file `cmd/rss-cli/article_fetch_cmd.go` and adding a database method `UpdateArticleContent()` to cache the fetched content.

Key design decision: `article fetch` always fetches fresh content from the network (explicit fetch = fresh), while the caching enables subsequent `article view` calls to show the full content without network requests.

## Recommendation

**Approach:** Create a dedicated `article fetch [id]` command that:
1. Retrieves the article from the database by ID to get the source URL
2. Calls `rss.FetchAndConvertArticle(url)` to fetch and convert
3. Updates the database `content` field with the markdown result
4. Displays the markdown output to the terminal

**Why:** This matches the user's explicit request for a separate command from `article view`. The separation keeps the mental model clean: `view` shows cached data, `fetch` retrieves fresh full content. Reusing the existing `content` field avoids schema migration (per decision D004).

## Implementation Landscape

### Key Files

- `cmd/rss-cli/article_fetch_cmd.go` — **NEW** CLI command file following the pattern of `article_cmd.go`. Implements the `article fetch [id]` subcommand with proper error handling, database caching, and markdown output.
- `pkg/database/article.go` — **MODIFY** Add `UpdateArticleContent(id int, content string) error` method to update the `content` field of an existing article. This is needed to cache the fetched markdown.
- `cmd/rss-cli/article_cmd.go` — **MODIFY** Register the new `articleFetchCmd` subcommand in the `init()` function.
- `pkg/rss/converter.go` — Already implemented in S01, provides `FetchAndConvertArticle(url)` function.

### Build Order

1. **First:** Add `UpdateArticleContent()` to `pkg/database/article.go` — this is a simple database operation that unblocks the command implementation
2. **Second:** Create `article_fetch_cmd.go` with the full command implementation
3. **Third:** Register the command in `article_cmd.go`
4. **Verify:** Build the CLI and run manual tests against real RSS feeds

### Verification Approach

```bash
# Build the CLI
go build -o rss-cli ./cmd/rss-cli

# Run all tests (should pass from S01)
go test ./...

# Manual verification against real feeds:
# 1. Fetch an article (first time = network request)
./rss-cli article fetch [id] --text

# 2. Verify database caching by checking the article view
./rss-cli article view [id] --text

# 3. Test error scenarios with invalid IDs or unreachable URLs
./rss-cli article fetch 99999
```

## Don't Hand-Roll

| Problem | Existing Solution | Why Use It |
|---------|------------------|------------|
| HTML to markdown conversion | `github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown` | Already used in S01, tested and working |
| HTTP fetching with error categorization | `rss.FetchAndConvertArticle()` from S01 | Reuse the existing function with proper error handling |
| CLI command structure | Existing `article_cmd.go` pattern | Consistent with existing commands, follows Cobra best practices |
| Database operations | `pkg/database/article.go` patterns | Uses existing connection, parameterized queries, error handling |
| JSON/text output | `ui.Printer` from `pkg/ui/output.go` | Consistent output formatting across all commands |

## Constraints

- **Database schema cannot change** — must reuse existing `content` field (per decision D004)
- **Must use existing HTTPError pattern** — for consistent error messages across the CLI
- **CGO dependency** — `html-to-markdown` library requires CGO-enabled build environment
- **Article must have a valid URL** — cannot fetch without a source link

## Common Pitfalls

- **Relative URLs in RSS feeds** — Some feeds may have relative URLs instead of absolute. The `FetchAndConvertArticle()` function expects absolute URLs. If the article's `Link` field is empty or relative, the fetch will fail. **Mitigation:** Check if `article.Link` is empty before fetching, return clear error message.
- **Database content update** — Forgetting to update the `content` field after fetching means the caching won't work. **Mitigation:** Make the database update part of the success path before displaying output.
- **Output mode handling** — The `--text` mode should display raw markdown, while `--json` mode should output structured JSON. **Mitigation:** Follow the existing pattern in `articleViewCmd` which handles both modes via `ui.Printer`.

## Open Risks

- **Paywalled or JavaScript-rendered articles** — Some articles may require JavaScript execution or have paywalls that prevent proper content extraction. The converter will return `ErrCategoryAccessDenied` or `ErrCategoryEmptyContent` in these cases, which is acceptable.
- **Article URLs that redirect or change** — Some RSS feed links may redirect to different URLs. The HTTP client handles up to 10 redirects, but excessive redirects will fail with `ErrCategoryTooManyRedirects`.

## Skills Discovered

No additional skills needed — this work uses existing patterns and libraries already in the codebase.

## Sources

- M002 Context document — architectural decisions and scope
- S01 Summary — `FetchAndConvertArticle()` implementation details
- Existing `article_cmd.go` — CLI command pattern reference
- `pkg/ui/output.go` — output formatting patterns
