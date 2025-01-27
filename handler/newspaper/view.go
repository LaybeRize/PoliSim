package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"log/slog"
	"net/http"
)

func GetSpecificPublicationPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	pubID := request.PathValue("id")
	found, err := database.GetPublicationForUser(pubID, acc.IsAtLeastPressAdmin())
	if !found || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.ViewPublicationPage{}
	var pub *database.Publication
	pub, page.Articles, err = database.GetPublication(pubID)
	if page.QueryError = err != nil; !page.QueryError {
		page.Publication = *pub
	} else {
		slog.Error(err.Error())
	}

	handler.MakeFullPage(writer, acc, page)
}

func PatchPublishPublication(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	pubID := request.PathValue("id")
	found, err := database.GetPublicationForUser(pubID, acc.IsAtLeastPressAdmin())
	if !found || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	err = database.PublishPublication(pubID)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: "Es ist ein Fehler beim Publizieren aufgetreten",
			IsError: true,
		})
	}

	page := &handler.ViewPublicationPage{}
	var pub *database.Publication
	pub, page.Articles, err = database.GetPublication(pubID)
	if page.QueryError = err != nil; !page.QueryError {
		page.Publication = *pub
	} else {
		slog.Error(err.Error())
	}

	handler.MakePage(writer, acc, page)
}

func DeleteArticle(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastPressAdmin() {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: "Fehler beim parsen der Informationen",
			IsError: true,
		})
		return
	}
	rejectionText := values.GetTrimmedString("rejection")
	if rejectionText == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: "Der Zurückweisungsgrund darf nicht leer sein",
			IsError: true,
		})
		return
	}

	transaction, err := database.RejectableArticle(request.PathValue("id"))
	if err != nil {
		slog.Debug("Possible error while trying to localte article", "error", err)
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: "Konnte keinen Artikel mit der angegeben ID finden, welcher noch nicht publiziert wurde",
			IsError: true,
		})
		return
	}
	err = transaction.DeleteArticle()
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: "Konnte den Artikel nicht löschen",
			IsError: true,
		})
		return
	}

	letter := &database.Letter{
		Title: fmt.Sprintf("Zurückweisung des Artikels '%s' geschrieben für %s",
			transaction.Article.Title, transaction.NewspaperName),
		Author:   loc.AdminstrationName,
		Flair:    "",
		Signable: false,
		Body: handler.MakeMarkdown(fmt.Sprintf("# Zurückweisungsgrund\n\n%s\n\n# Artikelinhalt\n\n```%s```",
			rejectionText, transaction.Article.RawBody)),
		Reader: []string{transaction.Article.Author},
	}
	letter.ID = helper.GetUniqueID(letter.Author)

	err = transaction.CreateLetter(letter)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: "Fehler beim erstellen des Briefs an den Autor des Artikels",
			IsError: true,
		})
		return
	}

	writer.Header().Add("HX-Redirect", "/check/newspapers")
	writer.WriteHeader(http.StatusFound)
}
