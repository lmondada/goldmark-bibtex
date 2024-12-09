package bibtex

import (
	"fmt"

	"github.com/jschaf/bibtex"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// CitationRenderer is a renderer.NodeRenderer implementation that renders Citation nodes.
type CitationRenderer struct {
	bibliography map[string]bibtex.Entry
}

// NewCitationRenderer returns a new CitationRenderer.
func NewCitationRenderer(bib []bibtex.Entry) renderer.NodeRenderer {
	bibMap := make(map[string]bibtex.Entry, len(bib))
	for _, b := range bib {
		bibMap[b.Key] = b
	}

	return &CitationRenderer{
		bibliography: bibMap,
	}
}

// RegisterFuncs implements renderer.NodeRenderer interface.
func (r *CitationRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(CitationKind, r.renderCitation)
}

func (r *CitationRenderer) renderCitation(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*Citation)
	entry, ok := r.bibliography[n.Key]
	if !ok {
		// Citation not found, render as is
		_, _ = w.WriteString("[?")
		_, _ = w.WriteString(n.Key)
		_, _ = w.WriteString("]")
		return ast.WalkContinue, nil
	}

	// Format the citation based on the entry type
	switch entry.Type {
	case bibtex.EntryArticle:
		r.renderArticle(w, &entry)
	case bibtex.EntryBook:
		r.renderBook(w, &entry)
	default:
		r.renderDefault(w, &entry)
	}

	return ast.WalkContinue, nil
}

func (r *CitationRenderer) renderArticle(w util.BufWriter, entry *bibtex.Entry) {
	authors := entry.Tags["author"]
	year := entry.Tags["year"]
	_, _ = w.WriteString("[")
	_, _ = w.WriteString(fmt.Sprintf("%s, %s", authors, year))
	_, _ = w.WriteString("]")
}

func (r *CitationRenderer) renderBook(w util.BufWriter, entry *bibtex.Entry) {
	authors := entry.Tags["author"]
	year := entry.Tags["year"]
	_, _ = w.WriteString("[")
	_, _ = w.WriteString(fmt.Sprintf("%s, %s", authors, year))
	_, _ = w.WriteString("]")
}

func (r *CitationRenderer) renderDefault(w util.BufWriter, entry *bibtex.Entry) {
	_, _ = w.WriteString("[")
	_, _ = fmt.Fprintf(w, "%s", entry.Tags["title"])
	_, _ = w.WriteString("]")
}
