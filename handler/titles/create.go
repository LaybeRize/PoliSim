package titles

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"net/http"
	"strings"
)

const (
	maxNameLength     = 600
	maxMainTypeLength = 400
	maxSubTypeLength
	maxFlairLength = 200
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
		page.Message = loc.ErrorLoadingAccountNames
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateTitlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	newTitle := &database.Title{
		Name:     values.GetTrimmedString("name"),
		MainType: values.GetTrimmedString("main-group"),
		SubType:  values.GetTrimmedString("sub-group"),
		Flair:    values.GetTrimmedString("flair"),
	}
	names := values.GetTrimmedArray("[]holder")

	if newTitle.Name == "" || newTitle.MainType == "" || newTitle.SubType == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.TitleGeneralInformationEmpty})
		return
	}

	if len([]rune(newTitle.Name)) > maxNameLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.TitleGeneralNameTooLong, maxNameLength)})
		return
	}

	if len([]rune(newTitle.MainType)) > maxMainTypeLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.TitleGeneralMainGroupTooLong, maxMainTypeLength)})
		return
	}

	if len([]rune(newTitle.SubType)) > maxSubTypeLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.TitleGeneralSubGroupTooLong, maxSubTypeLength)})
		return
	}

	if strings.Contains(newTitle.Flair, ",") ||
		strings.Contains(newTitle.Flair, ";") ||
		len([]rune(newTitle.Flair)) > maxFlairLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.TitleGeneralFlairContainsInvalidCharactersOrIsTooLong, maxFlairLength)})
		return
	}

	err = database.CreateTitle(newTitle, names)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.TitleErrorWhileCreating})
		return
	}

	page := &handler.CreateTitlePage{Holder: []string{""}}
	page.IsError = false
	page.Message = loc.TitleSuccessfullyCreated
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message = "\n" + loc.ErrorLoadingAccountNames
	}
	handler.MakeFullPage(writer, acc, page)
}
