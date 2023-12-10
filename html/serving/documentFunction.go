package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func InstallDocumentPager() {
	composition.PageTitleMap[composition.ViewDocument] = builder.Translation["documentViewListPageTitle"]
	composition.SidebarTitleMap[composition.ViewDocument] = builder.Translation["documentViewListSidebarText"]
	composition.GetHTMXFunctions[composition.ViewDocument] = GetViewDocumentsService
	composition.PatchHTMXFunctions[composition.ViewDocument] = PatchViewDocumentsService

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

	extraInfo := &extraction.DocumentQueryInfo{
		ViewAccountID: acc.ID,
		IsAdmin:       isAdmin,
	}
	extractURLFieldValues(extraInfo, r, composition.MinDocuments, int64(standardAmount), composition.MaxDocuments)

	html := composition.GetDocumentPage(extraInfo)
	viewDocumentsRenderRequest(w, r, acc, html)
}

func PatchViewDocumentsService(w http.ResponseWriter, r *http.Request) {
	acc, isAdmin := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	extraInfo := &extraction.DocumentQueryInfo{
		ViewAccountID: acc.ID,
		IsAdmin:       isAdmin,
	}

	err := extractFormValuesForFields(extraInfo, r, 0)
	if err != nil {
		ShowErrorPage(w, r, acc, builder.Translation["extractionError"])
		return
	}
	if extraInfo.Amount < composition.MinDocuments {
		extraInfo.Amount = composition.MinDocuments
	}
	if extraInfo.Amount > composition.MaxDocuments {
		extraInfo.Amount = composition.MaxDocuments
	}
	str := composition.GetExtraString(extraInfo)
	if strings.TrimSpace(r.PostFormValue("addWritten")) == "true" {
		str += "&written=" + extraInfo.Written
	} else {
		extraInfo.Written = ""
	}

	pushURL(w, fmt.Sprintf("/%s?amount=%d%s", string(composition.ViewDocument), extraInfo.Amount, str))
	html := composition.GetDocumentPage(extraInfo)
	viewDocumentsRenderRequest(w, r, acc, html)
}

var viewDocumentsRenderRequest = genericRenderer(composition.ViewDocument)
