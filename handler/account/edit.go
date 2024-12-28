package account

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
)

func GetEditAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	handler.MakeFullPage(writer, acc, &handler.EditAccountPage{})
}

func PostEditAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	handler.MakePage(writer, acc, &handler.EditAccountPage{})
}
