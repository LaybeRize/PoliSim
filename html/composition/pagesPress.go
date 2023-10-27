package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func GetCreatePressReleasePage(acc *extraction.AccountAuth, press *validation.CreateArticle, val validation.Message) Node {
	return getBasePageWrapper(
		getPageHeader(CreateUser),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateUser), CLASS("w-[800px]"),
			getUserDropdown(acc, press.Account, Translation["role"]),
			getSimpleTextInput("title", "title", press.Title, Translation["pressTitle"]),
			getSimpleTextInput("subtitle", "subtitle", press.Subtitle, Translation["pressSubtitle"]),
			getCheckBox("breakingNews", press.IsBreakingNews, false, "true", "breakingNews", Translation["pressBreakingNews"], nil),
			getSubmitButton(Translation["createArticleButton"])),
		GetMessage(val),
	)
}
