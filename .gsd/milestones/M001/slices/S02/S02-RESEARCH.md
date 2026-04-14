# S02 Research: Browser Open Flag

## Summary

The `article view` command does not currently exist in the codebase. The `article_cmd.go` file only contains `list` and `mark` subcommands. To implement S02, we need to create a new `view` subcommand that retrieves an article by ID and optionally opens it in the default browser using the `--open` flag.

The `github.com/toqueteos/webbrowser` library is not in go.mod and was just added during this research (v1.2.1). This is a simple, cross-platform library that provides a single `webbrowser.Open(url string)` function which opens the given URL in the system's default browser.

## Recommendation

Create a new `articleViewCmd` in `cmd/rss-cli/article_cmd.go` with the following structure:

1. Takes a required positional argument `[id]` (article ID)
2. Has an optional `--open` boolean flag
3. Retrieves the article from the database by ID
4. Outputs the article details in JSON format
5. If `--open` is true, calls `webbrowser.Open(article.Link)` to open the article URL

This approach:
- Follows the existing command pattern used in `feed_cmd.go` and `article_cmd.go`
- Uses the database's existing `GetArticles` method with a filter, or we can add a `GetArticleByID` helper
- Leverages the `ui.OutputJSON` pattern for consistent output
- Is cross-platform (webbrowser supports Windows, macOS, Linux)

## Implementation Landscape

### Key Files to Modify

1. **cmd/rss-cli/article_cmd.go** - Add `articleViewCmd` subcommand with `--open` flag
2. **pkg/database/article.go** - Add `GetArticleByID(id int) (*Article, error)` helper method (optional, or use existing GetArticles with filter)
3. **go.mod** - Already updated with `github.com/toqueteos/webbrowser v1.2.1`

### Build Order

1. Add `GetArticleByID` to `pkg/database/article.go` (cleaner than filtering GetArticles)
2. Add `articleViewCmd` to `cmd/rss-cli/article_cmd.go`:
   - Define the command with `Args: cobra.ExactArgs(1)`
   - Add `--open` flag in `init()`
   - Implement RunE to fetch article and optionally open browser
3. Register the command in `init()` with `articleCmd.AddCommand(articleViewCmd)`

### Verification Approach

1. Build the CLI: `go build -o rss-cli ./cmd/rss-cli`
2. Add a test feed if needed: `./rss-cli feed add <url>`
3. Update the feed: `./rss-cli feed update-all`
4. Test view without flag: `./rss-cli article view 1`
5. Test view with flag: `./rss-cli article view 1 --open` (should open browser)
6. Verify JSON output contains all article fields
7. Run tests: `go test ./...`

### Flag Pattern (from existing code)

```go
var articleViewCmd = &cobra.Command{
    Use:   "view [id]",
    Short: "View an article by ID",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        openInBrowser, _ := cmd.Flags().GetBool("open")
        // ... rest of implementation
    },
}

func init() {
    articleCmd.AddCommand(articleViewCmd)
    articleViewCmd.Flags().Bool("open", false, "Open article in default browser")
}
```

### webbrowser Library Usage

```go
import "github.com/toqueteos/webbrowser"

// Usage
err := webbrowser.Open("https://example.com")
if err != nil {
    return printer.Error(fmt.Sprintf("Failed to open browser: %v", err))
}
```
