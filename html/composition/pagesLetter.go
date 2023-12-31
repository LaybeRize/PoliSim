package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
	"slices"
	"strings"
)

func GetCreateNormalLetterPage(acc *database.AccountAuth, letter *validation.CreateLetter, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message = Translation["errorQueryingNames"] + "\n" + val.Message
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getPageHeader(CreateLetter),
		getFormStandardForm("form", POST, "/"+HTMXPreRouter+string(CreateLetter), CLASS("w-[800px]"),
			getUserDropdown(acc, letter.Account, Translation["accountLetter"]),
			getSimpleTextInput("title", "title", letter.Title, Translation["letterTitle"]),
			getCheckBoxWithHideScript(letter.NoSigning, "true", "noSigning", Translation["letterNoSigning"], "allHaveToAgree"),
			getStandardCheckBox(letter.AllHaveToAgree, "true", "allHaveToAgree", Translation["letterAllHaveToAgree"]),
			getTextArea("content", "content", letter.Content, Translation["letterContent"],
				MarkdownFormPage),
			getEditableList(letter.Reader, "reader", "displayNames", Translation["addLetterReaderButtonText"], "w-[800px]"),
			getSubmitButton("createLetterButton", Translation["createLetterButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}

func GetCreateModMailPage(letter *validation.CreateLetter, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message = Translation["errorQueryingNames"] + "\n" + val.Message
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getPageHeader(CreateModmail),
		getFormStandardForm("form", POST, "/"+HTMXPreRouter+string(CreateModmail), CLASS("w-[800px]"),
			getSimpleTextInput("authorAccount", "authorAccount", letter.Account, Translation["modMailAccount"]),
			getSimpleTextInput("flair", "flair", letter.Flair, Translation["modMailFlair"]),
			getSimpleTextInput("title", "title", letter.Title, Translation["modMailTitle"]),
			getCheckBoxWithHideScript(letter.NoSigning, "true", "noSigning", Translation["modMailNoSigning"], "allHaveToAgree"),
			getStandardCheckBox(letter.AllHaveToAgree, "true", "allHaveToAgree", Translation["modMailAllHaveToAgree"]),
			getTextArea("content", "content", letter.Content, Translation["modMailContent"],
				MarkdownFormPage),
			getEditableList(letter.Reader, "reader", "displayNames", Translation["addModMailReaderButtonText"], "w-[800px]"),
			getSubmitButton("createModMailButton", Translation["createModMailButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}

func GetLetterViewPersonalLetters(acc *database.AccountAuth, extra *logic.QueryInfo) Node {
	view, err := extra.GetLetter()
	if err != nil {
		return GetErrorPage(Translation["errorLoadingLetters"])
	}
	nodes := make([]Node, len(*view.LetterList))
	for i, item := range *view.LetterList {
		link := string(ViewLetterLink) + url.PathEscape(extra.ViewAccountName) + "/" + url.PathEscape(item.UUID)
		nodes[i] = getClickableLink("/"+HTMXPreRouter+link, "/"+link, Group(getStandardBoxClass,
			IfElse(item.Read, STYLE("--clr-border: rgb(40 51 69);"), STYLE("--clr-border: rgb(140 140 140);")),
			H1(CLASS("text-2xl"), Text(item.Title)),
			P(Text(Translation["authorShortFormLetter"], item.Author))))
	}
	if len(nodes) == 0 {
		nodes = []Node{
			DIV(CLASS("w-[800px] box box-e p-2 mt-2 flex items-center flex-col"),
				STYLE("--clr-border: rgb(40 51 69);"),
				P(CLASS("text-xl text-rose-600"), Text(Translation["noLettersFound"]))),
		}
	}
	return getBasePageWrapper(
		getPageHeader(ViewLetter),
		getUserModificationForLetters(acc, extra.ViewAccountName, Translation["selectedReaderLetter"]),
		Group(nodes...),
		pagerFooter(view.BeforeUUID, view.NextUUID,
			fmt.Sprintf("%s%s?uuid=%s&amount=%d&before=true", string(ViewLetterLink),
				url.PathEscape(extra.ViewAccountName), url.QueryEscape(view.BeforeUUID), extra.Amount),
			fmt.Sprintf("%s%s?uuid=%s&amount=%d", string(ViewLetterLink), url.PathEscape(extra.ViewAccountName),
				url.QueryEscape(view.NextUUID), extra.Amount)),
	)
}

func GetViewModmailList(acc *database.AccountAuth, extra *logic.QueryInfo) Node {
	view, err := extra.GetModMails()
	if err != nil {
		return GetErrorPage(Translation["errorLoadingLetters"])
	}
	nodes := make([]Node, len(*view.LetterList))
	for i, item := range *view.LetterList {
		link := string(ViewLetterLink) + url.PathEscape(acc.DisplayName) + "/" + url.PathEscape(item.UUID)
		nodes[i] = getClickableLink("/"+HTMXPreRouter+link, "/"+link, Group(getStandardBoxClass,
			STYLE("--clr-border: rgb(40 51 69);"),
			H1(CLASS("text-2xl"), Text(item.Title)),
			P(Text(Translation["authorShortFormLetter"], item.Author))))
	}

	return getBasePageWrapper(
		getPageHeader(ViewModMails),
		Group(nodes...),
		pagerFooter(view.BeforeUUID, view.NextUUID,
			fmt.Sprintf("%s?uuid=%s&amount=%d&before=true", string(ViewModMails),
				url.QueryEscape(view.BeforeUUID), extra.Amount),
			fmt.Sprintf("%s?uuid=%s&amount=%d", string(ViewModMails),
				url.QueryEscape(view.NextUUID), extra.Amount)),
	)
}

func getUserModificationForLetters(user *database.AccountAuth, selectedAccount string, labelText string) Node {
	return DIV(CLASS("mt-2 w-full"),
		LABEL(FOR("reader"), Text(labelText)),
		SELECT(ID("reader"), NAME("reader"), CLASS("bg-slate-700 appearance-none w-full py-2 px-3"),
			HXPATCH("/"+HTMXPreRouter+string(ChangeViewLetterAccount)), HXTRIGGER("change"),
			HXTARGET("#"+MainBodyID), HXSWAP("outerHTML"),
			getUserOptions(user, selectedAccount),
		),
		BUTTON(TYPE("submit"), CLASS(buttonClassAttribute+" mt-2 mr-2"), HXPATCH("/"+HTMXPreRouter+string(MarkAllLetterAccount)), HXTARGET("#"+MainBodyID),
			HXINCLUDE("[name='reader']"), HXSWAP("outerHTML"),
			ID("mark-all-read"), TEST("mark-all-read"), Text(Translation["markAllLettersRead"])),
		GetMessage(validation.Message{}),
	)
}

func GetSingLetterView(account *database.Account, letterUUID string, isMod bool, val validation.Message) Node {
	letter, err := extraction.GetLetterByIDOnlyWithAccount(letterUUID, account.ID, isMod)
	if err != nil {
		return GetErrorPage(Translation["errorWithSpecificLetter"])
	}
	go logic.SetLetterAsRead(letterUUID, account.ID)
	infoText := ""
	if letter.ModMessage {
		infoText = fmt.Sprintf(letter.Written.Format(Translation["authorModMessage"]), letter.Author)
	} else {
		infoText = fmt.Sprintf(letter.Written.Format(Translation["authorNormalLetter"]), letter.Author)
	}
	hasNotSigned := slices.Index(letter.Info.PeopleNotYetSigned, account.DisplayName) != -1
	var specialNode Node = nil
	if !letter.Info.NoSigning && !letter.Info.AllHaveToAgree {
		notYetSigned := strings.Join(letter.Info.PeopleNotYetSigned, ", ")
		signed := strings.Join(letter.Info.Signed, ", ")
		rejected := strings.Join(letter.Info.Rejected, ", ")
		specialNode = DIV(CLASS("w-[800px] mt-2"),
			If(notYetSigned != "", P(Text(Translation["peopleNotSigned"], notYetSigned))),
			If(signed != "", P(Text(Translation["peopleSigned"], signed))),
			If(rejected != "", P(Text(Translation["peopleRejectedLetter"], rejected))),
		)
	} else if !letter.Info.NoSigning && letter.Info.AllHaveToAgree {
		if len(letter.Info.Rejected) != 0 {
			specialNode = DIV(CLASS("w-[800px] mt-2"), P(Text(Translation["atLeastOneRejection"])))
		} else if len(letter.Info.PeopleNotYetSigned) != 0 {
			specialNode = DIV(CLASS("w-[800px] mt-2"), P(Text(Translation["noDecisionYet"])))
		} else {
			specialNode = DIV(CLASS("w-[800px] mt-2"), P(Text(Translation["everyoneSigned"])))
		}
	}
	return getBasePageWrapper(
		getPageHeader(ViewSingleLetter),
		DIV(CLASS("w-[800px]"),
			H1(CLASS("text-3xl underline decoration-2 underline-offset-2"), Text(letter.Title)),
			P(CLASS("my-2"), I(Text(infoText)),
				If(letter.Flair != "", Group(I(Text("; ")), Text(letter.Flair)))),
		),
		DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
			Raw(letter.HTMLContent),
		),
		DIV(CLASS("w-[800px]"),
			P(CLASS("mt-2"), Text(Translation["uuidLetterText"], letter.UUID)),
			P(Text(Translation["allViewerText"], strings.Join(append(append(letter.Info.Signed,
				letter.Info.PeopleNotYetSigned...),
				letter.Info.Rejected...), ", "))),
		),
		specialNode,
		If(hasNotSigned && !letter.Info.NoSigning, DIV(CLASS("w-[800px] flex flex-row"),
			updateLetterButton("/"+HTMXPreRouter+string(updateLetterLink)+
				url.PathEscape(account.DisplayName)+"/"+
				letterUUID+"/sign", Translation["signLetter"]),
			updateLetterButton("/"+HTMXPreRouter+string(updateLetterLink)+
				url.PathEscape(account.DisplayName)+"/"+
				letterUUID+"/reject", Translation["rejectLetter"]))),
		GetMessage(val),
	)
}

func updateLetterButton(link string, buttonText string) Node {
	return A(HXPATCH(link), HXTARGET("#"+MainBodyID),
		HXSWAP("outerHTML"),
		P(CLASS("bg-slate-700 text-white p-2 mr-4 mt-2 disable-selection"), Text(buttonText)),
	)
}
