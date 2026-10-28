## Why

Two notebook-wide settings have no home in the Linden Index Specification.

1. **The taxonomy list (§13-Q6).** Hugo reads it from its own site `config.yaml`. A
   standalone indexer such as `lindexer` has no notebook-level place to read it from;
   the old `L0-CONF-ROOT.yml` was removed in spec 0.2.0.
2. **Archived documents.** `linny.vim` writes `archive: true` into a document's front
   matter, and §5.2 already lists `archive` as an L2 term-config key, but nothing says
   what the flag means for a *document*. Today no indexer and no client acts on it, so
   an archived note keeps showing up everywhere. Clients also disagree on the field
   name: `linny.vim` writes `archive`, the `linny-mcp` `archive` tool writes `archived`.

Both are notebook policy, not per-user or per-client preference, so they belong in the
notebook's own `lindenConfig/`.

## What Changes

- New optional file `lindenConfig/notebook.yml`, specified in a new §5.4 of
  `docs/linden-index-spec.md`, with two keys:
  - `taxonomies`: the list of taxonomy names. A standalone indexer SHOULD use it when
    present (SHALL when present). Resolves Q6.
  - `archive`: `{ field, default_hidden }`. Defines which front-matter field marks a
    document as archived, and whether clients hide archived documents by default.
- The index stays complete: an indexer MUST NOT drop archived documents. Hiding is
  client behaviour.
- Without `notebook.yml` everything behaves exactly as in 0.3.0. This keeps the change
  offerable upstream as-is.
- Spec version 0.3.0 → 0.4.0 (draft): §4.3 front matter mentions the archive field,
  §5 points at §5.4, Q6 marked RESOLVED, version history row added.
- `index-format` capability gains one requirement describing `notebook.yml`.

## Capabilities

### Modified Capabilities
- `index-format`: adds the optional `lindenConfig/notebook.yml` contract.

## Impact

- Docs/spec only. No code in this repository changes; `lindexer` reading
  `taxonomies` and the `archive` tool honouring `archive.field` are follow-ups.
- Clients (`linny.vim`, the linny-web theme) can implement archive hiding against a
  written contract.
