package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
	"net/url"
)

func GetEditAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.EditAccountPage{Account: nil}
	values := helper.GetAdvancedURLValues(request)
	var err error

	if values.Has("name") {
		var ownerAccount *database.Account
		page.Account, ownerAccount, err = database.GetAccountAndOwnerByAccountName(values.GetTrimmedString("name"))

		if err != nil {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: "Der gesuchte Name ist mit keinem Account verbunden"})
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

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	page := &handler.EditAccountPage{Account: nil}
	role := database.AccountRole(values.GetInt("role"))

	var ownerAccount *database.Account
	page.Account, ownerAccount, err = database.GetAccountAndOwnerByAccountName(values.GetTrimmedString("name"))
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

	// First checks if the account is not a PressUser because then changing roles is not allowed
	// and then if the role is valid (lower boundary is set, upper boundary is set by modifying users role)
	if page.Account.Role != database.PressUser && role <= database.User && role > acc.Role {
		page.Account.Role = role
	} else if page.Account.Role != database.PressUser {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die ausgewählte Rolle ist nicht valide"})
		return
	}

	page.Account.Blocked = values.GetBool("blocked")
	err = database.UpdateAccount(page.Account)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Es ist ein Fehler beim updaten des Accounts aufgetreten"})
		return
	}

	if ownerAccount, err = database.GetAccountByName(values.GetTrimmedString("linked")); err == nil &&
		page.Account.Role == database.PressUser {
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

	if page.Account.Role == database.PressUser && values.GetTrimmedString("linked") == "" {
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

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	accountByName, nameErr := database.GetAccountByName(values.GetTrimmedString("name"))
	accountByUsername, usernameErr := database.GetAccountByUsername(values.GetTrimmedString("username"))
	var name string

	switch true {
	case nameErr != nil && usernameErr != nil:
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte keinen Account finden, der den Informationen entspricht"})
		return
	case accountByName.Exists():
		name = accountByName.Name
	case accountByUsername.Exists():
		name = accountByUsername.Name
	}

	writer.Header().Add("HX-Redirect", "/edit/account?name="+url.QueryEscape(name))
	writer.WriteHeader(http.StatusFound)
}
