package htmlServer

import (
	"PoliSim/componentHelper"
	"PoliSim/dataValidation"
	"PoliSim/database"
	"PoliSim/htmlComposition"
	"net/http"
)

func InstallAccountManagment() {
	htmlComposition.PageTitleMap[htmlComposition.CreateUser] = componentHelper.Translation["createUserTitle"]
	htmlComposition.SidebarTitleMap[htmlComposition.CreateUser] = componentHelper.Translation["createUserSidebarText"]
	htmlComposition.GetHTMXFunctions[htmlComposition.CreateUser] = GetCreateUserService
	htmlComposition.PostHTMXFunctions[htmlComposition.CreateUser] = PostCreateUserService
}

func GetCreateUserService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, componentHelper.Translation["notAllowedToViewThisPage"])
		return
	}
	html := htmlComposition.GetCreateAccountPage(&dataValidation.AccountModification{
		Role: int(database.User),
	}, dataValidation.ValidationMessage{})
	createUserRenderRequest(w, r, acc.Role, html)
}

func PostCreateUserService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, componentHelper.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := dataValidation.ValidationMessage{}

	create := &dataValidation.AccountModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = componentHelper.Translation["extractionError"]
		createUserOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	msg = create.ValidateAccountCreation(acc)
	if !msg.Positive {
		createUserOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	html := htmlComposition.GetCreateAccountPage(create, msg)
	createUserRenderRequest(w, r, acc.Role, html)
}

var createUserRenderRequest = genericRenderer(htmlComposition.CreateUser)
var createUserOnlySwapMessage = genericMessageSwapper(htmlComposition.CreateUser)