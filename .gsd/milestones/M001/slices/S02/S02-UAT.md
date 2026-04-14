# S02: Browser open flag — UAT

**Milestone:** M001
**Written:** 2026-04-14T14:52:27.683Z

# S02: Browser open flag — UAT

**Milestone:** M001
**Written:** 2026-04-14

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: This slice adds a CLI flag with straightforward behavior - either the flag is present and works, or it doesn't. The browser opening is a side effect that can't be automated in headless CI, but the flag presence and JSON output are fully testable.

## Preconditions

- CLI built: `go build -o rss-cli ./cmd/rss-cli`
- Database exists with at least one article

## Smoke Test

```bash
./rss-cli article view --help | grep -q '--open' && echo "PASS: --open flag exists"
```

## Test Cases

### 1. Help shows --open flag

1. Run `./rss-cli article view --help`
2. **Expected:** Output contains `--open   Open article URL in default browser`

### 2. View article returns JSON with article data

1. Run `./rss-cli article view 4419 --json`
2. **Expected:** JSON output contains `status: success` and `article` object with `id`, `title`, `link`, `content` fields

### 3. View non-existent article returns error

1. Run `./rss-cli article view 999999 --json`
2. **Expected:** JSON output contains error message about article not found (no panic, clean error)

### 4. View with invalid ID returns error

1. Run `./rss-cli article view abc --json`
2. **Expected:** JSON output contains error about invalid article ID

## Edge Cases

### Missing article ID argument

1. Run `./rss-cli article view`
2. **Expected:** Cobra shows error "requires at least 1 arg(s)"

## Failure Signals

- Build fails with compilation errors
- --open flag missing from help output
- article view command panics or crashes
- JSON output malformed or missing required fields

## Not Proven By This UAT

- Actual browser opening (requires GUI environment, can't test in headless CI)
- Cross-platform browser compatibility (tested on Linux only)

## Notes for Tester

- The --open flag uses github.com/toqueteos/webbrowser which handles platform-specific browser launching
- Browser opening is a side effect - the command still returns JSON output even when --open is used
- Testing actual browser opening requires manual verification on each target platform (Linux, macOS, Windows)
