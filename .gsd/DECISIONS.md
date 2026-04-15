# Decisions Register

<!-- Append-only. Never edit or remove existing rows.
     To reverse a decision, add a new row that supersedes it.
     Read this file at the start of any planning or research phase. -->

| # | When | Scope | Decision | Choice | Rationale | Revisable? | Made By |
|---|------|-------|----------|--------|-----------|------------|---------|
| D001 | M001 | library | HTML to Markdown library | github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown | High-performance (Rust-based with Go bindings), preserves article structure, well-maintained, simple API | Yes — if CGO dependency causes cross-compilation issues | collaborative |
| D002 | M001 | library | Browser opening library | github.com/toqueteos/webbrowser | Cross-platform, simple API, handles OS-specific browser detection, well-established | Yes — if webbrowser has issues on target platforms | collaborative |
| D003 | M001 | arch | Content fetching strategy | Fetch full article content on-demand when viewing, not during feed update | Avoids storage bloat, always shows current content, respects publisher bandwidth, simpler implementation | Yes — if users report slow viewing experience | collaborative |
| D004 | M002 | arch | Article content caching strategy | Cache full article markdown in database content field after first fetch via article fetch command | User explicitly requested caching on first fetch. Reusing existing content field avoids schema migration. Supersedes D003 (no caching) for the new article fetch command. | Yes | collaborative |
