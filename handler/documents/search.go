package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"fmt"
	"net/http"
)

func GetSearchDocumentsPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := helper.GetAdvancedURLValues(request)

	page := &handler.SearchDocumentsPage{
		Amount:      query.GetInt("amount"),
		Page:        query.GetInt("page"),
		ShowBlocked: query.GetBool("blocked"),
	}

	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	page.HasPrevious = page.Page > 1
	var err error
	page.Results, err = database.GetDocumentList(page.Amount+1, page.Page, acc, page.ShowBlocked)
	if err != nil {
		page.Results = make([]database.SmallDocument, 0)
	}
	if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
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
		Amount:      values.GetInt("amount"),
		Page:        values.GetInt("page"),
		ShowBlocked: values.GetBool("blocked"),
	}

	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	page.HasPrevious = page.Page > 1
	page.Results, err = database.GetDocumentList(page.Amount+1, page.Page, acc, page.ShowBlocked)
	if err != nil {
		page.Results = make([]database.SmallDocument, 0)
	}
	if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
	}
	writer.Header().Add("Hx-Push-Url", "/search/documents"+
		fmt.Sprintf("?amount=%d&page=%d", page.Amount, page.Page))
	handler.MakePage(writer, acc, page)
}
