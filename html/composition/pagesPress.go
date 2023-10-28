package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func GetCreatePressReleasePage(acc *extraction.AccountAuth, press *validation.CreateArticle, val validation.Message) Node {
	return getBasePageWrapper(
		getPageHeader(CreateUser),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreatePressRelease), CLASS("w-[800px]"),
			getUserDropdown(acc, press.Account, Translation["accountPressRelease"]),
			getSimpleTextInput("title", "title", press.Title, Translation["pressTitle"]),
			getSimpleTextInput("subtitle", "subtitle", press.Subtitle, Translation["pressSubtitle"]),
			getCheckBox("breakingNews", press.IsBreakingNews, false, "true", "breakingNews", Translation["pressBreakingNews"], nil),
			getTextArea("content", "content", press.Content, Translation["pressContent"], true),
			getSubmitButton(Translation["createArticleButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}

func GetViewOfHiddenNewspaper() Node {
	list, err := extraction.GetHiddenPublication()
	if err != nil {
		return GetErrorPage(Translation["errorRetrievingPublication"])
	}
	nodes := make([]Node, len(*list))
	for i, item := range *list {
		nodes[i] = H1(Text(item.CreateTime.String()))
	}
	return getBasePageWrapper(
		getPageHeader(ViewHiddenNewspaperList),
		Group(nodes...),
	)
}

func GetViewSingleHiddenNewspaper(uuid string) Node {
	return getBasePageWrapper(
		DIV(Text(uuid)),
	)
}
