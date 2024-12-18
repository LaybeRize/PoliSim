package account

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
)

func PostLoginAccount(writer http.ResponseWriter, request *http.Request) {
	_, loggedIn := database.RefreshSession(writer, request)
	info := handler.LoginInfo{AccountName: "", ErrorMessage: "Du bist bereits angemeldet"}
	if loggedIn {
		info.Execute(writer)
		return
	}

	err := request.ParseForm()
	if err != nil {
		info.ErrorMessage = "Fehler beim parsen der Informationen"
		info.Execute(writer)
		return
	}

	username := request.Form.Get("username")
	acc, accErr := database.GetAccountByUsername(username)
	correctPassword := database.VerifyPassword(acc.Password, request.Form.Get("password"))
	if accErr != nil || !correctPassword {
		info.ErrorMessage = "Nutzername oder Passwort falsch"
		info.Execute(writer)
		return
	}

	database.CreateSession(writer, acc)
	info.AccountName = acc.Name
	info.ErrorMessage = ""
	info.Execute(writer)
}

func PostLogOutAccount(writer http.ResponseWriter, request *http.Request) {
	_, loggedIn := database.RefreshSession(writer, request)
	info := handler.LoginInfo{AccountName: "", ErrorMessage: "Du bist nicht angemeldet"}
	if !loggedIn {
		info.Execute(writer)
		return
	}
	database.EndSession(writer, request)
	info.ErrorMessage = "Erfolgreich ausgeloggt"
	info.Execute(writer)
}
