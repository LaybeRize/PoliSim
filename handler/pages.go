package handler

import (
	"PoliSim/database"
	"html/template"
	"io"
	"log"
	"net/http"
)

type SpecialTemplate struct {
	Template *template.Template
	Name     string
}

type BaseInfo struct {
	Account  *database.Account
	LoggedIn bool
	Title    string
}

var baseTemplate = template.Must(template.ParseFiles("templates/templates.gohtml"))

func ParsePage(name string) *SpecialTemplate {
	return &SpecialTemplate{Template: template.Must(template.ParseFiles("templates/templates.gohtml", "templates/"+name+".gohtml")),
		Name: name + ".gohtml"}
}

func (t *SpecialTemplate) Execute(w io.Writer, data any) error {
	return t.Template.ExecuteTemplate(w, t.Name, data)
}

type LoginInfo struct {
	AccountName  string
	ErrorMessage string
}

func (data LoginInfo) Execute(w http.ResponseWriter) {
	err := baseTemplate.ExecuteTemplate(w, "login", data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
