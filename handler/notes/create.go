package notes

import (
	"PoliSim/database"
	"PoliSim/handler"
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
		page.Refrences = strings.Join(loaded, ",")
	}
	page.IsError = true
	page.Message = ""

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		page.Message = "Konnte nicht alle möglichen Autoren finden"
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

	err := request.ParseForm()
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	author, err := database.GetAccountByName(request.Form.Get("author"))
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

	references := handler.MakeCommaSeperatedStringToList(request.Form.Get("references"))
	note := &database.BlackboardNote{
		ID:       handler.GetUniqueID(author.Name),
		Title:    request.Form.Get("title"),
		Author:   author.Name,
		Flair:    flairString,
		PostedAt: time.Now().UTC(),
		Body:     handler.MakeMarkdown(request.Form.Get("markdown")),
		Removed:  false,
		Parents:  nil,
		Children: nil,
	}

	if note.Title == "" || string(note.Body) == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Titel oder Inhalt sind leer"})
		return
	}

	if len(note.Title) > 400 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Titel ist zu lang (400 Zeichen maximal)"})
		return
	}

	err = database.CreateNote(note, references)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim erstellen der Notiz aufgetreten"})
		return
	}

	page := &handler.CreateNotesPage{Refrences: strings.Join(append(references, note.ID), ","), Author: author.Name}
	page.IsError = false
	page.Message = "Notiz erfolgreich erstellt"

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		page.Message += "\nHinweis: Konnte nicht alle möglichen Autoren finden"
		arr = make([]string, 0, 1)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = acc.Name
	page.PossibleAuthors = arr

	handler.MakePage(writer, acc, page)
}
