package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

//Todo: make this

func GetSearchPublicationsPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	query := request.URL.Query()

	page := &handler.SearchNotesPage{
		Query: query.Get("query"),
	}
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
	page.Results, err = database.SearchForNotes(page.Amount+1, page.Page, page.Query)
	if err != nil {
		page.Results = make([]database.TruncatedBlackboardNotes, 0)
	}
	if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
	}
	handler.MakeFullPage(writer, acc, page)
}

func PutSearchPublicationPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if err := request.ParseForm(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := &handler.SearchNotesPage{}
	page.Query = helper.GetFormEntry(request, "query")
	page.Amount, _ = strconv.Atoi(helper.GetFormEntry(request, "amount"))
	page.Page, _ = strconv.Atoi(helper.GetFormEntry(request, "page"))
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	page.HasPrevious = page.Page > 1
	var err error
	page.Results, err = database.SearchForNotes(page.Amount+1, page.Page, page.Query)
	if err != nil {
		page.Results = make([]database.TruncatedBlackboardNotes, 0)
	}
	if len(page.Results) > page.Amount {
		page.HasNext = true
		page.Results = page.Results[:page.Amount]
	}
	writer.Header().Add("Hx-Push-Url", "/search/notes?query="+url.QueryEscape(page.Query)+
		fmt.Sprintf("&amount=%d&page=%d", page.Amount, page.Page))
	handler.MakePage(writer, acc, page)
}
