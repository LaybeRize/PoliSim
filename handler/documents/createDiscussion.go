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

const (
	addMin     = time.Hour * 24
	addMax     = time.Hour * 24 * 15
	runMinDays = 1
	runMaxDays = 14
)

func GetCreateDiscussionPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	locTime := time.Now().In(acc.TimeZone)
	page := &handler.CreateDiscussionPage{
		DateTime: locTime.Add(time.Hour * 48).Format("2006-01-02T15:04"),
		MinTime:  locTime.Add(addMin).Format("2006-01-02T15:04"),
		MaxTime:  locTime.Add(addMax).Format("2006-01-02T15:04"),
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

	page.Message = strings.TrimSpace(page.Message)
	handler.MakeFullPage(writer, acc, page)
}

func PostCreateDiscussionPage(writer http.ResponseWriter, request *http.Request) {
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
		Type:                database.DocTypeDiscussion,
		Organisation:        values.GetTrimmedString("organisation"),
		Title:               values.GetTrimmedString("title"),
		Author:              values.GetTrimmedString("author"),
		Body:                handler.MakeMarkdown(values.GetTrimmedString("markdown")),
		Public:              values.GetBool("public"),
		Removed:             false,
		MemberParticipation: values.GetBool("member"),
		AdminParticipation:  values.GetBool("admin"),
		Participants:        values.GetTrimmedArray("[]participants"),
		Reader:              values.GetTrimmedArray("[]reader"),
		End:                 values.GetTime("end-time", "2006-01-02T15:04", acc.TimeZone),
	}

	if doc.End.IsZero() {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentGeneralTimestampInvalid})
		return
	}

	locTime := time.Now().In(acc.TimeZone)
	if doc.End.Before(locTime.Add(addMin)) || doc.End.After(locTime.Add(addMax)) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTimeNotInAreaDiscussion})
		return
	}
	doc.End = doc.End.UTC()

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
			Message: loc.DocumentCreateDiscussionHasInvalidVisibility})
		return
	} else if errors.Is(err, database.NotAllowedError) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCreateDiscussionNotAllowedError})
		return
	} else if err != nil {
		slog.Info(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCreateDiscussionError})
		return
	}

	writer.Header().Add("HX-Redirect", fmt.Sprintf("/view/document/%s", doc.ID))
	writer.WriteHeader(http.StatusFound)
}
