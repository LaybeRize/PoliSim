package letter

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"net/http"
)

func GetAdminLetterSearchPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastAdmin() {
		handler.GetNotFoundPage(writer, request)
		return
	}

	handler.MakeFullPage(writer, acc, &handler.AdminSearchLetterPage{AccountNameToUse: loc.AdministrationName})
}

func GetLetterViewPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}
	id := request.PathValue("id")
	reader := helper.GetAdvancedURLValues(request).GetString("viewer")
	if allowed, err := database.IsAccountAllowedToPostWith(acc, reader); !checkValidSpecialAccounts(acc, reader) &&
		(!allowed || err != nil) {
		handler.GetNotFoundPage(writer, request)
		return
	}
	letter, err := database.GetLetterForReader(id, reader)
	if err != nil {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.ViewLetterPage{Letter: *letter}
	handler.MakeFullPage(writer, acc, page)
}

func PatchLetterViewPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}
	id := request.PathValue("id")
	query := helper.GetAdvancedURLValues(request)
	reader := query.GetTrimmedString("viewer")
	if allowed, err := database.IsAccountAllowedToPostWith(acc, reader); !checkValidSpecialAccounts(acc, reader) &&
		(!allowed || err != nil) {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	var err error
	switch query.GetTrimmedString("decision") {
	case "accept":
		err = database.UpdateSingatureStatus(id, reader, true)
	case "decline":
		err = database.UpdateSingatureStatus(id, reader, false)
	default:
		writer.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	letter, err := database.GetLetterForReader(id, reader)
	if err != nil {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	page := &handler.ViewLetterPage{Letter: *letter}
	handler.MakePage(writer, acc, page)
}

func checkValidSpecialAccounts(acc *database.Account, reader string) bool {
	return (reader == loc.AdministrationAccountName || reader == loc.AdministrationName) && acc.IsAtLeastAdmin()
}
