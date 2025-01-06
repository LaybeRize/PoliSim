package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"log/slog"
	"net/http"
	"strings"
)

func GetCreateArticlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.CreateArticlePage{}
	page.IsError = true
	page.Message = ""

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		slog.Debug(err.Error())
		page.Message = "Konnte nicht alle möglichen Autoren finden"
		arr = make([]string, 0)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = acc.Name
	page.PossibleAuthors = arr
	page.PossibleNewspaper, err = database.GetNewspaperNameListForAccount(acc.Name)
	if err != nil {
		slog.Debug(err.Error())
		page.Message = "\n" + "Konnte nicht alle möglichen Zeitungen für ausgewählten Account finden"
		page.Message = strings.TrimSpace(page.Message)
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateArticlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	err := request.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen"})
		return
	}

	article := &database.NewspaperArticle{
		Title:    helper.GetFormEntry(request, "title"),
		Subtitle: helper.GetFormEntry(request, "subtitle"),
		Author:   helper.GetFormEntry(request, "author"),
		RawBody:  helper.GetFormEntry(request, "markdown"),
	}
	isSpecial := helper.GetFormEntry(request, "special") == "true"
	newspaper := helper.GetFormEntry(request, "newspaper")

	if article.Title == "" || article.Body == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Titel oder Inhalt sind leer"})
		return
	}

	if len(article.Title) > 400 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Titel ist zu lang (400 Zeichen maximal)"})
		return
	}

	if len(article.Subtitle) > 600 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Untertitel ist zu lang (600 Zeichen maximal)"})
		return
	}

	allowed, err := database.CheckIfUserAllowedInNewspaper(acc, article.Author, newspaper)
	if !allowed || err != nil {
		if err != nil {
			slog.Debug(err.Error())
		}
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehlende Berechtigung um mit diesem Account in dieser Zeitung zu posten"})
		return
	}

	article.Flair, err = database.GetAccountFlairs(&database.Account{Name: article.Author})
	article.Body = handler.MakeMarkdown(article.RawBody)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim laden der Flairs für den Autor"})
		return
	}

	err = database.CreateArticle(article, isSpecial, newspaper)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim erstellen des Artikels"})
		return
	}

	page := &handler.CreateArticlePage{}
	page.IsError = false
	page.Message = "Artikel erfolgreich erstellt"

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		slog.Debug(err.Error())
		page.Message = "\n" + "Konnte nicht alle möglichen Autoren finden"
		arr = make([]string, 0)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = article.Author
	page.PossibleAuthors = arr
	page.PossibleNewspaper, err = database.GetNewspaperNameListForAccount(acc.Name)
	if err != nil {
		page.Message = "\n" + "Konnte nicht alle möglichen Zeitungen für ausgewählten Account finden"
	}

	handler.MakePage(writer, acc, page)

}

func GetFindNewspaperForAccountPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehlende Berechtigung"})
		return
	}

	err := request.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen"})
		return
	}

	baseAcc, owner, err := database.GetAccountAndOwnerByAccountName(helper.GetFormEntry(request, "author"))
	if !((baseAcc.Exists() && baseAcc.Name == acc.Name) || (owner.Exists() && owner.Name == acc.Name)) || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehlende Berechtigung um die Informationen für diesen Account anzufordern"})
		return
	}

	page := &handler.CreateArticlePage{}
	page.PossibleNewspaper, err = database.GetNewspaperNameListForAccount(baseAcc.Name)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte nicht alle möglichen Zeitungen für ausgewählten Account finden"})
		return
	}
	handler.MakeSpecialPagePart(writer, page)
}
