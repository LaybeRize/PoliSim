package account

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
	"net/url"
	"strconv"
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
		var ownerAccount *database.Account
		page.Account, ownerAccount, err = database.GetAccountAndOwnerByAccountName(accountName[0])

		page.IsError = true
		if err != nil {
			print(err.Error())
			page.Account = nil
			page.Message = "Der gesuchte Name ist mit keinem Account verbunden"
			page.AccountNames, page.AccountUsernames, err = database.GetNames()
			if err != nil {
				page.Message += "\nEs ist ein Fehler bei der Suche nach den Namenslisten aufgetreten"
			}
			handler.MakeFullPage(writer, acc, page)
			return
		}
		if page.Account.Role == database.PressUser {
			if ownerAccount.Exists() {
				page.LinkedAccountName = ownerAccount.Name
			}
			page.AccountNames, err = database.GetNamesForActiveUsers()
			if err != nil {
				page.Message = "Konnte Namen für mögliche Accountbesitzer nicht laden"
				handler.MakeFullPage(writer, acc, page)
				return
			}
		}

		page.IsError = false
		page.Message = "Gesuchten Account gefunden"
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

	err := request.ParseForm()
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	page := &handler.EditAccountPage{Account: nil}
	var role int

	var ownerAccount *database.Account
	page.Account, ownerAccount, err = database.GetAccountAndOwnerByAccountName(request.Form.Get("name"))
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte keinen Account zum modifizieren finden"})
		return
	}

	if ownerAccount.Exists() {
		page.LinkedAccountName = ownerAccount.Name
	}

	if page.Account.Role <= acc.Role {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Du besitzt nicht die Berechtigung diesen Account anzupassen"})
		return
	}

	role, err = strconv.Atoi(request.Form.Get("role"))
	// First checks if the account is not a PressUser because then changing roles is not allowed
	// and then if the role is valid (lower boundary is set, upper boundary is set by modifying users role)
	if page.Account.Role != database.PressUser && role <= int(database.User) && role > int(acc.Role) {
		page.Account.Role = database.AccountRole(role)
	} else if page.Account.Role != database.PressUser {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die ausgewählte Rolle ist nicht valide"})
		return
	}

	page.Account.Blocked = "true" == request.Form.Get("blocked")
	err = database.UpdateAccount(page.Account)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim updaten des Accounts aufgetreten"})
		return
	}

	if ownerAccount, err = database.GetAccountByName(request.Form.Get("linked")); err == nil && page.Account.Role == database.PressUser {
		if ownerAccount.Role == database.PressUser {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Ein Presse-Nutzer kann kein Besitzer eines anderen Presse-Nutzers sein"})
			return
		}

		err = database.RemoveOwner(page.Account.Name)
		if err != nil {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Es ist ein Fehler beim entferne des bisherigen Besitzers aufgetreten"})
			return
		}
		err = database.MakeOwner(ownerAccount.Name, page.Account.Name)
		if err != nil {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Es ist ein Fehler beim entferne des bisherigen Besitzers aufgetreten"})
			return
		}
		page.LinkedAccountName = ownerAccount.Name
	}

	if page.Account.Role == database.PressUser && request.Form.Get("linked") == "" {
		err = database.RemoveOwner(page.Account.Name)
		if err != nil {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Es ist ein Fehler beim entferne des bisherigen Besitzers aufgetreten"})
			return
		}
		page.LinkedAccountName = ""
	}

	page.IsError = false
	page.Message = "Account erfolgreich angepasst"
	handler.MakePage(writer, acc, page)
}

func PostEditSearchAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	page := &handler.EditAccountPage{MessageUpdate: handler.MessageUpdate{IsError: true}}

	err := request.ParseForm()
	if err != nil {
		page.Message = "Fehler beim parsen der Informationen"
		makeEditSearchPage(writer, acc, page)
		return
	}

	accountByName, nameErr := database.GetAccountByName(request.Form.Get("name"))
	accountByUsername, usernameErr := database.GetAccountByUsername(request.Form.Get("username"))
	var name string

	switch true {
	case nameErr != nil && usernameErr != nil:
		page.Message = "Konnte keinen Account finden, der den Informationen entspricht"
		makeEditSearchPage(writer, acc, page)
		return
	case accountByName.Exists():
		name = accountByName.Name
	case accountByUsername.Exists():
		name = accountByUsername.Name
	}

	writer.Header().Add("HX-Redirect", "/edit/account?name="+url.QueryEscape(name))
	writer.WriteHeader(http.StatusFound)
}

func makeEditSearchPage(writer http.ResponseWriter, acc *database.Account, page *handler.EditAccountPage) {
	var err error
	page.AccountNames, page.AccountUsernames, err = database.GetNames()
	if err != nil {
		page.Message += "\nEs ist ein Fehler bei der Suche nach den Namenslisten aufgetreten"
	}
	handler.MakePage(writer, acc, page)
}
