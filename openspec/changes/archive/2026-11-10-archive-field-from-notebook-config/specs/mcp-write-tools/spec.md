## ADDED Requirements

### Requirement: Archive tool writes the notebook's archive field

The `archive` tool SHALL set the front-matter key named by `archive.field` in
`lindenConfig/notebook.yml` to the YAML boolean `true`. When `notebook.yml` is absent
or does not parse, the tool SHALL set `archived: true`. When `notebook.yml` is present
without `archive.field`, the tool SHALL set `archive: true`.

#### Scenario: No notebook.yml

- **WHEN** a notebook has no `lindenConfig/notebook.yml` and a document is archived
- **THEN** the document's front matter contains `archived: true`

#### Scenario: notebook.yml without archive.field

- **WHEN** `lindenConfig/notebook.yml` exists but declares no `archive.field`
- **THEN** archiving a document sets `archive: true`

#### Scenario: Explicit archive.field

- **WHEN** `lindenConfig/notebook.yml` declares `archive: { field: hidden }`
- **THEN** archiving a document sets `hidden: true`
