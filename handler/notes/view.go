package notes

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
	"strings"
)

func GetNotesViewPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	length := len(request.URL.Query()["loaded"])
	page := &handler.NotesPage{
		LoadedNoteIDs: make([]string, 0, length),
		LoadedNotes:   make([]*database.BlackboardNote, 0, length),
	}

	for _, loadElement := range request.URL.Query()["loaded"] {
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
	element, err := database.GetNote(request.URL.Query().Get("request"))
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	loaded, exists := request.URL.Query()["loaded"]
	if !exists {
		loaded = make([]string, 0)
	}
	element.Viewer = acc
	arr := append(loaded, element.ID)
	url := "/notes?"
	for _, e := range arr {
		url += "loaded=" + e + "&"
	}
	writer.Header().Add("Hx-Push-Url", strings.TrimSuffix(url, "&"))
	handler.MakeSpecialPagePart(writer, &handler.NotesUpdate{BlackboardNote: *element})
}
