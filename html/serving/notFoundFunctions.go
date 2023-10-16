package serving

import (
	"PoliSim/html/composition"
	"net/http"
)

func NotFoundService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(w, r)
	html := composition.GetNotFoundPage()
	renderRequest(w, false, groupNodes(updateInformation(r, acc.Role, composition.NotFound),
		html))
}
