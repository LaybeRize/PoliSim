package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"log/slog"
	"net/http"
)

func GetFindOrganisationForAccountPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.UpdateOrganisationForUser{}
	author := values.GetTrimmedString("author")
	allowed, err := database.IsAccountAllowedToPostWith(acc, author)
	if !allowed || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissionForAccountInfo})
		return
	}

	page.PossibleOrganisations, err = database.GetOrganisationNamesAdminIn(author)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ErrorFindingAllOrganisationsForAccount})
		return
	}

	handler.MakeSpecialPagePart(writer, page)
}

func PatchFixUserList(writer http.ResponseWriter, request *http.Request) {
	_, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.ReaderAndParticipants{
		Participants: values.GetTrimmedArray("[]participants"),
		Reader:       values.GetTrimmedArray("[]reader"),
	}

	page.Reader, err = database.FilterNameListForNonBlocked(page.Reader, 1)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCloudNotFilterReaders})
		return
	}
	page.Participants, err = database.FilterNameListForNonBlocked(page.Participants, 1)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentCloudNotFilterParticipants})
		return
	}

	handler.MakeSpecialPagePart(writer, page)
}
