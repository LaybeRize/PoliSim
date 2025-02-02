package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"log/slog"
	"net/http"
	"strings"
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
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "Konnte Accountnamen nicht laden"
	}

	page.NewspaperNames, err = database.GetNewspaperNameList()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + "Konnte Zeitungsnamen nicht laden"
		page.Message = strings.TrimSpace(page.Message)
	}

	page.Publications, err = database.GetUnpublishedPublications()
	if err != nil {
		page.HadError = true
		slog.Debug(err.Error())
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateNewspaperPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastAdmin() {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	newspaper := &database.Newspaper{
		Name:    values.GetTrimmedString("name"),
		Authors: nil,
	}
	err = database.CreateNewspaper(newspaper)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Erstellen der Zeitung (überprüfe ob die Zeitung bereits existiert)"})
		return
	}

	page := &handler.ManageNewspaperPage{MessageUpdate: handler.MessageUpdate{IsError: false,
		Message: "Zeitung erfolgreich erstellt"}, HadError: false}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		slog.Debug(err.Error())
		page.Message = "\n" + "Konnte Accountnamen nicht laden"
	}

	page.Publications, err = database.GetUnpublishedPublications()
	if err != nil {
		page.HadError = true
		slog.Debug(err.Error())
	}

	handler.MakePage(writer, acc, page)
}

func PutSearchNewspaperPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastAdmin() {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.ManageNewspaperPage{}
	page.IsError = false

	newspaper, err := database.GetFullNewspaperInfo(values.GetTrimmedString("name"))
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler bei der Suche der Zeitung"})
		return
	}

	page.Message = "Zeitung gefunden"
	newspaper.Authors = append(newspaper.Authors, "")
	page.Newspaper = *newspaper

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + "Konnte Accountnamen nicht laden"
	}

	page.NewspaperNames, err = database.GetNewspaperNameList()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + "Konnte Zeitungsnamen nicht laden"
	}

	handler.MakeSpecialPagePart(writer, page)
}

func PatchUpdateNewspaperPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastAdmin() {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.ManageNewspaperPage{Newspaper: database.Newspaper{
		Name:    values.GetTrimmedString("name"),
		Authors: values.GetTrimmedArray("[]author"),
	}}
	page.IsError = false

	err = database.RemoveAccountsFromNewspaper(&page.Newspaper)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Anpassen der Zeitung"})
		return
	}
	err = database.UpdateNewspaper(&page.Newspaper)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim hinzufügen der neuen Autoren zur Zeitung"})
		return
	}

	if newspaper, err := database.GetFullNewspaperInfo(page.Newspaper.Name); err == nil {
		newspaper.Authors = append(newspaper.Authors, "")
		page.Newspaper = *newspaper
	}

	page.Message = "Zeitung angepasst"
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + "Konnte Accountnamen nicht laden"
	}

	page.NewspaperNames, err = database.GetNewspaperNameList()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + "Konnte Zeitungsnamen nicht laden"
	}

	handler.MakeSpecialPagePart(writer, page)
}
