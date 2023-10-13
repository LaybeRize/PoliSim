package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"net/http"
)

func InstallStart() {
	composition.PageTitleMap[composition.Start] = builder.Translation["startPageTitle"]
	composition.SidebarTitleMap[composition.Start] = builder.Translation["startSidebarText"]
	composition.GetHTMXFunctions[composition.Start] = GetStartService
	composition.PostHTMXFunctions[composition.Login] = PostLoginService
	composition.PostHTMXFunctions[composition.Logout] = PostLogoutService
}

func GetStartService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(w, r)
	html := composition.GetStartPage(acc, validation.Message{})
	startRenderRequest(w, r, acc.Role, html)
}

func PostLoginService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	msg := validation.Message{}
	if ok {
		msg.Message = builder.Translation["alreadyLoggedIn"]
		startOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	try := validation.LoginForm{}
	err := extractFormValuesForFields(&try, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		startOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	msg, loginAccount := try.TryLogin(w, r)
	if !msg.Positive {
		startOnlySwapMessage(w, r, msg, acc.Role)
		return
	}

	html := composition.GetStartPage(&extraction.AccountAuth{
		DisplayName: loginAccount.DisplayName,
		Role:        loginAccount.Role,
	}, msg)
	renderRequest(w, false, builder.Group(
		updateInformation(w, r, loginAccount.Role, composition.Start),
		html).Render)
}

func PostLogoutService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(w, r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	val := validation.Message{}
	if !ok {
		val.Message = builder.Translation["alreadyLoggedOut"]
		startOnlySwapMessage(w, r, val, acc.Role)
		return
	}
	cookie := validation.InvalidateAccountToken()
	w.Header().Set("Set-Cookie", cookie.String())

	val.Positive = true
	val.Message = builder.Translation["successfullyLoggedOut"]
	html := composition.GetStartPage(&extraction.AccountAuth{}, val)
	startRenderRequest(w, r, database.NotLoggedIn, html)
}

var startRenderRequest = genericRenderer(composition.Start)
var startOnlySwapMessage = genericMessageSwapper(composition.Start)
