---
name: news-summarizer
description: Fetches unread RSS articles from feeds and generates summaries. Use when you need to read and summarize news from RSS feeds stored in the rss-cli SQLite database.
---

# News Summarizer Skill

This skill helps you fetch unread news articles from RSS feeds and summarize them.

## Setup

No additional setup required. The skill uses the `rss-cli` tool that must be installed in the project.

## Usage

### Get Unread Articles

List all unread articles across all feeds:

```bash
./rss-cli article list --unread --json
```

### Get Unread Articles from Specific Feed

First, list feeds to get the feed ID:

```bash
./rss-cli feed list --json
```

Then get unread articles for that feed:

```bash
./rss-cli article list --unread --feed <feed_id> --json
```

### Update Feeds and Summarize

Update all feeds to get latest articles:

```bash
./rss-cli feed update-all
```

Then list unread articles:

```bash
./rss-cli article list --unread --json
```

### Mark Articles as Read

After reading/summarizing, mark articles as read:

```bash
./rss-cli article mark <article_id> read
```

## Output Format

All commands support `--json` flag for machine-readable output. Example output:

```json
{
  "articles": [
    {
      "id": 1,
      "feed_id": 1,
      "guid": "https://example.com/article-1",
      "title": "Sample Article",
      "content": "Article content here...",
      "url": "https://example.com/article-1",
      "published_at": "2024-01-01T12:00:00Z",
      "read": false
    }
  ],
  "count": 1
}
```
