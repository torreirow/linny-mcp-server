package index

import (
	"os"
	"path/filepath"
	"testing"
)

// fallbackStore indexes two notes: one with a hyphenated word, one without.
func fallbackStore(t *testing.T) *Store {
	t.Helper()
	root := t.TempDir()
	content := filepath.Join(root, "content")
	if err := os.MkdirAll(content, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(content, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("connector.md", "---\ntitle: Connector-test 2026-11-09\n---\nTemporary note for the Connector-test of the backup run.\n")
	write("other.md", "---\ntitle: Other\n---\nA backup schedule without connectors.\n")

	g, _, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	store, err := OpenStore(filepath.Join(t.TempDir(), "index.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Populate(g); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func hitFiles(hits []SearchHit) map[string]bool {
	files := map[string]bool{}
	for _, h := range hits {
		files[h.Filename] = true
	}
	return files
}

func TestSearchHyphenatedWordFallsBack(t *testing.T) {
	store := fallbackStore(t)
	for _, q := range []string{"Connector-test", "connector-test backup", `Connector-test"`, "note:"} {
		hits, err := store.Search(q, 10)
		if err != nil {
			t.Fatalf("Search(%q) errored: %v", q, err)
		}
		files := hitFiles(hits)
		if !files["connector.md"] || files["other.md"] {
			t.Fatalf("Search(%q) = %v, want only connector.md", q, files)
		}
	}
}

func TestSearchValidSyntaxUnchanged(t *testing.T) {
	store := fallbackStore(t)
	// A prefix query is valid FTS5; quoting it would turn `*` into literal text
	// and lose the match on "connectors".
	hits, err := store.Search("connector*", 10)
	if err != nil {
		t.Fatal(err)
	}
	files := hitFiles(hits)
	if !files["connector.md"] || !files["other.md"] {
		t.Fatalf("prefix query should match both notes, got %v", files)
	}
	// OR keeps its operator meaning.
	hits, err = store.Search("schedule OR temporary", 10)
	if err != nil {
		t.Fatal(err)
	}
	if files := hitFiles(hits); len(files) != 2 {
		t.Fatalf("OR query should match both notes, got %v", files)
	}
}

func TestLiteralQuery(t *testing.T) {
	cases := map[string]string{
		"Connector-test": `"Connector-test"`,
		"  a-b   c ":     `"a-b" "c"`,
		`say "hi"`:       `"say" """hi"""`,
		"":               "",
		"col:value":      `"col:value"`,
	}
	for in, want := range cases {
		if got := literalQuery(in); got != want {
			t.Errorf("literalQuery(%q) = %q, want %q", in, got, want)
		}
	}
}
