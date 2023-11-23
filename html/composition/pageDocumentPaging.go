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
			getStandardBoxClass, STYLE("--clr-border: rgb(40 51 69);"),
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
	extraStr := GetExtraString(extra)
	return getBasePageWrapper(
		getPageHeader(ViewDocument),
		DIV(CLASS("p-2.5 mt-3 w-[800px] flex items-center px-4 duration-300 cursor-pointer text-white hover:bg-blue-600"),
			HYPERSCRIPT("on click toggle .hidden on next <div/> from me then toggle .rotate-180 on last <span/> in first <div/> in me"),
			DIV(CLASS("flex justify-between items-center"),
				SPAN(CLASS("text-[15px] mr-4 text-gray-200 font-bold"), Text(Translation["advancedDocumentSearch"])),
				SPAN(CLASS("text-sm rotate-180"),
					I(CLASS("bi bi-chevron-down")),
				),
			),
		),
		DIV(ID("advanced-search-div"), CLASS("text-left text-sm mt-2 w-auto mx-auto text-gray-200 font-bold hidden"),
			getFormStandardForm("form", PATCH, "/"+APIPreRoute+string(ViewDocument), CLASS("w-[800px]"),
				getInput("amount", "amount", strconv.FormatInt(int64(extra.Amount), 10), Translation["documentSearchAmount"], "number",
					"", "", MIN(strconv.FormatInt(MinDocuments, 10)), MAX(strconv.FormatInt(MaxDocuments, 10))),
				getSimpleTextInput("organisation", "organisation", extra.Organisation, Translation["documentSearchOrganisation"]),
				getSimpleTextInput("author", "author", extra.Author, Translation["documentSearchAuthor"]),
				getSimpleTextInput("title", "title", extra.Title, Translation["documentSearchTitle"]),
				If(isAdmin, getCheckBox("hideblock", extra.HideBlock, false, "true", "hideblock", Translation["documentSearchHideBlock"], nil)),
				getCheckBox("text", extra.Text, false, "true", "text", Translation["documentSearchText"], nil),
				getCheckBox("discussion", extra.Discussion, false, "true", "discussion", Translation["documentSearchDiscussion"], nil),
				getCheckBox("votes", extra.Votes, false, "true", "votes", Translation["documentSearchVotes"], nil),
				getCheckBox("addWritten", false, false, "true", "addWritten", Translation["documentSearchAddWritten"], nil),
				getInput("written", "written", time.Now().Format("2006-01-02"), Translation["documentSearchWritten"], "date", "", ""),
				getSubmitButton("make-advanced-search-query", Translation["makeAdvancedSearch"]))),
		Group(nodes...),
		pagerFooter(view.BeforeUUID, view.NextUUID,
			before+extraStr, next+extraStr),
	)
}

func GetExtraString(extra *extraction.ExtraInfo) string {
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
