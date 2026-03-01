# AGENTS.md

Guidelines for agentic coding agents working in this repository.

## Build/Lint/Test Commands

```bash
# Build the CLI
go build -o rss-cli ./cmd/rss-cli

# Run the CLI directly
go run ./cmd/rss-cli [command]

# Run all tests
go test ./...

# Run tests for a specific package
go test ./pkg/database

# Run a single test by name
go test -run TestFunctionName ./pkg/database

# Run a single test file
go test ./pkg/database -run TestSpecificTest

# Run tests with verbose output
go test -v ./...

# Check for race conditions
go test -race ./...

# Lint (requires golangci-lint)
golangci-lint run

# Format code
go fmt ./...
```

## Project Structure

```
rss-cli/
├── cmd/rss-cli/       # CLI commands (Cobra)
│   ├── main.go        # Root command and initialization
│   ├── feed_cmd.go    # Feed management commands
│   ├── article_cmd.go # Article commands
│   └── import_cmd.go  # OPML import/export
├── pkg/
│   ├── database/      # SQLite database layer
│   ├── rss/           # RSS feed fetching/parsing
│   ├── opml/          # OPML import/export
│   └── ui/            # Output formatting (JSON)
├── go.mod
└── go.sum
```

## Code Style Guidelines

### Imports

Group imports in this order with blank lines between:
1. Standard library
2. External dependencies
3. Internal packages

```go
import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"rss-cli/pkg/database"
)
```

### Naming Conventions

- **Packages**: lowercase single word (e.g., `database`, `rss`, `opml`, `ui`)
- **Types**: PascalCase for exported, camelCase for unexported
- **Functions/Methods**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase, descriptive names preferred over short names
- **Constants**: PascalCase for exported, camelCase for unexported
- **Interfaces**: typically end with `-er` (e.g., `Reader`, `Writer`)

### Error Handling

- Return errors to callers; don't panic in library code
- Wrap errors with context when appropriate: `fmt.Errorf("operation failed: %w", err)`
- In CLI commands, use `ui.OutputError()` to format error JSON output
- Always check errors from database operations, file I/O, and HTTP requests

```go
// CLI command pattern
func (cmd *cobra.Command, args []string) error {
    db, err := database.NewDB(dbPath)
    if err != nil {
        return ui.OutputError(fmt.Sprintf("Failed to connect to database: %v", err))
    }
    defer db.Close()
    // ...
}
```

### Database Operations

- Use parameterized queries to prevent SQL injection
- Use transactions for operations that modify multiple tables
- Always defer `rows.Close()` when querying
- Use `ON CONFLICT` for upserts (see `article.go`)

### JSON Output

All CLI commands output JSON using `ui.OutputJSON()` or `ui.OutputError()`:
```go
return ui.OutputJSON(map[string]interface{}{
    "status": "success",
    "feed":   feed,
})
```

### Types and Structs

- Use pointers for optional fields (e.g., `*string`, `*int` for nullable DB columns)
- Add JSON tags to all struct fields for serialization
- Use `omitempty` for optional JSON fields

```go
type Feed struct {
    ID          int     `json:"id"`
    Title       string  `json:"title"`
    LastUpdated *string `json:"last_updated,omitempty"`
}
```

### CLI Commands

- Use `RunE` instead of `Run` for commands that can error
- Use `cobra.ExactArgs(n)` for required positional arguments
- Access flags via `cmd.Flags().GetString()`, `cmd.Flags().GetBool()`, etc.
- Register subcommands in `init()` functions

### Defer Patterns

- Always defer cleanup operations (Close, Rollback) immediately after acquisition
- Use `defer tx.Rollback()` before commit; the rollback becomes a no-op after commit

### Testing

When writing tests:
- Place test files in the same package with `_test.go` suffix
- Use table-driven tests for multiple test cases
- Test both success and error paths

```go
func TestAddFeed(t *testing.T) {
    tests := []struct {
        name    string
        title   string
        url     string
        wantErr bool
    }{
        {"valid feed", "Test", "https://example.com/feed", false},
        {"empty url", "Test", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

## Key Dependencies

- **Cobra**: CLI framework (`github.com/spf13/cobra`)
- **go-sqlite3**: SQLite driver (`github.com/mattn/go-sqlite3`) - requires CGO
- **gofeed**: RSS/Atom parsing (`github.com/mmcdole/gofeed`)

## Common Patterns

### Database Connection

```go
dbPath, _ := cmd.Flags().GetString("db-path")
db, err := database.NewDB(dbPath)
if err != nil {
    return ui.OutputError(fmt.Sprintf("Failed to connect to database: %v", err))
}
defer db.Close()
```

### Feed Update Pattern

```go
parsedFeed, err := rss.FetchAndParseFeed(url)
if err != nil {
    db.IncrementErrorCount(feed.ID)
    // handle error
}
```

### Filter Pattern with Pointers

```go
type ArticleFilter struct {
    FeedID *int
    Read   *bool
    Limit  *int
}
```

## Notes

- The database path defaults to `~/.rss-cli.db` (expanded in root command's PersistentPreRun)
- All output is JSON format for easy parsing by agents/scripts
- SQLite database is created automatically if it doesn't exist