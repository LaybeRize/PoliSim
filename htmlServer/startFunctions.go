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
	acc, _ := CheckUserPrivilges(w, r)
	html := htmlComposition.GetStartPage(acc, dataValidation.ValidationMessage{})
	renderRequest(w, false, updateInformation(r, acc.Role, htmlComposition.Start),
		html.Render)
}

func PostLoginService(w http.ResponseWriter, r *http.Request) {
	_, ok := CheckUserPrivilges(w, r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	if ok {
		onlySwapMessage(w, dataValidation.ValidationMessage{
			Message: componentHelper.Translation["alreadyLoggedIn"],
		})
		return
	}
	try := dataValidation.LoginForm{}
	err := extractValuesForFields(&try, r, 0)
	if err != nil {
		onlySwapMessage(w, dataValidation.ValidationMessage{
			Message: componentHelper.Translation["extractionError"],
		})
		return
	}
	msg, loginAccount, cookie := try.TryLogin()
	if !msg.Positive {
		onlySwapMessage(w, msg)
		return
	}

	http.SetCookie(w, cookie)

	html := htmlComposition.GetStartPage(&dataExtraction.AccountAuth{
		DisplayName: loginAccount.DisplayName,
		Role:        loginAccount.Role,
	}, msg)
	renderRequest(w, false, updateInformation(r, loginAccount.Role, htmlComposition.Start),
		html.Render)
}

func PostLogoutService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivilges(w, r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	val := dataValidation.ValidationMessage{}
	if !ok {
		val.Message = componentHelper.Translation["alreadyLoggedOut"]
		onlySwapMessage(w, val)
		return
	}
	err, cookie := dataValidation.InvalidateAccountToken(acc)
	if err != nil {
		val.Message = componentHelper.Translation["errorWhileTryingToLogYouOut"]
		onlySwapMessage(w, val)
		return
	}
	w.Header().Set("Set-Cookie", cookie.String())

	val.Positive = true
	val.Message = componentHelper.Translation["successfullyLoggedOut"]
	html := htmlComposition.GetStartPage(&dataExtraction.AccountAuth{}, val)
	renderRequest(w, false, updateInformation(r, database.NotLoggedIn, htmlComposition.Start),
		html.Render)
}
