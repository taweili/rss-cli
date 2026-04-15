---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Register article fetch command

Register articleFetchCmd subcommand in article_cmd.go init() function and add the command to the articleCmd.AddCommand() call.

## Inputs

- `cmd/rss-cli/article_cmd.go`
- `cmd/rss-cli/article_fetch_cmd.go`

## Expected Output

- `cmd/rss-cli/article_cmd.go`

## Verification

go build -o rss-cli ./cmd/rss-cli && ./rss-cli article --help | grep -q fetch
