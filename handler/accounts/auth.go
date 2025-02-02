package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"net/http"
)

func PostLoginAccount(writer http.ResponseWriter, request *http.Request) {
	_, loggedIn := database.RefreshSession(writer, request)

	if loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountNotLoggedIn})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	username := values.GetTrimmedString("username")
	loginAcc, accErr := database.GetAccountByUsername(username)
	if accErr != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountNameOrPasswordWrong})
		return
	}
	correctPassword := database.VerifyPassword(loginAcc.Password, values.GetString("password"))
	if !correctPassword || loginAcc.Role == database.PressUser {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountNameOrPasswordWrong})
		return
	}

	database.CreateSession(writer, loginAcc)
	page := &handler.HomePage{
		Account: loginAcc,
		MessageUpdate: handler.MessageUpdate{
			Message: loc.AccountSuccessfullyLoggedIn,
			IsError: false,
		},
	}
	handler.MakePage(writer, loginAcc, page)
}

func PostLogOutAccount(writer http.ResponseWriter, request *http.Request) {
	_, loggedIn := database.RefreshSession(writer, request)

	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountCurrentlyNotLoggedIn})
		return
	}
	database.EndSession(writer, request)
	page := &handler.HomePage{
		Account: nil,
		MessageUpdate: handler.MessageUpdate{
			Message: loc.AccountSuccessfullyLoggedOut,
			IsError: false,
		},
	}
	handler.MakePage(writer, nil, page)
}
