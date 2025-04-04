package letter

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func GetCreateLetterPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}
	page := &handler.CreateLetterPage{Recipients: []string{""}}
	page.IsError = true
	page.Author = acc.Name

	var err error
	page.PossibleAuthors, err = database.GetMyAccountNames(acc)
	if err != nil {
		page.PossibleAuthors = []string{acc.Name}
		page.Message = loc.CouldNotFindAllAuthors
	}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + loc.LetterErrorLoadingRecipients
	}

	if acc.IsAtLeastAdmin() {
		page.PossibleAuthors = append(page.PossibleAuthors, loc.AdministrationAccountName)
	}
	page.Message = strings.TrimSpace(page.Message)
	handler.MakeFullPage(writer, acc, page)
}

func PostCreateLetterPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	letter := &database.Letter{
		Title:    values.GetTrimmedString("title"),
		Author:   values.GetTrimmedString("author"),
		Body:     handler.MakeMarkdown(values.GetTrimmedString("markdown")),
		Signable: values.GetBool("signable"),
		Reader:   values.GetTrimmedArray("[]recipient"),
	}
	letter.ID = helper.GetUniqueID(letter.Author)
	letter.Flair, err = database.GetAccountFlairs(&database.Account{Name: letter.Author})
	if err != nil {
		slog.Info(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ErrorLoadingFlairInfoForAccount})
		return
	}

	if letter.Title == "" || letter.Body == "" {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ContentOrBodyAreEmpty})
		return
	}

	const maxTitleLength = 400
	if len([]rune(letter.Title)) > maxTitleLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ErrorTitleTooLong, maxTitleLength)})
		return
	}

	allowed, _ := database.IsAccountAllowedToPostWith(acc, letter.Author)
	if !allowed && !(acc.IsAtLeastAdmin() && letter.Author == loc.AdministrationAccountName) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.LetterNotAllowedToPostWithThatAccount})
		return
	}

	letter.Reader, err = database.FilterNameListForNonBlocked(letter.Reader, 0)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.LetterRecipientListUnvalidated})
		return
	}

	if len(letter.Reader) == 0 {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.LetterNeedAtLeastOneRecipient})
		return
	}

	err = database.CreateLetter(letter)
	if err != nil {
		slog.Error(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.LetterErrorWhileSending})
		return
	}

	page := &handler.CreateLetterPage{Author: letter.Author, Recipients: []string{""}}
	page.IsError = false
	page.Message = loc.LetterSuccessfullySendLetter

	page.PossibleAuthors, err = database.GetMyAccountNames(acc)
	if err != nil {
		page.PossibleAuthors = []string{acc.Name}
		page.Message += "\n" + loc.CouldNotFindAllAuthors
	}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + loc.LetterErrorLoadingRecipients
	}

	if acc.IsAtLeastAdmin() {
		page.PossibleAuthors = append(page.PossibleAuthors, loc.AdministrationAccountName)
	}
	handler.MakePage(writer, acc, page)
}

func PatchCheckCreateLetterPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.CreateLetterPage{
		Title:      values.GetTrimmedString("title"),
		Author:     values.GetTrimmedString("author"),
		Body:       values.GetTrimmedString("markdown"),
		Signable:   values.GetBool("signable"),
		Recipients: values.GetTrimmedArray("[]recipient"),
	}
	page.Information = handler.MakeMarkdown(page.Body)
	page.IsError = true

	if page.Title == "" || page.Body == "" {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ContentOrBodyAreEmpty})
		return
	}

	const maxTitleLength = 400
	if len([]rune(page.Title)) > maxTitleLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ErrorTitleTooLong, maxTitleLength)})
		return
	}

	allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Author)
	if !allowed && !(acc.IsAtLeastAdmin() && page.Author == loc.AdministrationAccountName) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.LetterNotAllowedToPostWithThatAccount})
		return
	}

	page.Recipients, err = database.FilterNameListForNonBlocked(page.Recipients, 0)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.LetterRecipientListUnvalidated})
		return
	}

	if len(page.Recipients) == 0 {
		page.Recipients = []string{""}
		page.Message = loc.LetterNeedAtLeastOneRecipient
	} else {
		page.IsError = false
		page.Message = loc.LetterAllowedToBeSent
	}

	page.PossibleAuthors, err = database.GetMyAccountNames(acc)
	if err != nil {
		page.PossibleAuthors = []string{acc.Name}
		page.Message += "\n" + loc.CouldNotFindAllAuthors
	}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + loc.LetterErrorLoadingRecipients
	}

	if acc.IsAtLeastAdmin() {
		page.PossibleAuthors = append(page.PossibleAuthors, loc.AdministrationAccountName)
	}
	page.Message = strings.TrimSpace(page.Message)
	handler.MakePage(writer, acc, page)
}
