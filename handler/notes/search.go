package notes

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"log/slog"
	"net/http"
)

func GetSearchNotePage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := helper.GetAdvancedURLValues(request)

	page := &handler.SearchNotesPage{
		Query:       query.GetTrimmedString("query"),
		Amount:      query.GetInt("amount"),
		ShowBlocked: query.GetBool("blocked"),
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	var backward bool
	page.PreviousItemTime, backward = query.GetUTCTime("backward", false)
	page.NextItemTime, _ = query.GetUTCTime("forward", true)

	var err error
	if backward {
		page.Results, err = database.SearchForNotesBackwards(acc, page.Amount, page.PreviousItemTime, page.Query, page.ShowBlocked)
	} else {
		page.Results, err = database.SearchForNotesForwards(acc, page.Amount, page.NextItemTime, page.Query, page.ShowBlocked)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.TruncatedBlackboardNotes, 0)
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
		page.NextItemTime = page.Results[page.Amount].PostedAt
		page.NextItemID = page.Results[page.Amount].ID
		page.Results = page.Results[:page.Amount]
	} else if backward && len(page.Results) > page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[1].PostedAt
		page.PreviousItemID = page.Results[1].ID
		page.Results = page.Results[1:]
	} else if backward && len(page.Results) > page.Amount {
		amt := len(page.Results) - page.Amount
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[amt].PostedAt
		page.PreviousItemID = page.Results[amt].ID
		page.Results = page.Results[amt:]
	}

	handler.MakeFullPage(writer, acc, page)
}

func PutSearchNotePage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := &handler.SearchNotesPage{
		Query:       values.GetTrimmedString("query"),
		Amount:      values.GetInt("amount"),
		ShowBlocked: values.GetBool("blocked"),
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}
	var backward bool
	page.PreviousItemTime, backward = values.GetUTCTime("backward", false)
	page.NextItemTime, _ = values.GetUTCTime("forward", true)

	if backward {
		page.Results, err = database.SearchForNotesBackwards(acc, page.Amount, page.PreviousItemTime, page.Query, page.ShowBlocked)
	} else {
		page.Results, err = database.SearchForNotesForwards(acc, page.Amount, page.NextItemTime, page.Query, page.ShowBlocked)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.TruncatedBlackboardNotes, 0)
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
		page.NextItemTime = page.Results[page.Amount].PostedAt
		page.NextItemID = page.Results[page.Amount].ID
		page.Results = page.Results[:page.Amount]
	} else if backward && len(page.Results) > page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[1].PostedAt
		page.PreviousItemID = page.Results[1].ID
		page.Results = page.Results[1:]
	} else if backward && len(page.Results) > page.Amount {
		amt := len(page.Results) - page.Amount
		page.HasPrevious = true
		page.PreviousItemTime = page.Results[amt].PostedAt
		page.PreviousItemID = page.Results[amt].ID
		page.Results = page.Results[amt:]
	}

	writer.Header().Add("Hx-Push-Url", "/search/notes?"+values.Encode())
	handler.MakePage(writer, acc, page)
}
