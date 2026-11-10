## 1. Implementation

- [x] 1.1 `internal/index/query.go`: `literalQuery` and a retry-on-error wrapper used by `Search`
- [x] 1.2 `internal/index/query_scoped.go`: `SearchScoped` uses the same wrapper
- [x] 1.3 Tests: hyphenated word finds hits (plain and scoped); valid FTS5 syntax unchanged; `literalQuery` quoting

## 2. Docs and verification

- [x] 2.1 `docs/tools.md` search row
- [x] 2.2 CHANGELOG `[Unreleased]` entry
- [x] 2.3 `go test ./...`, `nix build .#default` and `openspec validate search-literal-fallback --strict`
