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

type CreateAccountPage struct {
	NavInfo NavigationInfo
	Account database.Account
	Message string
	IsError bool
}

func (p *CreateAccountPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

type MarkdownBox struct {
	Information template.HTML
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
	templateForge = template.Must(template.New("").Funcs(template.FuncMap{
		"isUser": func(acc *database.Account) bool {
			return acc != nil
		},
		"isPressAdmin": func(acc *database.Account) bool {
			if acc == nil {
				return false
			}
			return acc.Role <= database.PRESS_ADMIN
		},
		"isAdmin": func(acc *database.Account) bool {
			if acc == nil {
				return false
			}
			return acc.Role <= database.ADMIN
		},
		"isHeadAdmin": func(acc *database.Account) bool {
			if acc == nil {
				return false
			}
			return acc.Role <= database.HEAD_ADMIN
		},
	}).ParseFiles(filenames...))
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
	case *CreateAccountPage:
		executeTemplate(w, "createAccount", data)
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
	case *CreateAccountPage:
		fullPage.Base.Title = "Nutzer erstellen"
		executeTemplate(w, "createAccountFull", fullPage)
	default:
		panic("Struct given to MakeFullPage() is not registered")
	}
}

type SpecialPage string

const (
	MARKDOWN SpecialPage = "markdownBox"
)

func MakeSpecialPagePart(w http.ResponseWriter, page SpecialPage, data any) {
	executeTemplate(w, string(page), data)
}

func executeTemplate(w http.ResponseWriter, name string, data any) {
	err := templateForge.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
