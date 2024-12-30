package apa

import (
	"fmt"

	"github.com/jschaf/bibtex"
	bibtexAst "github.com/jschaf/bibtex/ast"
)

// FormatAuthors formats a list of authors according to APA style
func FormatAuthors(authors bibtexAst.Authors) string {
	var authorList string
	for i, author := range authors {
		if i > 0 {
			if i == len(authors)-1 {
				authorList += ", & "
			} else {
				authorList += ", "
			}
		}
		authorList += fmt.Sprintf(`<span class="author">%s %s</span>`,
			author.Last.(*bibtexAst.Text).Value,
			author.First.(*bibtexAst.Text).Value)
	}
	return authorList
}

// FormatCitationKey formats a short citation key
func FormatCitationKey(entry *bibtex.Entry) string {
	authors := entry.Tags["author"].(bibtexAst.Authors)
	year := entry.Tags["year"].(*bibtexAst.Text)
	firstAuthor := authors[0]
	lastName := firstAuthor.Last.(*bibtexAst.Text).Value
	return fmt.Sprintf(`<span class="citation-key">%s, %s</span>`,
		TrimLastName(lastName),
		year.Value)
}

// TrimLastName trims an author's last name to 6 characters if it's longer
func TrimLastName(name string) string {
	if len(name) > 7 {
		return name[:6] + "."
	}
	return name
}

// FormatCitation formats a full citation in APA style
func FormatCitation(entry *bibtex.Entry) string {
	authors := entry.Tags["author"].(bibtexAst.Authors)
	year := entry.Tags["year"].(*bibtexAst.Text)
	authorList := FormatAuthors(authors)

	var citation string
	switch entry.Type {
	case bibtex.EntryArticle:
		citation = formatArticle(authorList, year.Value, entry)
	case bibtex.EntryInProceedings:
		citation = formatProceedings(authorList, year.Value, entry)
	case bibtex.EntryBook:
		citation = formatBook(authorList, year.Value, entry)
	default:
		citation = formatDefault(authorList, year.Value, entry)
	}

	return citation
}

func formatArticle(authors, year string, entry *bibtex.Entry) string {
	title := getFieldText(entry, "title")
	journal := getFieldText(entry, "journal")
	volume := getFieldText(entry, "volume")
	pages := getFieldText(entry, "pages")

	citation := fmt.Sprintf(`<span class="citation-full">%s (%s). <span class="title">%s</span>. <span class="journal">%s</span>`,
		authors, year, title, journal)
	if volume != "" {
		citation += fmt.Sprintf(`, <span class="volume">%s</span>`, volume)
	}
	if pages != "" {
		citation += fmt.Sprintf(`, <span class="pages">%s</span>`, pages)
	}
	citation += ".</span>"
	return citation
}

func formatProceedings(authors, year string, entry *bibtex.Entry) string {
	title := getFieldText(entry, "title")
	booktitle := getFieldText(entry, "booktitle")
	pages := getFieldText(entry, "pages")

	citation := fmt.Sprintf(`<span class="citation-full">%s (%s). <span class="title">%s</span>. In <span class="booktitle">%s</span>`,
		authors, year, title, booktitle)
	if pages != "" {
		citation += fmt.Sprintf(` (pp. <span class="pages">%s</span>)`, pages)
	}
	citation += ".</span>"
	return citation
}

func formatBook(authors, year string, entry *bibtex.Entry) string {
	title := getFieldText(entry, "title")
	publisher := getFieldText(entry, "publisher")

	citation := fmt.Sprintf(`<span class="citation-full">%s (%s). <span class="title"><em>%s</em></span>`,
		authors, year, title)
	if publisher != "" {
		citation += fmt.Sprintf(`. <span class="publisher">%s</span>`, publisher)
	}
	citation += ".</span>"
	return citation
}

func formatDefault(authors, year string, entry *bibtex.Entry) string {
	title := getFieldText(entry, "title")
	return fmt.Sprintf(`<span class="citation-full">%s (%s). <span class="title">%s</span>.</span>`,
		authors, year, title)
}

func getFieldText(entry *bibtex.Entry, field string) string {
	if v, ok := entry.Tags[field]; ok {
		return v.(*bibtexAst.Text).Value
	}
	return fmt.Sprintf("??%s??", field)
}
