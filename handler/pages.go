package handler

import (
	"PoliSim/database"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type FullPage struct {
	Language string
	Base     BaseInfo
	Content  PageStruct
}

type BaseInfo struct {
	Title string
	Icon  string
}

type NavigationInfo struct {
	Account *database.Account
}

type PageStruct interface {
	SetNavInfo(navInfo NavigationInfo)
	getPageName() string
}

type NotFoundPage struct {
	NavInfo NavigationInfo
}

func (p *NotFoundPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *NotFoundPage) getPageName() string {
	return "_notFound"
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

func (p *HomePage) getPageName() string {
	return "_home"
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

func (p *CreateAccountPage) getPageName() string {
	return "accountCreate"
}

type MyProfilePage struct {
	NavInfo  NavigationInfo
	Settings ModifyPersonalSettings
	Password ChangePassword
}

func (p *MyProfilePage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *MyProfilePage) getPageName() string {
	return "_profil"
}

type EditAccountPage struct {
	NavInfo           NavigationInfo
	Account           *database.Account
	LinkedAccountName string
	AccountNames      []string
	AccountUsernames  []string
	MessageUpdate
}

func (p *EditAccountPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *EditAccountPage) getPageName() string {
	return "accountEdit"
}

type NotesPage struct {
	NavInfo       NavigationInfo
	LoadedNoteIDs []string
	LoadedNotes   []*database.BlackboardNote
}

func (p *NotesPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *NotesPage) getPageName() string {
	return "notesView"
}

type CreateNotesPage struct {
	NavInfo         NavigationInfo
	Refrences       string
	Title           string
	Author          string
	PossibleAuthors []string
	Body            string
	MessageUpdate
	MarkdownBox
}

func (p *CreateNotesPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *CreateNotesPage) getPageName() string {
	return "noteCreate"
}

type SearchNotesPage struct {
	NavInfo     NavigationInfo
	Query       string
	Amount      int
	Page        int
	HasNext     bool
	HasPrevious bool
	Results     []database.TruncatedBlackboardNotes
}

func (p *SearchNotesPage) NextPage() int {
	return p.Page + 1
}

func (p *SearchNotesPage) PreviousPage() int {
	return p.Page - 1
}

func (p *SearchNotesPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *SearchNotesPage) getPageName() string {
	return "notesSearch"
}

type CreateTitlePage struct {
	NavInfo      NavigationInfo
	Title        database.Title
	Holder       []string
	AccountNames []string
	MessageUpdate
}

func (p *CreateTitlePage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *CreateTitlePage) getPageName() string {
	return "titleCreate"
}

type EditTitlePage struct {
	NavInfo      NavigationInfo
	Title        *database.Title
	Holder       []string
	AccountNames []string
	Titels       []string
	MessageUpdate
}

func (p *EditTitlePage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *EditTitlePage) getPageName() string {
	return "titleEdit"
}

type CreateOrganisationPage struct {
	NavInfo      NavigationInfo
	Organisation database.Organisation
	User         []string
	Admin        []string
	AccountNames []string
	MessageUpdate
}

func (p *CreateOrganisationPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *CreateOrganisationPage) getPageName() string {
	return "organisationCreate"
}

type EditOrganisationPage struct {
	NavInfo       NavigationInfo
	Organisation  *database.Organisation
	User          []string
	Admin         []string
	AccountNames  []string
	Organisations []string
	MessageUpdate
}

func (p *EditOrganisationPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *EditOrganisationPage) getPageName() string {
	return "organisationEdit"
}

type PartialStruct interface {
	getRenderInfo() (string, string) //first the templateForge key, then the definition name
}

type PartialRedirectStruct interface {
	getRenderInfo() (string, string) //first the templateForge key, then the definition name
	targetElement() string
}

type ChangePassword struct {
	OldPassword       string
	NewPassword       string
	RepeatNewPassword string
	Message           string
	IsError           bool
}

func (p *ChangePassword) getRenderInfo() (string, string) {
	return "_profil", "changeMyPassword"
}

type ModifyPersonalSettings struct {
	FontScaling int64
	Message     string
	IsError     bool
}

func (p *ModifyPersonalSettings) getRenderInfo() (string, string) {
	return "_profil", "changeMySettings"
}

type MarkdownBox struct {
	Information template.HTML
}

func (p *MarkdownBox) getRenderInfo() (string, string) {
	return "templates", "markdownBox"
}

type NotesUpdate struct {
	database.BlackboardNote
}

func (p *NotesUpdate) getRenderInfo() (string, string) {
	return "notesView", "singleNote"
}

type MessageUpdate struct {
	Message string
	IsError bool
}

func (p *MessageUpdate) getRenderInfo() (string, string) {
	return "templates", "message"
}

func (p *MessageUpdate) targetElement() string {
	return "#message-div"
}

//go:embed _pages/*
var pages embed.FS

//go:embed _templates/*
var templates embed.FS

var templateForge map[string]*template.Template = make(map[string]*template.Template)

var iconPath = "/public/fallback_icon.png"

func init() {
	_, _ = fmt.Fprintf(os.Stdout, "Reading All Templates\n")
	templateForge["templates"] = template.Must(template.ParseFS(templates, "*/*"))
	files, err := pages.ReadDir("_pages")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := strings.TrimSuffix(file.Name(), ".gohtml")
		page, pageErr := pages.ReadFile("_pages/" + file.Name())
		if pageErr != nil {
			panic(pageErr)
		}
		templateForge[name] = template.Must(template.Must(template.ParseFS(templates, "*/*")).Parse(string(page)))
	}
	_, _ = fmt.Fprintf(os.Stdout, "Successfully created the Template Forge\n")
	if os.Getenv("ICON_PATH") != "" {
		iconPath = os.Getenv("ICON_PATH")
	}
}

func MakePage(w http.ResponseWriter, acc *database.Account, data PageStruct) {
	navInfo := NavigationInfo{Account: acc}
	data.SetNavInfo(navInfo)

	currentTemplate, exists := templateForge[data.getPageName()]
	if !exists {
		panic("Could not find a template for data. Page required would be: " + data.getPageName())
	}

	err := currentTemplate.ExecuteTemplate(w, "page", data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func MakeFullPage(w http.ResponseWriter, acc *database.Account, data PageStruct) {
	navInfo := NavigationInfo{Account: acc}
	data.SetNavInfo(navInfo)

	fullPage := FullPage{
		Language: "de",
		Base: BaseInfo{
			Icon: iconPath,
		},
		Content: data,
	}

	switch data.(type) {
	case *HomePage:
		fullPage.Base.Title = "Home"
	case *NotFoundPage:
		fullPage.Base.Title = "Seite nicht gefunden"
	case *CreateAccountPage:
		fullPage.Base.Title = "Nutzer erstellen"
	case *MyProfilePage:
		fullPage.Base.Title = "Mein Profil"
	case *EditAccountPage:
		fullPage.Base.Title = "Accounts anpassen"
	case *NotesPage:
		fullPage.Base.Title = "Noitzen anschauen"
	case *CreateNotesPage:
		fullPage.Base.Title = "Notiz erstellen"
	case *SearchNotesPage:
		fullPage.Base.Title = "Notizen durchsuchen"
	case *CreateTitlePage:
		fullPage.Base.Title = "Titel erstellen"
	case *EditTitlePage:
		fullPage.Base.Title = "Titel bearbeiten"
	case *CreateOrganisationPage:
		fullPage.Base.Title = "Organisation erstellen"
	case *EditOrganisationPage:
		fullPage.Base.Title = "Organisation bearbeiten"
	default:
		panic("Struct given to MakeFullPage() is not registered")
	}

	currentTemplate, exists := templateForge[data.getPageName()]
	if !exists {
		panic("Could not find a template for data. Page required would be: " + data.getPageName())
	}
	err := currentTemplate.ExecuteTemplate(w, "fullPage", fullPage)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func MakeSpecialPagePart(w http.ResponseWriter, data PartialStruct) {
	pageName, templateName := data.getRenderInfo()

	currentTemplate, exists := templateForge[pageName]
	if !exists {
		panic("Could not find a template for partial page data. Page required would be: " + pageName)
	}
	err := currentTemplate.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func MakeSpecialPagePartWithRedirect(w http.ResponseWriter, data PartialRedirectStruct) {
	w.Header().Add("HX-Retarget", data.targetElement())
	pageName, templateName := data.getRenderInfo()

	currentTemplate, exists := templateForge[pageName]
	if !exists {
		panic("Could not find a template for partial page data. Page required would be: " + pageName)
	}
	err := currentTemplate.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func RedirectToErrorPage(w http.ResponseWriter) {
	w.Header().Add("HX-Redirect", "/page-not-found")
	w.WriteHeader(http.StatusSeeOther)
}
