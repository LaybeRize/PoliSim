package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
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
		getPageHeader(CreateTextDocument),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateTextDocument), CLASS("w-[800px]"),
			node,
			getSimpleTextInput("title", "title", document.Title, Translation["titleDiscussion"]),
			getSimpleTextInput("subtitle", "subtitle", document.Subtitle, Translation["subtitleDiscussion"]),
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

func ViewDiscussionPage(acc *extraction.AccountAuth, uuidStr string) Node {
	return getBasePageWrapper(
		Text(uuidStr),
		Text(acc.DisplayName),
	)
}
