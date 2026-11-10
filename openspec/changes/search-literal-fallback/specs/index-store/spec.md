## ADDED Requirements

### Requirement: Search falls back to a literal query

When FTS5 rejects a search query, `Search` and `SearchScoped` SHALL retry once with
each whitespace-separated token of the query as a quoted FTS5 string, joined by
spaces. A query that FTS5 accepts SHALL be run unchanged. When the literal retry also
fails, the original error SHALL be returned.

#### Scenario: Hyphenated word

- **WHEN** `Search("Connector-test", 10)` is called and a document contains "Connector-test"
- **THEN** that document is returned and no error is reported

#### Scenario: Valid FTS5 syntax is unchanged

- **WHEN** a query uses FTS5 syntax such as a prefix `back*`
- **THEN** it is run as given and matches as FTS5 defines
