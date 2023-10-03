package htmlServer

import (
	"PoliSim/componentHelper"
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
	acc, _ := CheckUserPrivilges(w, r)
	html := htmlComposition.GetStartPage(acc)
	renderRequest(w, false, html.Render)
}

func PostLoginService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivilges(w, r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	if ok {
		//tell the user he is already logged in
		return
	}
	try := dataValidation.LoginForm{}
	err := extractValuesForFields(&try, r, 0)
	if err != nil {
		//handel extraction error
		return
	}
	msg, loginAccount, cookie := try.TryLogin()
	if !msg.Positive {
		//login failed
		return
	}

	http.SetCookie(w, cookie)
	_ = acc
	_ = loginAccount
}

func PostLogoutService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivilges(w, r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	if !ok {
		//tell the user he is already logged out
		return
	}
	err, cookie := dataValidation.InvalidateAccountToken(acc)
	if err != nil {
		//report error
		return
	}
	w.Header().Set("Set-Cookie", cookie.String())
}
