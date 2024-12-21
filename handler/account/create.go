package account

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
	"strconv"
)

func GetCreateAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HEAD_ADMIN {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := handler.CreateAccountPage{}

	handler.MakeFullPage(writer, acc, &page)
}

func PostCreateAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HEAD_ADMIN {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}
	page := handler.CreateAccountPage{IsError: true}

	err := request.ParseForm()
	if err != nil {
		page.Message = "Fehler beim parsen der Informationen"
		handler.MakePage(writer, acc, &page)
		return
	}

	page.Account.Name = request.Form.Get("name")
	page.Account.Username = request.Form.Get("username")
	page.Account.Password = request.Form.Get("password")
	role, err := strconv.Atoi(request.Form.Get("role"))
	page.Account.Role = database.AccountRole(role)

	if page.Account.Name == "" || len(page.Account.Name) > 200 {
		page.Message = "Der Anzeigename des Accounts ist entweder leer oder überschreitet das 200 Zeichenlimit"
		handler.MakePage(writer, acc, &page)
		return
	}

	if page.Account.Username == "" || len(page.Account.Username) > 200 {
		page.Message = "Der Nutzername des Accounts ist entweder leer oder überschreitet das 200 Zeichenlimit"
		handler.MakePage(writer, acc, &page)
		return
	}

	if len(page.Account.Password) < 10 && page.Account.Role != database.PRESS_USER {
		page.Message = "Das Password hat weniger als 10 Zeichen"
		handler.MakePage(writer, acc, &page)
		return
	}

	if err != nil || page.Account.Role < database.HEAD_ADMIN || page.Account.Role > database.PRESS_USER {
		page.Message = "Die ausgewählte Rolle für den Nutzer ist nicht valide"
		handler.MakePage(writer, acc, &page)
		return
	}

	if page.Account.Role <= acc.Role {
		page.Message = "Du bist nicht berechtigt eine Account mit den selben oder höheren Berechtigungen zu erstellen"
		handler.MakePage(writer, acc, &page)
		return
	}

	newAccount := page.Account
	newAccount.Password, err = database.HashPassword(newAccount.Password)
	if err != nil {
		page.Message = "Es ist ein Fehler beim hashen des Passworts aufgetreten"
		handler.MakePage(writer, acc, &page)
		return
	}

	err = database.CreateAccount(&newAccount)
	if err != nil {
		page.Message = "Der Nutzer konnte nicht erstellt werden\nBitte überprüfe ob Anzeigename oder Nutzername einzigartig sind"
		handler.MakePage(writer, acc, &page)
		return
	}

	page = handler.CreateAccountPage{IsError: false, Message: "Account erfolgreich erstellt\nDer Nutzername ist: " + page.Account.Username + "\nDas Passwort ist: " + page.Account.Password}
	handler.MakePage(writer, acc, &page)
}
