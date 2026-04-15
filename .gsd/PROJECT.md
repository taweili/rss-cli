# Project

## What This Is

A command-line RSS reader for managing feeds and reading news articles in the terminal.

## Core Value

Read full news articles from the command line without leaving the terminal, with the option to open articles in a browser when needed.

## Current State

The tool supports:
- Feed management (add, list, remove, update)
- Article listing with filters (read/unread, by feed, by limit)
- Read status tracking
- OPML import
- SQLite backend for persistence
- **Article viewing** - Fetch and display full article content in terminal with HTML-to-Markdown conversion
- **Browser integration** - Open articles in default browser with `--open` flag or `article open` command
- **Comprehensive error handling** - Clear, actionable error messages for network failures, HTTP errors (404, 403, 500), timeouts, and parse failures

## Architecture / Key Patterns

- **Language:** Go
- **CLI Framework:** Cobra
- **Database:** SQLite (github.com/mattn/go-sqlite3)
- **RSS Parsing:** gofeed (github.com/mmcdole/gofeed)
- **HTML to Markdown:** github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown
- **Browser Opening:** github.com/toqueteos/webbrowser
- **Output:** JSON by default, optional plain text mode
- **Error Handling:** `ui.OutputError()` for CLI commands, error categorization with UserMessage() method

## Capability Contract

See `.gsd/REQUIREMENTS.md` for the explicit capability contract, requirement status, and coverage mapping.

All requirements for M001 validated:
- R001 — View full article content in terminal ✅
- R002 — Open article in default browser ✅
- R003 — HTML to Markdown conversion preserves structure ✅
- R004 — Clear error messages on fetch failure ✅
- R005 — Fetch content on-demand ✅

## Milestone Sequence

- [x] M001: Article viewing for daily news reading — Fetch and display full article content in terminal with browser open option
- [ ] M002: Article fetch with caching — New `article fetch [id]` command that fetches full article HTML, converts to markdown, caches in database
