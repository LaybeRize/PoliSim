package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"log/slog"
	"net/http"
)

func PostCreateComment(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Diese Funktion ist nicht verfügbar"})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	docId := request.PathValue("id")
	comment := &database.DocumentComment{
		Author: values.GetTrimmedString("author"),
		Body:   handler.MakeMarkdown(values.GetTrimmedString("markdown")),
	}
	comment.ID = helper.GetUniqueID(comment.Author)

	if comment.Body == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Inhalt ist leer"})
		return
	}

	allowed, err := database.IsAccountAllowedToPostWith(acc, comment.Author)
	if !allowed || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehlende Berechtigung um mit diesem Account ein Dokument zu erstellen"})
		return
	}

	comment.Flair, err = database.GetAccountFlairs(&database.Account{Name: comment.Author})
	if err != nil {
		slog.Info(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim laden der Flairs für den Autor"})
		return
	}

	err = database.CreateDocumentComment(docId, comment)
	if err != nil {
		slog.Info(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Speichern des Kommentars"})
		return
	}

	if obj := getDocumentPageObject(acc, request); obj != nil {
		handler.MakePage(writer, acc, obj)
	} else {
		handler.PartialGetNotFoundPage(writer, request)
	}
}
