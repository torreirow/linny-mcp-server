package index

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

// SearchHit is one ranked full-text search result.
type SearchHit struct {
	Filename string  `json:"filename"`
	Title    string  `json:"title"`
	Snippet  string  `json:"snippet"`
	Score    float64 `json:"score"` // FTS5 bm25; lower is more relevant
}

// DocRef names a document without its content: enough to list documents for a
// reader without a second lookup per row.
type DocRef struct {
	Filename string `json:"filename"`
	Title    string `json:"title"`
}

// Doc is a single document as returned by GetDoc.
type Doc struct {
	Filename string         `json:"filename"`
	Title    string         `json:"title"`
	Props    map[string]any `json:"props"`
	Body     string         `json:"body"`
}

const defaultSearchLimit = 20

// Search runs a full-text query over titles and bodies, returning results ranked
// by FTS5 bm25 (best first) with a snippet from the body. A query that matches
// nothing returns an empty slice and no error. A query FTS5 rejects is retried
// once as literal text (see searchWithFallback).
func (s *Store) Search(query string, limit int) ([]SearchHit, error) {
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	return searchWithFallback(query, func(q string) ([]SearchHit, error) {
		return s.scanHits(
			`SELECT filename, title, snippet(docs_fts, 2, '[', ']', '…', 8), bm25(docs_fts)
			 FROM docs_fts
			 WHERE docs_fts MATCH ?
			 ORDER BY bm25(docs_fts)
			 LIMIT ?`, q, limit)
	})
}

// searchWithFallback runs query; when FTS5 rejects it (e.g. "Connector-test"
// reads as a column filter) it retries once with literalQuery. A query that
// already works is never rewritten, and if the retry fails too the original
// error is returned.
func searchWithFallback(query string, run func(string) ([]SearchHit, error)) ([]SearchHit, error) {
	hits, err := run(query)
	if err == nil {
		return hits, nil
	}
	lit := literalQuery(query)
	if lit == "" || lit == query {
		return nil, err
	}
	if hits, litErr := run(lit); litErr == nil {
		return hits, nil
	}
	return nil, err
}

// literalQuery turns free text into an FTS5 query without operators: every
// whitespace-separated token becomes a quoted string (internal quotes doubled),
// joined by spaces, i.e. an implicit AND of the tokens as phrases.
func literalQuery(query string) string {
	tokens := strings.Fields(query)
	for i, t := range tokens {
		tokens[i] = `"` + strings.ReplaceAll(t, `"`, `""`) + `"`
	}
	return strings.Join(tokens, " ")
}

// scanHits runs a search statement whose columns are filename, title, snippet
// and bm25 score.
func (s *Store) scanHits(stmt string, args ...any) ([]SearchHit, error) {
	rows, err := s.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var hits []SearchHit
	for rows.Next() {
		var h SearchHit
		if err := rows.Scan(&h.Filename, &h.Title, &h.Snippet, &h.Score); err != nil {
			return nil, err
		}
		hits = append(hits, h)
	}
	return hits, rows.Err()
}

// ListTaxonomies returns the taxonomy names that occur in the index, sorted.
func (s *Store) ListTaxonomies() ([]string, error) {
	return s.queryStrings(`SELECT name FROM taxonomies ORDER BY name`)
}

// AllDocFilenames returns every document filename in the store, sorted.
func (s *Store) AllDocFilenames() ([]string, error) {
	return s.queryStrings(`SELECT filename FROM docs ORDER BY filename`)
}

// TermsForTaxonomy returns the terms of a taxonomy, sorted.
func (s *Store) TermsForTaxonomy(taxonomy string) ([]string, error) {
	return s.queryStrings(`SELECT term FROM terms WHERE taxonomy = ? ORDER BY term`, taxonomy)
}

// DocsByTerm returns the filenames tagged with (taxonomy, term), sorted.
func (s *Store) DocsByTerm(taxonomy, term string) ([]string, error) {
	return s.queryStrings(
		`SELECT filename FROM membership WHERE taxonomy = ? AND term = ? ORDER BY filename`,
		taxonomy, term)
}

// GetDoc returns a document by filename. ok is false when it is not indexed.
func (s *Store) GetDoc(filename string) (doc Doc, ok bool, err error) {
	var title, propsJSON string
	err = s.db.QueryRow(`SELECT title, props_json FROM docs WHERE filename = ?`, filename).
		Scan(&title, &propsJSON)
	if errors.Is(err, sql.ErrNoRows) {
		return Doc{}, false, nil
	}
	if err != nil {
		return Doc{}, false, err
	}

	var body string
	// Body lives in the FTS table; ignore a missing row defensively.
	if err := s.db.QueryRow(`SELECT body FROM docs_fts WHERE filename = ?`, filename).Scan(&body); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return Doc{}, false, err
	}

	props := map[string]any{}
	if propsJSON != "" {
		if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
			return Doc{}, false, err
		}
	}
	return Doc{Filename: filename, Title: title, Props: props, Body: body}, true, nil
}

// TermMembership is a (taxonomy, term) a document belongs to.
type TermMembership struct {
	Taxonomy string `json:"taxonomy"`
	Term     string `json:"term"`
}

// TermsOfDoc returns the (taxonomy, term) memberships of a document, sorted.
func (s *Store) TermsOfDoc(filename string) ([]TermMembership, error) {
	rows, err := s.db.Query(
		`SELECT taxonomy, term FROM membership WHERE filename = ? ORDER BY taxonomy, term`, filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []TermMembership{}
	for rows.Next() {
		var m TermMembership
		if err := rows.Scan(&m.Taxonomy, &m.Term); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) queryStrings(query string, args ...any) ([]string, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := []string{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
