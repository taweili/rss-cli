# RSS CLI

A minimal command-line RSS reader with SQLite backend and OPML support.

## Features

- Manage RSS feeds (add, list, remove, update)
- Track read/unread status of articles
- Import feeds from OPML files
- Clean JSON output for easy agent integration

## Installation

```bash
go install .
```

## Usage

```bash
# Add a feed
rss feed add https://example.com/feed.xml

# List all feeds
rss feed list

# Update a feed
rss feed update 1

# List articles
rss article list --unread

# Mark an article as read
rss article mark 1 read

# Import from OPML
rss import feed.opml
```

## Commands

### Feed Management

- `rss feed add <url>` - Add a new RSS feed
- `rss feed list` - List all feeds
- `rss feed remove <id>` - Remove a feed by ID
- `rss feed update <id>` - Update a specific feed 
- `rss feed update-all` - Update all feeds

### Article Management

- `rss article list` - List all articles (use --unread, --read to filter)
- `rss article mark <id> read|unread` - Mark an article's read status

### Import / Export

- `rss import <opml-file>` - Import feeds from OPML file
- `rss export <opml-file>` - Export feeds to OPML file

## Database

SQLite database located at `~/.rss-cli.db` by default (use `-d` to override).

## Output Format

All commands output clean JSON suitable for integration with agents and scripts.

## Database Schema

- `feeds` table: id, title, url, last_updated, error_count
- `articles` table: id, feed_id, guid, title, content, link, published_at, read