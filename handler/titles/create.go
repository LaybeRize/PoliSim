package titles

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
)

func GetCreateTitlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := handler.CreateTitlePage{}

	handler.MakeFullPage(writer, acc, &page)
}

func PostCreateTitlePage(writer http.ResponseWriter, request *http.Request) {

}
