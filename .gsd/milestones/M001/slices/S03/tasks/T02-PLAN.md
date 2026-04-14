---
estimated_steps: 9
estimated_files: 1
skills_used: []
---

# T02: Add article open command with browser detection

Implement the article open command that opens article URLs in the default browser with proper error handling.

Steps:
1. Add articleOpenCmd to cmd/rss-cli/article_cmd.go
2. Implement detectBrowser() function that checks platform-specific commands: Linux (xdg-open), macOS (open), Windows (start)
3. Implement openBrowser(url) function using exec.Command with fallback chain
4. Wire up error handling: no browser found, launch failure
5. Add --browser flag to article view command as alternative
6. Use printer.Error() for all error output
7. Handle edge cases: empty URL, invalid URL format

## Inputs

- ``cmd/rss-cli/article_cmd.go``
- ``pkg/ui/output.go``

## Expected Output

- ``cmd/rss-cli/article_cmd.go``

## Verification

go build -o rss-cli ./cmd/rss-cli && ./rss-cli article open --help

## Observability Impact

none
