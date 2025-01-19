package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	addMin = time.Hour * 24
	addMax = time.Hour * 24 * 15
)

func GetCreateDiscussionPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	locTime := time.Now().In(acc.TimeZone)
	page := &handler.CreateDiscussionPage{
		DateTime: locTime.Add(time.Hour * 48).Format("2006-01-02T15:04"),
		MinTime:  locTime.Add(addMin).Format("2006-01-02T15:04"),
		MaxTime:  locTime.Add(addMax).Format("2006-01-02T15:04"),
	}
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
	page.PossibleOrganisations, err = database.GetOrganisationNamesAdminIn(acc.Name)
	if err != nil {
		slog.Debug(err.Error())
		page.Message = "\n" + "Konnte nicht alle erlaubten Organisationen für ausgewählten Account finden"
	}
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		slog.Debug(err.Error())
		page.Message += "\n" + "Es ist ein Fehler bei der Suche nach der Accountnamensliste aufgetreten"
	}

	page.Message = strings.TrimSpace(page.Message)
	handler.MakeFullPage(writer, acc, page)
}

func PostCreateDiscussionPage(writer http.ResponseWriter, request *http.Request) {
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

	doc := &database.Document{
		Type:                database.DocTypeDiscussion,
		Organisation:        helper.GetFormEntry(request, "organisation"),
		Title:               helper.GetFormEntry(request, "title"),
		Author:              helper.GetFormEntry(request, "author"),
		Body:                handler.MakeMarkdown(helper.GetFormEntry(request, "markdown")),
		Public:              helper.GetBoolFormEntry(request, "public"),
		Removed:             false,
		MemberParticipation: helper.GetBoolFormEntry(request, "member"),
		AdminParticipation:  helper.GetBoolFormEntry(request, "admin"),
		Participants:        helper.GetFormList(request, "[]participants"),
		Reader:              helper.GetFormList(request, "[]reader"),
	}

	doc.End, err = time.ParseInLocation("2006-01-02T15:04", helper.GetFormEntry(request, "end-time"),
		acc.TimeZone)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der angegebene Zeitstempel für das Ende ist nicht gültig"})
		return
	}

	locTime := time.Now().In(acc.TimeZone)
	if doc.End.Before(locTime.Add(addMin)) || doc.End.After(locTime.Add(addMax)) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der angegebene Zeitstempel ist entweder in weniger als 24 Stunden oder in mehr als 15 Tagen"})
		return
	}
	doc.End = doc.End.UTC()

	if doc.Title == "" || doc.Body == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Titel oder Inhalt sind leer"})
		return
	}

	if len(doc.Title) > 400 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Titel ist zu lang (400 Zeichen maximal)"})
		return
	}

	allowed, err := database.IsAccountAllowedToPostWith(acc, doc.Author)
	if !allowed || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehlende Berechtigung um mit diesem Account ein Dokument zu erstellen"})
		return
	}

	doc.ID = helper.GetUniqueID(doc.Author)

	doc.Flair, err = database.GetAccountFlairs(&database.Account{Name: doc.Author})
	if err != nil {
		slog.Info(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim laden der Flairs für den Autor"})
		return
	}

	err = database.CreateDocument(doc, acc)
	if err != nil {
		slog.Info(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim erstellen des Dokuments"})
		return
	}

	writer.Header().Add("HX-Redirect", fmt.Sprintf("/view/document/%s", doc.ID))
	writer.WriteHeader(http.StatusFound)
}

func PatchFixUserList(writer http.ResponseWriter, request *http.Request) {
	_, loggedIn := database.RefreshSession(writer, request)
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

	page := &handler.ReaderAndParticipants{
		Participants: helper.GetFormList(request, "[]participants"),
		Reader:       helper.GetFormList(request, "[]reader"),
	}

	page.Reader, err = database.FilterNameListForNonBlocked(page.Reader, 1)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte Lesernamensliste nicht filtern"})
		return
	}
	page.Participants, err = database.FilterNameListForNonBlocked(page.Participants, 1)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte Teilnehmernamensliste nicht filtern"})
		return
	}

	handler.MakeSpecialPagePart(writer, page)
}
