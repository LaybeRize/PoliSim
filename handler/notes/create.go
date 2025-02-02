package notes

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"net/http"
	"strings"
	"time"
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
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	author, err := database.GetAccountByName(values.GetTrimmedString("author"))
	if err != nil || author.Blocked {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der ausgewählte Autor ist nicht valide"})
		return
	}

	ownerName, err := database.GetOwnerName(author)
	if author.Name != acc.Name && (err != nil || ownerName != acc.Name) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der ausgewählte Autor ist nicht valide"})
		return
	}

	flairString, err := database.GetAccountFlairs(author)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Beim Laden der Informationen für den Author ist ein Fehler aufgetreten"})
		return
	}

	references := values.GetCommaSeperatedArray("references")
	note := &database.BlackboardNote{
		ID:       helper.GetUniqueID(author.Name),
		Title:    values.GetTrimmedString("title"),
		Author:   author.Name,
		Flair:    flairString,
		PostedAt: time.Now().UTC(),
		Body:     handler.MakeMarkdown(values.GetTrimmedString("markdown")),
		Removed:  false,
		Parents:  nil,
		Children: nil,
	}

	if note.Title == "" || string(note.Body) == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ContentOrBodyAreEmpty})
		return
	}

	const maxTitleLength = 400
	if len([]rune(note.Title)) > maxTitleLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ErrorTitleTooLong, maxTitleLength)})
		return
	}

	err = database.CreateNote(note, references)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim erstellen der Notiz aufgetreten"})
		return
	}

	page := &handler.CreateNotesPage{References: strings.Join(append(references, note.ID), ", "), Author: author.Name}
	page.IsError = false
	page.Message = "Notiz erfolgreich erstellt"

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		page.Message += "\n" + "Hinweis: Konnte nicht alle möglichen Autoren finden"
		arr = make([]string, 0, 1)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = acc.Name
	page.PossibleAuthors = arr

	handler.MakePage(writer, acc, page)
}
