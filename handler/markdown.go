package handler

import (
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"bytes"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var extensions = parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
	parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.OrderedListStart |
	parser.BackslashLineBreak | parser.DefinitionLists | parser.EmptyLinesBreakList | parser.Footnotes |
	parser.SuperSubscript
var policy = bluemonday.NewPolicy()

func init() {
	policy.AllowElements("dl", "dt", "dd", "table", "th", "td", "tfoot", "h1", "h2", "h3", "h4", "h5", "h6",
		"pre", "code", "hr", "ul", "ol", "p", "a", "img", "mark", "blockquote", "details", "summary",
		"small", "li", "span", "tbody", "thead", "tr", "sub", "sup", "del", "strong", "em", "br")
	policy.AllowAttrs("class").Matching(bluemonday.SpaceSeparatedTokens).OnElements("span", "del", "mark")
	policy.AllowAttrs("href").OnElements("a")
	policy.AllowAttrs("colspan", "rowspan").Matching(bluemonday.Integer).OnElements("th", "td")
	policy.AllowAttrs("align").Matching(bluemonday.CellAlign).OnElements("th", "td")
	policy.AllowAttrs("start").Matching(bluemonday.Integer).OnElements("ol")

	policy.AllowStandardURLs()
	policy.AllowAttrs("align").Matching(bluemonday.ImageAlign).OnElements("img")
	policy.AllowAttrs("alt").Matching(bluemonday.Paragraph).OnElements("img")
	policy.AllowAttrs("style").Matching(regexp.MustCompile(`^ *(height|width) *: *[0-9]+(.[0-9]+)?(%|rem) *(; *)?$`)).OnElements("img")
	policy.AllowAttrs("src").OnElements("img")
	policy.AllowAttrs("referrerpolicy").Matching(regexp.MustCompile(`^no-referrer$`)).OnElements("img")
}

func MakeMarkdown(md string) template.HTML {
	if md == "" {
		return ""
	}
	htmlResult := markdown.NormalizeNewlines([]byte(md))
	htmlResult = markdown.ToHTML(htmlResult, parser.NewWithExtensions(extensions), getRenderer())
	htmlResult = bytes.ReplaceAll(htmlResult, []byte("<img"), []byte("<img referrerpolicy=\"no-referrer\" "))
	htmlResult = policy.SanitizeBytes(htmlResult)
	return template.HTML(htmlResult)
}

func PostMakeMarkdown(writer http.ResponseWriter, request *http.Request) {
	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		MakeSpecialPagePart(writer, &MarkdownBox{Information: MakeMarkdown(loc.MarkdownParseError)})
		return
	}
	MakeSpecialPagePart(writer, &MarkdownBox{Information: MakeMarkdown(values.GetTrimmedString("markdown"))})
}

// Stuff for removing the div

func getRenderer() *html.Renderer {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: myRenderHook,
	}
	return html.NewRenderer(opts)
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node.(type) {
	case *ast.List:
		List(w, node.(*ast.List), entering)
	case *ast.CodeBlock:
		CodeBlock(node.(*ast.CodeBlock))
		return ast.GoToNext, false
	default:
		return ast.GoToNext, false
	}
	return ast.GoToNext, true
}

// List writes ast.List node
func List(w io.Writer, list *ast.List, entering bool) {
	if entering {
		listEnter(w, list)
	} else {
		listExit(w, list)
	}
}

func listEnter(w io.Writer, nodeData *ast.List) {
	var attrs []string

	if nodeData.IsFootnotesList {
		_, _ = w.Write([]byte("\n<span class=\"footnotes\">\n\n<hr />"))
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
		_, _ = w.Write([]byte("\n</span>\n"))
	}
}

func CodeBlock(node *ast.CodeBlock) {
	maxLen := len(node.Literal)
	if node.Literal[maxLen-1] == ' ' && node.Literal[maxLen-2] == '\n' {
		maxLen -= 2
	}
	node.Literal = node.Literal[0:maxLen]
}
