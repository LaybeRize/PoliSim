package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func InstallDocumentDiscussion() {
	composition.PageTitleMap[composition.CreateDiscussionDocument] = builder.Translation["documentDiscussionCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateDiscussionDocument] = builder.Translation["documentDiscussionCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateDiscussionDocument] = GetDocumentDiscussionCreationService
	composition.PostHTMXFunctions[composition.CreateDiscussionDocument] = PostDocumentDiscussionCreationService
	composition.PageTitleMap[composition.ViewDiscussionDocument] = builder.Translation["documentDiscussionViewPageTitle"]
	composition.GetHTMXFunctions[composition.ViewDiscussionDocument] = GetDocumentDiscussionViewService
	composition.PostHTMXFunctions[composition.CommentDiscussion] = PostCommentDiscussionViewService
	composition.PatchHTMXFunctions[composition.ChangeCommentDocument] = PatchChangeCommentVisibilityService
}

func PatchChangeCommentVisibilityService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	docUUID := chi.URLParam(r, "doc")
	msg := validation.Message{Positive: false}
	exists, err := logic.ChangeVisibiltyComment(chi.URLParam(r, "comment"), docUUID)
	if err != nil {
		msg.Message = builder.Translation["errorProcessing"]
		viewDiscussionDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}
	if !exists {
		msg.Message = builder.Translation["commentDoesNotExists"]
		viewDiscussionDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg.Positive = true
	msg.Message = builder.Translation["changeCommentVisiblitySuccessfull"]
	html := composition.ViewDiscussionPage(acc, docUUID, true, msg)
	viewDiscussionDocumentRenderRequest(w, r, acc, html)
}

func PostCommentDiscussionViewService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)

	isAdmin := acc.Role == database.Admin || acc.Role == database.HeadAdmin
	uuidStr := chi.URLParam(r, "uuid")
	_, err := extraction.GetDocumentForUser(uuidStr, acc.ID, isAdmin, database.FinishedDiscussion, database.RunningDiscussion)

	if !ok || err != nil {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	create := &validation.AddComment{}
	msg := validation.Message{Positive: false}
	err = extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		viewDiscussionDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.AddComment(uuidStr, acc)

	if !msg.Positive {
		viewDiscussionDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.ViewDiscussionPage(acc, uuidStr, isAdmin, msg)
	viewDiscussionDocumentRenderRequest(w, r, acc, html)
}

func GetDocumentDiscussionViewService(w http.ResponseWriter, r *http.Request) {
	acc, admin := CheckUserPrivileges(r, database.Admin, database.HeadAdmin)

	html := composition.ViewDiscussionPage(acc, chi.URLParam(r, "uuid"), admin,
		validation.Message{})
	viewDiscussionDocumentRenderRequest(w, r, acc, html)
}

var viewDiscussionDocumentRenderRequest = genericRenderer(composition.ViewDiscussionDocument)
var viewDiscussionDocumentOnlySwapMessage = genericMessageSwapper(composition.ViewDiscussionDocument)

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
	html := composition.ViewDiscussionPage(acc, create.UUIDredirect,
		acc.Role == database.Admin || acc.Role == database.HeadAdmin,
		validation.Message{})
	viewDiscussionDocumentRenderRequest(w, r, acc, html)
}

func GetDocumentDiscussionCreationService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.CreateDiscussionPage(acc, &validation.CreateDiscussion{
		EndTime: time.Now().Add(time.Hour * 25).Format("2006-01-02T15:04"),
	}, validation.Message{})
	createDiscussionDocumentRenderRequest(w, r, acc, html)
}

var createDiscussionDocumentRenderRequest = genericRenderer(composition.CreateDiscussionDocument)
var createDiscussionDocumentOnlySwapMessage = genericMessageSwapper(composition.CreateDiscussionDocument)
