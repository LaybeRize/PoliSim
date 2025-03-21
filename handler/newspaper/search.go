package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"log/slog"
	"net/http"
)

// Todo use timestamps for paging in the future

func GetSearchPublicationsPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := helper.GetAdvancedURLValues(request)

	page := &handler.SearchPublicationsPage{
		Query:  query.GetTrimmedString("query"),
		Amount: query.GetInt("amount"),
		Page:   query.GetInt("page"),
	}

	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	page.HasPrevious = page.Page > 1
	var err error
	page.Results, err = database.GetPublishedNewspaper(page.Amount, page.Page, page.Query)
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.Publication, 0)
	}
	if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
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
		Page:   values.GetInt("page"),
	}

	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	page.HasPrevious = page.Page > 1
	page.Results, err = database.GetPublishedNewspaper(page.Amount, page.Page, page.Query)
	if err != nil {
		slog.Debug(err.Error())
		page.Results = make([]database.Publication, 0)
	}
	if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
	}
	writer.Header().Add("Hx-Push-Url", "/search/publications?"+values.Encode())
	handler.MakePage(writer, acc, page)
}
