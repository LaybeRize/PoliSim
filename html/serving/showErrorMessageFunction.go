package serving

import (
	"PoliSim/data/database"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"net/http"
)

func InstallErrorPage() {
	composition.PageTitleMap[composition.ErrorPage] = builder.Translation["errorPageTitle"]
}

func ShowErrorPage(w http.ResponseWriter, r *http.Request, acc *database.AccountAuth, errorText string) {
	renderRequest(w, updateInformation(w, r, acc, composition.ErrorPage),
		composition.GetErrorPage(errorText))
}
