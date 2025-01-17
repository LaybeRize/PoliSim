package organisations

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
	"net/url"
	"strings"
)

func GetEditOrgansationPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.EditOrganisationPage{Organisation: nil}
	var err error

	if orgName, exists := request.URL.Query()["name"]; exists {
		page.Organisation, page.User, page.Admin, err = database.GetFullOrganisationInfo(orgName[0])

		page.IsError = true
		if err != nil {
			page.Organisation = nil
			page.Message = "Der gesuchte Name ist mit keiner Organisation verbunden"
			page.Organisations, err = database.GetOrganisationNameList()
			if err != nil {
				page.Message += "\n" + "Es ist ein Fehler bei der Suche nach der Organisationsamensliste aufgetreten"
			}
			handler.MakeFullPage(writer, acc, page)
			return
		}

		page.User = append(page.User, "")
		page.Admin = append(page.Admin, "")
		page.IsError = false
		page.Message = "Gesuchte Organisation gefunden"

		page.AccountNames, err = database.GetNonBlockedNames()
		if err != nil {
			page.Message += "\n" + "Es ist ein Fehler bei der Suche nach der Accountnamensliste aufgetreten"
		}

		handler.MakeFullPage(writer, acc, page)
		return
	}

	page.Organisations, err = database.GetOrganisationNameList()
	if err != nil {
		page.IsError = true
		page.Message = "Es ist ein Fehler bei der Suche nach der Organisationsnamensliste aufgetreten"
	}

	handler.MakeFullPage(writer, acc, page)
}

func PatchEditOrganisationPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	err := request.ParseForm()
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	oldOrganisationName := helper.GetFormEntry(request, "oldName")
	organisationUpdate := &database.Organisation{}
	organisationUpdate.Name = helper.GetFormEntry(request, "name")
	database.GetIntegerFormEntry(request, "visiblity", &organisationUpdate.Visibility)
	organisationUpdate.MainType = helper.GetFormEntry(request, "main-group")
	organisationUpdate.SubType = helper.GetFormEntry(request, "sub-group")
	organisationUpdate.Flair = helper.GetFormEntry(request, "flair")
	userNames := helper.GetFormList(request, "[]user")
	adminNames := helper.GetFormList(request, "[]admin")

	if organisationUpdate.Name == "" || len(organisationUpdate.Name) > 400 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Organisationsname leer oder länger als 400 Zeichen"})
		return
	}

	if organisationUpdate.MainType == "" || len(organisationUpdate.MainType) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Hauptgruppe leer oder länger als 200 Zeichen"})
		return
	}

	if organisationUpdate.SubType == "" || len(organisationUpdate.SubType) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Untergruppe leer oder länger als 200 Zeichen"})
		return
	}

	if strings.Contains(organisationUpdate.Flair, ",") ||
		strings.Contains(organisationUpdate.Flair, ";") ||
		len(organisationUpdate.Flair) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Flair enthält ein Komma, Semikolon oder ist länger als 200 Zeichen"})
		return
	}

	if !organisationUpdate.VisibilityIsValid() {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die ausgewählte Sichtbarkeit ist nicht valide"})
		return
	}

	organisationUpdate.ClearInvalidFlair()

	err = database.UpdateOrganisation(oldOrganisationName, organisationUpdate)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim überarbeiten der Organisation aufgetreten"})
		return
	}

	err = database.AddOrganisationMember(organisationUpdate, userNames, adminNames)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte Organisationsmitglieder nicht erfolgreich updaten"})
		return
	}

	page := &handler.EditOrganisationPage{Organisation: organisationUpdate, User: userNames, Admin: adminNames}
	if _, actualUser, actualAdmins, err := database.GetFullOrganisationInfo(organisationUpdate.Name); err == nil {
		page.User = append(actualUser, "")
		page.Admin = append(actualAdmins, "")
	}
	page.IsError = false
	page.Message = "Organisation erfolgreich angepasst"
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + "Es ist ein Fehler bei der Suche nach der Accountnamensliste aufgetreten"
	}
	handler.MakePage(writer, acc, page)
}

func PutOrganisationSearchPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	err := request.ParseForm()
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	name := helper.GetFormEntry(request, "name")
	_, err = database.GetOrganisationByName(name)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte keine Organisation finden, die den Namen trägt"})
		return
	}

	writer.Header().Add("HX-Redirect", "/edit/organisation?name="+url.QueryEscape(name))
	writer.WriteHeader(http.StatusFound)
}
