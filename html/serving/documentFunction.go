package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InstallDocumentPager() {
	composition.PageTitleMap[composition.ViewDocument] = builder.Translation["documentViewListPageTitle"]
	composition.SidebarTitleMap[composition.ViewDocument] = builder.Translation["documentViewListSidebarText"]
	composition.GetHTMXFunctions[composition.ViewDocument] = GetViewDocumentsService

	composition.PatchHTMXFunctions[composition.BlockDocument] = PatchBlockDocumentService
}

func PatchBlockDocumentService(w http.ResponseWriter, r *http.Request) {
	_, isAdmin := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	val := validation.Message{}
	if !isAdmin {
		val.Message = builder.Translation["notAllowedToBlockDocuments"]
		renderRequest(w, composition.GetMessage(val))
		return
	}
	doc, err := logic.BlockDocument(chi.URLParam(r, "uuid"))
	if err != nil {
		val.Message = builder.Translation["errorChangingBlock"]
		renderRequest(w, composition.GetMessage(val))
	}

	val.Positive = true
	if doc.Blocked {
		val.Message = builder.Translation["successfullyBlockedDocument"]
	} else {
		val.Message = builder.Translation["successfullyUnblockedDocument"]
	}
	renderRequest(w, composition.GetMessage(val), composition.GetNewDocumentHeader(doc))
}

func GetViewDocumentsService(w http.ResponseWriter, r *http.Request) {
	acc, isAdmin := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)

	extraInfo := &extraction.ExtraInfo{
		ViewAccountID: acc.ID,
	}
	extractURLFieldValues(extraInfo, r, 5, int64(standardAmount), 50)

	html := composition.GetDocumentPage(isAdmin, extraInfo)
	viewDocumentsRenderRequest(w, r, acc, html)
}

var viewDocumentsRenderRequest = genericRenderer(composition.ViewDocument)
