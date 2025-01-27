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

func GetCreateVotePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	locTime := time.Now().In(acc.TimeZone)
	page := &handler.CreateVotePage{
		DateTime: locTime.Add(time.Hour * 48).Format("2006-01-02T15:04"),
		MinTime:  locTime.Add(addMin).Format("2006-01-02T15:04"),
		MaxTime:  locTime.Add(addMax).Format("2006-01-02T15:04"),
	}
	page.Reader = []string{""}
	page.Participants = []string{""}
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
	page.VoteChoice, err = database.GetVoteInfoList(acc)
	if err != nil {
		slog.Debug(err.Error())
		page.Message += "\n" + "Es ist ein Fehler bei der Suche nach den Abstimmung des Accounts aufgetreten"
	}

	page.Message = strings.TrimSpace(page.Message)
	handler.MakeFullPage(writer, acc, page)
}

func PostCreateVotePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen"})
		return
	}

	doc := &database.Document{
		Type:                database.DocTypeVote,
		Organisation:        values.GetTrimmedString("organisation"),
		Title:               values.GetTrimmedString("title"),
		Author:              values.GetTrimmedString("author"),
		Body:                handler.MakeMarkdown(values.GetTrimmedString("markdown")),
		Public:              values.GetBool("public"),
		Removed:             false,
		MemberParticipation: values.GetBool("member"),
		AdminParticipation:  values.GetBool("admin"),
		Participants:        values.GetTrimmedArray("[]participants"),
		Reader:              values.GetTrimmedArray("[]reader"),
		VoteIDs:             values.GetCommaSeperatedArray("votes"),
		End:                 values.GetTime("end-time", "2006-01-02T15:04", acc.TimeZone),
	}

	if doc.End.IsZero() {
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
