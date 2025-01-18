package letter

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

func GetPagePersonalLetter(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}
	query := request.URL.Query()

	var err error
	page := &handler.SearchLetterPage{
		Account: query.Get("account"),
	}
	page.Amount, _ = strconv.Atoi(query.Get("amount"))
	page.Page, _ = strconv.Atoi(query.Get("page"))
	page.PossibleAccounts, err = database.GetMyAccountNames(acc)

	if err != nil {
		page.PossibleAccounts = []string{acc.Name}
	}
	if acc.IsAtLeastAdmin() {
		page.PossibleAccounts = append(page.PossibleAccounts, loc.AdminstrationAccountName)
	}
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	accounts := page.PossibleAccounts

	if acc.IsAtLeastAdmin() && page.Account == loc.AdminstrationAccountName {
		accounts = []string{page.Account}
	} else if allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Account); allowed {
		accounts = []string{page.Account}
	}
	slog.Debug("Accounts", "list", accounts, "account", page.Account)

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

	err := request.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	page := &handler.SearchLetterPage{}
	page.Account = helper.GetFormEntry(request, "account")
	database.GetIntegerFormEntry(request, "amount", &page.Amount)
	database.GetIntegerFormEntry(request, "page", &page.Page)
	page.PossibleAccounts, err = database.GetMyAccountNames(acc)

	if err != nil {
		page.PossibleAccounts = []string{acc.Name}
	}
	if acc.IsAtLeastAdmin() {
		page.PossibleAccounts = append(page.PossibleAccounts, loc.AdminstrationAccountName)
	}
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	accounts := page.PossibleAccounts

	if acc.IsAtLeastAdmin() && page.Account == loc.AdminstrationAccountName {
		accounts = []string{page.Account}
	} else if allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Account); allowed {
		accounts = []string{page.Account}
	}
	slog.Debug("Accounts", "list", accounts, "account", page.Account)

	page.Results, err = database.GetLetterList(accounts, page.Amount+1, page.Page)
	page.HasPrevious = page.Page > 1
	if err != nil {
		slog.Error(err.Error())
		page.Results = make([]database.ReducedLetter, 0)
	} else if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
	}

	writer.Header().Add("Hx-Push-Url", "/my/letter?account="+url.QueryEscape(page.Account)+
		fmt.Sprintf("&amount=%d&page=%d", page.Amount, page.Page))
	handler.MakePage(writer, acc, page)
}
