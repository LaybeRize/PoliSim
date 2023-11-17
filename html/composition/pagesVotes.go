package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"strconv"
)

const voteContainerDiv = "vote-container-div"

func GetCreateVotePage(acc *extraction.AccountAuth, document *validation.CreateVote, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message = Translation["errorQueryingNames"] + "\n" + val.Message
	}
	node, err := UpdateUserOrganisations(acc, &extraction.AccountModification{ID: acc.ID,
		DisplayName: acc.DisplayName}, "", "user")
	if err != nil {
		val.Message = Translation["errorRetrievingOrganisationForAccount"] + "\n" + val.Message
	}

	return getBasePageWrapper(
		getDataList("displayNames", display),
		SCRIPT(SRC("/public/json-enc-custom.js")),
		getPageHeader(CreateVoteDocument),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateVoteDocument), CLASS("w-[800px]"),
			HXEXTEND("json-enc-custom"),
			node,
			getSimpleTextInput("title", "title", "", Translation["titleVoteDocument"]),
			getSimpleTextInput("subtitle", "subtitle", "", Translation["subtitleVoteDocument"]),
			getInput("endTime", "endTime", document.EndTime, Translation["endTimeVote"], "datetime-local", "", ""),
			getTextArea("content", "content", "", Translation["contentVoteDocument"],
				MarkdownJsonPage),
			getCheckBox("private", false, false, "", "private", Translation["privateVote"],
				HYPERSCRIPT("on click toggle .hidden on #anyoneCanComment")),
			getCheckBox("anyoneCanVote", false, false, "", "anyoneCanVote", Translation["anyoneCanVote"], nil),
			getCheckBox("membersCanVote", false, false, "", "membersCanVote", Translation["membersCanVote"], nil),
			getHiddenEmptyInput("attendents"), getHiddenEmptyInput("voter"),
			DIV(CLASS("flex flex-row"),
				getEditableList([]string{}, "attendents", "displayNames",
					Translation["addVoteAttendentButton"], "w-[400px]"),
				getEditableList([]string{}, "voter", "displayNames",
					Translation["addVoteVoterButton"], "w-[400px] ml-2"),
			),
			DIV(ID(voteContainerDiv),
				getPartialVote("1")),
			getPartialButton("2", false), BR(),
			getSubmitButton("createVoteButton", Translation["createVoteDocument"]),
		),
		GetMessage(val),
		getPreviewElement(),
	)
}

func getHiddenEmptyInput(name string) Node {
	return INPUT(TYPE("text"), NAME(name), VALUE(""), HIDDEN())
}

const votePartialButtonID = "vote-partial-button"

func GetVotePartial(partialNumber int64) Node {
	partial := strconv.FormatInt(partialNumber, 10)
	return Group(getPartialVote(partial),
		getPartialButton(strconv.FormatInt(partialNumber+1, 10), true),
	)
}

func getPartialVote(number string) Node {
	return DIV(CLASS("w-[800px] box box-e p-2 mt-2"), ID("vote-partial-"+number),
		STYLE("--clr-border: rgb(40 51 69);"),
		DIV(CLASS("w-full flex justify-between flex-row"),
			P(CLASS("text-xl"), Text(Translation["votePartialHeader"], number)),
			BUTTON(CLASS("bg-slate-700 text-white p-2 mt-2 hover:bg-rose-800"), TYPE("button"),
				ID("deleteVoteButton"), TEST("deleteVoteButton-"+number),
				HYPERSCRIPT("on click tell #vote-partial-"+number+" remove yourself"), Text(Translation["deleteVote"])),
		),
		getSimpleTextInput("question-questionText-"+number, "question["+number+"][questionText]",
			"", Translation["voteQuestionText"]),
		getDropDown("question["+number+"][type]", "question-type-"+
			number, Translation["voteQuestionType"], false,
			database.VoteTypes, database.VoteTranslation, database.SingleVote),
		getCheckBox("viewCountsWhileRunning-"+number, false, false, "",
			"question["+number+"][viewCountsWhileRunning]", Translation["viewCountsWhileRunning"], nil),
		getCheckBox("viewNamesWhileRunning-"+number, false, false, "",
			"question["+number+"][viewNamesWhileRunning]", Translation["viewNamesWhileRunning"],
			HYPERSCRIPT("on click toggle .hidden on #viewNamesAfterFinished-"+number)),
		getCheckBox("viewNamesAfterFinished-"+number, false, false, "",
			"question["+number+"][viewNamesAfterFinished]", Translation["viewNamesAfterFinished"], nil),
		getHiddenEmptyInput("question["+number+"][answers]"),
		getEditableList([]string{}, "question["+number+"][answers]", "",
			Translation["voteAddAnswersToQuestion"], "w-full"),
	)
}

func getPartialButton(number string, withSwap bool) Node {
	return BUTTON(If(withSwap, HXSWAPOOB("true")), ID(votePartialButtonID), TYPE("button"),
		HXTARGET("#"+voteContainerDiv), HXSWAP("beforeend"), HXPATCH("/"+APIPreRoute+string(requestVotePartialLink)+number),
		P(CLASS("bg-slate-700 text-white p-2 mt-2"), STYLE("text-align: center;"),
			Text(Translation["addnewVoteToPost"]),
		))
}

func GetVoteViewPage(acc *extraction.AccountAuth, uuidStr string, isAdmin bool, val validation.Message) Node {
	doc, err := extraction.GetDocumentForUser(uuidStr, acc.ID, isAdmin, database.FinishedVote, database.RunningVote)
	if err != nil {
		return GetErrorPage(Translation["documentDoesNotExists"])
	}

	if doc.Type == database.RunningVote {
		go logic.CloseVoteIfTimeIsUp(doc.Info.Finishing, doc.UUID)
	}
	votes, err := extraction.GetVotesForDocument(doc.UUID)
	if err != nil {
		return GetErrorPage(Translation["errorLoadingVotesForDocument"])
	}
	votesDivs := make([]Node, len(votes))
	for i, item := range votes {
		switch item.Info.VoteMethod {
		case database.SingleVote:
		case database.MultipleVotes:
		case database.RankedVotes:
		case database.ThreeCategoryVoting:
		}
		votesDivs[i] = DIV(Text(item.UUID))
	}

	return getBasePageWrapper(
		Group(votesDivs...),
	)
}
