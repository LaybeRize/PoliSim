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
		page.Message = loc.ErrorLoadingAccountNames
	}

	page.NewspaperNames, err = database.GetNewspaperNameList()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + loc.NewspaperErrorLoadingNewspaperNames
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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.NewspaperErrorWhileCreatingNewspaper})
		return
	}

	page := &handler.ManageNewspaperPage{MessageUpdate: handler.MessageUpdate{IsError: false,
		Message: loc.NewspaperSuccessfullyCreatedNewspaper}, HadError: false}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		slog.Debug(err.Error())
		page.Message = "\n" + loc.ErrorLoadingAccountNames
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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.ManageNewspaperPage{}
	page.IsError = false

	newspaper, err := database.GetFullNewspaperInfo(values.GetTrimmedString("name"))
	if err != nil {
		slog.Error(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.NewspaperErrorWhileSearchingNewspaper})
		return
	}

	page.Message = loc.NewspaperSuccessfullyFoundNewspaper
	newspaper.Authors = append(newspaper.Authors, "")
	page.Newspaper = *newspaper

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + loc.ErrorLoadingAccountNames
	}

	page.NewspaperNames, err = database.GetNewspaperNameList()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + loc.NewspaperErrorLoadingNewspaperNames
	}

	handler.MakeSpecialPagePart(writer, page)
}

func PatchUpdateNewspaperPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastAdmin() {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.NewspaperErrorWhileChangingNewspaper})
		return
	}
	err = database.UpdateNewspaper(&page.Newspaper)
	if err != nil {
		slog.Error(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.NewspaperErrorWhileAddingReporters})
		return
	}

	if newspaper, err := database.GetFullNewspaperInfo(page.Newspaper.Name); err == nil {
		newspaper.Authors = append(newspaper.Authors, "")
		page.Newspaper = *newspaper
	}

	page.Message = loc.NewspaperSuccessfullyChangedNewspaper
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + loc.ErrorLoadingAccountNames
	}

	page.NewspaperNames, err = database.GetNewspaperNameList()
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = "\n" + loc.NewspaperErrorLoadingNewspaperNames
	}

	handler.MakeSpecialPagePart(writer, page)
}
