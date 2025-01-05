package acm

import (
	"fmt"
	"strings"

	"github.com/jschaf/bibtex"
	bibtexAst "github.com/jschaf/bibtex/ast"
)

func FormatAuthor(author *bibtexAst.Author) (authorFmt string) {
	authorFmt += author.First.(*bibtexAst.Text).Value
	authorFmt += " "
	authorFmt += author.Prefix.(*bibtexAst.Text).Value
	authorFmt += " "
	authorFmt += author.Last.(*bibtexAst.Text).Value
	return
}

// Join list of strings with commas and "and" for the last one
func join(list []string) (s string) {
	for i, item := range list {
		if i > 0 {
			if i == len(list)-1 {
				s += " and "
			} else {
				s += ", "
			}
		}
		s += item
	}
	return
}

// FormatAuthors formats a list of authors according to ACM style
// ACM style uses full names and separates authors with commas, using "and" for the last author
func FormatAuthors(authors bibtexAst.Authors) string {
	authorList := make([]string, len(authors))
	for i, author := range authors {
		authorList[i] = FormatAuthor(author)
	}

	return fmt.Sprintf(`<span class="authors">%s</span>`, join(authorList))
}

func formatDoi(doi string) string {
	return fmt.Sprintf(`doi: <a href="https://doi.org/%s">%s</a>`, doi, doi)
}

func getArticleRef(entry *bibtex.Entry) articleRef {
	return articleRef{
		authors: FormatAuthors(entry.Tags["author"].(bibtexAst.Authors)),
		year:    getFieldText(entry, "year"),
		title:   getFieldText(entry, "title"),
		journal: getFieldText(entry, "journal"),
		number:  getFieldText(entry, "number"),
		volume:  getFieldText(entry, "volume"),
		pages:   getFieldText(entry, "pages"),
		doi:     getFieldText(entry, "doi"),
		month:   getFieldText(entry, "month"),
	}
}

func getProceedingsRef(entry *bibtex.Entry) proceedingsRef {
	return proceedingsRef{
		authors:   FormatAuthors(entry.Tags["author"].(bibtexAst.Authors)),
		year:      getFieldText(entry, "year"),
		title:     getFieldText(entry, "title"),
		booktitle: getFieldText(entry, "booktitle"),
		month:     getFieldText(entry, "month"),
		address:   getFieldText(entry, "address"),
		pages:     getFieldText(entry, "pages"),
		doi:       getFieldText(entry, "doi"),
		publisher: getFieldText(entry, "publisher"),
	}
}

func getBookRef(entry *bibtex.Entry) bookRef {
	return bookRef{
		authors:   FormatAuthors(entry.Tags["author"].(bibtexAst.Authors)),
		year:      getFieldText(entry, "year"),
		title:     getFieldText(entry, "title"),
		publisher: getFieldText(entry, "publisher"),
		address:   getFieldText(entry, "address"),
		edition:   getFieldText(entry, "edition"),
		doi:       getFieldText(entry, "doi"),
	}
}

func getArxivRef(entry *bibtex.Entry) arxivRef {
	return arxivRef{
		authors:      getFieldText(entry, "author"),
		year:         getFieldText(entry, "year"),
		title:        getFieldText(entry, "title"),
		eprint:       getFieldText(entry, "eprint"),
		primaryClass: getFieldText(entry, "primaryclass"),
	}
}

func getDefaultRef(entry *bibtex.Entry) defaultRef {
	return defaultRef{
		authors:      getFieldText(entry, "author"),
		year:         getFieldText(entry, "year"),
		month:        getFieldText(entry, "month"),
		title:        getFieldText(entry, "title"),
		howpublished: getFieldText(entry, "howpublished"),
		url:          getFieldText(entry, "url"),
	}
}

// FormatCitation formats a full citation in ACM style
func FormatCitation(entry *bibtex.Entry) string {
	archivePrefix := getFieldText(entry, "archiveprefix")

	switch strings.ToLower(entry.Type) {
	case "article":
		articleRef := getArticleRef(entry)
		return formatArticle(articleRef)
	case "inproceedings", "conference":
		proceedingsRef := getProceedingsRef(entry)
		return formatProceedings(proceedingsRef)
	case "book":
		bookRef := getBookRef(entry)
		return formatBook(bookRef)
	default:
		if !strings.EqualFold(archivePrefix, "arXiv") {
			defaultRef := getDefaultRef(entry)
			return formatDefault(defaultRef)
		}
		// Handle arXiv papers specially
		arxivRef := getArxivRef(entry)
		return formatArxiv(arxivRef)
	}
}

type articleRef struct {
	authors string
	year    string
	title   string
	journal string
	number  string
	volume  string
	pages   string
	doi     string
	month   string
}

// formatArticle formats an article citation in ACM style
// Example: Patricia S. Abril and Robert Plant. 2007. The patent holder's dilemma: Buy, sell, or troll? Commun. ACM 50, 1 (Jan. 2007), 36-44. https://doi.org/10.1145/1188913.1188915
func formatArticle(article articleRef) string {
	citation := fmt.Sprintf(
		`<span class="citation">%s. %s. %s. <em>%s</em>`,
		article.authors, article.year, article.title, article.journal,
	)

	if article.volume != "" {
		citation += " " + article.volume
		if article.number != "" {
			citation += ", " + article.number
		}
	}

	if article.month != "" || article.pages != "" {
		citation += " ("
		if article.month != "" {
			citation += article.month + " " + article.year
		}
		if article.pages != "" {
			if article.month != "" {
				citation += ", "
			}
			citation += article.pages
		}
		citation += ")"
	}

	if article.doi != "" {
		citation += fmt.Sprintf(". %s", formatDoi(article.doi))
	}

	citation += "</span>"
	return citation
}

type proceedingsRef struct {
	authors   string
	year      string
	title     string
	booktitle string
	pages     string
	doi       string
	address   string
	publisher string
	month     string
}

// formatProceedings formats a conference proceedings citation in ACM style
// Example: Sten Andler. 1979. Predicate path expressions. In Proceedings of the 6th. ACM SIGACT-SIGPLAN Symposium on Principles of Programming Languages (POPL '79), January 29 - 31, 1979, San Antonio, Texas. ACM Inc., New York, NY, 226-236. https://doi.org/10.1145/567752.567774
func formatProceedings(ref proceedingsRef) string {
	citation := fmt.Sprintf(`<span class="citation">%s. %s. %s. In <em>%s</em>`,
		ref.authors, ref.year, ref.title, ref.booktitle)

	if ref.month != "" || ref.address != "" {
		citation += ", "
		if ref.month != "" {
			citation += ref.month + " " + ref.year
		}
		if ref.address != "" {
			if ref.month != "" {
				citation += ", "
			}
			citation += ref.address
		}
	}

	if ref.publisher != "" {
		citation += ". " + ref.publisher
	}

	if ref.pages != "" {
		citation += ", " + ref.pages
	}

	if ref.doi != "" {
		citation += fmt.Sprintf(". %s", formatDoi(ref.doi))
	}

	citation += "</span>"
	return citation
}

type bookRef struct {
	authors   string
	year      string
	title     string
	publisher string
	address   string
	edition   string
	doi       string
}

// formatBook formats a book citation in ACM style
// Example: David Kosiur. 2001. Understanding Policy-Based Networking (2nd. ed.). Wiley, New York, NY.
func formatBook(ref bookRef) string {
	citation := fmt.Sprintf(`<span class="citation">%s. %s. <em>%s</em>`,
		ref.authors, ref.year, ref.title)

	if ref.edition != "" {
		citation += fmt.Sprintf(" (%s ed.)", ref.edition)
	}

	if ref.publisher != "" {
		citation += ". " + ref.publisher
		if ref.address != "" {
			citation += ", " + ref.address
		}
	}

	if ref.doi != "" {
		citation += fmt.Sprintf(". %s", formatDoi(ref.doi))
	}

	citation += "</span>"
	return citation
}

type arxivRef struct {
	authors      string
	year         string
	title        string
	eprint       string
	primaryClass string
}

// formatArxiv formats an arXiv paper citation in ACM style
// Example: "Ali Javadi-Abhari et al. 2024. Quantum computing with Qiskit. arXiv: 2405.08810 [quant-ph]"
func formatArxiv(ref arxivRef) string {
	citation := fmt.Sprintf(`<span class="citation">%s. %s. %s`,
		ref.authors, ref.year, ref.title)

	if ref.eprint != "" {
		citation += fmt.Sprintf(". arXiv: <a href=\"https://arxiv.org/abs/%s\">%s", ref.eprint, ref.eprint)
		if ref.primaryClass != "" {
			citation += fmt.Sprintf(" [%s]", ref.primaryClass)
		}
		citation += "</a>"
	}

	citation += "</span>"
	return citation
}

type defaultRef struct {
	authors      string
	year         string
	month        string
	title        string
	howpublished string
	url          string
}

// formatDefault formats other types of citations in ACM style
// This is a basic formatter that includes the essential elements of a citation:
// authors, year, title, and URL/DOI if available
func formatDefault(ref defaultRef) string {
	citation := fmt.Sprintf(`<span class="citation">%s. %s. %s`,
		ref.authors, ref.year, ref.title)

	if ref.month != "" {
		citation += fmt.Sprintf(". (%s %s)", ref.month, ref.year)
	}

	if ref.url != "" || ref.howpublished != "" {
		citation += ". "
		if ref.url != "" {
			citation += "Retrieved "
		}
		if ref.howpublished != "" {
			citation += ref.howpublished
			citation += " "
		}
		if ref.url != "" {
			citation += fmt.Sprintf("from <a href=\"%s\">%s</a>", ref.url, ref.url)
		}
	}

	citation += "</span>"
	return citation
}

// getFieldText safely gets the text of a BibTeX field
func getFieldText(entry *bibtex.Entry, field string) string {
	if f := entry.Tags[field]; f != nil {
		return f.(*bibtexAst.Text).Value
	}
	return ""
}
