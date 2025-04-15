package handler

import (
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"net/http"
)

func PostMakeMarkdown(writer http.ResponseWriter, request *http.Request) {
	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		MakeSpecialPagePart(writer, &MarkdownBox{Information: helper.MakeMarkdown(loc.MarkdownParseError)})
		return
	}
	MakeSpecialPagePart(writer, &MarkdownBox{Information: helper.MakeMarkdown(values.GetTrimmedString("markdown"))})
}
