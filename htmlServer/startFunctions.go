package htmlServer

import (
	"PoliSim/htmlComposition"
	"net/http"
)

func InstallStart() {
	htmlComposition.HandlerList[htmlComposition.Start] = &htmlComposition.HttpHandling{
		TitleText:          "Startseite",
		SidebarButtonText:  nil,
		SidebarSubMenuText: nil,
		GetFunction:        GetStartService,
		PostFunction:       PostStartService,
		PatchFunction:      nil,
		DeleteFunction:     nil,
	}
}

func GetStartService(w http.ResponseWriter, r *http.Request) {
	html := htmlComposition.GetStartPage(r.URL.RawQuery)
	renderRequest(w, false, html.Render)
}

func PostStartService(w http.ResponseWriter, r *http.Request) {

}
