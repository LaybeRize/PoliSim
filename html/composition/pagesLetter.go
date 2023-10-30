package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func GetCreateNormalLetterPage(acc *extraction.AccountAuth, letter *validation.CreateLetter, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message = Translation["errorQueryingNames"] + "\n" + val.Message
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getPageHeader(CreateLetter),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateLetter), CLASS("w-[800px]"),
			getUserDropdown(acc, letter.Account, Translation["accountLetter"]),
			getSimpleTextInput("title", "title", letter.Title, Translation["letterTitle"]),
			getCheckBox("noSigning", letter.NoSigning, false, "true", "noSigning", Translation["letterNoSigning"],
				HYPERSCRIPT("on click toggle .hidden on #allHaveToAgree")),
			getCheckBox("allHaveToAgree", letter.AllHaveToAgree, false, "true", "allHaveToAgree", Translation["letterAllHaveToAgree"], nil),
			getTextArea("content", "content", letter.Content, Translation["letterContent"], true),
			getEditableList(letter.Reader, "reader", "displayNames", Translation["addLetterReaderButtonText"], "w-[800px]"),
			getSubmitButton(Translation["createLetterButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}
