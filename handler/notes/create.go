package notes

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

func GetNoteCreatePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	loaded, exists := request.URL.Query()["loaded"]
	page := &handler.CreateNotesPage{}
	if exists {
		page.References = strings.Join(loaded, ",")
	}
	page.IsError = true
	page.Message = ""

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		page.Message = loc.CouldNotFindAllAuthors
		arr = make([]string, 0, 1)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = acc.Name
	page.PossibleAuthors = arr

	handler.MakeFullPage(writer, acc, page)
}

func PostNoteCreatePage(writer http.ResponseWriter, request *http.Request) {
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

	author := values.GetTrimmedString("author")
	allowed, err := database.IsAccountAllowedToPostWith(acc, author)
	if !allowed || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.NoteAuthorIsInvalid})
		return
	}

	flairString, err := database.GetAccountFlairs(&database.Account{Name: author})
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ErrorLoadingFlairInfoForAccount})
		return
	}

	references := values.GetCommaSeperatedArray("references")
	note := &database.BlackboardNote{
		ID:       helper.GetUniqueID(author),
		Title:    values.GetTrimmedString("title"),
		Author:   author,
		Flair:    flairString,
		Body:     handler.MakeMarkdown(values.GetTrimmedString("markdown")),
		Removed:  false,
		Parents:  nil,
		Children: nil,
	}

	if note.Title == "" || string(note.Body) == "" {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ContentOrBodyAreEmpty})
		return
	}

	const maxTitleLength = 400
	if len([]rune(note.Title)) > maxTitleLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ErrorTitleTooLong, maxTitleLength)})
		return
	}

	err = database.CreateNote(note, references)
	if err != nil {
		slog.Error(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.NoteErrorWhileCreatingNote})
		return
	}

	page := &handler.CreateNotesPage{References: strings.Join(append(references, note.ID), ", "), Author: author}
	page.IsError = false
	page.Message = loc.NoteSuccessfullyCreatedNote

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		page.Message += "\n" + loc.CouldNotFindAllAuthors
		arr = make([]string, 0, 1)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = acc.Name
	page.PossibleAuthors = arr

	handler.MakePage(writer, acc, page)
}
