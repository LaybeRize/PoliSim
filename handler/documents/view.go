package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
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

func GetVoteView(writer http.ResponseWriter, request *http.Request) {
	// Todo make this
}

func PostVote(writer http.ResponseWriter, request *http.Request) {
	// Todo: make this
}
