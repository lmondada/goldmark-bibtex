package bibtex

import (
	"fmt"

	"github.com/lmondada/bibtex"
	bibtexAst "github.com/lmondada/bibtex/ast"

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
		// Citation not found, render as question mark
		_, _ = w.WriteString("[?]")
		return ast.WalkContinue, nil
	}

	// Format the citation based on the entry type
	switch entry.Type {
	case bibtex.EntryArticle, bibtex.EntryInProceedings, bibtex.EntryBook:
		r.renderAuthorYear(w, &entry)
	default:
		r.renderDefault(w, &entry)
	}

	return ast.WalkContinue, nil
}

func fmtAuthorYear(authorsExpr, yearExpr bibtexAst.Expr) string {
	authors := authorsExpr.(bibtexAst.Authors)
	firstAuthor := authors[0]
	lastName := firstAuthor.Last.(*bibtexAst.Text).Value
	year := yearExpr.(*bibtexAst.Text).Value
	return fmt.Sprintf("[%s, %s]", lastName, year)
}

func (r *CitationRenderer) renderAuthorYear(w util.BufWriter, entry *bibtex.Entry) {
	authors := entry.Tags["author"]
	year := entry.Tags["year"]
	_, _ = w.WriteString(fmtAuthorYear(authors, year))
}

func (r *CitationRenderer) renderBook(w util.BufWriter, entry *bibtex.Entry) {
	authors := entry.Tags["author"]
	year := entry.Tags["year"]
	_, _ = w.WriteString(fmtAuthorYear(authors, year))
}

func (r *CitationRenderer) renderDefault(w util.BufWriter, entry *bibtex.Entry) {
	_, _ = w.WriteString("[")
	_, _ = w.WriteString(entry.Tags["title"].(*bibtexAst.Text).Value)
	_, _ = w.WriteString("]")
}
