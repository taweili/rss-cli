---
estimated_steps: 10
estimated_files: 1
skills_used: []
---

# T04: Manual verification and integration testing

Perform manual verification of all error scenarios and edge cases.

Steps:
1. Test article view with invalid ID (database error)
2. Test article view with non-existent article (empty result)
3. Test article view with 404 URL (simulate with httpbin.org)
4. Test article view with 403 URL (simulate with httpbin.org)
5. Test article open command with valid article
6. Test article open command with empty URL
7. Verify all error messages are clear and actionable
8. Run full test suite: go test ./...

## Inputs

- ``rss-cli` binary`

## Expected Output

- `verification log in task summary`

## Verification

go test ./... && echo 'All tests passed'

## Observability Impact

none
