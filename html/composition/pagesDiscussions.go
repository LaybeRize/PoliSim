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
			getInput("subtitle", "endTime", document.EndTime, Translation["endTimeDiscussion"], "datetime-local", "", ""),
			getTextArea("content", "content", document.Content, Translation["contentDiscussion"], true),
			getCheckBox("private", document.Private, false, "true", "private", Translation["privateDiscussion"],
				HYPERSCRIPT("on click toggle .hidden on #anyoneCanComment")),
			getCheckBox("anyoneCanComment", document.AnyoneCanComment, false, "true", "anyoneCanComment", Translation["anyoneCanCommentDiscussion"], nil),
			getCheckBox("membersCanComment", document.MembersCanComment, false, "true", "membersCanComment", Translation["membersCanCommentDiscussion"], nil),
			DIV(CLASS("flex flex-row"),
				getEditableList(document.Reader, "reader", "displayNames",
					Translation["addDiscussionReaderButton"], "w-[400px]"),
				getEditableList(document.Writer, "writer", "displayNames",
					Translation["addDiscussionWriterButton"], "w-[400px] ml-2"),
			),
			getSubmitButton(Translation["createDiscussionButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
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
		if disc.Hidden && !isAdmin {
			comments[i] = DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
				P(Text(Translation["commentHasBeenHidden"])),
			)
			continue
		}
		comments[i] = DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
			If(disc.Hidden, P(CLASS("text-rose-600"), Text(Translation["commentCurrentlyHidden"]))),
			P(CLASS("mb-2"), I(Text(disc.Written.Format(Translation["commentWrittenAuthor"]), disc.Author)),
				If(disc.Flair != "", Group(I(Text("; ")), Text(disc.Flair)))),
			Raw(disc.HTMLContent),
			If(isAdmin, getCustomRequestClickable(HXPATCH, "/"+APIPreRoute+string(ChangeCommentDocumentLink)+
				url.PathEscape(doc.UUID)+"/"+url.PathEscape(disc.UUID), "", P(CLASS("bg-slate-700 text-white p-2 mt-2"),
				STYLE("text-align: center;"), IfElse(!disc.Hidden, Text(Translation["hideCommentDiscussion"]),
					Text(Translation["showCommentDiscussion"]))),
			)),
		)
	}

	return getBasePageWrapper(
		getPageHeader(ViewDiscussionDocument),
		DIV(CLASS("w-[800px]"),
			H1(CLASS("text-3xl underline decoration-2 underline-offset-2"), Text(doc.Title)),
			If(doc.Subtitle.Valid, H1(CLASS("text-2xl"), Text(doc.Subtitle.String))),
			P(CLASS("my-2"), I(Text(fmt.Sprintf(doc.Written.Format(Translation["authorDiscussionDocument"]), doc.Organisation, doc.Author))),
				If(doc.Flair != "", Group(I(Text("; ")), Text(doc.Flair)))),
		),
		DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
			Raw(doc.HTMLContent),
		),
		DIV(CLASS("w-[800px] mt-2"),
			IfElse(doc.Type == database.FinishedDiscussion,
				P(I(Text(doc.Info.Finishing.Format(Translation["discussionFinished"])))),
				P(I(Text(doc.Info.Finishing.Format(Translation["discussionRunning"]))))),
			If(doc.Private,
				P(Text(Translation["discussionIsPrivate"]))),
			If(len(doc.Viewer) != 0 && doc.Private,
				P(Text(Translation["peopleAllowedToView"], reduceAccountsToString(doc.Viewer)))),
			If(len(doc.Poster) != 0 && !doc.AnyPosterAllowed,
				P(Text(Translation["peopleAllowedToComment"], reduceAccountsToString(doc.Poster)))),
			If(doc.AnyPosterAllowed,
				P(Text(Translation["anyPosterAllowed"]))),
			If(doc.OrganisationPosterAllowed && !doc.AnyPosterAllowed,
				P(Text(Translation["onlyOrganisationMemberAllowed"]))),
		),
		Group(comments...),
		If(doc.Type == database.RunningDiscussion && acc.ID != 0, Group(
			getFormStandardForm("form", POST, "/"+APIPreRoute+string(CommentDiscussionLink)+url.PathEscape(doc.UUID), CLASS("mt-2 w-[800px]"),
				getUserDropdown(acc, "", Translation["discussionCommentAuthor"]),
				getTextArea("content", "content", "", Translation["discussionCommentContent"], true),
				getSubmitButton(Translation["addCommentButton"])),
			GetMessage(val),
			getPreviewElement(),
		)),
	)
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
