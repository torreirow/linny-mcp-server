## 1. Implementation

- [x] 1.1 `internal/index/meta.go`: `ArchivePolicy{Field, Hidden}` via `LoadArchivePolicy`; `ArchiveField` derives from it
- [x] 1.2 `internal/index/query_scoped.go`: `ExcludeArchived` narrows a readable subquery (`json_type(...) IS NOT 'true'`)
- [x] 1.3 `internal/mcp/tools.go`: `listScope(includeArchived)`; `include_archived` on `search`, `docs_by_term`, `starred_docs`, `changed_since`; descriptions
- [x] 1.4 Tests: policy table (incl. `default_hidden`); four tools hide by default and list with `include_archived`; string `"true"` not archived; custom field; `get_doc` unaffected; no change without `default_hidden`

## 2. Docs and verification

- [x] 2.1 `docs/tools.md` rows and archive paragraph
- [x] 2.2 CHANGELOG `[Unreleased]` entry
- [x] 2.3 `go test ./...`, golangci-lint, `nix flake check`, `openspec validate hide-archived-by-default --strict`
