# S03: Tests and verification — UAT

**Milestone:** M002
**Written:** 2026-04-15T02:29:29.803Z

# S03: Tests and verification — UAT

**Milestone:** M002
**Written:** 2026-04-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: This slice validates existing functionality through automated tests and manual verification against real feeds. No new runtime features were added - only verification of the article fetch command implemented in S02.

## Preconditions

- CLI binary built: `go build -o rss-cli ./cmd/rss-cli`
- Database initialized with feeds and articles
- Network access available for fetching article content

## Smoke Test

Run `./rss-cli article fetch [id]` with a valid article ID - should return JSON with status "success" and markdown content.

## Test Cases

### 1. Automated test suite passes

1. Run `go test ./...`
2. **Expected:** All packages pass (cmd/rss-cli, pkg/database, pkg/rss)
3. Run `go build -o rss-cli ./cmd/rss-cli`
4. **Expected:** Binary builds without errors

### 2. Feed update populates articles

1. Run `./rss-cli feed update-all`
2. **Expected:** Feeds update successfully (some may fail with expected errors like timeouts)
3. Run `./rss-cli article list --limit 5`
4. **Expected:** Returns JSON array with at least 5 articles containing id, title, link fields

### 3. Article fetch converts and caches content

1. Run `./rss-cli article fetch [id]` with valid article ID
2. **Expected:** Returns JSON with status "success", id, title, link, and content fields
3. Content field contains markdown-formatted text with headings, paragraphs, links
4. Run same command again
5. **Expected:** Returns cached content (faster response, same content)

### 4. Error handling for invalid articles

1. Run `./rss-cli article fetch [id]` with non-existent ID
2. **Expected:** Returns error JSON with clear message "Article not found"

## Edge Cases

### Network failure during fetch

1. Test with article from site that times out or returns errors
2. **Expected:** Returns error JSON with category-specific message (Timeout, NotFound, AccessDenied)

### Empty article content

1. Fetch article that converts to empty/no meaningful content
2. **Expected:** Returns error with "EmptyContent" category and user message

## Failure Signals

- Test failures in `go test ./...` output
- Build errors from `go build`
- Article fetch returns error for valid articles
- Content field missing or empty in fetch output
- Database caching fails (content not persisted)

## Not Proven By This UAT

- Long-term cache eviction behavior (not implemented)
- Performance under heavy concurrent fetch load
- Articles with complex JavaScript-rendered content

## Notes for Tester

- Some feeds may consistently fail (McKinsey had network issues during testing) - this is site-specific
- The Register feeds work reliably for testing
- First fetch takes longer (network request), subsequent fetches use cache
