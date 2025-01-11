package letter

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"fmt"
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
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	accounts := page.PossibleAccounts

	allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Account)
	if allowed {
		accounts = []string{page.Account}
	}

	page.Results, err = database.GetLetterList(accounts, page.Amount+1, page.Page)
	page.HasPrevious = page.Page > 1
	if err != nil {
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

	var err error
	page := &handler.SearchLetterPage{}
	page.Account = helper.GetFormEntry(request, "account")
	page.Amount, _ = strconv.Atoi(helper.GetFormEntry(request, "amount"))
	page.Page, _ = strconv.Atoi(helper.GetFormEntry(request, "page"))
	page.PossibleAccounts, err = database.GetMyAccountNames(acc)

	if err != nil {
		page.PossibleAccounts = []string{acc.Name}
	}
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	accounts := page.PossibleAccounts

	allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Account)
	if allowed {
		accounts = []string{page.Account}
	}

	page.Results, err = database.GetLetterList(accounts, page.Amount+1, page.Page)
	page.HasPrevious = page.Page > 1
	if err != nil {
		page.Results = make([]database.ReducedLetter, 0)
	} else if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
	}

	writer.Header().Add("Hx-Push-Url", "/my/letter?account="+url.QueryEscape(page.Account)+
		fmt.Sprintf("&amount=%d&page=%d", page.Amount, page.Page))
	handler.MakePage(writer, acc, page)
}
