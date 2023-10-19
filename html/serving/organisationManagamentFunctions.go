package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
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
}

func GetOrganisationCreateService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetCreateOrganisationPage(&validation.OrganisationModification{}, validation.Message{})
	createOrganisationRenderRequest(w, r, acc.Role, html)
}

func PostOrganisationCreateService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	create := &validation.OrganisationModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createOrganisationOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	msg = create.CreateOrganisation()
	if !msg.Positive {
		createOrganisationOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	html := composition.GetCreateOrganisationPage(create, msg)
	createOrganisationRenderRequest(w, r, acc.Role, html)
}

var createOrganisationRenderRequest = genericRenderer(composition.CreateOrganisation)
var createOrganisationOnlySwapMessage = genericMessageSwapper(composition.CreateOrganisation)

func GetOrganisationEditService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetModifyOrganisationPage(&validation.OrganisationModification{}, validation.Message{})
	editOrganisationRenderRequest(w, r, acc.Role, html)
}

func PostOrganisationEditService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	create := &validation.OrganisationModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		editOrganisationOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	msg = create.ModifyOrganisation()
	if !msg.Positive {
		editOrganisationOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	html := composition.GetModifyOrganisationPage(create, msg)
	editOrganisationRenderRequest(w, r, acc.Role, html)
}

func PatchOrganisationSearchService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin, database.Admin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	create := &validation.OrganisationModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		editOrganisationOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	msg = create.SearchOrganisation()
	if !msg.Positive {
		editOrganisationOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	html := composition.GetModifyOrganisationPage(create, msg)
	editOrganisationRenderRequest(w, r, acc.Role, html)
}

var editOrganisationRenderRequest = genericRenderer(composition.EditOrganisation)
var editOrganisationOnlySwapMessage = genericMessageSwapper(composition.EditOrganisation)
