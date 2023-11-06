package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InstallDocumentDiscussion() {
	composition.PageTitleMap[composition.CreateDiscussionDocument] = builder.Translation["documentDiscussionCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateDiscussionDocument] = builder.Translation["documentDiscussionCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateDiscussionDocument] = GetDocumentDiscussionCreationService
	composition.PostHTMXFunctions[composition.CreateDiscussionDocument] = PostDocumentDiscussionCreationService
	composition.PageTitleMap[composition.ViewDiscussionDocument] = builder.Translation["documentDiscussionViewPageTitle"]
	composition.GetHTMXFunctions[composition.ViewDiscussionDocument] = GetDocumentDiscussionViewService
}

func GetDocumentDiscussionViewService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(r)

	html := composition.ViewDiscussionPage(acc, chi.URLParam(r, "uuid"))
	viewDiscussionDocumentRenderRequest(w, r, acc, html)
}

var viewDiscussionDocumentRenderRequest = genericRenderer(composition.ViewDiscussionDocument)

func PostDocumentDiscussionCreationService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	create := &validation.CreateDiscussion{}
	msg := validation.Message{Positive: false}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createDiscussionDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.CreateDiscussion(acc.ID)

	if !msg.Positive {
		createDiscussionDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	w.Header().Set("HX-Push-Url", "/"+string(composition.ViewDiscussionDocumentLink)+create.UUIDredirect)
	html := composition.ViewDiscussionPage(acc, create.UUIDredirect)
	viewDiscussionDocumentRenderRequest(w, r, acc, html)
}

func GetDocumentDiscussionCreationService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.CreateDiscussionPage(acc, &validation.CreateDiscussion{}, validation.Message{})
	createDiscussionDocumentRenderRequest(w, r, acc, html)
}

var createDiscussionDocumentRenderRequest = genericRenderer(composition.CreateDiscussionDocument)
var createDiscussionDocumentOnlySwapMessage = genericMessageSwapper(composition.CreateDiscussionDocument)
