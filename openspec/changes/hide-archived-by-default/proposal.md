## Why

Spec 0.4.0 (§5.4) lets a notebook declare `archive.default_hidden: true`: clients
should then hide archived documents from document listings by default and offer a
way to show them. The MCP server ignores the setting, so an agent searching or
listing a notebook keeps getting documents its owner archived, while `linny.vim` and
the web theme are meant to hide them.

## What Changes

- The server reads `archive.default_hidden` from `lindenConfig/notebook.yml`, next to
  `archive.field`, through a new `index.LoadArchivePolicy`.
- With `default_hidden: true`, `search`, `docs_by_term`, `starred_docs` and
  `changed_since` leave out documents whose archive field is the YAML boolean `true`.
  The filter is a SQL predicate on the caller's readable subquery, like scope.
- Each of those tools takes `include_archived` (bool) to list archived documents too.
- `get_doc`, `history` and `diff` are unaffected: opening a document directly is not
  a listing (§5.4).
- Without `notebook.yml`, or with `default_hidden` absent or `false`, nothing changes.

## Capabilities

### Modified Capabilities
- `mcp-read-tools`: listing tools honour `archive.default_hidden` and accept
  `include_archived`.

## Impact

- `internal/index/meta.go` (`ArchivePolicy`), `internal/index/query_scoped.go`
  (`ExcludeArchived`), `internal/mcp/tools.go`.
- `docs/tools.md`, CHANGELOG.
- Search also hides archived documents. §5.4 names term, starred and recent listings;
  for an agent, search is its main listing, and `include_archived` keeps the archive
  reachable.
- Stacked on `feat/archive-field` (needs the `notebook.yml` reader).
