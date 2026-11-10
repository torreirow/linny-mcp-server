## Why

`search` passes the query verbatim to FTS5 `MATCH`. Ordinary text that is not valid
FTS5 syntax fails with an SQL error instead of returning results: `Connector-test`
gives `SQL logic error: no such column: test`, and the same happens for any
hyphenated word, a stray quote or a trailing colon. A user or agent searching for a
word as it appears in a note gets an error, not hits.

## What Changes

- When the query is valid FTS5, `Search` and `SearchScoped` behave exactly as today
  (operators, phrases and prefixes keep working).
- When FTS5 rejects the query, the store retries once with a literal form of it: each
  whitespace-separated token as a quoted FTS5 string (internal `"` doubled), joined
  by spaces (implicit AND). `Connector-test` becomes `"Connector-test"`, which the
  tokenizer reads as the phrase `connector test`.
- If the literal form fails as well, the original error is returned.

## Capabilities

### Modified Capabilities
- `index-store`: full-text search falls back to a literal query on FTS5 syntax errors.

## Impact

- `internal/index/query.go`, `internal/index/query_scoped.go`.
- No behaviour change for queries that already worked, so the branch can be offered
  upstream unchanged.
