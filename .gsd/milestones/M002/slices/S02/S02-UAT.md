# S02: Article fetch command — UAT

**Milestone:** M002
**Written:** 2026-04-15T02:21:10.959Z

# S02: Article fetch command — UAT

**Milestone:** M002
**Written:** 2026-04-15

## UAT Type

- UAT mode: live-runtime
- Why this mode is sufficient: The slice implements a CLI command that can be directly invoked and verified. The command integrates with existing database and RSS fetching infrastructure, all of which have unit test coverage.

## Preconditions

1. RSS CLI built: `go build -o rss-cli ./cmd/rss-cli`
2. Database exists with at least one article (from feed update or manual insertion)
3. Article has a valid URL pointing to an accessible web page

## Smoke Test

Verify the fetch command is registered and accessible:
```bash
./rss-cli article --help | grep -q fetch
```

## Test Cases

### 1. Fetch command help

1. Run `./rss-cli article fetch --help`
2. **Expected:** Command displays usage information with description "Fetch full article content from source URL and convert to markdown"

### 2. Fetch article with text output

1. Run `./rss-cli article fetch [id] --text` where [id] is an existing article with a URL
2. **Expected:** Command outputs markdown-formatted article content with headings, paragraphs, and links preserved

### 3. Fetch article with JSON output

1. Run `./rss-cli article fetch [id] --json` where [id] is an existing article
2. **Expected:** Command outputs JSON with article metadata and content field containing markdown

### 4. Fetch non-existent article

1. Run `./rss-cli article fetch 99999` with non-existent ID
2. **Expected:** Error message indicating article not found

### 5. Fetch article without URL

1. Run `./rss-cli article fetch [id]` where article has NULL or empty URL
2. **Expected:** Error message indicating article has no URL

## Edge Cases

### Network failure

1. Run `./rss-cli article fetch [id]` where article URL points to non-existent domain
2. **Expected:** Clear error message about network failure or connection error

### HTTP 404 error

1. Run `./rss-cli article fetch [id]` where article URL returns 404
2. **Expected:** Error message "Article not found (404)" with appropriate error category

### Empty article content

1. Run `./rss-cli article fetch [id]` where URL returns empty HTML
2. **Expected:** Error message indicating empty content or parse failure

## Failure Signals

- Build fails: `go build -o rss-cli ./cmd/rss-cli` returns non-zero exit code
- Command not registered: `./rss-cli article --help` does not show 'fetch' command
- Tests fail: `go test ./...` returns non-zero exit code
- Panic or crash when running fetch command

## Not Proven By This UAT

- Long-term cache eviction behavior (not implemented in this slice)
- Performance under high load (single-user CLI tool)
- Database storage limits with many cached articles

## Notes for Tester

- The fetch command caches content to the database on successful fetch. Subsequent fetches of the same article will re-fetch from the source (not read from cache).
- Content is stored in the articles.content field as markdown.
- The command reuses the existing rss.FetchAndConvertArticle() function, which has comprehensive error handling for network errors, HTTP errors, and parse failures.
