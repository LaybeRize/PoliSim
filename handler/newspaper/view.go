package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"log/slog"
	"net/http"
)

func GetSpecificPublicationPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	pubID := request.PathValue("id")
	found, err := database.GetPublicationForUser(pubID, acc.IsAtLeastPressAdmin())
	if !found || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.ViewPublicationPage{}
	var pub *database.Publication
	pub, page.Articles, err = database.GetPublication(pubID)
	if page.QueryError = err != nil; !page.QueryError {
		page.Publication = *pub
	} else {
		slog.Error(err.Error())
	}

	handler.MakeFullPage(writer, acc, page)
}
