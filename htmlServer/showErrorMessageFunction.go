package htmlServer

import (
	"PoliSim/componentHelper"
	"PoliSim/dataExtraction"
	"PoliSim/htmlComposition"
	"net/http"
)

func InstallErrorPage() {
	htmlComposition.PageTitleMap[htmlComposition.ErrorPage] = componentHelper.Translation["errorPageTitle"]
}

func ShowErrorPage(w http.ResponseWriter, r *http.Request, acc *dataExtraction.AccountAuth, errorText string) {
	renderRequest(w, false, groupNodes(updateInformation(w, r, acc.Role, htmlComposition.ErrorPage),
		htmlComposition.GetErrorPage(errorText)))
}
