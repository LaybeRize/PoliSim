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
		Query: &database.LetterSearch{
			Title:            query.GetTrimmedString("title"),
			ExactTitleMatch:  query.GetBool("match-title"),
			Author:           query.GetTrimmedString("author"),
			ExactAuthorMatch: query.GetBool("match-author"),
			ShowOnlyUnread:   query.GetBool("only-unread"),
		},
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
	recName := query.GetTrimmedString("rec-name")

	if backward {
		page.Results, err = database.GetLetterListBackwards(accounts, page.Amount, page.PreviousItemTime, recName, page.Query)
	} else {
		page.Results, err = database.GetLetterListForwards(accounts, page.Amount, page.NextItemTime, recName, page.Query)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.ReducedLetter, 0)
	}

	if len(page.Results) > 0 {
		id := query.GetTrimmedString("id")
		if !backward && id == page.Results[0].ID && page.Results[0].Recipient == recName {
			page.HasPrevious = true
			page.PreviousItemTime = page.NextItemTime
			page.PreviousItemID = id
			page.PreviousItemRec = page.Results[0].Recipient
		} else if lst := len(page.Results) - 1; backward && id == page.Results[lst].ID && page.Results[lst].Recipient == recName {
			page.HasNext = true
			page.NextItemTime = page.PreviousItemTime
			page.NextItemID = id
			page.NextItemRec = page.Results[lst].Recipient
			page.Results = page.Results[:lst]
		}
	}

	if !backward && len(page.Results) > page.Amount {
		page.HasNext = true
		page.NextItemTime = page.Results[page.Amount].Written
		page.NextItemID = page.Results[page.Amount].ID
		page.NextItemRec = page.Results[page.Amount].Recipient
		page.Results = page.Results[:page.Amount]
	} else if backward && len(page.Results) > page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[1].Written
		page.PreviousItemID = page.Results[1].ID
		page.PreviousItemRec = page.Results[1].Recipient
		page.Results = page.Results[1:]
	} else if backward && len(page.Results) > page.Amount {
		amt := len(page.Results) - page.Amount
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[amt].Written
		page.PreviousItemID = page.Results[amt].ID
		page.PreviousItemRec = page.Results[amt].Recipient
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
		Query: &database.LetterSearch{
			Title:            values.GetTrimmedString("title"),
			ExactTitleMatch:  values.GetBool("match-title"),
			Author:           values.GetTrimmedString("author"),
			ExactAuthorMatch: values.GetBool("match-author"),
			ShowOnlyUnread:   values.GetBool("only-unread"),
		},
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
	recName := values.GetTrimmedString("rec-name")

	if backward {
		page.Results, err = database.GetLetterListBackwards(accounts, page.Amount, page.PreviousItemTime, recName, page.Query)
	} else {
		page.Results, err = database.GetLetterListForwards(accounts, page.Amount, page.NextItemTime, recName, page.Query)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.ReducedLetter, 0)
	}

	if len(page.Results) > 0 {
		id := values.GetTrimmedString("id")
		if !backward && id == page.Results[0].ID && page.Results[0].Recipient == recName {
			page.HasPrevious = true
			page.PreviousItemTime = page.NextItemTime
			page.PreviousItemID = id
			page.PreviousItemRec = page.Results[0].Recipient
		} else if lst := len(page.Results) - 1; backward && id == page.Results[lst].ID && page.Results[lst].Recipient == recName {
			page.HasNext = true
			page.NextItemTime = page.PreviousItemTime
			page.NextItemID = id
			page.NextItemRec = page.Results[lst].Recipient
			page.Results = page.Results[:lst]
		}
	}

	if !backward && len(page.Results) > page.Amount {
		page.HasNext = true
		page.NextItemTime = page.Results[page.Amount].Written
		page.NextItemID = page.Results[page.Amount].ID
		page.NextItemRec = page.Results[page.Amount].Recipient
		page.Results = page.Results[:page.Amount]
	} else if backward && len(page.Results) > page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[1].Written
		page.PreviousItemID = page.Results[1].ID
		page.PreviousItemRec = page.Results[1].Recipient
		page.Results = page.Results[1:]
	} else if backward && len(page.Results) > page.Amount {
		amt := len(page.Results) - page.Amount
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[amt].Written
		page.PreviousItemID = page.Results[amt].ID
		page.PreviousItemRec = page.Results[amt].Recipient
		page.Results = page.Results[amt:]
	}

	values.DeleteEmptyFields([]string{"title", "author", "account"})
	writer.Header().Add("Hx-Push-Url", "/my/letter?"+values.Encode())
	handler.MakePage(writer, acc, page)
}
