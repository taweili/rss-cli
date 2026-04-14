# Requirements

This file is the explicit capability and coverage contract for the project.

Use it to track what is actively in scope, what has been validated by completed work, what is intentionally deferred, and what is explicitly out of scope.

## Active

### R001 — View full article content in terminal
- Class: core-capability
- Status: active
- Description: User can run `article view [id]` to fetch and display full article content in the terminal
- Why it matters: Enables actual news reading without leaving the terminal
- Source: user
- Primary owning slice: M001/S01
- Supporting slices: none
- Validation: unmapped
- Notes: Fetches from stored link URL on-demand, converts HTML to Markdown

### R002 — Open article in default browser
- Class: core-capability
- Status: active
- Description: User can run `article view [id] --open` to open the article URL in their default browser
- Why it matters: Fallback for articles that don't render well in terminal, or when user prefers browser
- Source: user
- Primary owning slice: M001/S02
- Supporting slices: none
- Validation: unmapped
- Notes: Uses cross-platform browser opening library

### R003 — HTML to Markdown conversion preserves structure
- Class: quality-attribute
- Status: active
- Description: Conversion preserves headings, paragraphs, lists, links, code blocks, and basic formatting
- Why it matters: Articles must remain readable and structured, not just plain text
- Source: user
- Primary owning slice: M001/S01
- Supporting slices: none
- Validation: unmapped
- Notes: Uses github.com/kreuzberg-dev/html-to-markdown library

### R004 — Clear error messages on fetch failure
- Class: failure-visibility
- Status: active
- Description: Network errors, HTTP errors (404, 403), and parse failures show clear, actionable error messages
- Why it matters: User needs to know why an article didn't load and whether to retry
- Source: user
- Primary owning slice: M001/S03
- Supporting slices: none
- Validation: unmapped
- Notes: Timeout, 404, 403/paywall, parse errors all have distinct messages

### R005 — Fetch content on-demand
- Class: constraint
- Status: active
- Description: Full article content is fetched from the source URL when viewing, not cached in the database
- Why it matters: Avoids storage bloat, always shows current content, respects publisher bandwidth
- Source: user
- Primary owning slice: M001/S01
- Supporting slices: none
- Validation: unmapped
- Notes: No database schema changes needed

## Out of Scope

### R006 — OPML export
- Class: anti-feature
- Status: out-of-scope
- Description: Export feeds to OPML file format
- Why it matters: Explicitly excluded to prevent scope creep
- Source: user
- Primary owning slice: none
- Supporting slices: none
- Validation: n/a
- Notes: User confirmed not needed for this milestone

## Traceability

| ID | Class | Status | Primary owner | Supporting | Proof |
|---|---|---|---|---|---|
| R001 | core-capability | active | M001/S01 | none | unmapped |
| R002 | core-capability | active | M001/S02 | none | unmapped |
| R003 | quality-attribute | active | M001/S01 | none | unmapped |
| R004 | failure-visibility | active | M001/S03 | none | unmapped |
| R005 | constraint | active | M001/S01 | none | unmapped |
| R006 | anti-feature | out-of-scope | none | none | n/a |

## Coverage Summary

- Active requirements: 5
- Mapped to slices: 5
- Validated: 0
- Unmapped active requirements: 0
