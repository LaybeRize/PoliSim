package organisations

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
	"strings"
)

func GetCreateOrganisationPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.CreateOrganisationPage{
		Admin: []string{""},
		User:  []string{""},
	}
	var err error
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.IsError = true
		page.Message = "Konnte Accountnamen nicht laden"
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateOrganisationPage(writer http.ResponseWriter, request *http.Request) {
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

	newOrganisation := &database.Organisation{}
	newOrganisation.Name = helper.GetFormEntry(request, "name")
	newOrganisation.Visibility = database.OrganisationVisibility(helper.GetFormEntry(request, "visiblity"))
	newOrganisation.MainType = helper.GetFormEntry(request, "main-group")
	newOrganisation.SubType = helper.GetFormEntry(request, "sub-group")
	newOrganisation.Flair = helper.GetFormEntry(request, "flair")
	userNames := helper.GetFormList(request, "[]user")
	adminNames := helper.GetFormList(request, "[]admin")

	if newOrganisation.Name == "" || len(newOrganisation.Name) > 400 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Organisationsname leer oder länger als 400 Zeichen"})
		return
	}

	if newOrganisation.MainType == "" || len(newOrganisation.MainType) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Hauptgruppe leer oder länger als 200 Zeichen"})
		return
	}

	if newOrganisation.SubType == "" || len(newOrganisation.SubType) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Untergruppe leer oder länger als 200 Zeichen"})
		return
	}

	if strings.Contains(newOrganisation.Flair, ",") ||
		strings.Contains(newOrganisation.Flair, ";") ||
		len(newOrganisation.Flair) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Flair enthält ein Komma, Semikolon oder ist länger als 200 Zeichen"})
		return
	}

	if !newOrganisation.VisibilityIsValid() {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die ausgewählte Sichtbarkeit ist nicht valide"})
		return
	}

	newOrganisation.ClearInvalidFlair()

	err = database.CreateOrganisation(newOrganisation, userNames, adminNames)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim erstellen der Organisation aufgetreten (Überprüf ob der Name der " +
				"Organisation einzigartig ist)"})
		return
	}

	page := &handler.CreateOrganisationPage{
		Admin: []string{""},
		User:  []string{""},
	}
	page.IsError = false
	page.Message = "Organisation erfolgreich erstellt"
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message = "\nKonnte Accountnamen nicht laden"
	}
	handler.MakeFullPage(writer, acc, page)
}
