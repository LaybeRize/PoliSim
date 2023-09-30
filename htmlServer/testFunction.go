package htmlServer

import (
	"PoliSim/htmlComposition"
	"net/http"
)

func ServeTestGet(w http.ResponseWriter, r *http.Request) {
	html := htmlComposition.GetBasePage("Test Page")
	renderRequest(w, true, html.Render)
}
