package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
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

const (
	maxLengthNames    = 200
	minLengthPassword = 10
)

func PostCreateAccount(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	newAccount := &database.Account{
		Name:     values.GetTrimmedString("name"),
		Username: values.GetTrimmedString("username"),
		Password: values.GetString("password"),
		Role:     database.AccountRole(values.GetInt("role")),
	}

	if newAccount.Name == "" || len(newAccount.Name) > maxLengthNames {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.AccountDisplayNameTooLongOrNotAtAll, maxLengthNames)})
		return
	}

	if newAccount.Username == "" || len(newAccount.Username) > maxLengthNames {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.AccountUsernameTooLongOrNotAtAll, maxLengthNames)})
		return
	}

	if len(newAccount.Password) < minLengthPassword && newAccount.Role != database.PressUser {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.AccountPasswordTooShort, minLengthPassword)})
		return
	}

	if newAccount.Role < database.HeadAdmin || newAccount.Role > database.PressUser {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountSelectedInvalidRole})
		return
	}

	if newAccount.Role <= acc.Role {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountNotAllowedToCreateAccountOfThatRank})
		return
	}

	clearTextPassword := newAccount.Password
	newAccount.Password, err = database.HashPassword(newAccount.Password)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountPasswordHashingError})
		return
	}

	err = database.CreateAccount(newAccount)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountCreationError})
		return
	}

	page := &handler.CreateAccountPage{MessageUpdate: handler.MessageUpdate{
		IsError: false,
		Message: fmt.Sprintf(loc.AccountSuccessfullyCreated, newAccount.Username, clearTextPassword),
	}}
	handler.MakePage(writer, acc, page)
}
