package handler

import (
	"PoliSim/database"
	"log"
	"net/http"
)

type HomePage struct {
	Base BaseInfo
	Info LoginInfo
}

var HomeTemplate = ParsePage("home")

func (p *HomePage) Execute(w http.ResponseWriter) {
	p.Base.Title = "Home"
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

func (p *NotFoundPage) Execute(w http.ResponseWriter) {
	p.Base.Title = "Seite nicht gefunden"
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
	base := BaseInfo{Account: acc, LoggedIn: loggedIn}

	page := HomePage{Base: base, Info: LoginInfo{ErrorMessage: ""}}
	if acc != nil {
		page.Info.AccountName = acc.Name
	}
	page.Execute(writer)
}

func GetNotFoundPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	base := BaseInfo{Account: acc, LoggedIn: loggedIn}

	page := NotFoundPage{Base: base}
	page.Execute(writer)
}
