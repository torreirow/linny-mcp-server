## Why

The `archive` tool writes `archived: true`, while `linny.vim` writes `archive: true`
and spec 0.4.0 (§5.4) names `archive` as the default archive field. A document
archived through the MCP server is therefore not archived for any other client of the
same notebook.

## What Changes

- The `archive` tool reads the archive field from `lindenConfig/notebook.yml`
  (`archive.field`, §5.4).
- Without `notebook.yml` the tool keeps writing `archived: true`, exactly as today.
- With `notebook.yml` but no `archive.field`, the tool writes the spec default
  `archive`.
- A `notebook.yml` that does not parse is treated as absent (best-effort, like the
  other lindenConfig readers).
- Tool description and `docs/tools.md` describe the configurable field.

## Capabilities

### Modified Capabilities
- `mcp-write-tools`: the `archive` tool writes the notebook's archive field.

## Impact

- `internal/index/meta.go`: `NotebookMeta` gains `ArchiveField`.
- `internal/mcp/write.go`: `archive` uses it.
- No behaviour change for notebooks without `notebook.yml`, so the branch can be
  offered upstream unchanged.
