## ADDED Requirements

### Requirement: Optional notebook-level config `lindenConfig/notebook.yml`

A notebook MAY contain `lindenConfig/notebook.yml`. When the file or any key in it is
absent, indexers and clients SHALL behave as specified for version 0.3.0.

#### Scenario: No notebook.yml

- **WHEN** a notebook has no `lindenConfig/notebook.yml`
- **THEN** the index files and client listings are identical to version 0.3.0
  behaviour, including documents with `archive: true`

### Requirement: `notebook.yml` declares the taxonomy list

`notebook.yml` MAY declare `taxonomies`, a list of taxonomy names. When it is present,
a standalone indexer SHALL use it as the notebook's taxonomy list.

#### Scenario: Taxonomy list declared

- **WHEN** `notebook.yml` contains `taxonomies: [tag, project, customer]`
- **THEN** a standalone indexer treats exactly `tag`, `project` and `customer` as
  taxonomies

### Requirement: `notebook.yml` declares the archive policy

`notebook.yml` MAY declare `archive` with `field` (string, default `archive`) and
`default_hidden` (boolean, default `false`). A document is archived when its
front-matter value for `archive.field` is the boolean `true`. Indexers SHALL NOT omit
archived documents from any index file. When `default_hidden` is `true`, clients SHOULD
hide archived documents from document listings by default and SHOULD offer a way to
show them.

#### Scenario: Archived document stays in the index

- **WHEN** `notebook.yml` sets `archive: { default_hidden: true }` and `week-3.md` has
  `archive: true`
- **THEN** `week-3.md` is still present in `_index_docs_with_props.json` and in every
  L2 index it belongs to

#### Scenario: Custom archive field

- **WHEN** `notebook.yml` sets `archive: { field: archived, default_hidden: true }`
- **THEN** a document with `archived: true` is archived and a document with only
  `archive: true` is not

#### Scenario: Non-boolean value does not archive

- **WHEN** a document has `archive: "true"` (a string)
- **THEN** it is not archived
