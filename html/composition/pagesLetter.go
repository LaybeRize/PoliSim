package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
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

func GetLetterViewPersonalLetters(acc *extraction.AccountAuth, extra *logic.ExtraInfo) Node {
	view, err := extra.GetLetter()
	if err != nil {
		return GetErrorPage(Translation["errorLoadingLetters"])
	}
	nodes := make([]Node, len(*view.LetterList))
	for i, item := range *view.LetterList {
		link := string(ViewLetterLink) + "/" + url.PathEscape(extra.ViewAccountName) + "/" + url.PathEscape(item.UUID)
		nodes[i] = getClickableLink("/"+APIPreRoute+link, "/"+link, Group(
			CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
			H1(CLASS("text-2xl"), Text(item.Title)),
			P(Text(item.Author))))
	}
	beforeLink, nextLink := generateLetterLink(extra, view)
	return getBasePageWrapper(
		getPageHeader(ViewLetter),
		getUserDropdownForLetter(acc, extra.ViewAccountName, Translation["selectedReaderLetter"]),
		Group(nodes...),
		DIV(CLASS("w-[800px] flex justify-between flex-row"),
			DIV(If(view.BeforeUUID != "", getClickableLink("/"+APIPreRoute+beforeLink, "/"+beforeLink,
				P(CLASS("bg-slate-700 text-white p-2 mt-2"), Text("test")),
			))),
			DIV(If(view.NextUUID != "", getClickableLink("/"+APIPreRoute+nextLink, "/"+nextLink,
				P(CLASS("bg-slate-700 text-white p-2 mt-2"), Text("test")),
			))),
		),
	)
}

func generateLetterLink(extra *logic.ExtraInfo, view *logic.ViewLetter) (beforeLink string, nextLink string) {
	beforeLink = fmt.Sprintf("%s/%s?uuid=%s&amount=%d&before=true", string(ViewLetterLink), url.PathEscape(extra.ViewAccountName),
		url.QueryEscape(view.BeforeUUID), extra.Amount)
	nextLink = fmt.Sprintf("%s/%s?uuid=%s&amount=%d", string(ViewLetterLink), url.PathEscape(extra.ViewAccountName),
		url.QueryEscape(view.NextUUID), extra.Amount)
	return
}

func getUserDropdownForLetter(user *extraction.AccountAuth, selectedAccount string, labelText string) Node {
	return DIV(CLASS("mt-2 w-full"),
		LABEL(FOR("reader"), Text(labelText)),
		SELECT(ID("reader"), NAME("reader"), CLASS("bg-slate-700 appearance-none w-full py-2 px-3"),
			HXPATCH("/"+APIPreRoute+string(ChangeViewLetterAccount)), HXTRIGGER("change"),
			HXTARGET("#"+MainBodyID), HXSWAP("outerHTML"),
			getUserOptions(user, selectedAccount),
		),
	)
}
