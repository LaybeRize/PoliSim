package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"log/slog"
	"net/http"
)

func GetDocumentViewPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)

	if obj := getDocumentPageObject(acc, request); obj != nil {
		handler.MakeFullPage(writer, acc, obj)
	} else {
		handler.GetNotFoundPage(writer, request)
	}
}

func getDocumentPageObject(acc *database.Account, request *http.Request) *handler.DocumentViewPage {
	id := request.PathValue("id")
	var err error
	page := &handler.DocumentViewPage{ColorPalettes: database.ColorPaletteMap}
	page.Document, page.Commentator, err = database.GetDocumentForUser(id, acc)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	return page
}

const elementID = "tag-message"

func PostNewDocumentTagPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Diese Funktion ist nicht verfügbar", ElementID: elementID})
		return
	}

	err := request.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen", ElementID: elementID})
		return
	}

	tag := &database.DocumentTag{
		ID:              helper.GetUniqueID(acc.Name),
		Text:            helper.GetFormEntry(request, "text"),
		BackgroundColor: helper.GetFormEntry(request, "background-color"),
		TextColor:       helper.GetFormEntry(request, "text-color"),
		LinkColor:       helper.GetFormEntry(request, "link-color"),
		Links:           helper.GetCommaListFormEntry(request, "links"),
	}

	if tag.Text == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der Tag-Text ist leer", ElementID: elementID})
		return
	}

	if len(tag.Text) > 400 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der Tag-Text ist länger als 400 Zeichen", ElementID: elementID})
		return
	}

	if !helper.StringIsAColor(tag.BackgroundColor) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die Farbe für den Hintergrund ist nicht valide", ElementID: elementID})
		return
	}

	if !helper.StringIsAColor(tag.TextColor) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die Farbe für den Text ist nicht valide", ElementID: elementID})
		return
	}

	if !helper.StringIsAColor(tag.LinkColor) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die Farbe für die Links ist nicht valide", ElementID: elementID})
		return
	}

	err = database.CreateTagForDocument(request.PathValue("id"), acc, tag)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Erstellen des Tags", ElementID: elementID})
		return
	}

	if obj := getDocumentPageObject(acc, request); obj != nil {
		handler.MakePage(writer, acc, obj)
	} else {
		handler.GetNotFoundPage(writer, request)
	}
}

func GetVoteView(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)

	page := &handler.ViewVotePage{}
	var err error
	page.VoteInstance, page.VoteResults, err = database.GetVoteForUser(request.PathValue("id"), acc)

	if err != nil {
		handler.GetNotFoundPage(writer, request)
	}

	if acc.Exists() {
		page.Voter, err = database.GetOwnedAccountNames(acc)
		page.Voter = append([]string{acc.Name}, page.Voter...)
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostVote(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Diese Funktion ist nicht verfügbar", ElementID: elementID})
		return
	}
	// Todo finish this
}
