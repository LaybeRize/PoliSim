package titles

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"net/http"
	"net/url"
	"strings"
)

func GetEditTitlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.EditTitlePage{Title: nil}
	var err error

	query := helper.GetAdvancedURLValues(request)
	if query.Has("name") {
		page.Title, page.Holder, err = database.GetTitleAndHolder(query.GetTrimmedString("name"))

		page.IsError = true
		if err != nil {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Der gesuchte Name ist mit keinem Titel verbunden"})
			return
		}

		page.Holder = append(page.Holder, "")
		page.IsError = false
		page.Message = "Gesuchter Titel gefunden"

		page.AccountNames, err = database.GetNonBlockedNames()
		if err != nil {
			page.Message += "\n" + "Es ist ein Fehler bei der Suche nach der Accountnamensliste aufgetreten"
		}

		handler.MakeFullPage(writer, acc, page)
		return
	}

	page.Titels, err = database.GetTitleNameList()
	if err != nil {
		page.IsError = true
		page.Message = "Es ist ein Fehler bei der Suche nach der Titelnamensliste aufgetreten"
	}

	handler.MakeFullPage(writer, acc, page)
}

func PatchEditTitlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	oldTitleName := values.GetTrimmedString("oldName")
	titleUpdate := &database.Title{
		Name:     values.GetTrimmedString("name"),
		MainType: values.GetTrimmedString("main-group"),
		SubType:  values.GetTrimmedString("sub"),
		Flair:    values.GetTrimmedString("flair"),
	}
	names := values.GetTrimmedArray("[]holder")

	if titleUpdate.Name == "" || len(titleUpdate.Name) > 400 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Titelname leer oder länger als 400 Zeichen"})
		return
	}

	if titleUpdate.MainType == "" || len(titleUpdate.MainType) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Hauptgruppe leer oder länger als 200 Zeichen"})
		return
	}

	if titleUpdate.SubType == "" || len(titleUpdate.SubType) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Untergruppe leer oder länger als 200 Zeichen"})
		return
	}

	if strings.Contains(titleUpdate.Flair, ",") ||
		strings.Contains(titleUpdate.Flair, ";") ||
		len(titleUpdate.Flair) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Flair enthält ein Komma, Semikolon oder ist länger als 200 Zeichen"})
		return
	}

	err = database.UpdateTitle(oldTitleName, titleUpdate)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim überarbeiten des Titels aufgetreten"})
		return
	}

	err = database.AddTitleHolder(titleUpdate, names)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte Titel-Halter nicht erfolgreich updaten"})
		return
	}

	page := &handler.EditTitlePage{Title: titleUpdate, Holder: names}
	if _, actualHolder, err := database.GetTitleAndHolder(titleUpdate.Name); err == nil {
		page.Holder = append(actualHolder, "")
	}
	page.IsError = false
	page.Message = "Titel erfolgreich angepasst"
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + "Es ist ein Fehler bei der Suche nach der Accountnamensliste aufgetreten"
	}
	handler.MakePage(writer, acc, page)
}

func PutTitleSearchPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	name := values.GetTrimmedString("name")
	_, err = database.GetTitleByName(name)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte keinen Titel finden, der den Namen trägt"})
		return
	}

	writer.Header().Add("HX-Redirect", "/edit/title?name="+url.QueryEscape(name))
	writer.WriteHeader(http.StatusFound)
}
