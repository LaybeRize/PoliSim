package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"fmt"
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

func ViewDiscussionPage(acc *extraction.AccountAuth, uuidStr string, isAdmin bool) Node {
	doc, err := extraction.GetDocumentForUser(uuidStr, acc.ID, isAdmin, database.FinishedDiscussion, database.RunningDiscussion)
	if err != nil {
		return GetErrorPage(Translation["documentDoesNotExists"])
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
			If(len(doc.Viewer) != 0 && doc.Private,
				P(Text(Translation["peopleAllowedToView"], reduceAccountsToString(doc.Viewer)))),
			If(len(doc.Poster) != 0 && !doc.AnyPosterAllowed,
				P(Text(Translation["peopleAllowedToComment"], reduceAccountsToString(doc.Poster)))),
			If(doc.AnyPosterAllowed,
				P(Text(Translation["anyPosterAllowed"]))),
			If(doc.OrganisationPosterAllowed && !doc.AnyPosterAllowed,
				P(Text(Translation["onlyOrganisationMemberAllowed"]))),
		),
	)
}

func reduceAccountsToString(accs []database.Account) string {
	result := accs[0].DisplayName
	for _, acc := range accs[1:] {
		result += ", " + acc.DisplayName
	}
	return result
}
