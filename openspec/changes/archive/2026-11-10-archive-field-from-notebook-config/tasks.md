## 1. Implementation

- [x] 1.1 `internal/index/meta.go`: read `lindenConfig/notebook.yml`, expose `ArchiveField`
- [x] 1.2 `internal/mcp/write.go`: `archive` writes `ArchiveField`; tool description
- [x] 1.3 Tests: no file → `archived`; file without field → `archive`; explicit field; invalid YAML → `archived`

## 2. Docs and verification

- [x] 2.1 `docs/tools.md` archive row
- [x] 2.2 CHANGELOG `[Unreleased]` entry
- [x] 2.3 `go test ./...` and `openspec validate archive-field-from-notebook-config --strict`
