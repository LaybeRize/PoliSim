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
	"strings"
	"time"
)

func GetCreateVotePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	year, month, day := time.Now().UTC().Date()
	page := &handler.CreateVotePage{
		DateTime: time.Date(year, month, day+runMinDays+1, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
		MinTime:  time.Date(year, month, day+runMinDays, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
		MaxTime:  time.Date(year, month, day+runMaxDays, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
	}
	page.Reader = []string{""}
	page.Participants = []string{""}
	page.IsError = true
	page.Message = ""

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		slog.Debug(err.Error())
		page.Message = loc.CouldNotFindAllAuthors
		arr = make([]string, 0)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = acc.Name
	page.PossibleAuthors = arr
	page.PossibleOrganisations, err = database.GetOrganisationNamesAdminIn(acc.Name)
	if err != nil {
		slog.Debug(err.Error())
		page.Message = "\n" + loc.ErrorFindingAllOrganisationsForAccount
	}
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		slog.Debug(err.Error())
		page.Message += "\n" + loc.ErrorSearchingForAccountNames
	}
	page.VoteChoice, err = database.GetVoteInfoList(acc)
	if err != nil {
		slog.Debug(err.Error())
		page.Message += "\n" + loc.DocumentSearchErrorVotes
	}

	page.Message = strings.TrimSpace(page.Message)
	handler.MakeFullPage(writer, acc, page)
}

func PostCreateVotePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	doc := &database.Document{
		Type:                database.DocTypeVote,
		Organisation:        values.GetTrimmedString("organisation"),
		Title:               values.GetTrimmedString("title"),
		Author:              values.GetTrimmedString("author"),
		Body:                helper.MakeMarkdown(values.GetTrimmedString("markdown")),
		Public:              values.GetBool("public"),
		Removed:             false,
		MemberParticipation: values.GetBool("member"),
		AdminParticipation:  values.GetBool("admin"),
		Participants:        values.GetTrimmedArray("[]participants"),
		Reader:              values.GetTrimmedArray("[]reader"),
		VoteIDs:             values.GetCommaSeperatedArray("votes"),
		End:                 values.GetTime("end-time", "2006-01-02", time.UTC),
	}

	if doc.End.IsZero() {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentGeneralTimestampInvalid})
		return
	}

	doc.End = doc.End.Add((23 * time.Hour) + (50 * time.Minute))
	currTime := time.Now().UTC()
	if doc.End.Before(currTime.Add(addMin)) || doc.End.After(currTime.Add(addMax)) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTimeNotInAreaVote})
		return
	}

	if doc.Title == "" || doc.Body == "" {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ContentOrBodyAreEmpty})
		return
	}

	const maxTitleLength = 400
	if len([]rune(doc.Title)) > maxTitleLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ErrorTitleTooLong, maxTitleLength)})
		return
	}

	allowed, err := database.IsAccountAllowedToPostWith(acc, doc.Author)
	if !allowed || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentGeneralMissingPermissionForDocumentCreation})
		return
	}

	doc.ID = helper.GetUniqueID(doc.Author)

	doc.Flair, err = database.GetAccountFlairs(&database.Account{Name: doc.Author})
	if err != nil {
		slog.Info(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ErrorLoadingFlairInfoForAccount})
		return
	}

	err = database.CreateDocument(doc, acc)
	if errors.Is(err, database.DocumentHasInvalidVisibility) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCreateVoteHasInvalidVisibility})
		return
	} else if errors.Is(err, database.NotAllowedError) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCreateVoteNotAllowedError})
		return
	} else if errors.Is(err, database.DocumentHasNoAttachedVotes) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCreateVoteHasNoAttachedVotes})
		return
	} else if err != nil {
		slog.Info(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCreateVoteError})
		return
	}

	writer.Header().Add("HX-Redirect", fmt.Sprintf("/view/document/%s", doc.ID))
	writer.WriteHeader(http.StatusFound)
}
