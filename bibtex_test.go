package bibtex

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jschaf/bibtex"
	"github.com/jschaf/bibtex/ast"
	"github.com/yuin/goldmark"
)

const testBibContent = `@InProceedings{Albert1989,
  author = {Albert, Luc},
  title = {Average Case Complexity Analysis of {RETE} Pattern-Match Algorithm and Average Size of Join in Database},
  publisher = {Springer},
  year = {1989},
  pages = {223--241},
  booktitle = {Foundations of Software Technology and Theoretical Computer Science, Ninth Conference, Bangalore, India, December 19-21, 1989, Proceedings},
}`

func createTempBibFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test*.bib")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpfile.Name()) })

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpfile.Name()
}

func verifyBibliography(t *testing.T, bibExtender *Extender) {
	t.Helper()
	// Find entry
	var entry *bibtex.Entry
	key := "Albert1989"
	for _, e := range bibExtender.Bibliography {
		if e.Key == key {
			if entry != nil {
				t.Fatalf("Expected to find unique entry with key '%s'", key)
			}
			entry = &e
		}
	}
	if entry == nil {
		t.Fatalf("Expected to find entry with key '%s'", key)
	}

	// Verify the entry type
	if entry.Type != "inproceedings" {
		t.Errorf("Expected entry type inproceedings, got %s", entry.Type)
	}

	// Test the parsed fields
	tests := []struct {
		field    string
		expected string
	}{
		{"author", "Luc Albert"},
		{"title", "Average Case Complexity Analysis of RETE Pattern-Match Algorithm and Average Size of Join in Database"},
		{"booktitle", "Foundations of Software Technology and Theoretical Computer Science, Ninth Conference, Bangalore, India, December 19-21, 1989, Proceedings"},
		{"publisher", "Springer"},
		{"year", "1989"},
		{"pages", "223--241"},
	}

	for _, tt := range tests {
		got := entry.Tags[tt.field]
		switch got.Kind() {
		case ast.KindAuthors:
			got := got.(ast.Authors)
			if len(got) != 1 {
				t.Errorf("Expected 1 author, got %d", len(got))
			}
			a := got[0]
			firstLastName := a.First.(*ast.Text).Value + " " + a.Last.(*ast.Text).Value
			if firstLastName != tt.expected {
				t.Errorf("Expected name to be %s, got %q", tt.expected, firstLastName)
			}
		case ast.KindText:
			got := got.(*ast.Text)
			if got.Value != tt.expected {
				t.Errorf("Field %s = %q; want %s", tt.field, got.Value, tt.expected)
			}
		default:
			t.Errorf("Unexpected type for field %s: %T", tt.field, got)
		}
	}
}

func verifyMarkdownConversion(t *testing.T, bibExtender *Extender) {
	t.Helper()
	markdown := goldmark.New(
		goldmark.WithExtensions(bibExtender),
	)

	source := []byte("As shown in @Albert1989, the results are significant.")
	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		t.Fatal(err)
	}

	citationExp := `<span class="citation-key">Albert, 1989</span><span class="citation-full"><span class="author">Albert Luc</span> (1989). <span class="title">Average Case Complexity Analysis of RETE Pattern-Match Algorithm and Average Size of Join in Database</span>. In <span class="booktitle">Foundations of Software Technology and Theoretical Computer Science, Ninth Conference, Bangalore, India, December 19-21, 1989, Proceedings</span> (pp. <span class="pages">223--241</span>).</span>`
	expected := fmt.Sprintf("<p>As shown in %s, the results are significant.</p>\n", citationExp)
	if got := buf.String(); got != expected {
		t.Errorf("Markdown conversion = %s; want %s", got, expected)
	}
}

func TestBibTeXParsingFromString(t *testing.T) {
	bibFile := createTempBibFile(t, testBibContent)
	bibExtender, err := New(bibFile)
	if err != nil {
		t.Fatal(err)
	}

	verifyBibliography(t, bibExtender)
	verifyMarkdownConversion(t, bibExtender)
}

func TestBibTeXParsingFromFile(t *testing.T) {
	// Create a test file in the testdata directory
	testdataDir := "testdata"
	bibFile := filepath.Join(testdataDir, "refs.bib")

	bibExtender, err := New(bibFile)
	if err != nil {
		t.Fatal(err)
	}

	verifyBibliography(t, bibExtender)
	verifyMarkdownConversion(t, bibExtender)
}
