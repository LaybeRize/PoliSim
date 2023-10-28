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
