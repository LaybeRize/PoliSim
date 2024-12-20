package handler

import (
	"PoliSim/database"
	"net/http"
)

func GetHomePage(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		GetNotFoundPage(writer, request)
		return
	}

	acc, _ := database.RefreshSession(writer, request)
	page := HomePage{
		Account: acc,
		Message: "",
		IsError: false,
	}
	MakeFullPage(writer, acc, &page)
}

func GetNotFoundPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	page := NotFoundPage{}
	MakeFullPage(writer, acc, &page)
}

func PartialGetNotFoundPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	page := NotFoundPage{}
	MakePage(writer, acc, &page)
}
