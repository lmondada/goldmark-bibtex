package bibtex

import (
	"os"

	"github.com/lmondada/bibtex"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extender is a goldmark extension for rendering BibTeX citations.
type Extender struct {
	Bibliography []bibtex.Entry
}

// New creates a new BibTeX extender with the given bibliography file.
func New(bibFile string) (*Extender, error) {
	f, err := os.Open(bibFile)
	if err != nil {
		return nil, err
	}
	bib := bibtex.New(
		bibtex.WithResolvers(
			// NewAuthorResolver creates a resolver for the "author" field that parses
			// author names into an ast.Authors node.
			bibtex.NewAuthorResolver("author"),
			// SimplifyEscapedTextResolver replaces ast.TextEscaped nodes with a plain
			// ast.Text containing the value that was escaped. Meaning, `\&` is converted to
			// `&`.
			bibtex.ResolverFunc(bibtex.SimplifyEscapedTextResolver),
			// RenderParsedTextResolver replaces ast.ParsedText with a simplified rendering
			// of ast.Text.
			bibtex.NewRenderParsedTextResolver(),
		),
	)
	file, err := bib.Parse(f)
	if err != nil {
		return nil, err
	}
	entries, err := bib.Resolve(file)
	if err != nil {
		panic(err.Error())
	}

	return &Extender{
		Bibliography: entries,
	}, nil
}

// Extend implements goldmark.Extender interface.
func (e *Extender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewCitationParser(), 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewCitationRenderer(e.Bibliography), 100),
		),
	)
}
