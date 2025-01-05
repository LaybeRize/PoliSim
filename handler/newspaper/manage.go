package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
)

func GetManageNewspaperPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastPressAdmin() {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.ManageNewspaperPage{}
	var err error

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.IsError = true
		page.Message = "Konnte Accountnamen nicht laden"
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateNewspaperPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastAdmin() {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true, Message: "Fehlende Berechtigung"})
		return
	}

	err := request.ParseForm()
	if err != nil {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	newspaper := &database.Newspaper{
		Name:    helper.GetFormEntry(request, "name"),
		Authors: nil,
	}
	err = database.CreateNewspaper(newspaper)
	if err != nil {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Erstellen der Zeitung (überprüfe ob die Zeitung bereits existiert)"})
		return
	}

	handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: false, Message: "Zeitung erfolgreich erstellt"})
}
