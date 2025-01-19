package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"fmt"
	"net/http"
	"strconv"
)

func GetSearchDocumentsPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := request.URL.Query()

	page := &handler.SearchDocumentsPage{}
	page.Amount, _ = strconv.Atoi(query.Get("amount"))
	page.Page, _ = strconv.Atoi(query.Get("page"))
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	page.HasPrevious = page.Page > 1
	var err error
	page.Results, err = database.GetDocumentList(page.Amount+1, page.Page, acc)
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
	if err := request.ParseForm(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := &handler.SearchDocumentsPage{}
	database.GetIntegerFormEntry(request, "amount", &page.Amount)
	database.GetIntegerFormEntry(request, "page", &page.Page)
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	page.HasPrevious = page.Page > 1
	var err error
	page.Results, err = database.GetDocumentList(page.Amount+1, page.Page, acc)
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
