package serving

import (
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi"
	"net/http"
)

func InstallTitlePages() {
	composition.PageTitleMap[composition.ViewTitles] = builder.Translation["titleViewPageTitle"]
	composition.SidebarTitleMap[composition.ViewTitles] = builder.Translation["titleViewSidebarText"]
	composition.GetHTMXFunctions[composition.ViewTitles] = GetTitleViewService
	composition.GetHTMXFunctions[composition.GetTitleSubGroup] = GetSubGroupHTMLElement
}

func GetTitleViewService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(w, r)
	html := composition.GetViewTitelPage()
	viewTitleRenderRequest(w, r, acc.Role, html)
}

func GetSubGroupHTMLElement(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(w, r)
	mainGroup := chi.URLParam(r, "mainGroup")
	subGroup := chi.URLParam(r, "subGroup")
	html := composition.GetViewSubGroupOfTitles(mainGroup, subGroup)
	viewTitleRenderRequest(w, r, acc.Role, html)
}

var viewTitleRenderRequest = genericRenderer(composition.ViewTitles)
