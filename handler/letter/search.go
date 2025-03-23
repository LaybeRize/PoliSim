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
	}

	page.PossibleAccounts, err = database.GetMyAccountNames(acc)

	if err != nil {
		page.PossibleAccounts = []string{acc.Name}
	}
	if acc.IsAtLeastAdmin() {
		page.PossibleAccounts = append(page.PossibleAccounts, loc.AdministrationAccountName)
	}
	accounts := page.PossibleAccounts

	if acc.IsAtLeastAdmin() && page.Account == loc.AdministrationAccountName {
		accounts = []string{page.Account}
	} else if allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Account); allowed {
		accounts = []string{page.Account}
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	var backward bool
	page.PreviousItemTime, backward = query.GetUTCTime("backward", false)
	page.NextItemTime, _ = query.GetUTCTime("forward", true)

	if backward {
		page.Results, err = database.GetLetterListBackwards(accounts, page.Amount, page.PreviousItemTime)
	} else {
		page.Results, err = database.GetLetterListForwards(accounts, page.Amount, page.NextItemTime)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.ReducedLetter, 0)
	}

	if len(page.Results) > 0 {
		id := query.GetTrimmedString("id")
		if !backward && id == page.Results[0].ID {
			page.HasPrevious = true
			page.PreviousItemTime = page.NextItemTime
			page.PreviousItemID = id
		} else if backward && id == page.Results[len(page.Results)-1].ID {
			page.HasNext = true
			page.NextItemTime = page.PreviousItemTime
			page.NextItemID = id
			page.Results = page.Results[:len(page.Results)-1]
		}
	}

	if !backward && len(page.Results) > page.Amount {
		page.HasNext = true
		page.NextItemTime = page.Results[page.Amount].Written
		page.NextItemID = page.Results[page.Amount].ID
		page.Results = page.Results[:page.Amount]
	} else if backward && len(page.Results) > page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[1].Written
		page.PreviousItemID = page.Results[1].ID
		page.Results = page.Results[1:]
	} else if backward && len(page.Results) > page.Amount {
		amt := len(page.Results) - page.Amount
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[amt].Written
		page.PreviousItemID = page.Results[amt].ID
		page.Results = page.Results[amt:]
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
	}

	page.PossibleAccounts, err = database.GetMyAccountNames(acc)

	if err != nil {
		page.PossibleAccounts = []string{acc.Name}
	}
	if acc.IsAtLeastAdmin() {
		page.PossibleAccounts = append(page.PossibleAccounts, loc.AdministrationAccountName)
	}
	accounts := page.PossibleAccounts

	if acc.IsAtLeastAdmin() && page.Account == loc.AdministrationAccountName {
		accounts = []string{page.Account}
	} else if allowed, _ := database.IsAccountAllowedToPostWith(acc, page.Account); allowed {
		accounts = []string{page.Account}
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	var backward bool
	page.PreviousItemTime, backward = values.GetUTCTime("backward", false)
	page.NextItemTime, _ = values.GetUTCTime("forward", true)

	if backward {
		page.Results, err = database.GetLetterListBackwards(accounts, page.Amount, page.PreviousItemTime)
	} else {
		page.Results, err = database.GetLetterListForwards(accounts, page.Amount, page.NextItemTime)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.ReducedLetter, 0)
	}

	if len(page.Results) > 0 {
		id := values.GetTrimmedString("id")
		if !backward && id == page.Results[0].ID {
			page.HasPrevious = true
			page.PreviousItemTime = page.NextItemTime
			page.PreviousItemID = id
		} else if backward && id == page.Results[len(page.Results)-1].ID {
			page.HasNext = true
			page.NextItemTime = page.PreviousItemTime
			page.NextItemID = id
			page.Results = page.Results[:len(page.Results)-1]
		}
	}

	if !backward && len(page.Results) > page.Amount {
		page.HasNext = true
		page.NextItemTime = page.Results[page.Amount].Written
		page.NextItemID = page.Results[page.Amount].ID
		page.Results = page.Results[:page.Amount]
	} else if backward && len(page.Results) > page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[1].Written
		page.PreviousItemID = page.Results[1].ID
		page.Results = page.Results[1:]
	} else if backward && len(page.Results) > page.Amount {
		amt := len(page.Results) - page.Amount
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[amt].Written
		page.PreviousItemID = page.Results[amt].ID
		page.Results = page.Results[amt:]
	}

	writer.Header().Add("Hx-Push-Url", "/my/letter?"+values.Encode())
	handler.MakePage(writer, acc, page)
}
