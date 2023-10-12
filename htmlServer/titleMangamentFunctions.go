package htmlServer

import (
	"PoliSim/htmlComposition"
	"github.com/go-chi/chi"
	"net/http"
)

func InstallTitlePages() {

}

func GetSubGroupHTMLElement(w http.ResponseWriter, r *http.Request) {
	// url: title/get-sub-group/{mainGroup}/{subGroup}
	mainGroup := chi.URLParam(r, "mainGroup")
	subGroup := chi.URLParam(r, "subGroup")
	html := htmlComposition.GetViewSubGroupOfTitles(mainGroup, subGroup)
	//do the rendering you know
	_ = html
}
