package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
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
		page.Message = "Konnte die ausgewählte Abstimmung nicht laden"
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

	err := request.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen"})
		return
	}
	page := &handler.CreateVoteElementPage{VoteNumbers: voteArray}
	database.GetIntegerFormEntry(request, "number", &page.CurrNumber)
	if page.CurrNumber < voteArray[0] || page.CurrNumber > voteArray[len(voteArray)-1] {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die ausgewählte Nummer für die Abstimmung ist nicht zulässig"})
		return
	}

	page.Vote = &database.VoteInstance{
		ID:                    helper.GetUniqueID(acc.Name),
		Question:              helper.GetFormEntry(request, "question"),
		Answers:               helper.GetFormList(request, "[]answers"),
		ShowVotesDuringVoting: helper.GetBoolFormEntry(request, "show-during"),
		Anonymous:             helper.GetBoolFormEntry(request, "anonymous"),
	}
	database.GetIntegerFormEntry(request, "type", &page.Vote.Type)
	database.GetIntegerFormEntry(request, "max-votes", &page.Vote.MaxVotes)

	if !page.Vote.HasValidType() {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der ausgewählte Abstimmungstyp für die Abstimmung ist nicht zulässig"})
		return
	}

	if page.Vote.MaxVotes < 1 && page.Vote.Type == database.VoteShares {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die maximale Stimmenzahl pro Nutzer darf nicht kleiner als 1 sein für den ausgewählten Abstimmungstypen"})
		return
	}

	if len(page.Vote.Answers) < 1 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es muss mindestens eine Antwort zur Abstimmung stehen"})
		return
	}

	if page.Vote.Question == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die Abstimmung muss eine Frage haben, über die abgestimmt wird"})
		return
	}

	if len(page.Vote.Question) > 1000 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die Abstimmungsfrage darf nicht länger als 1000 Zeichen sein"})
		return
	}

	err = database.CreateOrUpdateVote(page.Vote, acc, page.CurrNumber)
	if err != nil {
		slog.Error(err.Error())
		page.IsError = true
		page.Message = "Es ist ein Fehler beim speichern der Abstimmung aufgetreten"
	} else {
		page.IsError = false
		page.Message = "Abstimmung erfolgreich gespeichert"
	}

	handler.MakePage(writer, acc, page)
}

func PatchGetVoteElementPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	err := request.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen"})
		return
	}

	page := &handler.CreateVoteElementPage{VoteNumbers: voteArray}
	database.GetIntegerFormEntry(request, "number", &page.CurrNumber)
	if page.CurrNumber < voteArray[0] || page.CurrNumber > voteArray[len(voteArray)-1] {
		page.CurrNumber = voteArray[0]
	}

	page.Vote, err = database.GetVote(acc, page.CurrNumber)
	if err != nil {
		page.IsError = true
		page.Message = "Konnte die ausgewählte Abstimmung nicht laden"
	}

	if page.Vote == nil {
		page.Vote = &referenceVote
	} else {
		page.Vote.ConvertToAnswer()
	}

	handler.MakePage(writer, acc, page)
}
