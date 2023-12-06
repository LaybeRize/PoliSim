package serving

import (
	"PoliSim/data/database"
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
	acc, _ := CheckUserPrivileges(r)
	html := composition.GetStartPage(acc, validation.Message{})
	startRenderRequest(w, r, acc, html)
}

func PostLoginService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	msg := validation.Message{}
	if ok {
		msg.Message = builder.Translation["alreadyLoggedIn"]
		startOnlySwapMessage(w, r, msg, acc)
		return
	}

	try := validation.LoginForm{}
	err := extractFormValuesForFields(&try, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		startOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg, loginAccount := try.TryLogin(w, r)
	if !msg.Positive {
		startOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetStartPage(&database.AccountAuth{
		DisplayName: loginAccount.DisplayName,
		Role:        loginAccount.Role,
	}, msg)
	renderRequest(w, updateInformation(w, r, &database.AccountAuth{
		ID:          loginAccount.ID,
		DisplayName: loginAccount.DisplayName,
		Suspended:   loginAccount.Suspended,
		Role:        loginAccount.Role,
		HasLetters:  loginAccount.HasLetters,
		Session:     acc.Session,
	}, composition.Start),
		html)
}

func PostLogoutService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	val := validation.Message{}
	if !ok {
		val.Message = builder.Translation["alreadyLoggedOut"]
		startOnlySwapMessage(w, r, val, acc)
		return
	}
	cookie := validation.InvalidateAccountToken()

	val.Positive = true
	val.Message = builder.Translation["successfullyLoggedOut"]
	html := composition.GetStartPage(&database.AccountAuth{}, val)
	acc.Role = database.NotLoggedIn
	update := updateInformation(w, r, acc, composition.Start)
	w.Header().Set("Set-Cookie", cookie.String())
	renderRequest(w, update, html)
}

var startRenderRequest = genericRenderer(composition.Start)
var startOnlySwapMessage = genericMessageSwapper(composition.Start)
