package handler

import (
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"net/http"
	"strings"
)

var extensions = parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
	parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.OrderedListStart |
	parser.BackslashLineBreak | parser.DefinitionLists | parser.EmptyLinesBreakList | parser.Footnotes |
	parser.SuperSubscript
var policy = bluemonday.NewPolicy().AllowElements("mark", "details", "summary", "small")

func MakeMarkdown(md string) template.HTML {
	if md == "" {
		return ""
	}
	intermediate := markdown.NormalizeNewlines(policy.SanitizeBytes([]byte(md)))
	htmlResult := markdown.ToHTML(intermediate, parser.NewWithExtensions(extensions), nil)
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
