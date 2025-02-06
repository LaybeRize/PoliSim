package letter

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"log/slog"
	"net/http"
)

func GetPagePersonalLetter(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}
	query := helper.GetAdvancedURLValues(request)

	var err error
	page := &handler.SearchLetterPage{
		Account: query.GetTrimmedString("account"),
		Amount:  query.GetInt("amount"),
		Page:    query.GetInt("page"),
	}

	page.PossibleAccounts, err = database.GetMyAccountNames(acc)

	if err != nil {
		page.PossibleAccounts = []string{acc.Name}
	}
	if acc.IsAtLeastAdmin() {
		page.PossibleAccounts = append(page.PossibleAccounts, loc.AdministrationAccountName)
	}
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	accounts := page.PossibleAccounts

	if acc.IsAtLeastAdmin() && page.Account == loc.AdministrationAccountName {
		accounts = []string{page.Account}
	} else if allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Account); allowed {
		accounts = []string{page.Account}
	}

	page.Results, err = database.GetLetterList(accounts, page.Amount+1, page.Page)
	page.HasPrevious = page.Page > 1
	if err != nil {
		slog.Error(err.Error())
		page.Results = make([]database.ReducedLetter, 0)
	} else if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
	}

	handler.MakeFullPage(writer, acc, page)
}

func PutPagePersonalLetter(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	page := &handler.SearchLetterPage{
		Account: values.GetTrimmedString("account"),
		Amount:  values.GetInt("amount"),
		Page:    values.GetInt("page"),
	}

	page.PossibleAccounts, err = database.GetMyAccountNames(acc)

	if err != nil {
		page.PossibleAccounts = []string{acc.Name}
	}
	if acc.IsAtLeastAdmin() {
		page.PossibleAccounts = append(page.PossibleAccounts, loc.AdministrationAccountName)
	}
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	accounts := page.PossibleAccounts

	if acc.IsAtLeastAdmin() && page.Account == loc.AdministrationAccountName {
		accounts = []string{page.Account}
	} else if allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Account); allowed {
		accounts = []string{page.Account}
	}

	page.Results, err = database.GetLetterList(accounts, page.Amount+1, page.Page)
	page.HasPrevious = page.Page > 1
	if err != nil {
		slog.Error(err.Error())
		page.Results = make([]database.ReducedLetter, 0)
	} else if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
	}

	writer.Header().Add("Hx-Push-Url", "/my/letter?"+values.Encode())
	handler.MakePage(writer, acc, page)
}
