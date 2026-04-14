# S01: Article view command — UAT

**Milestone:** M001
**Written:** 2026-04-14T14:28:41.516Z

# UAT: Article View Command

## Test Cases

### 1. View article in JSON mode (default)
```bash
./rss-cli article view [id]
```
**Expected:** JSON output with id, title, link, content (markdown), and status fields.

### 2. View article in text mode
```bash
./rss-cli article view [id] --text
```
**Expected:** Human-readable output with title, link, and markdown content displayed.

### 3. Invalid article ID
```bash
./rss-cli article view invalid
```
**Expected:** Error message "Invalid article ID"

### 4. Non-existent article
```bash
./rss-cli article view 99999
```
**Expected:** Error message "Failed to retrieve article: sql: no rows in result set"

## Manual Test Results

- ✅ Test 1: Verified with article 17494 from Slashdot feed
- ✅ Test 2: Verified with article 17494 using --text flag
- ✅ Test 3: Verified with "invalid" as article ID
- ✅ Test 4: Verified with non-existent article ID 99999
