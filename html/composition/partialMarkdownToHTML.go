package composition

import (
	"PoliSim/helper"
	. "PoliSim/html/builder"
)

func getPreviewElement() Node {
	return DIV(CLASS("w-[800px] mt-2"),
		H1(CLASS("text-2xl text-white mb-2"), Text(Translation["previewTitle"])),
		DIV(ID(DisplayID), CLASS("")))
}

func GetUpdatePreviewElement(content string) Node {
	return DIV(ID(DisplayID), CLASS("w-full box box-e p-2"), STYLE("--clr-border: rgb(40 51 69);"),
		Raw(helper.CreateHTML(content)))
}
