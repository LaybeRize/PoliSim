package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
)

func PostLoginAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	page := handler.HomePage{
		Account: acc,
		Message: "Du bist bereits angemeldet",
		IsError: true,
	}
	if loggedIn {
		handler.MakePage(writer, acc, &page)
		return
	}

	err := request.ParseForm()
	if err != nil {
		page.Message = "Fehler beim parsen der Informationen"
		handler.MakePage(writer, acc, &page)
		return
	}

	username := helper.GetFormEntry(request, "username")
	loginAcc, accErr := database.GetAccountByUsername(username)
	page.Message = "Nutzername oder Passwort falsch"
	if accErr != nil {
		handler.MakePage(writer, acc, &page)
		return
	}
	correctPassword := database.VerifyPassword(loginAcc.Password, helper.GetPureFormEntry(request, "password"))
	if !correctPassword || loginAcc.Role == database.PressUser {
		handler.MakePage(writer, acc, &page)
		return
	}

	database.CreateSession(writer, loginAcc)
	page = handler.HomePage{
		Account: loginAcc,
		Message: "Erfolgreich angemeldet",
		IsError: false,
	}
	handler.MakePage(writer, loginAcc, &page)
}

func PostLogOutAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	page := handler.HomePage{
		Account: acc,
		Message: "Du bist nicht angemeldet",
		IsError: true,
	}
	if !loggedIn {
		handler.MakePage(writer, acc, &page)
		return
	}
	database.EndSession(writer, request)
	page.Account = nil
	page.Message = "Erfolgreich ausgeloggt"
	page.IsError = false
	handler.MakePage(writer, nil, &page)
}
