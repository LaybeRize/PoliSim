package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"log/slog"
	"net/http"
)

func GetSearchPublicationsPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := helper.GetAdvancedURLValues(request)

	page := &handler.SearchPublicationsPage{
		Query:  query.GetTrimmedString("query"),
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
		page.Results, err = database.GetPublishedNewspaperBackwards(page.Amount, page.PreviousItemTime, page.Query)
	} else {
		page.Results, err = database.GetPublishedNewspaperForwards(page.Amount, page.NextItemTime, page.Query)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.Publication, 0)
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
		page.NextItemTime = page.Results[page.Amount].PublishedDate
		page.NextItemID = page.Results[page.Amount].ID
		page.Results = page.Results[:page.Amount]
	} else if backward && len(page.Results) == page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[0].PublishedDate
		page.PreviousItemID = page.Results[0].ID
	}

	handler.MakeFullPage(writer, acc, page)
}

func PutSearchPublicationPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := &handler.SearchPublicationsPage{
		Query:  values.GetTrimmedString("query"),
		Amount: values.GetInt("amount"),
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	var backward bool
	page.PreviousItemTime, backward = values.GetUTCTime("backward", false)
	page.NextItemTime, _ = values.GetUTCTime("forward", true)

	if backward {
		page.Results, err = database.GetPublishedNewspaperBackwards(page.Amount, page.PreviousItemTime, page.Query)
	} else {
		page.Results, err = database.GetPublishedNewspaperForwards(page.Amount, page.NextItemTime, page.Query)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.Publication, 0)
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
		page.NextItemTime = page.Results[page.Amount].PublishedDate
		page.NextItemID = page.Results[page.Amount].ID
		page.Results = page.Results[:page.Amount]
	} else if backward && len(page.Results) == page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[0].PublishedDate
		page.PreviousItemID = page.Results[0].ID
	}

	writer.Header().Add("Hx-Push-Url", "/search/publications?"+values.Encode())
	handler.MakePage(writer, acc, page)
}
