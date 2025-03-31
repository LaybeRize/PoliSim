package handler

import (
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var extensions = parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
	parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.OrderedListStart |
	parser.BackslashLineBreak | parser.DefinitionLists | parser.EmptyLinesBreakList | parser.Footnotes |
	parser.SuperSubscript
var policy = bluemonday.NewPolicy().AllowElements("dl", "dt", "dd", "table", "th", "td", "tfoot", "h1", "h2", "h3", "h4", "h5", "h6", "pre", "code", "hr", "ul", "ol", "p", "a", "img", "mark", "blockquote", "details", "summary",
	"small", "li", "span", "tbody", "thead", "tr", "sub", "sup", "del", "strong", "em", "br").AllowAttrs("class").OnElements("span").AllowAttrs("alt",
	"src").OnElements("img").AllowAttrs("href").OnElements("a").AllowAttrs("colspan", "rowspan").OnElements("th", "td")

func init() {
	policy.AllowStandardURLs()
}

func MakeMarkdown(md string) template.HTML {
	if md == "" {
		return ""
	}
	htmlResult := policy.SanitizeBytes(markdown.ToHTML(markdown.NormalizeNewlines([]byte(md)), parser.NewWithExtensions(extensions), getRenderer()))
	return template.HTML(strings.ReplaceAll(strings.ReplaceAll(string(htmlResult), "<code>\n", "<code>"), "\n</code>", "</code>"))
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
	if para, ok := node.(*ast.List); ok {
		List(w, para, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
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
