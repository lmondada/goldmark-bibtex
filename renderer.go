package bibtex

import (
	"github.com/lmondada/bibtex"
	"github.com/lmondada/goldmark-bibtex/apa"
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
	reg.Register(CitationKind, r.Render)
}

func (r *CitationRenderer) Render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

	r.renderCitation(w, &entry)

	return ast.WalkContinue, nil
}

func (r *CitationRenderer) renderCitation(w util.BufWriter, entry *bibtex.Entry) {
	_, _ = w.WriteString(`<span class="citation">`)
	_, _ = w.WriteString(apa.FormatCitationKey(entry))
	_, _ = w.WriteString(apa.FormatCitation(entry))
	_, _ = w.WriteString(`</span>`)
}
