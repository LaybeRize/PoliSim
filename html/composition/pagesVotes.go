package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"fmt"
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
		items := []Node{}
		switch item.Info.VoteMethod {
		case database.SingleVote:
			items = getSingleVote(&item, doc.UUID)
		case database.MultipleVotes:
			items = getMultipleVote(&item, doc.UUID)
		case database.RankedVotes:
			items = getRankedVote(&item, doc.UUID)
		case database.ThreeCategoryVoting:
			items = getThreeCategoryVote(&item, doc.UUID)
		}
		votesDivs[i] = compactVoteView("vote-"+item.UUID,
			fmt.Sprintf(Translation["votingPartialHeader"], i+1), items)
	}

	return getBasePageWrapper(
		getPageHeader(ViewVoteDocument),
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
			P(I(CLASS("bi bi-calendar")), IfElse(doc.Type == database.FinishedVote,
				I(Text(" "+doc.Info.Finishing.Format(Translation["voteFinished"]))),
				I(Text(" "+doc.Info.Finishing.Format(Translation["voteRunning"]))))),
			If(doc.Private,
				P(I(CLASS("bi bi-person-fill-lock")), Text(" "+Translation["voteIsPrivate"]))),
			If(doc.AnyPosterAllowed,
				P(I(CLASS("bi bi-people-fill")), Text(" "+Translation["anyVoterAllowed"]))),
			If(doc.OrganisationPosterAllowed && !doc.AnyPosterAllowed,
				P(I(CLASS("bi bi-people-fill")), Text(" "+Translation["onlyOrganisationMemberAllowedToVote"]))),
			If(len(doc.Viewer) != 0 && doc.Private,
				P(Text(Translation["peopleAllowedToView"], reduceAccountsToString(doc.Viewer)))),
			If(len(doc.Poster) != 0 && !doc.AnyPosterAllowed,
				P(Text(Translation["peopleAllowedToComment"], reduceAccountsToString(doc.Poster)))),
		),
		Group(votesDivs...),
	)
}

func compactVoteView(id string, text string, children []Node) Node {
	return DIV(ID(id), TEST(id),
		DIV(CLASS("p-2.5 mt-3 w-[800px] flex items-center px-4 duration-300 cursor-pointer text-white hover:bg-blue-600"),
			HYPERSCRIPT("on click toggle .hidden on next <div/> from me then toggle .rotate-180 on last <span/> in first <div/> in me"),
			DIV(CLASS("flex justify-between items-center"),
				SPAN(CLASS("text-[15px] mr-4 text-gray-200 font-bold"), Text(text)),
				SPAN(CLASS("text-sm rotate-180"),
					I(CLASS("bi bi-chevron-down")),
				),
			),
		),
		DIV(ID(id+"-content"), Group(children...), CLASS("text-left text-sm mt-2 w-4/5 mx-auto text-gray-200 font-bold hidden")),
	)
}

func getSingleVote(item *database.Votes, docUUID string) []Node {
	return []Node{}
}
func getMultipleVote(item *database.Votes, docUUID string) []Node {
	return []Node{}
}
func getRankedVote(item *database.Votes, docUUID string) []Node {
	return []Node{}
}
func getThreeCategoryVote(item *database.Votes, docUUID string) []Node {
	return []Node{}
}
