package handler

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var extensions = parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
	parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.OrderedListStart |
	parser.BackslashLineBreak | parser.DefinitionLists | parser.EmptyLinesBreakList | parser.Footnotes |
	parser.SuperSubscript
var policy = bluemonday.UGCPolicy().AllowAttrs("class").OnElements("div")

func MakeMarkdown(md string) template.HTML {
	intermediate := markdown.NormalizeNewlines([]byte(md))
	maybeUnsafeHTML := markdown.ToHTML(intermediate, parser.NewWithExtensions(extensions), nil)
	htmlResult := string(policy.SanitizeBytes(maybeUnsafeHTML))
	return template.HTML(strings.ReplaceAll(strings.ReplaceAll(htmlResult, "<code>\n", "<code>"), "\n</code>", "</code>"))
}

func PostMakeMarkdown(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		MakeSpecialPagePart(writer, &MarkdownBox{Information: MakeMarkdown("`Anfrage konnte nicht geparsed werden`")})
	}
	MakeSpecialPagePart(writer, &MarkdownBox{Information: MakeMarkdown(request.Form.Get("markdown"))})
}

func GetUniqueID(author string) string {
	authorRunes := []rune(author)
	sum := 0
	for i, singleRune := range authorRunes {
		if i > 3 {
			break
		}
		sum += int(singleRune)
	}
	return fmt.Sprintf("%x-%x", sum, time.Now().UnixNano()/1000000)
}

func MakeCommaSeperatedStringToList(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return make([]string, 0)
	}
	arr := strings.Split(input, ",")
	result := make([]string, 0, len(arr))
	for _, element := range arr {
		element = strings.TrimSpace(element)
		if element != "" {
			result = append(result, element)
		}
	}
	return result
}
