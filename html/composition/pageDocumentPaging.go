package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	MinDocuments = 5
	MaxDocuments = 50
)

var (
	textSymbol = I(CLASS("bi bi-card-text")).PureNode()
	discSymbol = I(CLASS("bi bi-chat-fill")).PureNode()
	voteSymbol = I(CLASS("bi bi-check2-square")).PureNode()
	runSymbol  = I(CLASS("bi bi-hourglass-split")).PureNode()
)

func GetDocumentPage(extra *extraction.DocumentQueryInfo) Node {
	view, err := logic.GetDocuments(extra)
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
			classification = textSymbol
		case database.FinishedDiscussion:
			link = string(ViewDiscussionDocumentLink) + link
			classification = discSymbol
		case database.RunningDiscussion:
			link = string(ViewDiscussionDocumentLink) + link
			classification = Group(discSymbol, Text(" "), runSymbol)
		case database.FinishedVote:
			link = string(ViewVoteDocumentLink) + link
			classification = voteSymbol
		case database.RunningVote:
			link = string(ViewVoteDocumentLink) + link
			classification = Group(voteSymbol, Text(" "), runSymbol)
		}
		nodes[i] = getClickableLink("/"+HTMXPreRouter+link, "/"+link, Group(
			getStandardBoxClass, STYLE("--clr-border: rgb(40 51 69);"),
			H1(CLASS("text-2xl"), classification, Text(" "+item.Title)),
			P(I(Text(item.Written.Format(Translation["documentWrittenDate"])))),
			P(Text(Translation["authorOrganisationShortFormDocument"], item.Author, item.Organisation))))
	}
	if len(nodes) == 0 {
		nodes = []Node{
			DIV(CLASS("w-[800px] box box-e p-2 mt-2 flex items-center flex-col"),
				STYLE("--clr-border: rgb(40 51 69);"),
				P(CLASS("text-xl text-rose-600"), Text(Translation["noDocumentsFound"]))),
		}
	}
	before := fmt.Sprintf("%s?uuid=%s&amount=%d&before=true", string(ViewDocument),
		url.QueryEscape(view.BeforeUUID), extra.Amount)
	next := fmt.Sprintf("%s?uuid=%s&amount=%d", string(ViewDocument),
		url.QueryEscape(view.NextUUID), extra.Amount)
	extraStr := GetExtraString(extra)
	return getBasePageWrapper(
		getPageHeader(ViewDocument),
		toggleVisiblityOfNextDiv(Translation["advancedDocumentSearch"]),
		getAdvancedSearch(extra.IsAdmin, extra),
		Group(nodes...),
		pagerFooter(view.BeforeUUID, view.NextUUID,
			before+extraStr, next+extraStr),
	)
}

func getAdvancedSearch(isAdmin bool, extra *extraction.DocumentQueryInfo) Node {
	return DIV(ID("advanced-search-div"),
		CLASS("text-left text-sm mt-2 w-auto mx-auto text-gray-200 font-bold hidden"),
		getFormStandardForm("form", PATCH, "/"+HTMXPreRouter+string(ViewDocument), CLASS("w-[800px]"),
			getInput("amount", "amount", strconv.FormatInt(int64(extra.Amount), 10),
				Translation["documentSearchAmount"], "number", "", "",
				MIN(strconv.FormatInt(MinDocuments, 10)), MAX(strconv.FormatInt(MaxDocuments, 10))),
			getSimpleTextInput("organisation", "organisation", extra.Organisation,
				Translation["documentSearchOrganisation"]),
			getSimpleTextInput("author", "author", extra.Author, Translation["documentSearchAuthor"]),
			getSimpleTextInput("title", "title", extra.Title, Translation["documentSearchTitle"]),
			If(isAdmin, getStandardCheckBox(extra.HideBlock, "true", "hideblock",
				Translation["documentSearchHideBlock"])),
			getStandardCheckBox(extra.Text, "true", "text", Translation["documentSearchText"]),
			getStandardCheckBox(extra.Discussion, "true", "discussion",
				Translation["documentSearchDiscussion"]),
			getStandardCheckBox(extra.Votes, "true", "votes", Translation["documentSearchVotes"]),
			getCheckBoxWithHideScript(false, "true", "addWritten",
				Translation["documentSearchAddWritten"], "writtenDiv"),
			getInput("written", "written", time.Now().Format("2006-01-02"),
				Translation["documentSearchWritten"], "date", "", "hidden"),
			getSubmitButton("make-advanced-search-query", Translation["makeAdvancedSearch"])),
	)
}

func GetExtraString(extra *extraction.DocumentQueryInfo) string {
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
	if extra.Organisation != "" {
		result += "&organisation=" + url.QueryEscape(extra.Organisation)
	}
	if extra.Author != "" {
		result += "&author=" + url.QueryEscape(extra.Author)
	}
	if extra.Title != "" {
		result += "&title=" + url.QueryEscape(extra.Author)
	}
	return result
}
