package notes

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
)

func GetSearchNotePage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	handler.MakeFullPage(writer, acc, &handler.SearchNotesPage{})
}

func PutSearchNotePage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	handler.MakePage(writer, acc, &handler.SearchNotesPage{})
}
