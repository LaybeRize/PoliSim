package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"fmt"
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
			Message: "Diese Funktion ist nicht verfügbar"})
		return
	}

	err := request.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen"})
		return
	}

	voter := helper.GetFormEntry(request, "voter")
	allowed, err := database.IsAccountAllowedToPostWith(acc, voter)
	if !allowed || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehlende Berechtigung um mit diesem Account abzustimmen"})
		return
	}
	id := request.PathValue("id")
	answers, voteType, maxVotes, err := database.GetAnswersAndTypeForVote(id, acc)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Für diese Abstimmung kann keine Stimme abgegeben werden"})
		return
	}

	var votesCasted []int
	if helper.GetBoolFormEntry(request, "invalid") {
		votesCasted = nil

	} else if voteType.IsSingleVote() {
		votesCasted = make([]int, len(answers))
		var pos int
		database.GetIntegerFormEntry(request, "vote", &pos)

		if pos <= 0 || pos > len(answers) {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Die Abgegebene Stimme ist invalide" + "\n" + "Die ausgewählte Position der Antwort ist nicht gültig"})
			return
		}
		votesCasted[pos-1] = 1

	} else if voteType.IsMultipleVotes() {
		votesCasted = make([]int, len(answers))
		for i := range len(answers) {
			if helper.GetBoolFormEntry(request, fmt.Sprintf("vote-%d", i+1)) {
				votesCasted[i] = 1
			}
		}

	} else if voteType.IsVoteSharing() {
		votesCasted = make([]int, len(answers))
		var amount int64
		sum := int64(0)

		for i := range len(answers) {
			database.GetIntegerFormEntry(request, fmt.Sprintf("vote-%d", i+1), &amount)
			if amount < 0 {
				handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
					Message: "Die Abgegebene Stimme ist invalide" + "\n" + "Die Anzahl an Stimmen pro Antwort darf nicht kleiner als 0 sein"})
				return
			}
			votesCasted[i] = int(amount)
			sum += amount
		}

		if sum > maxVotes {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Die Abgegebene Stimme ist invalide" + "\n" + "Die Summe aller abgegebenen Stimmen überschreitet das festgelegte Maximum"})
			return
		}

	} else if voteType.IsRankedVoting() {
		votesCasted = make([]int, len(answers))
		var pos int
		lookUpMap := make(map[int]interface{})

		for i := range len(answers) {
			database.GetIntegerFormEntry(request, fmt.Sprintf("vote-%d", i+1), &pos)
			if pos < 0 {
				votesCasted[i] = -1
			} else if pos > len(answers) {
				handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
					Message: "Die Abgegebene Stimme ist invalide" + "\n" + "Einer der Ränge ist größer als maximal erlaubt"})
				return
			} else {
				_, exists := lookUpMap[pos]
				if exists {
					handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
						Message: "Die Abgegebene Stimme ist invalide" + "\n" + "Der selbe Rang darf nicht doppelt vergeben werden"})
					return
				}
				votesCasted[i] = pos
			}
		}

	}

	err = database.CastVoteWithAccount(voter, id, votesCasted)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Versuch die Stimme abzugeben\nÜberprüfe ob der Account stimmberechtigt ist"})
		return
	}

	page := &handler.ViewVotePage{}
	page.VoteInstance, page.VoteResults, err = database.GetVoteForUser(request.PathValue("id"), acc)
	if err != nil {
		handler.PartialGetNotFoundPage(writer, request)
	}

	page.IsError = false
	page.Message = "Stimme erfolgreich abgegeben"

	if acc.Exists() {
		page.Voter, err = database.GetOwnedAccountNames(acc)
		page.Voter = append([]string{acc.Name}, page.Voter...)
	}

	handler.MakePage(writer, acc, page)
}
