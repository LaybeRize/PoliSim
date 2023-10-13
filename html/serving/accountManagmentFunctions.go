package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"net/http"
)

func InstallAccountManagment() {
	composition.PageTitleMap[composition.CreateUser] = builder.Translation["createUserTitle"]
	composition.SidebarTitleMap[composition.CreateUser] = builder.Translation["createUserSidebarText"]
	composition.GetHTMXFunctions[composition.CreateUser] = GetCreateUserService
	composition.PostHTMXFunctions[composition.CreateUser] = PostCreateUserService
	composition.PageTitleMap[composition.EditUser] = builder.Translation["editUserTitle"]
	composition.SidebarTitleMap[composition.EditUser] = builder.Translation["editUserSidebarText"]
	composition.GetHTMXFunctions[composition.EditUser] = GetEditUserService
	composition.PostHTMXFunctions[composition.EditUser] = PostEditUserService
	composition.PatchHTMXFunctions[composition.SearchUser] = PatchSearchUserService
	composition.PageTitleMap[composition.ViewUser] = builder.Translation["viewUserTitle"]
	composition.SidebarTitleMap[composition.ViewUser] = builder.Translation["viewUserSidebarText"]
	composition.GetHTMXFunctions[composition.ViewUser] = GetViewUserService
}

func GetViewUserService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}
	html := composition.GetViewAccountList(r.URL.Query().Get("id"))
	viewUserRenderRequest(w, r, acc.Role, html)
}

func PatchSearchUserService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{}

	create := &validation.AccountModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		editUserOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	msg = create.RequestAccount()
	if !msg.Positive {
		editUserOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	html := composition.GetModifyAccount(create, msg)
	editUserRenderRequest(w, r, acc.Role, html)
}

func PostEditUserService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{}

	create := &validation.AccountModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		editUserOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	msg = create.ValidateAccountModification(acc)
	if !msg.Positive {
		editUserOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	html := composition.GetModifyAccount(create, msg)
	editUserRenderRequest(w, r, acc.Role, html)
}

func GetEditUserService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}
	html := composition.GetModifyAccount(&validation.AccountModification{
		Role: int(database.User),
	}, validation.Message{})
	editUserRenderRequest(w, r, acc.Role, html)
}

func GetCreateUserService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}
	html := composition.GetCreateAccountPage(&validation.AccountModification{
		Role: int(database.User),
	}, validation.Message{})
	createUserRenderRequest(w, r, acc.Role, html)
}

func PostCreateUserService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.HeadAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{}

	create := &validation.AccountModification{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createUserOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	msg = create.ValidateAccountCreation(acc)
	if !msg.Positive {
		createUserOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	html := composition.GetCreateAccountPage(create, msg)
	createUserRenderRequest(w, r, acc.Role, html)
}

var createUserRenderRequest = genericRenderer(composition.CreateUser)
var createUserOnlySwapMessage = genericMessageSwapper(composition.CreateUser)
var editUserRenderRequest = genericRenderer(composition.EditUser)
var editUserOnlySwapMessage = genericMessageSwapper(composition.EditUser)
var viewUserRenderRequest = genericRenderer(composition.ViewUser)
