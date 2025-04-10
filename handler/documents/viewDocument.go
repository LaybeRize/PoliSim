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

func GetDocumentVoteResults(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	csv, err := database.GetVoteCSVForDocument(request.PathValue("id"), acc)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.Header().Add("Content-Type", "text/csv;charset=utf-8")
	writer.Header().Add("Content-Disposition", `attachment; filename="results.csv"`)
	_, _ = writer.Write([]byte(csv))
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

func PatchRemoveComment(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastAdmin() {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	database.RemoveRestoreComment(request.PathValue("comment"))

	if obj := getDocumentPageObject(acc, request); obj != nil {
		handler.MakePage(writer, acc, obj)
	} else {
		handler.GetNotFoundPage(writer, request)
	}
}
