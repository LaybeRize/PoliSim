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

	page := &handler.EditAccountPage{Account: nil}
	var err error

	if accountName, exists := request.URL.Query()["name"]; exists {
		page.Account, err = database.GetAccountByName(accountName[0])

		if err != nil {
			page.Account = nil
			page.IsError = true
			page.Message = "Der gesuchte Name ist mit keinem Account verbunden"
			page.AccountNames, page.AccountUsernames, err = database.GetNames()
			if err != nil {
				page.Message += "\nEs ist ein Fehler bei der Suche nach den Namenslisten aufgetreten"
			}
			handler.MakeFullPage(writer, acc, page)
			return
		}

		handler.MakeFullPage(writer, acc, page)
		return
	}

	page.AccountNames, page.AccountUsernames, err = database.GetNames()
	if err != nil {
		page.IsError = true
		page.Message = "Es ist ein Fehler bei der Suche nach den Namenslisten aufgetreten"
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostEditAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	handler.MakePage(writer, acc, &handler.EditAccountPage{})
}

func PostEditSearchAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	handler.MakePage(writer, acc, &handler.EditAccountPage{})
}
