package htmlServer

import (
	"PoliSim/componentHelper"
	"PoliSim/htmlComposition"
	"github.com/go-chi/chi"
	"net/http"
)

func InstallTitlePages() {
	htmlComposition.PageTitleMap[htmlComposition.ViewTitles] = componentHelper.Translation["titleViewPageTitle"]
	htmlComposition.SidebarTitleMap[htmlComposition.ViewTitles] = componentHelper.Translation["titleViewSidebarText"]
	htmlComposition.GetHTMXFunctions[htmlComposition.ViewTitles] = GetTitleViewService
	htmlComposition.GetHTMXFunctions[htmlComposition.GetTitleSubGroup] = GetSubGroupHTMLElement
}

func GetTitleViewService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(w, r)
	html := htmlComposition.GetViewTitelPage()
	viewTitleRenderRequest(w, r, acc.Role, html)
}

func GetSubGroupHTMLElement(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(w, r)
	mainGroup := chi.URLParam(r, "mainGroup")
	subGroup := chi.URLParam(r, "subGroup")
	html := htmlComposition.GetViewSubGroupOfTitles(mainGroup, subGroup)
	viewTitleRenderRequest(w, r, acc.Role, html)
}

var viewTitleRenderRequest = genericRenderer(htmlComposition.ViewTitles)
