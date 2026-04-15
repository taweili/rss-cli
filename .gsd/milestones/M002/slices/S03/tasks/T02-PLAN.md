---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Manual verification against real RSS feeds

Manually test the article fetch command against 2-3 real RSS feeds to prove end-to-end functionality: fetching HTML, converting to markdown, caching to database, and displaying formatted output.

## Inputs

- `rss-cli (built binary)`
- `Database with articles from feed update-all`

## Expected Output

- `Manual verification evidence captured in task summary with sample output`

## Verification

./rss-cli feed update-all && ./rss-cli article list --limit 5 && ./rss-cli article fetch [id]
