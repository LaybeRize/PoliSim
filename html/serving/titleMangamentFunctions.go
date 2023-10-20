package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi"
	"net/http"
)

func InstallTitlePages() {
	composition.PageTitleMap[composition.ViewTitles] = builder.Translation["titleViewPageTitle"]
	composition.SidebarTitleMap[composition.ViewTitles] = builder.Translation["titleViewSidebarText"]
	composition.GetHTMXFunctions[composition.ViewTitles] = GetTitleViewService
	composition.GetHTMXFunctions[composition.GetTitleSubGroup] = GetSubGroupHTMLElement

	composition.PageTitleMap[composition.EditTitle] = builder.Translation["titleEditPageTitle"]
	composition.SidebarTitleMap[composition.EditTitle] = builder.Translation["titleEditSidebarText"]
	composition.GetHTMXFunctions[composition.EditTitle] = GetTitleEditService
	composition.PatchHTMXFunctions[composition.SearchTitle] = PatchSearchTitleService
	composition.PostHTMXFunctions[composition.EditTitle] = PostEditTitleService
	composition.PatchHTMXFunctions[composition.DeleteTitle] = DeleteTitleService

	composition.PageTitleMap[composition.CreateTitle] = builder.Translation["titleCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateTitle] = builder.Translation["titleCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateTitle] = GetTitleCreateService
	composition.PostHTMXFunctions[composition.CreateTitle] = PostTitleCreateService
}

func GetTitleViewService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(r)
	html := composition.GetViewTitelPage()
	viewTitleRenderRequest(w, r, acc, html)
}

func GetSubGroupHTMLElement(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(r)
	mainGroup := chi.URLParam(r, "mainGroup")
	subGroup := chi.URLParam(r, "subGroup")
	html := composition.GetViewSubGroupOfTitles(mainGroup, subGroup)
	viewTitleRenderRequest(w, r, acc, html)
}

var viewTitleRenderRequest = genericRenderer(composition.ViewTitles)

func GetTitleCreateService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetCreateTitlePage(&validation.TitleModification{}, validation.Message{})
	createTitleRenderRequest(w, r, acc, html)
}

func PostTitleCreateService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{}

	create := &validation.TitleModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createTitleOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.CreateTitle()
	if !msg.Positive {
		createTitleOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetCreateTitlePage(create, msg)
	createTitleRenderRequest(w, r, acc, html)
}

var createTitleRenderRequest = genericRenderer(composition.CreateTitle)
var createTitleOnlySwapMessage = genericMessageSwapper(composition.CreateTitle)

func GetTitleEditService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetModifyTitlePage(&validation.TitleModification{}, validation.Message{})
	editTitleRenderRequest(w, r, acc, html)
}

func PatchSearchTitleService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{}

	create := &validation.TitleModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		editTitleOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.SearchTitle()
	if !msg.Positive {
		editTitleOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetModifyTitlePage(create, msg)
	editTitleRenderRequest(w, r, acc, html)
}

func PostEditTitleService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{}

	create := &validation.TitleModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		editTitleOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.ModifyTitle()
	if !msg.Positive {
		editTitleOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetModifyTitlePage(create, msg)
	editTitleRenderRequest(w, r, acc, html)
}

func DeleteTitleService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{}

	create := &validation.TitleModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		editTitleOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.DeleteTitle()
	if !msg.Positive {
		editTitleOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetModifyTitlePage(create, msg)
	editTitleRenderRequest(w, r, acc, html)
}

var editTitleRenderRequest = genericRenderer(composition.EditTitle)
var editTitleOnlySwapMessage = genericMessageSwapper(composition.EditTitle)
