package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func CreateDocumentPage(acc *extraction.AccountAuth, document *validation.CreateDocument, val validation.Message) Node {
	return getBasePageWrapper(
		getPageHeader(CreateTextDocument),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateTextDocument), CLASS("w-[800px]"),
			getUserDropdown(acc, document.Account, Translation["accountTextDocument"]),
			getSimpleTextInput("title", "title", document.Title, Translation["titleTextDocument"]),
			getSimpleTextInput("subtitle", "subtitle", document.Subtitle, Translation["subtitleTextDocument"]),
			getTextArea("content", "content", document.Content, Translation["contentTextDocument"], true),
			getSubmitButton(Translation["createTextDocumentButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}

func ViewDocumentPage(acc *extraction.AccountAuth, uuidStr string) Node {
	return getBasePageWrapper(
		Text(uuidStr),
		Text(acc.DisplayName),
	)
}
