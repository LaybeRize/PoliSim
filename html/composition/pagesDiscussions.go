package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
	"time"
)

func CreateDiscussionPage(acc *extraction.AccountAuth, document *validation.CreateDiscussion, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message = Translation["errorQueryingNames"] + "\n" + val.Message
	}
	node, err := UpdateUserOrganisations(acc, &extraction.AccountModification{ID: acc.ID,
		DisplayName: acc.DisplayName}, document.Organisation, "user")
	if err != nil {
		val.Message = Translation["errorRetrievingOrganisationForAccount"] + "\n" + val.Message
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getPageHeader(CreateDiscussionDocument),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateDiscussionDocument), CLASS("w-[800px]"),
			node,
			getSimpleTextInput("title", "title", document.Title, Translation["titleDiscussion"]),
			getSimpleTextInput("subtitle", "subtitle", document.Subtitle, Translation["subtitleDiscussion"]),
			getInput("endTime", "endTime", document.EndTime, Translation["endTimeDiscussion"], "datetime-local", "", ""),
			getTextArea("content", "content", document.Content, Translation["contentDiscussion"],
				MarkdownFormPage),
			getCheckBox("private", document.Private, false, "true", "private", Translation["privateDiscussion"],
				HYPERSCRIPT("on click toggle .hidden on #anyoneCanComment")),
			getCheckBox("anyoneCanComment", document.AnyoneCanParticipate, false, "true", "anyoneCanComment", Translation["anyoneCanCommentDiscussion"], nil),
			getCheckBox("membersCanComment", document.MembersCanParticipate, false, "true", "membersCanComment", Translation["membersCanCommentDiscussion"], nil),
			DIV(CLASS("flex flex-row"),
				getEditableList(document.Onlooker, "reader", "displayNames",
					Translation["addDiscussionReaderButton"], "w-[400px]"),
				getEditableList(document.Participants, "writer", "displayNames",
					Translation["addDiscussionWriterButton"], "w-[400px] ml-2"),
			),
			getSubmitButton("createDiscussionButton", Translation["createDiscussionButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}

const (
	CommentSingleDivID = "comment-div-id-%s"
	AdditionDiv        = "add-comment-div"
	AllCommentDiv      = "all-comments-in-one-div"
	CommentContentDiv  = "comment-content-div"
)

func GetCommentRendered(docUUID string, disc *database.Discussions, isAdmin bool) Node {
	id := ID(fmt.Sprintf(CommentSingleDivID, disc.UUID))
	if disc.Hidden && !isAdmin {
		return DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"), id,
			P(Text(Translation["commentHasBeenHidden"])),
		)
	}
	return DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"), id,
		If(disc.Hidden, P(CLASS("text-rose-600"), Text(Translation["commentCurrentlyHidden"]))),
		P(CLASS("mb-2"), I(Text(disc.Written.Format(Translation["commentWrittenAuthor"]), disc.Author)),
			If(disc.Flair != "", Group(I(Text("; ")), Text(disc.Flair)))),
		Raw(disc.HTMLContent),
		If(isAdmin, getCustomRequestClickable(HXPATCH, "/"+APIPreRoute+string(ChangeCommentDocumentLink)+
			url.PathEscape(docUUID)+"/"+url.PathEscape(disc.UUID), "", P(CLASS("bg-slate-700 text-white p-2 mt-2"),
			STYLE("text-align: center;"), IfElse(!disc.Hidden, Text(Translation["hideCommentDiscussion"]),
				Text(Translation["showCommentDiscussion"]))),
		)),
	)
}

func GetNewAdditionSSEDiv() Node {
	return DIV(ID(AdditionDiv))
}

func ViewDiscussionPage(acc *extraction.AccountAuth, uuidStr string, isAdmin bool, val validation.Message) Node {
	doc, err := extraction.GetDocumentForUser(uuidStr, acc.ID, isAdmin, database.FinishedDiscussion, database.RunningDiscussion)
	if err != nil {
		return GetErrorPage(Translation["documentDoesNotExists"])
	}

	if doc.Type == database.RunningDiscussion {
		go logic.CloseDiscussionIfTimeIsUp(doc.Info.Finishing, doc.UUID)
	}
	comments := make([]Node, len(doc.Info.Discussion))
	for i, disc := range doc.Info.Discussion {
		comments[i] = GetCommentRendered(doc.UUID, &disc, isAdmin)
	}

	return getBasePageWrapper(
		getPageHeader(ViewDiscussionDocument),
		getDocumentHead(doc, isAdmin),
		getDocumentBody(doc),
		DIV(CLASS("w-[800px] mt-2"),
			getTimeDiscussionInfo(doc, nil),
			If(doc.Private,
				P(I(CLASS("bi bi-person-fill-lock")), Text(" "+Translation["discussionIsPrivate"]))),
			If(doc.AnyPosterAllowed,
				P(I(CLASS("bi bi-people-fill")), Text(" "+Translation["anyPosterAllowed"]))),
			If(doc.OrganisationPosterAllowed && !doc.AnyPosterAllowed,
				P(I(CLASS("bi bi-people-fill")), Text(" "+Translation["onlyOrganisationMemberAllowed"]))),
			If(len(doc.Viewer) != 0 && doc.Private,
				P(Text(Translation["peopleAllowedToView"], reduceAccountsToString(doc.Viewer)))),
			If(len(doc.Poster) != 0 && !doc.AnyPosterAllowed,
				P(Text(Translation["peopleAllowedToComment"], reduceAccountsToString(doc.Poster)))),
		),
		DIV(ID(AllCommentDiv),
			Group(comments...),
			If(doc.Type == database.RunningDiscussion, GetNewAdditionSSEDiv()),
		),
		getDiscussionScript(doc.UUID),
		GetMessage(val),
		If(doc.Type == database.RunningDiscussion && acc.ID != 0,
			DIV(ID("discussion-comment-div"),
				getFormStandardForm("form", POST, "/"+APIPreRoute+string(CommentDiscussionLink)+url.PathEscape(doc.UUID), CLASS("mt-2 w-[800px]"),
					getUserDropdown(acc, "", Translation["discussionCommentAuthor"]),
					getTextArea(CommentContentDiv, "content", "", Translation["discussionCommentContent"],
						MarkdownFormPage),
					getSubmitButton("submitCommentButton", Translation["addCommentButton"])),
				getPreviewElement(),
			)),
		If(doc.Type == database.RunningDiscussion, scriptForUpdateOnEnd(doc, DiscussionUpdateDocumentLink)),
	)
}

func getDiscussionScript(uuid string) Node {
	return SCRIPT(Raw(`
const es = new EventSource("/` + APIPreRoute + string(sseReaderDiscussionLink) + uuid + `");
es.addEventListener("change", (event) => {
    const parsedData = JSON.parse(event.data);
    const el = document.getElementById(parsedData.id);
    el.outerHTML = parsedData.data;
	htmx.process(document.getElementById(parsedData.replace));
	if (parsedData.id !== parsedData.replace) {
        document.getElementById("` + CommentContentDiv + `").value = "";
    }
});`))
}

func GetDiscussionViewPageUpdate(acc *extraction.AccountAuth, uuidStr string, isAdmin bool) Node {
	doc, err := extraction.GetDocumentForUser(uuidStr, acc.ID, isAdmin, database.FinishedDiscussion, database.RunningDiscussion)
	if err != nil {
		return GetMessage(validation.Message{Message: Translation["documentDoesNotExists"]})
	}

	if doc.Type == database.RunningDiscussion {
		if doc.Info.Finishing.Before(time.Now()) {
			logic.CloseDiscussionIfTimeIsUp(doc.Info.Finishing, doc.UUID)
		} else {
			return GetMessage(validation.Message{Message: Translation["discussionIsStillRunning"]})
		}
	}
	doc.Type = database.FinishedDiscussion
	comments := make([]Node, len(doc.Info.Discussion))
	for i, disc := range doc.Info.Discussion {
		comments[i] = GetCommentRendered(doc.UUID, &disc, isAdmin)
	}
	return Group(GetMessage(validation.Message{Message: Translation["discussionClosedJustNow"], Positive: true}),
		DIV(ID("discussion-comment-div"), HXSWAPOOB("true")),
		DIV(ID(AllCommentDiv), HXSWAPOOB("true"), Group(comments...)),
		getTimeDiscussionInfo(doc, HXSWAPOOB("true")))
}

func getTimeDiscussionInfo(doc *database.Document, extra Node) Node {
	return P(ID("timer-p-element"), extra, I(CLASS("bi bi-calendar")), IfElse(doc.Type == database.FinishedDiscussion,
		I(Text(" "+doc.Info.Finishing.Format(Translation["discussionFinished"]))),
		I(Text(" "+doc.Info.Finishing.Format(Translation["discussionRunning"])))))
}

func reduceAccountsToString(accs []database.Account) string {
	if len(accs) == 0 {
		return ""
	}
	result := accs[0].DisplayName
	for _, acc := range accs[1:] {
		result += ", " + acc.DisplayName
	}
	return result
}
