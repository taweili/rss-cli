# M002/S03: Tests and verification — Research

**Date:** 2026-04-15

## Summary

Slice S03 is a **light verification slice** — all implementation work (S01 converter, S02 fetch command) is complete with comprehensive test coverage. This slice focuses on running the full test suite and performing manual verification against real RSS feeds to prove end-to-end functionality.

The codebase already has extensive test coverage:
- **converter_test.go**: 12 successful HTML-to-markdown conversion scenarios covering headings, paragraphs, lists, links, code blocks, blockquotes, images, and nested structures
- **error_test.go**: 25+ tests for HTTP error categorization (404, 403, 500, 503, 410, timeouts, empty content, network errors, too many redirects)
- **article_test.go**: Database caching tests for `UpdateArticleContent()`
- **fetcher_test.go**: RSS feed fetching and parsing tests with mock HTTP servers

All tests currently pass (`go test ./...` succeeds with 0 failures).

## Recommendation

**Proceed directly to verification** — no additional implementation needed. Execute:

1. **Automated verification**: Run `go test ./...` and `go build -o rss-cli ./cmd/rss-cli` to confirm all tests pass and the CLI builds successfully
2. **Manual verification**: Test `article fetch [id]` against 2-3 real RSS feeds (e.g., Slashdot, Hacker News) to prove:
   - First fetch retrieves HTML, converts to markdown, and caches to database
   - Command outputs readable markdown with preserved formatting (headings, lists, links)
   - Error messages are clear for unreachable URLs or conversion failures

This is straightforward verification work — the implementation is complete and well-tested.

## Implementation Landscape

### Key Files

- `pkg/rss/converter.go` — `FetchAndConvertArticle(url)` function (implemented in S01)
- `pkg/rss/converter_test.go` — 12+ unit tests for HTML-to-markdown conversion
- `pkg/rss/fetcher.go` — HTTP client pattern with error categorization
- `pkg/rss/error_test.go` — 25+ error handling tests
- `pkg/database/article.go` — `UpdateArticleContent()` method for caching
- `pkg/database/article_test.go` — Database caching tests
- `cmd/rss-cli/article_fetch_cmd.go` — `article fetch [id]` CLI command
- `cmd/rss-cli/article_cmd.go` — Command registration (init function)

### Build Order

Already complete. S01 and S02 are done. S03 only requires:
1. Run existing tests (`go test ./...`)
2. Build CLI (`go build -o rss-cli ./cmd/rss-cli`)
3. Manual verification with real feeds

### Verification Approach

**Automated:**
```bash
# Run all tests
go test ./...

# Build the CLI
go build -o rss-cli ./cmd/rss-cli

# Verify command registration
./rss-cli article --help  # Should show 'fetch' command
```

**Manual (requires database with articles):**
```bash
# Update a feed to get articles
./rss-cli feed update-all

# List articles to get an ID
./rss-cli article list --limit 5

# Fetch an article (first time = network request + cache)
./rss-cli article fetch [id]

# Verify markdown output has proper formatting
# (headings, paragraphs, lists, links preserved)
```

## Constraints

- Must use existing `htmltomarkdown.Convert()` function (CGO dependency)
- Must reuse existing `HTTPError` pattern for error categorization
- Database schema cannot change (reuses existing `content` field)
- Tests must pass without modification (they already do)

## Common Pitfalls

- **Headless environment**: Browser-related tests may skip or fail in CI/headless environments — this is expected and handled via test skips
- **Real feed variability**: Some RSS feeds may have relative URLs, JavaScript-rendered content, or paywalls — manual testing should use known-working feeds like Slashdot
- **Network flakiness**: Integration tests with real URLs may fail due to temporary network issues — use mock HTTP servers for unit tests

## Verification Evidence Required

To mark S03 complete, capture:
1. `go test ./...` output showing all tests pass
2. `go build -o rss-cli ./cmd/rss-cli` succeeds
3. `./rss-cli article --help` shows `fetch` command
4. (Optional) Screenshot or output from manual `article fetch [id]` against a real feed
