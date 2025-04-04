package titles

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
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
			handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
				Message: loc.TitleNoTitleWithThatName})
			return
		}

		page.Holder = append(page.Holder, "")
		page.IsError = false
		page.Message = loc.TitleFoundTitle

		page.AccountNames, err = database.GetNonBlockedNames()
		if err != nil {
			page.Message += "\n" + loc.ErrorSearchingForAccountNames
		}

		handler.MakeFullPage(writer, acc, page)
		return
	}

	page.Titels, err = database.GetTitleNameList()
	if err != nil {
		page.IsError = true
		page.Message = loc.TitleErrorSearchingForTitleList
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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	oldTitleName := values.GetTrimmedString("oldName")
	titleUpdate := &database.Title{
		Name:     values.GetTrimmedString("name"),
		MainType: values.GetTrimmedString("main-group"),
		SubType:  values.GetTrimmedString("sub-group"),
		Flair:    values.GetTrimmedString("flair"),
	}
	names := values.GetTrimmedArray("[]holder")

	if titleUpdate.Name == "" || titleUpdate.MainType == "" || titleUpdate.SubType == "" {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.TitleGeneralInformationEmpty})
		return
	}

	if len([]rune(titleUpdate.Name)) > maxNameLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.TitleGeneralNameTooLong, maxNameLength)})
		return
	}

	if len([]rune(titleUpdate.MainType)) > maxMainTypeLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.TitleGeneralMainGroupTooLong, maxMainTypeLength)})
		return
	}

	if len([]rune(titleUpdate.SubType)) > maxSubTypeLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.TitleGeneralSubGroupTooLong, maxSubTypeLength)})
		return
	}

	if strings.Contains(titleUpdate.Flair, ",") ||
		strings.Contains(titleUpdate.Flair, ";") ||
		len([]rune(titleUpdate.Flair)) > maxFlairLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.TitleGeneralFlairContainsInvalidCharactersOrIsTooLong, maxFlairLength)})
		return
	}

	err = database.UpdateTitle(oldTitleName, titleUpdate)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.TitleErrorWhileUpdatingTitle})
		return
	}

	err = database.AddTitleHolder(titleUpdate, names)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.TitleErrorWhileUpdatingTitleHolder})
		return
	}

	page := &handler.EditTitlePage{Title: titleUpdate, Holder: names}
	if _, actualHolder, err := database.GetTitleAndHolder(titleUpdate.Name); err == nil {
		page.Holder = append(actualHolder, "")
	}
	page.IsError = false
	page.Message = loc.TitleSuccessfullyUpdated
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + loc.ErrorSearchingForAccountNames
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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	name := values.GetTrimmedString("name")
	_, err = database.GetTitleByName(name)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.TitleNotFoundByName})
		return
	}

	writer.Header().Add("HX-Redirect", "/edit/title?name="+url.QueryEscape(name))
	writer.WriteHeader(http.StatusFound)
}
