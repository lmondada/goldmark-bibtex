# Goldmark BibTeX

A [Goldmark](https://github.com/yuin/goldmark) extension for rendering BibTeX citations in markdown, using the [bibtex](https://github.com/jschaf/bibtex) parser.

## Usage

First, create a BibTeX file (e.g., `references.bib`) with your citations:

```bibtex
@article{smith2023,
  author = {Smith, John},
  title = {An Important Discovery},
  journal = {Journal of Important Things},
  year = {2023}
}
```

Then use the extension in your Go code:

```go
bibExtender, err := bibtex.New("references.bib")
if err != nil {
    log.Fatal(err)
}

markdown := goldmark.New(
    goldmark.WithExtensions(bibExtender),
)

var buf bytes.Buffer
if err := markdown.Convert([]byte(source), &buf); err != nil {
    log.Fatal(err)
}
```

Cite references in your markdown using the @ symbol followed by the citation key:

```markdown
As shown in @smith2023, the results are significant.
```

This will be rendered as: "As shown in [Smith, 2023], the results are significant."

## Features

- Inline citations using @key format
- Support for different BibTeX entry types (article, book, etc.)
- Customizable citation formatting
- Integration with standard BibTeX files
- Simple integration with Goldmark markdown parser

## Installation

```bash
go get github.com/lmondada/goldmark-bibtex
```

## License

MIT License
