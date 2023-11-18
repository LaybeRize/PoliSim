package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/logic"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"net/http"
)

func InstallDocumentPager() {
	composition.PageTitleMap[composition.ViewDocument] = builder.Translation["documentViewListPageTitle"]
	composition.SidebarTitleMap[composition.ViewDocument] = builder.Translation["documentViewListSidebarText"]
	composition.GetHTMXFunctions[composition.ViewDocument] = GetViewDocumentsService
}

func GetViewDocumentsService(w http.ResponseWriter, r *http.Request) {
	acc, isAdmin := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)

	extraInfo := &logic.ExtraInfo{
		ViewAccountID: acc.ID,
	}
	extractURLFieldValues(extraInfo, r, 5, int64(standardAmount), 50)

	html := composition.GetDocumentPage(isAdmin, extraInfo)
	viewDocumentsRenderRequest(w, r, acc, html)
}

var viewDocumentsRenderRequest = genericRenderer(composition.ViewDocument)
