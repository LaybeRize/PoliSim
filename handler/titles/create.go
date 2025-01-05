package titles

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
	"strings"
)

func GetCreateTitlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.CreateTitlePage{Holder: []string{""}}
	var err error
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.IsError = true
		page.Message = "Konnte Accountnamen nicht laden"
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateTitlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	err := request.ParseForm()
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	newTitle := &database.Title{}
	newTitle.Name = helper.GetFormEntry(request, "name")
	newTitle.MainType = helper.GetFormEntry(request, "main-group")
	newTitle.SubType = helper.GetFormEntry(request, "sub-group")
	newTitle.Flair = helper.GetFormEntry(request, "flair")
	names := helper.GetFormList(request, "[]holder")

	if newTitle.Name == "" || len(newTitle.Name) > 400 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Titelname leer oder länger als 400 Zeichen"})
		return
	}

	if newTitle.MainType == "" || len(newTitle.MainType) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Hauptgruppe leer oder länger als 200 Zeichen"})
		return
	}

	if newTitle.SubType == "" || len(newTitle.SubType) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Untergruppe leer oder länger als 200 Zeichen"})
		return
	}

	if strings.Contains(newTitle.Flair, ",") ||
		strings.Contains(newTitle.Flair, ";") ||
		len(newTitle.Flair) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Flair enthält ein Komma, Semikolon oder ist länger als 200 Zeichen"})
		return
	}

	err = database.CreateTitle(newTitle, names)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim erstellen des Titels aufgetreten (Überprüf ob der Name des Titel " +
				"einzigartig ist)"})
		return
	}

	page := &handler.CreateTitlePage{Holder: []string{""}}
	page.IsError = false
	page.Message = "Titel erfolgreich erstellt"
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message = "\nKonnte Accountnamen nicht laden"
	}
	handler.MakeFullPage(writer, acc, page)
}
