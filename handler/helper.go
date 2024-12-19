package handler

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"net/http"
	"strings"
)

var extensions = parser.NewWithExtensions(parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
	parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.OrderedListStart |
	parser.BackslashLineBreak | parser.DefinitionLists | parser.EmptyLinesBreakList | parser.Footnotes |
	parser.SuperSubscript)

func MakeMarkdown(md string) template.HTML {
	intermediate := markdown.NormalizeNewlines([]byte(md))
	maybeUnsafeHTML := markdown.ToHTML(intermediate, extensions, nil)
	htmlResult := string(bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML))
	return template.HTML(strings.ReplaceAll(strings.ReplaceAll(htmlResult, "<code>\n", "<code>"), "\n</code>", "</code>"))
}

func PostMakeMarkdown(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		MakeSpecialPagePart(writer, MARKDOWN, MarkdownBox{Information: MakeMarkdown("`Anfrage konnte nicht geparsed werden`")})
	}
	MakeSpecialPagePart(writer, MARKDOWN, MarkdownBox{Information: MakeMarkdown(request.Form.Get("markdown"))})
}
