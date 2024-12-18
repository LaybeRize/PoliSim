package handler

import (
	"PoliSim/database"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type FullPage struct {
	Base    BaseInfo
	Content any
}

type BaseInfo struct {
	Title    string
	Language string
	Icon     string
}

type NavigationInfo struct {
	Account  *database.Account
	LoggedIn bool
}

type PageStruct interface {
	SetNavInfo(navInfo NavigationInfo)
}

type NotFoundPage struct {
	NavInfo NavigationInfo
}

func (p *NotFoundPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

type HomePage struct {
	NavInfo NavigationInfo
	Account *database.Account
	Message string
	IsError bool
}

func (p *HomePage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

var templateForge *template.Template = nil

func init() {
	_, _ = fmt.Fprintf(os.Stdout, "Reading All Templates\n")
	files, err := os.ReadDir("templates")
	if err != nil {
		panic(err)
	}
	filenames := make([]string, len(files))
	for i, file := range files {
		filenames[i] = "templates/" + file.Name()
	}
	templateForge = template.Must(template.ParseFiles(filenames...))
	_, _ = fmt.Fprintf(os.Stdout, "Successfully created the Template Forge\n")
}

func MakePage(w http.ResponseWriter, acc *database.Account, data PageStruct) {
	navInfo := NavigationInfo{
		Account:  acc,
		LoggedIn: acc != nil,
	}
	data.SetNavInfo(navInfo)
	switch data.(type) {
	case *HomePage:
		executeTemplate(w, "home", data)
	case *NotFoundPage:
		executeTemplate(w, "notFound", data)
	default:
		panic("Struct given to MakePage() is not registered")
	}
}

func MakeFullPage(w http.ResponseWriter, acc *database.Account, data PageStruct) {
	navInfo := NavigationInfo{
		Account:  acc,
		LoggedIn: acc != nil,
	}
	data.SetNavInfo(navInfo)

	fullPage := FullPage{
		Base: BaseInfo{
			Title:    "",
			Language: "de",
			Icon:     "fallback_icon.png",
		},
		Content: data,
	}

	switch data.(type) {
	case *HomePage:
		fullPage.Base.Title = "Home"
		executeTemplate(w, "homeFull", fullPage)
	case *NotFoundPage:
		fullPage.Base.Title = "Seite nicht gefunden"
		executeTemplate(w, "notFoundFull", fullPage)
	default:
		panic("Struct given to MakeFullPage() is not registered")
	}
}

func executeTemplate(w http.ResponseWriter, name string, data any) {
	err := templateForge.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
