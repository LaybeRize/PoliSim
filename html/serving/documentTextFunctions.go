package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"
)

func InstallDocumentText() {
	composition.PageTitleMap[composition.CreateTextDocument] = builder.Translation["documentTextCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateTextDocument] = builder.Translation["documentTextCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateTextDocument] = GetDocumentTextCreationService
	composition.PostHTMXFunctions[composition.CreateTextDocument] = PostDocumentTextCreationService
	composition.PageTitleMap[composition.ViewTextDocument] = builder.Translation["documentTextViewPageTitle"]
	composition.GetHTMXFunctions[composition.ViewTextDocument] = GetDocumentTextViewService
	composition.PatchHTMXFunctions[composition.UpdateUserSelection] = PatchUserSelectionService

	composition.GetHTMXFunctions[composition.AddTagDocument] = GetAddTagService
	composition.PatchHTMXFunctions[composition.AddTagDocument] = PatchAddTagService
	composition.PatchHTMXFunctions[composition.ChangeTagDocument] = PatchChangeTagService
}

var tagManipulationMutex = sync.Mutex{}

func PatchChangeTagService(w http.ResponseWriter, r *http.Request) {
	tagManipulationMutex.Lock()
	defer tagManipulationMutex.Unlock()

	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	uuidDoc := chi.URLParam(r, "doc")
	uuidTag := chi.URLParam(r, "tag")
	doc, err := extraction.GetDocumentIfNotPrivate(database.LegislativeText, uuidDoc)
	if err != nil {
		viewTextDocumentOnlySwapMessage(w, r, validation.Message{
			Message: builder.Translation["documentDoesNotExistsOrNoPremissions"],
		}, acc)
		return
	}
	if !ok {
		viewTextDocumentOnlySwapMessage(w, r, validation.Message{
			Message: builder.Translation["documentDoesNotExistsOrNoPremissions"],
		}, acc)
		return
	}
	msg := validation.FlipTagHidden(uuidTag, doc)

	if !msg.Positive {
		viewTextDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.ViewDocumentPage(uuidDoc)
	viewTextDocumentRenderRequest(w, r, acc, html)
}

func PatchAddTagService(w http.ResponseWriter, r *http.Request) {
	tagManipulationMutex.Lock()
	defer tagManipulationMutex.Unlock()
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	uuidStr := chi.URLParam(r, "uuid")
	doc, err := extraction.GetDocumentIfNotPrivate(database.LegislativeText, uuidStr)
	if err != nil {
		viewTextDocumentOnlySwapMessage(w, r, validation.Message{
			Message: builder.Translation["documentDoesNotExistsOrNoPremissions"],
		}, acc)
		return
	}
	if !ok && extraction.HasAdminAccountInOrganisation(acc.ID, doc.Organisation) != nil {
		viewTextDocumentOnlySwapMessage(w, r, validation.Message{
			Message: builder.Translation["documentDoesNotExistsOrNoPremissions"],
		}, acc)
		return
	}

	create := &validation.AddTag{}
	msg := validation.Message{Positive: false}
	err = extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		viewTextDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.AddTagToDocument(doc)

	if !msg.Positive {
		viewTextDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.ViewDocumentPage(uuidStr)
	viewTextDocumentRenderRequest(w, r, acc, html)
}

func GetAddTagService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok && extraction.HasAdminAccountInOrganisation(acc.ID, r.URL.Query().Get("org")) != nil {
		renderRequest(w, builder.DIV())
		return
	}
	renderRequest(w, composition.GetTagAdminPanel(chi.URLParam(r, "uuid"), ok))
}

func PatchUserSelectionService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(r)

	err := r.ParseForm()
	name := ""
	if err == nil {
		name = r.PostFormValue("authorAccount")
	}
	account, ok, err := validation.IsAccountValidForUser(acc.ID, name)
	if !ok || err != nil {
		w.Header().Set("HX-Retarget", "#"+composition.MessageID)
		html := composition.GetMessage(validation.Message{
			Message:  builder.Translation["notAllowedToUseAccount"],
			Positive: false,
		})
		renderRequest(w, html)
		return
	}

	html, err := composition.UpdateUserOrganisations(acc, account, "", chi.URLParam(r, "isAdmin"))
	var extraNode builder.Node = nil
	if err != nil {
		extraNode = composition.GetMessageOOB(validation.Message{
			Message:  builder.Translation["errorRetrievingOrganisationForAccount"],
			Positive: false,
		})
	}
	renderRequest(w, builder.Group(html, extraNode))
}

func GetDocumentTextViewService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(r)

	html := composition.ViewDocumentPage(chi.URLParam(r, "uuid"))
	viewTextDocumentRenderRequest(w, r, acc, html)
}

var viewTextDocumentRenderRequest = genericRenderer(composition.ViewTextDocument)
var viewTextDocumentOnlySwapMessage = genericMessageSwapper(composition.ViewTextDocument)

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
	html := composition.ViewDocumentPage(create.UUIDredirect)
	viewTextDocumentRenderRequest(w, r, acc, html)
}

func GetDocumentTextCreationService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.CreateDocumentPage(acc, &validation.CreateDocument{TagColor: "#008000"}, validation.Message{})
	createTextDocumentRenderRequest(w, r, acc, html)
}

var createTextDocumentRenderRequest = genericRenderer(composition.CreateTextDocument)
var createTextDocumentOnlySwapMessage = genericMessageSwapper(composition.CreateTextDocument)