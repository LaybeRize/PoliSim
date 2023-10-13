package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

//TODO add all the new Translation to the DE.json

func GetCreateTitlePage(title *validation.TitleModification, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message += "\n" + Translation["errorQueryingNames"]
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getPageHeader(CreateTitle),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateTitle), CLASS("w-[800px]"),
			getSimpleTextInput("name", "name", title.Name, Translation["name"]),
			getEditableList([]string{}, "holder", "displayNames", Translation["addTitleHolderButtonText"], "w-[800px]"),
		),
	)
}

func GetModifyTitlePage(title *validation.TitleModification, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message += "\n" + Translation["errorQueryingNames"]
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getPageHeader(EditTitle),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(EditTitle), CLASS("w-[800px]"),
			getSubmitButtonOverwriteURL(Translation["searchTitleButton"], PATCH, "/"+APIPreRoute+string(SearchTitle)),
			getSimpleTextInput("newName", "newName", title.NewName, Translation["newName"]),
			getEditableList(title.Holder, "holder", "displayNames", Translation["addTitleHolderButtonText"], "w-[800px]"),
		),
	)
}
