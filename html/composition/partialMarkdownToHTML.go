package composition

import (
	"PoliSim/helper"
	. "PoliSim/html/builder"
)

func getPreviewElement() Node {
	return Group(
		DIV(ID(DisplayID), CLASS("")))
}

func GetUpdatePreviewElement(content string) Node {
	return DIV(ID(DisplayID), CLASS(""),
		Raw(helper.CreateHTML(content)))
}
