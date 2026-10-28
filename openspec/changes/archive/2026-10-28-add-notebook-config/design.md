## Context

The torrlinny notebook is read by three clients (linny.vim, linny-web, linny-mcp) with
two indexers (Hugo, `lindexer`). Archive hiding has to behave the same in all of them,
and the change must be acceptable upstream later, so it is opt-in.

## Decisions

**One file, `lindenConfig/notebook.yml`.** `lindenConfig/` is already where every
client and indexer looks; a new top-level file there needs no new path setting. The
name matches the scope: settings for the whole notebook.

**`taxonomies` is a list of names.** Each entry is both the front-matter key and the
taxonomy key in the index (`_index_taxonomies.json`, `<taxonomy>/`), i.e. the Hugo
`taxonomies:` map with key == value, which is how the reference notebook uses it.
A singular/plural mapping is not added until a notebook needs it.

**Archive is a front-matter boolean, named by `archive.field`.** Default field name is
`archive`, the name `linny.vim` already writes and §5.2 already uses for terms. Only the
YAML boolean `true` marks a document archived; strings like `"true"` do not, so a typo
cannot hide a document.

**The index stays complete; clients filter.** Indexers emit archived documents like any
other. Hiding in the indexer would make the Hugo and `lindexer` outputs diverge
(`verify --hugo`) and would make "show archived" impossible without a second index.
`_index_docs_with_props.json` already carries the field, so clients have what they
need.

**`default_hidden` defaults to `false`.** Absent file or absent key = 0.3.0 behaviour.

**Clients SHOULD offer a way to show archived documents.** How is left to the client.

## Non-goals

- Archived *terms* (L2 `archive: true`): unchanged, still presentation metadata.
- Changing any index file.
