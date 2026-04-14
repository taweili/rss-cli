# Project

## What This Is

A command-line RSS reader for managing feeds and reading news articles in the terminal.

## Core Value

Read full news articles from the command line without leaving the terminal, with the option to open articles in a browser when needed.

## Current State

The tool already supports:
- Feed management (add, list, remove, update)
- Article listing with filters (read/unread, by feed, by limit)
- Read status tracking
- OPML import
- SQLite backend for persistence

What's missing for daily use:
- Article viewing (fetch and display full content)
- Browser integration (open articles in default browser)

## Architecture / Key Patterns

- **Language:** Go
- **CLI Framework:** Cobra
- **Database:** SQLite (github.com/mattn/go-sqlite3)
- **RSS Parsing:** gofeed (github.com/mmcdole/gofeed)
- **Output:** JSON by default, optional plain text mode
- **Error Handling:** `ui.OutputError()` for CLI commands

## Capability Contract

See `.gsd/REQUIREMENTS.md` for the explicit capability contract, requirement status, and coverage mapping.

## Milestone Sequence

- [x] M001: Article viewing for daily news reading — Fetch and display full article content in terminal with browser open option
