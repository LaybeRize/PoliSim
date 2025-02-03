package notes

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
	"net/url"
)

func GetNotesViewPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := helper.GetAdvancedURLValues(request)

	length := len(query.GetArray("loaded"))
	page := &handler.NotesPage{
		LoadedNoteIDs: make([]string, 0, length),
		LoadedNotes:   make([]*database.BlackboardNote, 0, length),
	}

	for _, loadElement := range query.GetArray("loaded") {
		element, err := database.GetNote(loadElement)
		if err != nil {
			continue
		}
		element.Viewer = acc
		page.LoadedNoteIDs = append(page.LoadedNoteIDs, loadElement)
		page.LoadedNotes = append(page.LoadedNotes, element)
	}

	handler.MakeFullPage(writer, acc, page)
}

func RequestNote(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := helper.GetAdvancedURLValues(request)

	element, err := database.GetNote(query.GetTrimmedString("request"))
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	loaded := query.GetArray("loaded")
	element.Viewer = acc
	arr := append(loaded, element.ID)
	redirectUrl := "/notes?" + url.Values(map[string][]string{"loaded": arr}).Encode()

	writer.Header().Add("Hx-Push-Url", redirectUrl)
	handler.MakeSpecialPagePart(writer, &handler.NotesUpdate{BlackboardNote: *element})
}

func UnBlockNote(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	element, err := database.GetNote(request.PathValue("id"))
	if !acc.IsAtLeastAdmin() || err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	element.Removed = !element.Removed
	err = database.UpdateNoteRemovedStatus(element)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	element.Viewer = acc
	handler.MakeSpecialPagePart(writer, &handler.NotesUpdate{BlackboardNote: *element})
}
