package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InstallDocumentText() {
	composition.PageTitleMap[composition.CreateTextDocument] = builder.Translation["documentTextCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateTextDocument] = builder.Translation["documentTextCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateTextDocument] = GetDocumentTextCreationService
	composition.PostHTMXFunctions[composition.CreateTextDocument] = PostDocumentTextCreationService
	composition.PageTitleMap[composition.ViewTextDocument] = builder.Translation["documentTextViewPageTitle"]
	composition.GetHTMXFunctions[composition.ViewTextDocument] = GetDocumentTextViewService
}

func GetDocumentTextViewService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(r)

	html := composition.ViewDocumentPage(acc, chi.URLParam(r, "uuid"))
	viewTextDocumentRenderRequest(w, r, acc, html)
}

var viewTextDocumentRenderRequest = genericRenderer(composition.ViewTextDocument)

func PostDocumentTextCreationService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	create := &validation.CreateDocument{}
	msg := validation.Message{Positive: false}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createTextDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.CreateDocument(acc.ID)

	if !msg.Positive {
		createTextDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	w.Header().Set("HX-Push-Url", "/"+string(composition.ViewTextDocumentLink)+create.UUIDredirect)
	html := composition.ViewDocumentPage(acc, create.UUIDredirect)
	viewTextDocumentRenderRequest(w, r, acc, html)
}

func GetDocumentTextCreationService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.CreateDocumentPage(acc, &validation.CreateDocument{}, validation.Message{})
	createTextDocumentRenderRequest(w, r, acc, html)
}

var createTextDocumentRenderRequest = genericRenderer(composition.CreateTextDocument)
var createTextDocumentOnlySwapMessage = genericMessageSwapper(composition.CreateTextDocument)
