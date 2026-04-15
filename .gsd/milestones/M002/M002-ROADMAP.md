# M002: Article Fetch with Caching

## Vision
Add a new `article fetch [id]` command that fetches full article HTML from source URL, converts to markdown, displays in terminal, and caches the result in the database for subsequent views.

## Slice Overview
| ID | Slice | Risk | Depends | Done | After this |
|----|-------|------|---------|------|------------|
| S01 | S01 | medium | — | ✅ | FetchAndConvertArticle(url) fetches HTML from URL, converts to markdown, returns error with category for failures |
| S02 | S02 | medium | — | ✅ | article fetch [id] displays markdown in terminal and caches to database |
| S03 | S03 | low | — | ✅ | All tests pass, manual verification against real feeds succeeds |
