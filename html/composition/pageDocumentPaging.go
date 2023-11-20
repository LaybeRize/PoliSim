package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
)

func GetDocumentPage(isAdmin bool, extra *extraction.ExtraInfo) Node {
	view, err := logic.GetDocuments(isAdmin, extra)
	if err != nil {
		return GetErrorPage(Translation["errorLoadingDocuments"])
	}
	nodes := make([]Node, len(*view.DocumentList))
	for i, item := range *view.DocumentList {
		link := url.PathEscape(item.UUID)
		var classification Node
		switch item.Type {
		case database.LegislativeText:
			link = string(ViewTextDocumentLink) + link
			classification = I(CLASS("bi bi-card-text"))
		case database.FinishedDiscussion:
			link = string(ViewDiscussionDocumentLink) + link
			classification = I(CLASS("bi bi-chat-fill"))
		case database.RunningDiscussion:
			link = string(ViewDiscussionDocumentLink) + link
			classification = Group(I(CLASS("bi bi-chat-fill")), Text(" "), I(CLASS("bi bi-hourglass-split")))
		case database.FinishedVote:
			link = string(ViewVoteDocumentLink) + link
			classification = I(CLASS("bi bi-check2-square"))
		case database.RunningVote:
			link = string(ViewVoteDocumentLink) + link
			classification = Group(I(CLASS("bi bi-check2-square")), Text(" "), I(CLASS("bi bi-hourglass-split")))
		}
		nodes[i] = getClickableLink("/"+APIPreRoute+link, "/"+link, Group(
			CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
			H1(CLASS("text-2xl"), classification, Text(" "+item.Title)),
			P(I(Text(item.Written.Format(Translation["documentWrittenDate"])))),
			P(Text(Translation["authorOrganisationShortFormDocument"], item.Author, item.Organisation))))
	}
	if len(nodes) == 0 {
		nodes = []Node{
			DIV(CLASS("w-[800px] box box-e p-2 mt-2 flex items-center flex-col"), STYLE("--clr-border: rgb(40 51 69);"),
				P(CLASS("text-xl text-rose-600"), Text(Translation["noDocumentsFound"]))),
		}
	}
	before := fmt.Sprintf("%s?uuid=%s&amount=%d&before=true", string(ViewDocument),
		url.QueryEscape(view.BeforeUUID), extra.Amount)
	next := fmt.Sprintf("%s?uuid=%s&amount=%d", string(ViewDocument),
		url.QueryEscape(view.NextUUID), extra.Amount)
	extraStr := getExtraString(extra)
	return getBasePageWrapper(
		getPageHeader(ViewDocument),
		Group(nodes...),
		pagerFooter(view.BeforeUUID, view.NextUUID,
			before+extraStr, next+extraStr),
	)
}

func getExtraString(extra *extraction.ExtraInfo) string {
	result := ""
	if extra.HideBlock {
		result += "&hideblock=true"
	}
	if extra.Text {
		result += "&text=true"
	}
	if extra.Discussion {
		result += "&discussion=true"
	}
	if extra.Votes {
		result += "&votes=true"
	}
	return result
}
