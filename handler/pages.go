package handler

import (
	"PoliSim/database"
	loc "PoliSim/localisation"
	"embed"
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
	return "personalProfil"
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
	References      string
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

type ViewTitlePage struct {
	NavInfo        NavigationInfo
	TitleHierarchy map[string]map[string][]string
}

func (p *ViewTitlePage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *ViewTitlePage) getPageName() string {
	return "titleView"
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

type ViewOrganisationPage struct {
	NavInfo   NavigationInfo
	Hierarchy map[string]map[string][]database.Organisation
	HadError  bool
}

func (p *ViewOrganisationPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *ViewOrganisationPage) getPageName() string {
	return "organisationView"
}

type ManageNewspaperPage struct {
	NavInfo        NavigationInfo
	Newspaper      database.Newspaper
	AccountNames   []string
	NewspaperNames []string
	Publications   []database.Publication
	HadError       bool
	MessageUpdate
}

func (p *ManageNewspaperPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *ManageNewspaperPage) getPageName() string {
	return "newspaperManage"
}

func (p *ManageNewspaperPage) getRenderInfo() (string, string) {
	return "newspaperManage", "updateNewspaper"
}

type CreateArticlePage struct {
	NavInfo           NavigationInfo
	Title             string
	Subtitle          string
	Author            string
	PossibleAuthors   []string
	PossibleNewspaper []string
	Special           bool
	Body              string
	MessageUpdate
	MarkdownBox
}

func (p *CreateArticlePage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *CreateArticlePage) getPageName() string {
	return "newspaperCreate"
}

func (p *CreateArticlePage) getRenderInfo() (string, string) {
	return "newspaperCreate", "newspaperDropdown"
}

type ViewPublicationPage struct {
	NavInfo     NavigationInfo
	QueryError  bool
	Publication database.Publication
	Articles    []database.NewspaperArticle
	MessageUpdate
}

func (p *ViewPublicationPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *ViewPublicationPage) getPageName() string {
	return "newspaperPubView"
}

type SearchPublicationsPage struct {
	NavInfo     NavigationInfo
	Query       string
	Amount      int
	Page        int
	HasNext     bool
	HasPrevious bool
	Results     []database.Publication
}

func (p *SearchPublicationsPage) NextPage() int {
	return p.Page + 1
}

func (p *SearchPublicationsPage) PreviousPage() int {
	return p.Page - 1
}

func (p *SearchPublicationsPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *SearchPublicationsPage) getPageName() string {
	return "newspaperSearch"
}

type SearchLetterPage struct {
	NavInfo          NavigationInfo
	Account          string
	PossibleAccounts []string
	Amount           int
	Page             int
	HasNext          bool
	HasPrevious      bool
	Results          []database.ReducedLetter
}

func (p *SearchLetterPage) NextPage() int {
	return p.Page + 1
}

func (p *SearchLetterPage) PreviousPage() int {
	return p.Page - 1
}

func (p *SearchLetterPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *SearchLetterPage) getPageName() string {
	return "letterPersonalSearch"
}

type CreateLetterPage struct {
	NavInfo         NavigationInfo
	Title           string
	Author          string
	PossibleAuthors []string
	Recipients      []string
	AccountNames    []string
	Signable        bool
	Body            string
	MarkdownBox
	MessageUpdate
}

func (p *CreateLetterPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *CreateLetterPage) getPageName() string {
	return "letterCreate"
}

type ViewLetterPage struct {
	NavInfo NavigationInfo
	database.Letter
}

func (p *ViewLetterPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *ViewLetterPage) getPageName() string {
	return "letterView"
}

type AdminSearchLetterPage struct {
	NavInfo          NavigationInfo
	AccountNameToUse string
}

func (p *AdminSearchLetterPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *AdminSearchLetterPage) getPageName() string {
	return "letterAdminSearch"
}

type DocumentViewPage struct {
	NavInfo       NavigationInfo
	Document      *database.Document
	Commentator   []string
	ColorPalettes map[string]database.ColorPalette
	MessageUpdate
	MarkdownBox
}

func (p *DocumentViewPage) ColorIterator() func(func(int, database.ColorPalette) bool) {
	return func(yield func(int, database.ColorPalette) bool) {
		pos := 0
		for _, palette := range p.ColorPalettes {
			pos += 1
			if !yield(pos, palette) {
				return
			}
		}
	}
}

func (p *DocumentViewPage) CanComment() bool {
	return len(p.Commentator) != 0 && !p.Document.Ended()
}

func (p *DocumentViewPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *DocumentViewPage) getPageName() string {
	return "documentView"
}

type ViewVotePage struct {
	NavInfo      NavigationInfo
	VoteInstance *database.VoteInstance
	VoteResults  *database.AccountVotes
	Voter        []string
	MessageUpdate
}

func (p *ViewVotePage) CanVote() bool {
	return len(p.Voter) != 0 && !p.VoteInstance.Ended()
}

func (p *ViewVotePage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *ViewVotePage) getPageName() string {
	return "documentViewVote"
}

type CreateDocumentPage struct {
	NavInfo         NavigationInfo
	Title           string
	Author          string
	PossibleAuthors []string
	Body            string
	UpdateOrganisationForUser
	MessageUpdate
	MarkdownBox
}

func (p *CreateDocumentPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *CreateDocumentPage) getPageName() string {
	return "documentCreate"
}

type CreateDiscussionPage struct {
	NavInfo              NavigationInfo
	Title                string
	Author               string
	DateTime             string
	MaxTime              string
	MinTime              string
	PossibleAuthors      []string
	AccountNames         []string
	Public               bool
	LetMemberParticipate bool
	LetAdminsParticipate bool
	Body                 string
	ReaderAndParticipants
	UpdateOrganisationForUser
	MessageUpdate
	MarkdownBox
}

func (p *CreateDiscussionPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *CreateDiscussionPage) getPageName() string {
	return "documentCreateDiscussion"
}

type CreateVotePage struct {
	NavInfo              NavigationInfo
	Title                string
	Author               string
	DateTime             string
	MaxTime              string
	MinTime              string
	PossibleAuthors      []string
	AccountNames         []string
	Public               bool
	LetMemberParticipate bool
	LetAdminsParticipate bool
	Body                 string
	VoteChoice           []database.VoteInfo
	VoteChoiceList       string
	ReaderAndParticipants
	UpdateOrganisationForUser
	MessageUpdate
	MarkdownBox
}

func (p *CreateVotePage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *CreateVotePage) getPageName() string {
	return "documentCreateVote"
}

type CreateVoteElementPage struct {
	NavInfo     NavigationInfo
	VoteNumbers []int
	CurrNumber  int
	Vote        *database.VoteInstance
	MessageUpdate
}

func (p *CreateVoteElementPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *CreateVoteElementPage) getPageName() string {
	return "documentCreateVoteElement"
}

type SearchDocumentsPage struct {
	NavInfo     NavigationInfo
	Amount      int
	Page        int
	HasNext     bool
	HasPrevious bool
	Results     []database.SmallDocument
}

func (p *SearchDocumentsPage) NextPage() int {
	return p.Page + 1
}

func (p *SearchDocumentsPage) PreviousPage() int {
	return p.Page - 1
}

func (p *SearchDocumentsPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *SearchDocumentsPage) getPageName() string {
	return "documentSearch"
}

type EditColorPage struct {
	NavInfo         NavigationInfo
	AllowedToCreate bool
	AllowedToDelete bool
	Color           database.ColorPalette
	ColorPalettes   map[string]database.ColorPalette
	MessageUpdate
}

func (p *EditColorPage) SetNavInfo(navInfo NavigationInfo) {
	p.NavInfo = navInfo
}

func (p *EditColorPage) getPageName() string {
	return "documentColorEdit"
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
	return "personalProfil", "changeMyPassword"
}

type ModifyPersonalSettings struct {
	FontScaling int64
	TimeZone    string
	Message     string
	IsError     bool
}

func (p *ModifyPersonalSettings) getRenderInfo() (string, string) {
	return "personalProfil", "changeMySettings"
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

type SingleTitelUpdate struct {
	Found  bool
	Title  string
	Flair  string
	Holder string
}

func (p *SingleTitelUpdate) HasFlair() bool {
	return p.Flair != ""
}

func (p *SingleTitelUpdate) getRenderInfo() (string, string) {
	return "titleView", "singleTitle"
}

type SingleOrganisationUpdate struct {
	Name         string
	Organisation *database.Organisation
	User         string
	Admin        string
}

func (p *SingleOrganisationUpdate) getRenderInfo() (string, string) {
	return "organisationView", "singleOrganisation"
}

type UpdateOrganisationForUser struct {
	PossibleOrganisations []string
}

func (p *UpdateOrganisationForUser) getRenderInfo() (string, string) {
	return "templates", "organisationDropdown"
}

type ReaderAndParticipants struct {
	Reader       []string
	Participants []string
}

func (p *ReaderAndParticipants) getRenderInfo() (string, string) {
	return "templates", "readerAndParticipants"
}

type MessageUpdate struct {
	ElementID string
	Message   string
	IsError   bool
}

func (p *MessageUpdate) getRenderInfo() (string, string) {
	return "templates", "message"
}

func (p *MessageUpdate) targetElement() string {
	if p.ElementID != "" {
		return "#" + p.ElementID
	}
	return "#message-div"
}

func (p *MessageUpdate) GetElementID() string {
	if p.ElementID != "" {
		return p.ElementID
	}
	return "message-div"
}

//go:embed _pages/*
var pages embed.FS

//go:embed _templates/*
var templates embed.FS

var templateForge = make(map[string]*template.Template)

var iconPath = "/public/fallback_icon.png"

func init() {
	log.Println("Reading All Templates")
	templateString := getTemplatesAsSingleString()
	templateForge["templates"] = template.Must(template.New("").Parse(templateString))
	files, err := pages.ReadDir("_pages")
	if err != nil {
		log.Fatalf("Pages read directory error: %v", err)
	}
	for _, file := range files {
		name := strings.TrimSuffix(file.Name(), ".gohtml")
		page, pageErr := pages.ReadFile("_pages/" + file.Name())
		if pageErr != nil {
			log.Fatalf("page read content error: %v", pageErr)
		}
		templateForge[name] = template.Must(
			template.Must(template.New("").Parse(templateString)).Parse(
				loc.LocaliseTemplateString(page)))
	}
	log.Println("Successfully created the Template Forge")
	if os.Getenv("ICON_PATH") != "" {
		iconPath = os.Getenv("ICON_PATH")
	}
}

func getTemplatesAsSingleString() string {
	files, err := templates.ReadDir("_templates")
	if err != nil {
		log.Fatalf("Templates read directory error: %v", err)
	}
	arr := make([]string, len(files))
	for i, file := range files {
		temp, templateErr := templates.ReadFile("_templates/" + file.Name())
		if templateErr != nil {
			log.Fatalf("template read content error: %v", templateErr)
		}
		arr[i] = loc.LocaliseTemplateString(temp)
	}
	return strings.Join(arr, "\n")
}

func MakePage(w http.ResponseWriter, acc *database.Account, data PageStruct) {
	navInfo := NavigationInfo{Account: acc}
	data.SetNavInfo(navInfo)

	currentTemplate, exists := templateForge[data.getPageName()]
	if !exists {
		log.Fatalf("Could not find a template for data. Page required would be: %s", data.getPageName())
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
		fullPage.Base.Title = "Notizen anschauen"
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
	case *ViewTitlePage:
		fullPage.Base.Title = "Titelübersicht"
	case *ViewOrganisationPage:
		fullPage.Base.Title = "Organisationsübersicht"
	case *ManageNewspaperPage:
		fullPage.Base.Title = "Zeitungen verwalten"
	case *CreateArticlePage:
		fullPage.Base.Title = "Artikel erstellen"
	case *ViewPublicationPage:
		fullPage.Base.Title = "Zeitung"
	case *SearchPublicationsPage:
		fullPage.Base.Title = "Zeitungen durchsuchen"
	case *SearchLetterPage:
		fullPage.Base.Title = "Briefe durchsuchen"
	case *CreateLetterPage:
		fullPage.Base.Title = "Brief erstellen"
	case *AdminSearchLetterPage:
		fullPage.Base.Title = "Briefsuche mit ID"
	case *ViewLetterPage:
		fullPage.Base.Title = "Briefansicht"
	case *DocumentViewPage:
		fullPage.Base.Title = "Dokumentansicht"
	case *CreateDocumentPage:
		fullPage.Base.Title = "Dokument erstellen"
	case *CreateDiscussionPage:
		fullPage.Base.Title = "Diskussion erstellen"
	case *CreateVoteElementPage:
		fullPage.Base.Title = "Abstimmungen verwalten"
	case *CreateVotePage:
		fullPage.Base.Title = "Abstimmungsdokument erstellen"
	case *SearchDocumentsPage:
		fullPage.Base.Title = "Dokumente durchsuchen"
	case *ViewVotePage:
		fullPage.Base.Title = "Abstimmungsansicht"
	case *EditColorPage:
		fullPage.Base.Title = "Farbpaletten anpassen"
	default:
		log.Fatalf("Struct of type %T given to MakeFullPage() is not registered", data)
	}

	currentTemplate, exists := templateForge[data.getPageName()]
	if !exists {
		log.Fatalf("Could not find a template for data. Page required would be: %s", data.getPageName())
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
		log.Fatalf("Could not find a template for data. Page required would be: %s", pageName)
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
		log.Fatalf("Could not find a template for data. Page required would be: %s", pageName)
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
