package bibtex

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type citationParser struct{}

// NewCitationParser returns a new inline parser for citations.
func NewCitationParser() parser.InlineParser {
	return &citationParser{}
}

// Trigger implements parser.InlineParser interface.
func (s *citationParser) Trigger() []byte {
	return []byte{'@'}
}

// Parse implements parser.InlineParser interface.
func (s *citationParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	fmt.Println("Entering citationParser.Parse")
	line, _ := block.PeekLine()
	if len(line) <= 1 {
		fmt.Println("Exiting Parse: line too short")
		return nil
	}

	// Check if it's a citation
	if line[0] != '@' {
		fmt.Println("Exiting Parse: not a citation")
		return nil
	}

	// Find the citation key
	var i int
	for i = 1; i < len(line); i++ {
		if !isValidCitationChar(line[i]) {
			fmt.Printf("Stopped reading at index %d, encountered %c\n", i, line[i])
			break
		}
	}

	if i == 1 {
		fmt.Println("Exiting Parse: no valid citation key")
		return nil
	}

	block.Advance(i)
	citation := &Citation{
		BaseInline: ast.BaseInline{},
		Key:        string(line[1:i]),
		RawText:    string(line[:i]),
	}
	fmt.Printf("Exiting Parse: created citation with key %s\n", citation.Key)
	return citation
}

func isValidCitationChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_' || c == '-' || c == ':'
}
