package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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
		var items []Node
		switch item.Info.VoteMethod {
		case database.SingleVote:
			items = getSingleVote(acc, &item, doc.UUID)
		case database.MultipleVotes:
			items = getMultipleVote(acc, &item, doc.UUID)
		case database.RankedVotes:
			items = getRankedVote(acc, &item, doc.UUID)
		case database.ThreeCategoryVoting:
			items = getThreeCategoryVote(acc, &item, doc.UUID)
		}
		votesDivs[i] = compactVoteView("vote-"+item.UUID,
			fmt.Sprintf(Translation["votingPartialHeader"], i+1), items)
	}

	return getBasePageWrapper(
		SCRIPT(SRC("/public/json-enc-custom.js")),
		getPageHeader(ViewVoteDocument),
		getDocumentHead(doc, isAdmin),
		getDocumentBody(doc),
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
		GetMessage(val),
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
		DIV(ID(id+"-content"), Group(children...), CLASS("text-left text-sm mt-2 w-auto mx-auto text-gray-200 font-bold hidden")),
	)
}

func getSingleVote(acc *extraction.AccountAuth, item *database.Votes, docUUID string) []Node {
	RadioButtons := make([]Node, len(item.Info.Options))
	for i, str := range item.Info.Options {
		RadioButtons[i] = getRadioButton("vote-option-"+item.UUID+"-"+strconv.FormatInt(int64(i), 10),
			i == 0, false, strconv.FormatInt(int64(i), 10), "answerSingle", str, CONVERTTO("number"))
	}
	return []Node{P(CLASS("text-xl my-2"), Text(item.Question)),
		If(acc.Role != database.NotLoggedIn && !item.Finished,
			getFormStandardForm("form", PATCH, "/"+APIPreRoute+string(MakeVoteLink)+url.PathEscape(docUUID)+"/"+
				url.PathEscape(item.UUID)+"/"+url.PathEscape(string(database.SingleVote)), CLASS("w-[800px]"),
				HXEXTEND("json-enc-custom"),
				getUserDropdown(acc, "", Translation["giveVoteAuthorText"]),
				getCheckBox("invalidateVote", false, false, "", "invalidateVote", Translation["invalidateVoteCheckbox"],
					HYPERSCRIPT("on click toggle .hidden on #vote-card-div-"+item.UUID)),
				DIV(ID("vote-card-div-"+item.UUID), CLASS("w-full mt-4"),
					Group(RadioButtons...)),
				getSubmitButton("sendVote-"+item.UUID, Translation["sendVote"]),
			)),
		GetInfoStandardView(item, false),
	}
}
func getMultipleVote(acc *extraction.AccountAuth, item *database.Votes, docUUID string) []Node {
	checkBoxes := make([]Node, len(item.Info.Options))
	for i, str := range item.Info.Options {
		checkBoxes[i] = getCheckBox("vote-option-"+item.UUID+"-"+strconv.FormatInt(int64(i), 10),
			false, false, "", "answerMultiple["+strconv.FormatInt(int64(i), 10)+"]", str, nil)
	}
	return []Node{P(CLASS("text-xl my-2"), Text(item.Question)),
		If(acc.Role != database.NotLoggedIn && !item.Finished,
			getFormStandardForm("form", PATCH, "/"+APIPreRoute+string(MakeVoteLink)+url.PathEscape(docUUID)+"/"+
				url.PathEscape(item.UUID)+"/"+url.PathEscape(string(database.MultipleVotes)), CLASS("w-[800px]"),
				HXEXTEND("json-enc-custom"),
				getUserDropdown(acc, "", Translation["giveVoteAuthorText"]),
				getCheckBox("invalidateVote", false, false, "", "invalidateVote", Translation["invalidateVoteCheckbox"],
					HYPERSCRIPT("on click toggle .hidden on #vote-card-div-"+item.UUID)),
				DIV(ID("vote-card-div-"+item.UUID), CLASS("w-full mt-4"),
					Group(checkBoxes...),
				),
				getSubmitButton("sendVote-"+item.UUID, Translation["sendVote"]),
			)),
		GetInfoStandardView(item, false),
	}
}

func getRankedVote(acc *extraction.AccountAuth, item *database.Votes, docUUID string) []Node {
	rankings := make([]Node, len(item.Info.Options))
	for i, str := range item.Info.Options {
		rankings[i] = getInput("vote-option-"+item.UUID+"-"+strconv.FormatInt(int64(i), 10),
			"answerRanked["+strconv.FormatInt(int64(i), 10)+"]", "0", str, "number",
			"", "", MIN("0"), MAX(strconv.FormatInt(int64(len(item.Info.Options)), 10)),
			CONVERTTO("number"))
	}
	return []Node{P(CLASS("text-xl my-2"), Text(item.Question)),
		If(acc.Role != database.NotLoggedIn && !item.Finished,
			getFormStandardForm("form", PATCH, "/"+APIPreRoute+string(MakeVoteLink)+url.PathEscape(docUUID)+"/"+
				url.PathEscape(item.UUID)+"/"+url.PathEscape(string(database.RankedVotes)), CLASS("w-[800px]"),
				HXEXTEND("json-enc-custom"),
				getUserDropdown(acc, "", Translation["giveVoteAuthorText"]),
				getCheckBox("invalidateVote", false, false, "", "invalidateVote", Translation["invalidateVoteCheckbox"],
					HYPERSCRIPT("on click toggle .hidden on #vote-card-div-"+item.UUID)),
				DIV(ID("vote-card-div-"+item.UUID), CLASS("w-full mt-4"),
					Group(rankings...)),
				getSubmitButton("sendVote-"+item.UUID, Translation["sendVote"]),
			)),
		GetInfoRankedView(item, false),
	}
}

func getThreeCategoryVote(acc *extraction.AccountAuth, item *database.Votes, docUUID string) []Node {
	threeCategory := make([]Node, len(item.Info.Options))
	for i, str := range item.Info.Options {
		threeCategory[i] = Group(P(CLASS("mt-2"), Text(str)),
			DIV(CLASS("w-full grid grid-flow-col grid-cols-3 justify-stretch"),
				getRadioButton("vote-option-"+item.UUID+"-"+strconv.FormatInt(int64(i), 10)+"-for",
					false, false, "1", "answerThree["+strconv.FormatInt(int64(i), 10)+"]",
					Translation["voteFor"], CONVERTTO("number")),
				getRadioButton("vote-option-"+item.UUID+"-"+strconv.FormatInt(int64(i), 10)+"-neutral",
					true, false, "0", "answerThree["+strconv.FormatInt(int64(i), 10)+"]",
					Translation["voteNeutral"], CONVERTTO("number")),
				getRadioButton("vote-option-"+item.UUID+"-"+strconv.FormatInt(int64(i), 10)+"-against",
					false, false, "-1", "answerThree["+strconv.FormatInt(int64(i), 10)+"]",
					Translation["voteAgainst"], CONVERTTO("number")),
			))
	}
	return []Node{P(CLASS("text-xl my-2"), Text(item.Question)),
		If(acc.Role != database.NotLoggedIn && !item.Finished,
			getFormStandardForm("form", PATCH, "/"+APIPreRoute+string(MakeVoteLink)+url.PathEscape(docUUID)+"/"+
				url.PathEscape(item.UUID)+"/"+url.PathEscape(string(database.ThreeCategoryVoting)), CLASS("w-[800px]"),
				HXEXTEND("json-enc-custom"),
				getUserDropdown(acc, "", Translation["giveVoteAuthorText"]),
				getCheckBox("invalidateVote", false, false, "", "invalidateVote", Translation["invalidateVoteCheckbox"],
					HYPERSCRIPT("on click toggle .hidden on #vote-card-div-"+item.UUID)),
				DIV(CLASS("my-2")),
				DIV(ID("vote-card-div-"+item.UUID), CLASS("w-full mt-4"),
					Group(threeCategory...)),
				getSubmitButton("sendVote-"+item.UUID, Translation["sendVote"]),
			)),
		GetInfoStandardView(item, false),
	}
}

func GetInfoStandardView(item *database.Votes, oobSwap bool) Node {
	children := make([]Node, len(item.Info.Options)+1)
	for i, str := range item.Info.Options {
		children[i] = TR(
			getTableElement(StartPos, 1, Text(str)),
			getTableElement(EndPos, 1, Text(strconv.FormatInt(item.Info.Summary.Sums[str], 10))),
		)
	}
	children[len(children)-1] = TR(
		getTableElement(StartPos, 1, Text(Translation["validVoteVotes"])),
		getTableElement(EndPos, 1, Text(strconv.FormatInt(int64(len(item.Info.Summary.InvalidVotes)), 10))),
	)
	moreInfo := []Node{nil}
	generateMoreInfo := (item.ShowNamesAfterVoting && item.Finished) || (item.ShowNamesWhileVoting && !item.Finished)
	if generateMoreInfo && len(item.Info.VoteOrder) != 0 {
		moreInfo = make([]Node, len(item.Info.VoteOrder)+1)

		row := make([]Node, len(item.Info.Options)+1)
		row[0] = getTableHeader(StartPos, -1, Translation["questionExtensionName"])
		for i, str := range item.Info.Options {
			if i == len(item.Info.Options)-1 {
				row[i+1] = getTableHeader(EndPos, -1, str)
				continue
			}
			row[i+1] = getTableHeader(MiddlePos, -1, str)
		}
		moreInfo[0] = TR(Group(clone(row)...))
		for j, str := range item.Info.VoteOrder {
			row[0] = getTableElement(StartPos, 1, Text(str))
			for i, strOpt := range item.Info.Options {
				if i == len(item.Info.Options)-1 {
					row[i+1] = getTableElement(EndPos, 1, Text(strconv.FormatInt(item.Info.Results[str].Votes[strOpt], 10)))
					continue
				}
				row[i+1] = getTableElement(MiddlePos, 1, Text(strconv.FormatInt(item.Info.Results[str].Votes[strOpt], 10)))
			}
			moreInfo[j+1] = TR(Group(clone(row)...))
		}
	}
	return DIV(ID("vote-info-for-"+item.UUID), If(oobSwap, HXSWAPOOB("true")), CLASS("w-full"),
		If((item.ShowNumbersWhileVoting && !item.Finished) || item.Finished,
			TABLE(ID("vote-info-table-summary-"+item.UUID), CLASS("table-auto mt-4"),
				TR(
					getTableHeader(StartPos, -1, Translation["questionColumnSummary"]),
					getTableHeader(EndPos, -1, Translation["numbersColumnSummary"]),
				),
				Group(children...))),
		If(generateMoreInfo && len(item.Info.Summary.InvalidVotes) != 0, P(Text(Translation["invalidateVoteNames"], strings.Join(item.Info.Summary.InvalidVotes, ", ")))),
		If(generateMoreInfo && len(item.Info.VoteOrder) != 0, TABLE(ID("vote-info-table-summary-"+item.UUID), CLASS("table-auto mt-4"),
			Group(moreInfo...))))
}

func GetInfoRankedView(item *database.Votes, oobSwap bool) Node {
	moreInfo := []Node{nil}
	generateMoreInfo := item.Finished || item.ShowNumbersWhileVoting
	if generateMoreInfo && len(item.Info.VoteOrder) != 0 {
		showName := (item.Finished && item.ShowNamesAfterVoting) || (!item.Finished && item.ShowNamesWhileVoting)
		moreInfo = make([]Node, len(item.Info.VoteOrder)+1)

		row := make([]Node, len(item.Info.Options)+1)
		row[0] = getTableHeader(StartPos, -1, Translation["questionExtensionName"])
		for i, str := range item.Info.Options {
			if i == len(item.Info.Options)-1 {
				row[i+1] = getTableHeader(EndPos, -1, str)
				continue
			}
			row[i+1] = getTableHeader(MiddlePos, -1, str)
		}
		moreInfo[0] = TR(Group(clone(row)...))
		for j, str := range item.Info.VoteOrder {
			voter := str
			if !showName {
				voter = fmt.Sprintf(Translation["replaceNameWithNumber"], j+1)
			}
			row[0] = getTableElement(StartPos, 1, Text(voter))
			for i, strOpt := range item.Info.Options {
				if i == len(item.Info.Options)-1 {
					row[i+1] = getTableElement(EndPos, 1, Text(strconv.FormatInt(item.Info.Results[str].Votes[strOpt], 10)))
					continue
				}
				row[i+1] = getTableElement(MiddlePos, 1, Text(strconv.FormatInt(item.Info.Results[str].Votes[strOpt], 10)))
			}
			moreInfo[j+1] = TR(Group(clone(row)...))
		}
	}

	return DIV(ID("vote-info-for-"+item.UUID), If(oobSwap, HXSWAPOOB("true")), CLASS("w-full"),
		If((item.ShowNumbersWhileVoting && !item.Finished) || item.Finished, Group(
			P(CLASS("mt-4"), Text(Translation["invalidateVoteRankedText"], len(item.Info.Summary.InvalidVotes))),
			If(((item.ShowNamesAfterVoting && item.Finished) || (item.ShowNamesWhileVoting && !item.Finished)) && len(item.Info.Summary.InvalidVotes) != 0,
				P(Text(Translation["invalidateVoteNames"], strings.Join(item.Info.Summary.InvalidVotes, ", ")))),
			If(generateMoreInfo && len(item.Info.VoteOrder) != 0, TABLE(ID("vote-info-table-summary-"+item.UUID), CLASS("table-auto mt-4"),
				Group(moreInfo...))),
		)))
}

func clone(n []Node) []Node {
	newNodes := make([]Node, len(n))
	copy(newNodes, n)
	return newNodes
}
