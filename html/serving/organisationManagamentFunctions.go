package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi"
	"net/http"
)

func InstallOrganisationPages() {
	composition.PageTitleMap[composition.CreateOrganisation] = builder.Translation["organisationCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateOrganisation] = builder.Translation["organisationCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateOrganisation] = GetOrganisationCreateService
	composition.PostHTMXFunctions[composition.CreateOrganisation] = PostOrganisationCreateService

	composition.PageTitleMap[composition.EditOrganisation] = builder.Translation["organisationEditPageTitle"]
	composition.SidebarTitleMap[composition.EditOrganisation] = builder.Translation["organisationEditSidebarText"]
	composition.GetHTMXFunctions[composition.EditOrganisation] = GetOrganisationEditService
	composition.PostHTMXFunctions[composition.EditOrganisation] = PostOrganisationEditService
	composition.PatchHTMXFunctions[composition.SearchOrganisation] = PatchOrganisationSearchService

	composition.PageTitleMap[composition.ViewOrganisations] = builder.Translation["organisationViewPageTitle"]
	composition.SidebarTitleMap[composition.ViewOrganisations] = builder.Translation["organisationViewSidebarText"]
	composition.GetHTMXFunctions[composition.ViewOrganisations] = GetOrganisationViewService
	composition.GetHTMXFunctions[composition.GetOrganisationSubGroup] = GetSubGroupOrganisationHTMLElement

	composition.PageTitleMap[composition.ViewHiddenOrganisations] = builder.Translation["hiddenOrganisationViewPageTitle"]
	composition.SidebarTitleMap[composition.ViewHiddenOrganisations] = builder.Translation["hiddenOrganisationViewSidebarText"]
	composition.GetHTMXFunctions[composition.ViewHiddenOrganisations] = GetHiddenOrganisationViewService
}

func GetOrganisationCreateService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetCreateOrganisationPage(&validation.OrganisationModification{}, validation.Message{})
	createOrganisationRenderRequest(w, r, acc, html)
}

func PostOrganisationCreateService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	create := &validation.OrganisationModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createOrganisationOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.CreateOrganisation()
	if !msg.Positive {
		createOrganisationOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetCreateOrganisationPage(create, msg)
	createOrganisationRenderRequest(w, r, acc, html)
}

var createOrganisationRenderRequest = genericRenderer(composition.CreateOrganisation)
var createOrganisationOnlySwapMessage = genericMessageSwapper(composition.CreateOrganisation)

func GetOrganisationEditService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetModifyOrganisationPage(&validation.OrganisationModification{}, validation.Message{})
	editOrganisationRenderRequest(w, r, acc, html)
}

func PostOrganisationEditService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	create := &validation.OrganisationModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		editOrganisationOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.ModifyOrganisation()
	if !msg.Positive {
		editOrganisationOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetModifyOrganisationPage(create, msg)
	editOrganisationRenderRequest(w, r, acc, html)
}

func PatchOrganisationSearchService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	create := &validation.OrganisationModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		editOrganisationOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.SearchOrganisation()
	if !msg.Positive {
		editOrganisationOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetModifyOrganisationPage(create, msg)
	editOrganisationRenderRequest(w, r, acc, html)
}

var editOrganisationRenderRequest = genericRenderer(composition.EditOrganisation)
var editOrganisationOnlySwapMessage = genericMessageSwapper(composition.EditOrganisation)

func GetOrganisationViewService(w http.ResponseWriter, r *http.Request) {
	acc, isAdmin := CheckUserPrivileges(r, database.Admin, database.HeadAdmin)
	html := composition.GetViewOrganisationPage(acc.ID, isAdmin)
	viewOrganisationRenderRequest(w, r, acc, html)
}

func GetSubGroupOrganisationHTMLElement(w http.ResponseWriter, r *http.Request) {
	acc, isAdmin := CheckUserPrivileges(r, database.Admin, database.HeadAdmin)
	mainGroup := chi.URLParam(r, "mainGroup")
	subGroup := chi.URLParam(r, "subGroup")
	html := composition.GetViewSubGroupOfOrganisations(acc.ID, isAdmin, mainGroup, subGroup)
	viewOrganisationRenderRequest(w, r, acc, html)
}

var viewOrganisationRenderRequest = genericRenderer(composition.ViewOrganisations)

func GetHiddenOrganisationViewService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}
	html := composition.GetViewHiddenOrganisationPage()
	viewHiddenOrganisationRenderRequest(w, r, acc, html)
}

var viewHiddenOrganisationRenderRequest = genericRenderer(composition.ViewTitles)
