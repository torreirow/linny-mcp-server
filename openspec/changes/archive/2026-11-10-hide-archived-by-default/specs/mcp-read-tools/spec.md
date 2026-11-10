## ADDED Requirements

### Requirement: Listings hide archived documents by notebook policy

When `lindenConfig/notebook.yml` sets `archive.default_hidden: true`, the `search`,
`docs_by_term`, `starred_docs` and `changed_since` tools SHALL leave out documents whose
front-matter value for `archive.field` (default `archive`) is the YAML boolean `true`,
unless the call sets `include_archived: true`. `get_doc`, `history` and `diff` SHALL NOT
be affected. Without that setting, the tools SHALL return the same documents as before.

#### Scenario: Archived document hidden by default

- **WHEN** `notebook.yml` sets `archive.default_hidden: true` and a document has `archive: true`
- **THEN** `search`, `docs_by_term`, `starred_docs` and `changed_since` do not list it

#### Scenario: include_archived lists the archive

- **WHEN** the same tools are called with `include_archived: true`
- **THEN** the archived document is listed

#### Scenario: Direct access unaffected

- **WHEN** `get_doc` is called for an archived document
- **THEN** the document is returned

#### Scenario: A string value does not archive

- **WHEN** a document has `archive: "true"` (a string)
- **THEN** it is listed by default

#### Scenario: No default_hidden

- **WHEN** `notebook.yml` is absent or does not set `archive.default_hidden: true`
- **THEN** archived documents are listed as before
