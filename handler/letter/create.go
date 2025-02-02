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
		page.Message = "Konnte nicht alle möglichen Autoren laden"
	}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + "Konnte mögliche Empfängernamen nicht laden"
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
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
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
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ErrorLoadingFlairInfoForAccount})
		return
	}

	if letter.Title == "" || letter.Body == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ContentOrBodyAreEmpty})
		return
	}

	const maxTitleLength = 400
	if len([]rune(letter.Title)) > maxTitleLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ErrorTitleTooLong, maxTitleLength)})
		return
	}

	allowed, _ := database.IsAccountAllowedToPostWith(acc, letter.Author)
	if !allowed && !(acc.IsAtLeastAdmin() && letter.Author == loc.AdministrationAccountName) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der Brief darf nicht mit dem angegeben Account als Autor verschickt werden"})
		return
	}

	letter.Reader, err = database.FilterNameListForNonBlocked(letter.Reader, 0)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte Empfängerliste nicht validieren"})
		return
	}

	if len(letter.Reader) == 0 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die Anzahl an Empfängern für den Brief darf nicht 0 sein"})
		return
	}

	err = database.CreateLetter(letter)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim erstellen des Briefs aufgetreten"})
		return
	}

	page := &handler.CreateLetterPage{Author: letter.Author, Recipients: []string{""}}
	page.IsError = false
	page.Message = "Brief erfolgreich erstellt"

	page.PossibleAuthors, err = database.GetMyAccountNames(acc)
	if err != nil {
		page.PossibleAuthors = []string{acc.Name}
		page.Message += "\n" + "Konnte nicht alle möglichen Autoren laden"
	}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + "Konnte mögliche Empfängernamen nicht laden"
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
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
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
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ContentOrBodyAreEmpty})
		return
	}

	const maxTitleLength = 400
	if len([]rune(page.Title)) > maxTitleLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ErrorTitleTooLong, maxTitleLength)})
		return
	}

	allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Author)
	if !allowed && !(acc.IsAtLeastAdmin() && page.Author == loc.AdministrationAccountName) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der Brief darf nicht mit dem angegeben Account als Autor verschickt werden"})
		return
	}

	page.Recipients, err = database.FilterNameListForNonBlocked(page.Recipients, 0)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte Empfängerliste nicht validieren"})
		return
	}

	if len(page.Recipients) == 0 {
		page.Recipients = []string{""}
		page.Message = "Die Anzahl an Empfängern für den Brief darf nicht 0 sein"
	} else {
		page.IsError = false
		page.Message = "Der Brief darf so versendet werden"
	}

	page.PossibleAuthors, err = database.GetMyAccountNames(acc)
	if err != nil {
		page.PossibleAuthors = []string{acc.Name}
		page.Message += "\n" + "Konnte nicht alle möglichen Autoren laden"
	}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + "Konnte mögliche Empfängernamen nicht laden"
	}

	if acc.IsAtLeastAdmin() {
		page.PossibleAuthors = append(page.PossibleAuthors, loc.AdministrationAccountName)
	}
	page.Message = strings.TrimSpace(page.Message)
	handler.MakePage(writer, acc, page)
}
