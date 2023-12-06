package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
)

const (
	MinZwitscher = 5
	MaxZwitscher = 50
)

func GetZwitschers(acc *database.AccountAuth, extra *extraction.ExtraZwitscherInfo) Node {
	view, err := logic.GetZwitschers(extra)
	if err != nil {
		return GetErrorPage(Translation["errorLoadingZwitscher"])
	}
	nodes := make([]Node, len(*view.ZwitscherList))
	for i, item := range *view.ZwitscherList {
		link := string(ViewSingleZwitscherLink) + url.PathEscape(item.UUID)
		nodes[i] = getClickableLink("/"+HTMXPreRouter+link, "/"+link, Group(getStandardBoxClass,
			IfElse(item.Blocked, STYLE("--clr-border: rgb(159 18 57);"), STYLE("--clr-border: rgb(40 51 69);")),
			H1(CLASS("text-2xl"), Text(item.Author), If(item.Flair != "", Group(Text("; "), I(Text(item.Flair))))),
			P(I(Text(item.Written.Format(Translation["zwitscherWrittenTime"])))),
			Raw(item.HTMLContent)))
	}
	if len(nodes) == 0 {
		nodes = []Node{
			DIV(CLASS("w-[800px] box box-e p-2 mt-2 flex items-center flex-col"),
				STYLE("--clr-border: rgb(40 51 69);"),
				P(CLASS("text-xl text-rose-600"), Text(Translation["noZwitscherFound"]))),
		}
	}
	before := fmt.Sprintf("%s?uuid=%s&amount=%d&before=true", string(ViewZwitscher),
		url.QueryEscape(view.BeforeUUID), extra.Amount)
	next := fmt.Sprintf("%s?uuid=%s&amount=%d", string(ViewZwitscher),
		url.QueryEscape(view.NextUUID), extra.Amount)
	extraStr := GetExtraStringForZwitscher(extra)
	return getBasePageWrapper(
		getPageHeader(ViewZwitscher),
		If(acc.Role != database.NotLoggedIn, toggleVisiblityOfNextDiv(Translation["showTweetInterface"])),
		GetZwitscherInterface(acc, "", true, validation.Message{}),
		Group(nodes...),
		pagerFooter(view.BeforeUUID, view.NextUUID,
			before+extraStr, next+extraStr),
	)
}

func GetSingleZwitscher(acc *database.AccountAuth, isAdmin bool, uuid string) Node {
	zwt, err := extraction.GetZwitscher(uuid)
	if err != nil {
		return GetErrorPage(Translation["tweetWithUUIDDoesNotExist"])
	}
	return getBasePageWrapper(DIV(Text(zwt.HTMLContent)))
}

const ZwitscherInterfaceID = "zwitscher-interface-id"

func GetZwitscherInterface(acc *database.AccountAuth, responseTo string, hideInterface bool, msg validation.Message) Node {
	extra := ""
	if responseTo != "" {
		extra = "?zwitscher=" + url.QueryEscape(responseTo)
	}
	return DIV(CLASS("w-[800px]"), If(hideInterface, HIDDEN()), ID(ZwitscherInterfaceID),
		getFormStandardForm("createZwitscher", POST, "/"+HTMXPreRouter+string(CreateZwitscher)+extra,
			getUserDropdown(acc, "", Translation["zwitscherPerson"]),
			getTextArea("zwitscherContent", "content", "", Translation["zwitscherContent"],
				MarkdownFormPage),
			getSubmitButton("submitCommentButton", Translation["postZwitscher"]),
		),
		getPreviewElement(),
		GetMessage(msg))
}

func GetAuthorQueryString(author string) string {
	return "author=" + url.QueryEscape(author)
}

func GetExtraStringForZwitscher(extra *extraction.ExtraZwitscherInfo) string {
	result := ""
	if extra.HideBlock {
		result += "&hideblock=true"
	}
	if extra.ShowOnlyReplies {
		result += "&onlyreplies=true"
	}
	if extra.ShowOnlyZwitscher {
		result += "&onlyzwitscher=true"
	}
	if extra.Author != "" {
		result += "&" + GetAuthorQueryString(extra.Author)
	}
	return result
}
