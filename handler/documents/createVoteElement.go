package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"log/slog"
	"net/http"
)

var voteArray = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var referenceVote = database.VoteInstance{
	Question:              "",
	Answers:               []string{""},
	Type:                  database.SingleVote,
	MaxVotes:              100,
	ShowVotesDuringVoting: false,
	Anonymous:             true,
}

func GetCreateVoteElementPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}
	page := &handler.CreateVoteElementPage{VoteNumbers: voteArray, CurrNumber: voteArray[0]}
	var err error
	page.Vote, err = database.GetVote(acc, page.CurrNumber)
	if err != nil {
		page.IsError = true
		page.Message = loc.DocumentCouldNotLoadPersonalVote
	}

	if page.Vote == nil {
		page.Vote = &referenceVote
	} else {
		page.Vote.ConvertToAnswer()
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateVoteElementPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}
	page := &handler.CreateVoteElementPage{
		VoteNumbers: voteArray,
		CurrNumber:  values.GetInt("current-number"),
	}

	if page.CurrNumber < voteArray[0] || page.CurrNumber > voteArray[len(voteArray)-1] {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentInvalidVoteNumber})
		return
	}

	page.Vote = &database.VoteInstance{
		ID:                    helper.GetUniqueID(acc.Name),
		Question:              values.GetTrimmedString("question"),
		Answers:               values.GetFilteredArray("[]answers"),
		ShowVotesDuringVoting: values.GetBool("show-during"),
		Anonymous:             values.GetBool("anonymous"),
		Type:                  database.VoteType(values.GetInt("type")),
		MaxVotes:              values.GetInt("max-votes"),
	}

	if !page.Vote.HasValidType() {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentInvalidVoteType})
		return
	}

	if page.Vote.MaxVotes < 1 && page.Vote.Type == database.VoteShares {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentInvalidNumberMaxVotes})
		return
	}

	if len(page.Vote.Answers) < 1 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentAmountAnswersTooSmall})
		return
	}

	if page.Vote.Question == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentVoteMustHaveAQuestion})
		return
	}

	const maxQuestionLength = 1000
	if len([]rune(page.Vote.Question)) > maxQuestionLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.DocumentVoteQuestionTooLong, maxQuestionLength)})
		return
	}

	err = database.CreateOrUpdateVote(page.Vote, acc, page.CurrNumber)
	if err != nil {
		slog.Error(err.Error())
		page.IsError = true
		page.Message = loc.DocumentErrorSavingUserVote
	} else {
		page.IsError = false
		page.Message = loc.DocumentSuccessfullySavedUserVote
	}

	handler.MakePage(writer, acc, page)
}

func PatchGetVoteElementPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.CreateVoteElementPage{
		VoteNumbers: voteArray,
		CurrNumber:  values.GetInt("number"),
	}

	if page.CurrNumber < voteArray[0] || page.CurrNumber > voteArray[len(voteArray)-1] {
		page.CurrNumber = voteArray[0]
	}

	page.Vote, err = database.GetVote(acc, page.CurrNumber)
	if err != nil {
		page.IsError = true
		page.Message = loc.DocumentCouldNotLoadPersonalVote
	}

	if page.Vote == nil {
		page.Vote = &referenceVote
	} else {
		page.Vote.ConvertToAnswer()
	}

	handler.MakePage(writer, acc, page)
}
