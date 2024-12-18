package handler

import (
	"PoliSim/database"
	"log"
	"net/http"
)

type BaseInfo struct {
	Account  *database.Account
	LoggedIn bool
	Title    string
}

type HomePage struct {
	Base BaseInfo
}

var HomeTemplate = ParsePage("home")

func (p *HomePage) execute(w http.ResponseWriter) {
	err := HomeTemplate.Execute(w, p)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type NotFoundPage struct {
	Base BaseInfo
}

var NotFoundTemplate = ParsePage("notFound")

func (p *NotFoundPage) execute(w http.ResponseWriter) {
	err := NotFoundTemplate.Execute(w, p)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetHomePage(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		GetNotFoundPage(writer, request)
		return
	}

	acc, loggedIn := database.RefreshSession(writer, request)
	base := BaseInfo{Account: acc, LoggedIn: loggedIn, Title: "Home"}

	page := HomePage{Base: base}
	page.execute(writer)
}

func GetNotFoundPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	base := BaseInfo{Account: acc, LoggedIn: loggedIn, Title: "Seite nicht gefunden"}

	page := NotFoundPage{Base: base}
	page.execute(writer)
}
