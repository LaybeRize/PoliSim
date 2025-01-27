package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
)

func GetCreateAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := handler.CreateAccountPage{}

	handler.MakeFullPage(writer, acc, &page)
}

func PostCreateAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	newAccount := &database.Account{
		Name:     values.GetTrimmedString("name"),
		Username: values.GetTrimmedString("username"),
		Password: values.GetString("password"),
		Role:     database.AccountRole(values.GetInt("role")),
	}

	if newAccount.Name == "" || len(newAccount.Name) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der Anzeigename des Accounts ist entweder leer oder überschreitet das 200 Zeichenlimit"})
		return
	}

	if newAccount.Username == "" || len(newAccount.Username) > 200 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der Nutzername des Accounts ist entweder leer oder überschreitet das 200 Zeichenlimit"})
		return
	}

	if len(newAccount.Password) < 10 && newAccount.Role != database.PressUser {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Das Password hat weniger als 10 Zeichen"})
		return
	}

	if newAccount.Role < database.HeadAdmin || newAccount.Role > database.PressUser {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die ausgewählte Rolle für den Nutzer ist nicht valide"})
		return
	}

	if newAccount.Role <= acc.Role {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Du bist nicht berechtigt eine Account mit den selben oder höheren Berechtigungen zu erstellen"})
		return
	}

	newAccount.Password, err = database.HashPassword(newAccount.Password)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim Hashen des Passworts aufgetreten"})
		return
	}

	err = database.CreateAccount(newAccount)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Der Nutzer konnte nicht erstellt werden\nBitte überprüfe ob Anzeigename oder Nutzername einzigartig sind"})
		return
	}

	page := &handler.CreateAccountPage{MessageUpdate: handler.MessageUpdate{
		IsError: false,
		Message: "Account erfolgreich erstellt\nDer Nutzername ist: " + newAccount.Username + "\nDas Passwort ist: " + newAccount.Password,
	}}
	handler.MakePage(writer, acc, page)
}
