package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"log/slog"
	"net/http"
)

func GetSearchDocumentsPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := helper.GetAdvancedURLValues(request)

	page := &handler.SearchDocumentsPage{
		Amount: query.GetInt("amount"),
		Query: &database.DocumentSearch{
			ShowBlocked:             query.GetBool("blocked"),
			DocumentTitle:           query.GetTrimmedString("title"),
			ExactDocumentTitleMatch: query.GetBool("match-title"),
			Organisation:            query.GetTrimmedString("organisation"),
			ExactOrganisationMatch:  query.GetBool("match-organisation"),
			OrganisationName:        query.GetTrimmedString("organisation-name"),
			AuthorName:              query.GetTrimmedString("author"),
			ExactAuthorMatch:        query.GetBool("match-author"),
		},
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	var backward bool
	page.PreviousItemTime, backward = query.GetUTCTime("backward", false)
	page.NextItemTime, _ = query.GetUTCTime("forward", true)

	var err error
	if backward {
		page.Results, err = database.GetDocumentListBackwards(page.Amount, page.PreviousItemTime, acc, page.Query)
	} else {
		page.Results, err = database.GetDocumentListForwards(page.Amount, page.NextItemTime, acc, page.Query)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.SmallDocument, 0)
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

func PutSearchDocumentsPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := &handler.SearchDocumentsPage{
		Amount: values.GetInt("amount"),
		Query: &database.DocumentSearch{
			ShowBlocked:             values.GetBool("blocked"),
			DocumentTitle:           values.GetTrimmedString("title"),
			ExactDocumentTitleMatch: values.GetBool("match-title"),
			Organisation:            values.GetTrimmedString("organisation"),
			ExactOrganisationMatch:  values.GetBool("match-organisation"),
			OrganisationName:        values.GetTrimmedString("organisation-name"),
			AuthorName:              values.GetTrimmedString("author"),
			ExactAuthorMatch:        values.GetBool("match-author"),
		},
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	var backward bool
	page.PreviousItemTime, backward = values.GetUTCTime("backward", false)
	page.NextItemTime, _ = values.GetUTCTime("forward", true)

	if backward {
		page.Results, err = database.GetDocumentListBackwards(page.Amount, page.PreviousItemTime, acc, page.Query)
	} else {
		page.Results, err = database.GetDocumentListForwards(page.Amount, page.NextItemTime, acc, page.Query)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.SmallDocument, 0)
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

	values.DeleteEmptyFields([]string{"title", "organisation", "author"})
	writer.Header().Add("Hx-Push-Url", "/search/documents?"+values.Encode())
	handler.MakePage(writer, acc, page)
}

func GetPersonalSearchDocumentsPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := helper.GetAdvancedURLValues(request)

	page := &handler.SearchPersonalDocumentsPage{
		Amount: query.GetInt("amount"),
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	var backward bool
	page.PreviousItemTime, backward = query.GetUTCTime("backward", false)
	page.NextItemTime, _ = query.GetUTCTime("forward", true)

	var err error
	if backward {
		page.Results, err = database.GetPersonalDocumentListBackwards(page.Amount, page.PreviousItemTime, acc)
	} else {
		page.Results, err = database.GetPersonalDocumentListForwards(page.Amount, page.NextItemTime, acc)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.SmallDocument, 0)
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

func PutPersonalSearchDocumentsPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := &handler.SearchPersonalDocumentsPage{
		Amount: values.GetInt("amount"),
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	var backward bool
	page.PreviousItemTime, backward = values.GetUTCTime("backward", false)
	page.NextItemTime, _ = values.GetUTCTime("forward", true)

	if backward {
		page.Results, err = database.GetPersonalDocumentListBackwards(page.Amount, page.PreviousItemTime, acc)
	} else {
		page.Results, err = database.GetPersonalDocumentListForwards(page.Amount, page.NextItemTime, acc)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.SmallDocument, 0)
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

	writer.Header().Add("Hx-Push-Url", "/my/documents?"+values.Encode())
	handler.MakePage(writer, acc, page)
}
