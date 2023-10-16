package serving

import (
	"PoliSim/data/extraction"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"net/http"
)

func InstallErrorPage() {
	composition.PageTitleMap[composition.ErrorPage] = builder.Translation["errorPageTitle"]
}

func ShowErrorPage(w http.ResponseWriter, r *http.Request, acc *extraction.AccountAuth, errorText string) {
	renderRequest(w, false, groupNodes(updateInformation(r, acc.Role, composition.ErrorPage),
		composition.GetErrorPage(errorText)))
}
