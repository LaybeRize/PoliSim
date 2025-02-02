package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"errors"
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

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError, ElementID: elementID})
		return
	}

	tag := &database.DocumentTag{
		ID:              helper.GetUniqueID(acc.Name),
		Text:            values.GetTrimmedString("text"),
		BackgroundColor: values.GetTrimmedString("background-color"),
		TextColor:       values.GetTrimmedString("text-color"),
		LinkColor:       values.GetTrimmedString("link-color"),
		Links:           values.GetCommaSeperatedArray("links"),
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

func PatchRemoveDocument(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastAdmin() {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	database.RemoveRestoreDocument(request.PathValue("id"))

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

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	voter := values.GetTrimmedString("voter")
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
	if values.GetBool("invalid") {
		votesCasted = nil

	} else if voteType.IsSingleVote() {
		votesCasted = make([]int, len(answers))
		pos := values.GetInt("vote")

		if pos <= 0 || pos > len(answers) {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Die Abgegebene Stimme ist invalide" + "\n" + "Die ausgewählte Position der Antwort ist nicht gültig"})
			return
		}
		votesCasted[pos-1] = 1

	} else if voteType.IsMultipleVotes() {
		votesCasted = make([]int, len(answers))
		for i := range len(answers) {
			if values.GetBool(fmt.Sprintf("vote-%d", i+1)) {
				votesCasted[i] = 1
			}
		}

	} else if voteType.IsVoteSharing() {
		votesCasted = make([]int, len(answers))
		sum := 0

		for i := range len(answers) {
			amount := values.GetInt(fmt.Sprintf("vote-%d", i+1))
			if amount < 0 {
				handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
					Message: "Die Abgegebene Stimme ist invalide" + "\n" + "Die Anzahl an Stimmen pro Antwort darf nicht kleiner als 0 sein"})
				return
			}
			votesCasted[i] = amount
			sum += amount
		}

		if sum > maxVotes {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Die Abgegebene Stimme ist invalide" + "\n" + "Die Summe aller abgegebenen Stimmen überschreitet das festgelegte Maximum"})
			return
		}

	} else if voteType.IsRankedVoting() {
		votesCasted = make([]int, len(answers))
		lookUpMap := make(map[int]interface{})

		for i := range len(answers) {
			pos := values.GetInt(fmt.Sprintf("vote-%d", i+1))
			if pos <= 0 {
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
				lookUpMap[pos] = struct{}{}
				votesCasted[i] = pos
			}
		}

	}

	err = database.CastVoteWithAccount(voter, id, votesCasted)
	if errors.Is(err, database.AlreadyVoted) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Mit dem Account wurde bereits abgestimmt"})
		return
	} else if err != nil {
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
