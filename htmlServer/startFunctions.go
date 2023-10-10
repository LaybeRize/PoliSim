package htmlServer

import (
	"PoliSim/componentHelper"
	"PoliSim/dataExtraction"
	"PoliSim/dataValidation"
	"PoliSim/database"
	"PoliSim/htmlComposition"
	"net/http"
)

func InstallStart() {
	htmlComposition.PageTitleMap[htmlComposition.Start] = componentHelper.Translation["startPageTitle"]
	htmlComposition.SidebarTitleMap[htmlComposition.Start] = componentHelper.Translation["startSidebarText"]
	htmlComposition.GetHTMXFunctions[htmlComposition.Start] = GetStartService
	htmlComposition.PostHTMXFunctions[htmlComposition.Login] = PostLoginService
	htmlComposition.PostHTMXFunctions[htmlComposition.Logout] = PostLogoutService
}

func GetStartService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(w, r)
	html := htmlComposition.GetStartPage(acc, dataValidation.ValidationMessage{})
	startRenderRequest(w, r, acc.Role, html)
}

func PostLoginService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	msg := dataValidation.ValidationMessage{}
	if ok {
		msg.Message = componentHelper.Translation["alreadyLoggedIn"]
		startOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	try := dataValidation.LoginForm{}
	err := extractFormValuesForFields(&try, r, 0)
	if err != nil {
		msg.Message = componentHelper.Translation["extractionError"]
		startOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	msg, loginAccount := try.TryLogin(w, r)
	if !msg.Positive {
		startOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	html := htmlComposition.GetStartPage(&dataExtraction.AccountAuth{
		DisplayName: loginAccount.DisplayName,
		Role:        loginAccount.Role,
	}, msg)
	renderRequest(w, false, componentHelper.Group(
		updateInformation(w, r, loginAccount.Role, htmlComposition.Start),
		html).Render)
}

func PostLogoutService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	val := dataValidation.ValidationMessage{}
	if !ok {
		val.Message = componentHelper.Translation["alreadyLoggedOut"]
		startOnlySwapMessage(w, r, val, acc.Role)
		return
	}
	cookie := dataValidation.InvalidateAccountToken()
	w.Header().Set("Set-Cookie", cookie.String())

	val.Positive = true
	val.Message = componentHelper.Translation["successfullyLoggedOut"]
	html := htmlComposition.GetStartPage(&dataExtraction.AccountAuth{}, val)
	startRenderRequest(w, r, database.NotLoggedIn, html)
}

var startRenderRequest = genericRenderer(htmlComposition.Start)
var startOnlySwapMessage = genericMessageSwapper(htmlComposition.Start)
