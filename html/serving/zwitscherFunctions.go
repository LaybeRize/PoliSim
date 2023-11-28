package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"net/http"
)

func InstallZwitscher() {
	composition.PageTitleMap[composition.ViewZwitscher] = builder.Translation["zwitscherViewPageTitle"]
	composition.SidebarTitleMap[composition.ViewZwitscher] = builder.Translation["zwitscherViewSidebarText"]
	composition.GetHTMXFunctions[composition.ViewZwitscher] = GetZwitscherViewService
}

func GetZwitscherViewService(w http.ResponseWriter, r *http.Request) {
	acc, isAdmin := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)

	extraInfo := &extraction.ExtraZwitscherInfo{
		IsAdmin: isAdmin,
	}
	extractURLFieldValues(extraInfo, r, composition.MinZwitscher, int64(standardAmount), composition.MaxZwitscher)

	html := composition.GetZwitschers(acc, extraInfo)
	viewZwitscherRenderRequest(w, r, acc, html)
}

var viewZwitscherRenderRequest = genericRenderer(composition.ViewZwitscher)
