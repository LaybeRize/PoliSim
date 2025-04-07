package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
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
	query := helper.GetAdvancedURLValues(request)
	var err error

	if query.Has("name") {
		var ownerAccount *database.Account
		page.Account, ownerAccount, err = database.GetAccountAndOwnerByAccountName(query.GetTrimmedString("name"))

		if err != nil {
			handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
				Message: loc.AccountSearchedNameDoesNotCorrespond})
			return
		}
		if page.Account.Role == database.PressUser {
			if ownerAccount.Exists() {
				page.LinkedAccountName = ownerAccount.Name
			}
			page.AccountNames, err = database.GetNamesForActiveUsers()
			if err != nil {
				page.Message = loc.AccountErrorFindingNamesForOwner
				handler.MakeFullPage(writer, acc, page)
				return
			}
		}

		page.IsError = false
		page.Message = loc.AccountFoundSearchedName
		handler.MakeFullPage(writer, acc, page)
		return
	}

	page.AccountNames, page.AccountUsernames, err = database.GetNames()
	if err != nil {
		page.IsError = true
		page.Message = loc.AccountErrorSearchingNameList
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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.EditAccountPage{Account: nil}
	role := database.AccountRole(values.GetInt("role"))

	var ownerAccount *database.Account
	page.Account, ownerAccount, err = database.GetAccountAndOwnerByAccountName(values.GetTrimmedString("name"))
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountErrorNoAccountToModify})
		return
	}

	if ownerAccount.Exists() {
		page.LinkedAccountName = ownerAccount.Name
	}

	if page.Account.Role <= acc.Role {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountNoPermissionToEdit})
		return
	}

	// First checks if the account is not a PressUser because then changing roles is not allowed
	// and then if the role is valid (lower boundary is set, upper boundary is set by modifying users role)
	if page.Account.Role != database.PressUser && role <= database.User && role > acc.Role {
		page.Account.Role = role
	} else if page.Account.Role != database.PressUser {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountRoleIsNotAllowed})
		return
	}

	page.Account.Blocked = values.GetBool("blocked")
	err = database.UpdateAccount(page.Account)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountErrorWhileUpdating})
		return
	}

	if ownerAccount, err = database.GetAccountByName(values.GetTrimmedString("linked")); err == nil &&
		page.Account.Role == database.PressUser && !page.Account.Blocked {
		if ownerAccount.Role == database.PressUser {
			handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
				Message: loc.AccountPressUserOwnerIsPressUser})
			return
		}

		err = database.UpdateOwner(ownerAccount.Name, page.Account.Name)
		if err != nil {
			handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
				Message: loc.AccountPressUserOwnerChangeError})
			return
		}
		page.LinkedAccountName = ownerAccount.Name
	}

	if page.Account.Role == database.PressUser && values.GetTrimmedString("linked") == "" && !page.Account.Blocked {
		err = database.UpdateOwner(page.Account.Name, page.Account.Name)
		if err != nil {
			handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
				Message: loc.AccountPressUserOwnerRemovingError})
			return
		}
		page.LinkedAccountName = ""
	} else if page.Account.Blocked {
		page.LinkedAccountName = ""
	}

	page.IsError = false
	page.Message = loc.AccountSuccessfullyUpdated
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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	accountByName, nameErr := database.GetAccountByName(values.GetTrimmedString("name"))
	accountByUsername, usernameErr := database.GetAccountByUsername(values.GetTrimmedString("username"))
	var name string

	switch true {
	case nameErr != nil && usernameErr != nil:
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountSearchedNamesDoesNotCorrespond})
		return
	//goland:noinspection ALL
	case accountByName.Exists():
		name = accountByName.Name
	//goland:noinspection ALL
	case accountByUsername.Exists():
		name = accountByUsername.Name
	}

	writer.Header().Add("HX-Redirect", "/edit/account?name="+url.QueryEscape(name))
	writer.WriteHeader(http.StatusFound)
}
