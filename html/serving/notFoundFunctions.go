package serving

import (
	"PoliSim/html/composition"
	"net/http"
)

func NotFoundService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(r)
	html := composition.GetNotFoundPage()
	renderRequest(w, updateInformation(w, r, acc, composition.NotFound),
		html)
}
