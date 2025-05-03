package helper

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"io"
	"log"
	"regexp"
	"strings"
)

var extensions = parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
	parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.OrderedListStart |
	parser.BackslashLineBreak | parser.DefinitionLists | parser.EmptyLinesBreakList | parser.Footnotes |
	parser.SuperSubscript
var policy = bluemonday.NewPolicy()

func setupMarkdown() {
	log.Println("Setting up Markdown Rules")
	policy.AllowElements("dl", "dt", "dd", "table", "th", "td", "tfoot", "h1", "h2", "h3", "h4", "h5", "h6",
		"pre", "code", "hr", "ul", "ol", "p", "a", "img", "mark", "blockquote", "details", "summary",
		"small", "li", "span", "tbody", "thead", "tr", "sub", "sup", "del", "strong", "em", "br", "figure")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^(footnotes|added|removed)$`)).OnElements("figure", "del", "mark")
	policy.AllowAttrs("style").Matching(regexp.MustCompile(`^ *font-family *: *(serif|sans-serif|monospace) *(; *)?$`)).OnElements("span", "p")
	policy.AllowAttrs("style").Matching(regexp.MustCompile(`^ *list-style-type *: *[a-z-]* *(; *)?$`)).OnElements("ol", "ul")
	policy.AllowAttrs("href").OnElements("a")
	policy.AllowAttrs("colspan", "rowspan").Matching(bluemonday.Integer).OnElements("th", "td")
	policy.AllowAttrs("align").Matching(bluemonday.CellAlign).OnElements("th", "td")
	policy.AllowAttrs("start").Matching(bluemonday.Integer).OnElements("ol")
	policy.AllowAttrs("id").Matching(regexp.MustCompile(`^md-[a-zA-Z0-9-]*$`)).OnElements("span", "p", "h1", "h2", "h3", "h4", "h5", "h6")

	policy.AllowStandardURLs()
	policy.AllowAttrs("align").Matching(bluemonday.ImageAlign).OnElements("img")
	policy.AllowAttrs("alt").Matching(bluemonday.Paragraph).OnElements("img")
	policy.AllowAttrs("style").Matching(regexp.MustCompile(`^ *(height|width) *: *[0-9]+(.[0-9]+)?(%|rem) *(; *)?$`)).OnElements("img")
	policy.AllowAttrs("src").OnElements("img")
	policy.AllowAttrs("referrerpolicy").Matching(regexp.MustCompile(`^no-referrer$`)).OnElements("img", "a")
}

func MakeMarkdown(md string) template.HTML {
	if md == "" {
		return ""
	}
	htmlResult := markdown.NormalizeNewlines([]byte(md))
	htmlResult = markdown.ToHTML(htmlResult, parser.NewWithExtensions(extensions), getRenderer())
	htmlResult = bytes.ReplaceAll(htmlResult, []byte("<img"), []byte("<img referrerpolicy=\"no-referrer\" "))
	htmlResult = policy.SanitizeBytes(htmlResult)
	htmlResult = addAutoIDForHeadings(htmlResult)
	return template.HTML(htmlResult)
}

func getRenderer() *html.Renderer {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: myRenderHook,
	}
	return html.NewRenderer(opts)
}

func addAutoIDForHeadings(htmlInput []byte) []byte {
	result := bytes.NewBuffer([]byte{})
	scanner := bufio.NewScanner(bytes.NewReader(htmlInput))

	counter := 0
	transform := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		if len(data) >= 4 && data[0] == '<' && data[1] == 'h' &&
			data[2] >= 49 && data[2] <= 49+5 &&
			data[3] == '>' {
			//data starts with a heading that can be extended with an automatic generated ID
			counter += 1
			return 4, []byte{'<', 'h', data[2],
				' ', 'i', 'd', '=', '"', 'a', 'u', 't', 'o', '-', 'm', 'd', '-',
				byte(48 + ((counter / 100) % 10)),
				byte(48 + ((counter / 10) % 10)),
				byte(48 + (counter % 10)),
				'"', '>'}, nil
		} else if len(data) >= 4 && data[0] == '<' {
			//data starts with a < but is not a heading or one that can not be processed
			if i := bytes.Index(data[1:], []byte("<")); i > 0 {
				return i + 1, data[:i+1], nil
			}
		} else if i := bytes.Index(data, []byte("<")); i > 0 {
			// data does not start with a < and can safely be returned to the user
			return i, data[:i], nil
		}
		//returning all data that has not been read yet but there are no more headings
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		//end the process
		if atEOF {
			return 0, nil, io.EOF
		}
		if len(data) > (bufio.MaxScanTokenSize / 2) {
			return len(data), data, nil
		}
		//request more data if none of the above have handled the data
		return 0, nil, nil
	}

	scanner.Split(transform)
	for scanner.Scan() {
		_, _ = result.Write(scanner.Bytes())
	}
	return result.Bytes()
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node.(type) {
	case *ast.List:
		list(w, node.(*ast.List), entering)
	case *ast.CodeBlock:
		codeBlock(node.(*ast.CodeBlock))
		return ast.GoToNext, false
	default:
		return ast.GoToNext, false
	}
	return ast.GoToNext, true
}

func list(w io.Writer, list *ast.List, entering bool) {
	if entering {
		listEnter(w, list)
	} else {
		listExit(w, list)
	}
}

func listEnter(w io.Writer, nodeData *ast.List) {
	var attrs []string

	if nodeData.IsFootnotesList {
		_, _ = w.Write([]byte("\n<figure class=\"footnotes\">\n\n<hr />"))
	}
	if html.IsListItem(nodeData.Parent) {
		grand := nodeData.Parent.GetParent()
		if html.IsListTight(grand) {
			_, _ = w.Write([]byte("\n"))
		}
	}

	openTag := "<ul"
	if nodeData.ListFlags&ast.ListTypeOrdered != 0 {
		if nodeData.Start > 0 {
			attrs = append(attrs, fmt.Sprintf(`start="%d"`, nodeData.Start))
		}
		openTag = "<ol"
	}
	if nodeData.ListFlags&ast.ListTypeDefinition != 0 {
		openTag = "<dl"
	}
	attrs = append(attrs, html.BlockAttrs(nodeData)...)
	_, _ = w.Write([]byte(openTag + " " + strings.Join(attrs, " ") + " >\n"))
}

func listExit(w io.Writer, list *ast.List) {
	closeTag := "</ul>"
	if list.ListFlags&ast.ListTypeOrdered != 0 {
		closeTag = "</ol>"
	}
	if list.ListFlags&ast.ListTypeDefinition != 0 {
		closeTag = "</dl>"
	}
	_, _ = w.Write([]byte(closeTag))

	if list.IsFootnotesList {
		_, _ = w.Write([]byte("\n</figure>\n"))
	}
}

func codeBlock(node *ast.CodeBlock) {
	maxLen := len(node.Literal)
	if node.Literal[maxLen-1] == ' ' && node.Literal[maxLen-2] == '\n' {
		maxLen -= 2
	}
	node.Literal = node.Literal[0:maxLen]
}
