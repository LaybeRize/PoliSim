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

func GetCreateDocumentPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.CreateDocumentPage{}
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
		page.Message = strings.TrimSpace(page.Message)
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateDocumentPage(writer http.ResponseWriter, request *http.Request) {
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
		Type:                database.DocTypePost,
		Organisation:        values.GetTrimmedString("organisation"),
		Title:               values.GetTrimmedString("title"),
		Author:              values.GetTrimmedString("author"),
		Body:                helper.MakeMarkdown(values.GetTrimmedString("markdown")),
		Public:              true,
		Removed:             false,
		MemberParticipation: false,
		AdminParticipation:  false,
		End:                 time.Time{},
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
			Message: loc.DocumentCreatePostHasInvalidVisibility})
		return
	} else if errors.Is(err, database.NotAllowedError) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCreatePostNotAllowedError})
		return
	} else if err != nil {
		slog.Info(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCreatePostError})
		return
	}

	writer.Header().Add("HX-Redirect", fmt.Sprintf("/view/document/%s", doc.ID))
	writer.WriteHeader(http.StatusFound)
}
