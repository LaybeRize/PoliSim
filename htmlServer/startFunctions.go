package htmlServer

import (
	"PoliSim/database"
	"PoliSim/htmlComposition"
	"net/http"
)

func InstallStart() {
	htmlComposition.HandlerList[htmlComposition.Start] = &htmlComposition.HttpHandling{
		TitleText:         "Startseite",
		SidebarButtonText: "Start",
		HasSidebarButton:  true,
	}
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
	}
	_ = acc.ID
}

func PostLogoutService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivilges(w, r, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	if !ok {
		//tell the user he is already logged out
	}
	_ = acc.ID
}
