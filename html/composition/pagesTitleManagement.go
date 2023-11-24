package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func GetCreateTitlePage(title *validation.TitleModification, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message = Translation["errorQueryingNames"] + "\n" + val.Message
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getDataList("mainGroupNames", extraction.TitleMainGroupList),
		getDataListFromMap("subGroupNames", extraction.TitleSubGroupNameMap),
		getPageHeader(CreateTitle),
		getFormStandardForm("form", POST, "/"+HTMXPreRouter+string(CreateTitle), CLASS("w-[800px]"),
			getSimpleTextInput("name", "name", title.Name, Translation["titleName"]),
			getSimpleTextInput("flair", "flair", title.Flair, Translation["flair"]),
			getInput("mainGroup", "mainGroup", title.MainGroup, Translation["mainGroup"], "text", "mainGroupNames", ""),
			getInput("subGroup", "subGroup", title.SubGroup, Translation["subGroup"], "text", "subGroupNames", ""),
			getEditableList(title.Holder, "holder", "displayNames", Translation["addTitleHolderButtonText"], "w-[800px]"),
			getSubmitButton("createTitleButton", Translation["createTitleButton"]),
		),
		GetMessage(val),
	)
}

func GetModifyTitlePage(title *validation.TitleModification, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message = Translation["errorQueryingNames"] + "\n" + val.Message
	}
	titleNames, err := extraction.GetAllTitleNames()
	if err != nil {
		val.Message = Translation["errorQueryingTitleNames"] + "\n" + val.Message
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getDataList("titleNames", titleNames),
		getDataList("mainGroupNames", extraction.TitleMainGroupList),
		getDataListFromMap("subGroupNames", extraction.TitleSubGroupNameMap),
		getPageHeader(EditTitle),
		getFormStandardForm("form", POST, "/"+HTMXPreRouter+string(EditTitle), CLASS("w-[800px]"),
			getInput("name", "name", title.Name, Translation["titleName"], "text", "titleNames", ""),
			getSubmitButtonOverwriteURL(Translation["searchTitleButton"], PATCH, "/"+HTMXPreRouter+string(SearchTitle)),
			getSimpleTextInput("newName", "newName", title.NewName, Translation["newTitleName"]),
			getInput("mainGroup", "mainGroup", title.MainGroup, Translation["mainGroup"], "text", "mainGroupNames", ""),
			getInput("subGroup", "subGroup", title.SubGroup, Translation["subGroup"], "text", "subGroupNames", ""),
			getSimpleTextInput("flair", "flair", title.Flair, Translation["flair"]),
			getEditableList(title.Holder, "holder", "displayNames", Translation["addTitleHolderButtonText"], "w-[800px]"),
			getSubmitButton("modifyTitleButton", Translation["changeTitleButton"]),
			getSubmitButtonOverwriteURL(Translation["deleteTitleButton"], PATCH, "/"+HTMXPreRouter+string(DeleteTitle)),
		),
		GetMessage(val),
	)
}
