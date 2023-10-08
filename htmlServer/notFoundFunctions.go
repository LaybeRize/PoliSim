package htmlServer

import (
	"PoliSim/htmlComposition"
	"net/http"
)

func GetNotFoundService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivilges(w, r)
	html := htmlComposition.GetNotFoundPage()
	renderRequest(w, false, groupNodes(updateInformation(w, r, acc.Role, htmlComposition.NotFound),
		html))
}
