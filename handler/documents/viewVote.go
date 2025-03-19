package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

func GetVoteView(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)

	page := &handler.ViewVotePage{}
	var err error
	page.VoteInstance, page.VoteResults, err = database.GetVoteForUser(request.PathValue("id"), acc)

	if err != nil {
		handler.GetNotFoundPage(writer, request)
	}

	if acc.Exists() {
		page.Voter, err = database.GetOwnedAccountNames(acc)
		page.Voter = append([]string{acc.Name}, page.Voter...)
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostVote(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentGeneralFunctionNotAvailable})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	voter := values.GetTrimmedString("voter")
	allowed, err := database.IsAccountAllowedToPostWith(acc, voter)
	if !allowed || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentNotAllowedToVoteWithThatAccount})
		return
	}
	id := request.PathValue("id")
	answers, voteType, maxVotes, err := database.GetAnswersAndTypeForVote(id, acc)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentNotAllowedToVoteOnThis})
		return
	}

	var votesCasted []int
	if values.GetBool("invalid") {
		votesCasted = nil

	} else if voteType.IsSingleVote() {
		votesCasted = make([]int, len(answers))
		pos := values.GetInt("vote")

		if pos <= 0 || pos > len(answers) {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: loc.DocumentVoteIsInvalid + "\n" + loc.DocumentVotePositionInvalid})
			return
		}
		votesCasted[pos-1] = 1

	} else if voteType.IsMultipleVotes() {
		votesCasted = make([]int, len(answers))
		for i := range len(answers) {
			if values.GetBool(fmt.Sprintf("vote-%d", i+1)) {
				votesCasted[i] = 1
			}
		}

	} else if voteType.IsVoteSharing() {
		votesCasted = make([]int, len(answers))
		sum := 0

		for i := range len(answers) {
			amount := values.GetInt(fmt.Sprintf("vote-%d", i+1))
			if amount < 0 {
				handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
					Message: loc.DocumentVoteIsInvalid + "\n" + loc.DocumentVoteShareNotSmallerZero})
				return
			}
			votesCasted[i] = amount
			sum += amount
		}

		if sum > maxVotes {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: loc.DocumentVoteIsInvalid + "\n" + loc.DocumentVoteSumTooBig})
			return
		}

	} else if voteType.IsRankedVoting() {
		votesCasted = make([]int, len(answers))
		lookUpMap := make(map[int]interface{})

		for i := range len(answers) {
			pos := values.GetInt(fmt.Sprintf("vote-%d", i+1))
			if pos <= 0 {
				votesCasted[i] = -1
			} else if pos > len(answers) {
				handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
					Message: loc.DocumentVoteIsInvalid + "\n" + loc.DocumentVoteRankTooBig})
				return
			} else {
				_, exists := lookUpMap[pos]
				if exists {
					handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
						Message: loc.DocumentVoteIsInvalid + "\n" + loc.DocumentVoteInvalidDoubleRank})
					return
				}
				lookUpMap[pos] = struct{}{}
				votesCasted[i] = pos
			}
		}

	}

	err = database.CastVoteWithAccount(voter, id, votesCasted)
	if errors.Is(err, database.AlreadyVoted) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentAlreadyVotedWithThatAccount})
		return
	} else if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentErrorWhileVoting})
		return
	}

	page := &handler.ViewVotePage{}
	page.VoteInstance, page.VoteResults, err = database.GetVoteForUser(request.PathValue("id"), acc)
	if err != nil {
		handler.PartialGetNotFoundPage(writer, request)
	}

	page.IsError = false
	page.Message = loc.DocumentSuccessfullyVoted

	if acc.Exists() {
		page.Voter, err = database.GetOwnedAccountNames(acc)
		page.Voter = append([]string{acc.Name}, page.Voter...)
	}

	handler.MakePage(writer, acc, page)
}
