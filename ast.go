package bibtex

import (
	"github.com/yuin/goldmark/ast"
)

// Citation represents a citation node in the AST.
type Citation struct {
	ast.BaseInline
	Key     string
	RawText string
}

var CitationKind = ast.NewNodeKind("Citation")

func (n *Citation) Kind() ast.NodeKind {
	return CitationKind
}

// Dump implements Node.Dump.
func (n *Citation) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
