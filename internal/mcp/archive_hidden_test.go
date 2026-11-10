package mcp

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// archiveFixture adds three documents that share a search word, a tag and the
// starred flag: one archived, one visible, and one whose archive value is the
// string "true" (which spec §5.4 says does not archive). notebook, when not
// empty, is written as lindenConfig/notebook.yml. Everything is committed so
// changed_since sees it.
func archiveFixture(t *testing.T, notebook string) *writeFixture {
	t.Helper()
	f := newWriteFixture(t, "read:*", "write:*")
	if notebook != "" {
		if err := os.WriteFile(filepath.Join(f.root, "lindenConfig", "notebook.yml"), []byte(notebook), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	f.putDoc(t, "old_plan.md", "---\ntitle: Old Plan\nstarred: true\narchive: true\ntags:\n  - zebrafish\n---\n\nzebrafish plan\n")
	f.putDoc(t, "new_plan.md", "---\ntitle: New Plan\nstarred: true\ntags:\n  - zebrafish\n---\n\nzebrafish plan\n")
	f.putDoc(t, "quoted_plan.md", "---\ntitle: Quoted Plan\nstarred: true\narchive: \"true\"\ntags:\n  - zebrafish\n---\n\nzebrafish plan\n")
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-qm", "archive fixture"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = f.root
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=T", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=T", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return f
}

// listed runs the four listing tools and returns, per tool, the set of the
// fixture documents each one returned.
func listed(t *testing.T, f *writeFixture, includeArchived bool) map[string]map[string]bool {
	t.Helper()
	ctx := context.Background()
	rd := f.w.reader()
	fixture := map[string]bool{"old_plan.md": true, "new_plan.md": true, "quoted_plan.md": true}
	res := map[string]map[string]bool{}
	add := func(tool, slug string) {
		if !fixture[slug] {
			return
		}
		if res[tool] == nil {
			res[tool] = map[string]bool{}
		}
		res[tool][slug] = true
	}

	_, s, err := rd.search(ctx, nil, searchIn{Query: "zebrafish", IncludeArchived: includeArchived})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	for _, h := range s.Hits {
		add("search", h.Filename)
	}
	_, d, err := rd.docsByTerm(ctx, nil, docsByTermIn{Taxonomy: "tags", Term: "zebrafish", IncludeArchived: includeArchived})
	if err != nil {
		t.Fatalf("docs_by_term: %v", err)
	}
	for _, slug := range d.Docs {
		add("docs_by_term", slug)
	}
	_, st, err := rd.starredDocs(ctx, nil, starredDocsIn{IncludeArchived: includeArchived})
	if err != nil {
		t.Fatalf("starred_docs: %v", err)
	}
	for _, doc := range st.Docs {
		add("starred_docs", doc.Filename)
	}
	_, c, err := rd.changedSince(ctx, nil, changedSinceIn{Since: "2000-01-01", IncludeArchived: includeArchived})
	if err != nil {
		t.Fatalf("changed_since: %v", err)
	}
	for _, slug := range c.Docs {
		add("changed_since", slug)
	}
	return res
}

var listingTools = []string{"search", "docs_by_term", "starred_docs", "changed_since"}

func TestListingsHideArchivedWhenNotebookSaysSo(t *testing.T) {
	f := archiveFixture(t, "archive:\n  field: archive\n  default_hidden: true\n")

	hidden := listed(t, f, false)
	for _, tool := range listingTools {
		got := hidden[tool]
		if got["old_plan.md"] {
			t.Errorf("%s lists the archived document by default: %v", tool, got)
		}
		if !got["new_plan.md"] || !got["quoted_plan.md"] {
			t.Errorf("%s hides a document that is not archived (a string \"true\" does not archive): %v", tool, got)
		}
	}

	all := listed(t, f, true)
	for _, tool := range listingTools {
		if got := all[tool]; len(got) != 3 {
			t.Errorf("%s with include_archived = %v, want all three documents", tool, got)
		}
	}

	// Opening an archived document directly is unaffected.
	_, doc, err := f.w.reader().getDoc(context.Background(), nil, getDocIn{Slug: "old_plan.md"})
	if err != nil {
		t.Fatal(err)
	}
	if !doc.Found {
		t.Fatal("get_doc must still return an archived document")
	}
}

func TestListingsHideArchivedUnderCustomField(t *testing.T) {
	f := archiveFixture(t, "archive:\n  field: archived\n  default_hidden: true\n")
	// The fixture marks old_plan.md with `archive`, not the configured `archived`,
	// so nothing is archived under this policy.
	for _, tool := range listingTools {
		if got := listed(t, f, false)[tool]; len(got) != 3 {
			t.Errorf("%s = %v; only the configured field archives", tool, got)
		}
	}
}

func TestListingsUnchangedWithoutDefaultHidden(t *testing.T) {
	for name, notebook := range map[string]string{
		"no notebook.yml":          "",
		"default_hidden not set":   "archive:\n  field: archive\n",
		"default_hidden set false": "archive:\n  default_hidden: false\n",
	} {
		t.Run(name, func(t *testing.T) {
			f := archiveFixture(t, notebook)
			got := listed(t, f, false)
			for _, tool := range listingTools {
				if len(got[tool]) != 3 {
					t.Errorf("%s = %v, want all three documents", tool, got[tool])
				}
			}
		})
	}
}
