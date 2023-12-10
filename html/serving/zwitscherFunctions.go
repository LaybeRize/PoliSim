package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InstallZwitscher() {
	composition.PageTitleMap[composition.ViewZwitscher] = builder.Translation["zwitscherViewPageTitle"]
	composition.SidebarTitleMap[composition.ViewZwitscher] = builder.Translation["zwitscherViewSidebarText"]
	composition.GetHTMXFunctions[composition.ViewZwitscher] = GetZwitscherViewService
	composition.PostHTMXFunctions[composition.CreateZwitscher] = CreateZwitscherService
	composition.PageTitleMap[composition.ViewSingleZwitscher] = builder.Translation["singleZwitscherViewPageTitle"]
	composition.GetHTMXFunctions[composition.ViewSingleZwitscher] = GetSingleZwitscherViewService
}

func CreateZwitscherService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	parent := r.URL.Query().Get("zwitscher")
	create := &validation.CreateZwitscher{
		ParentZwitscher: parent,
	}
	html := builder.Node(nil)
	err := extractFormValuesForFields(create, r, 0)
	msg := validation.Message{Positive: false}
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
	} else {
		msg = create.CreateZwitscher(acc)
	}

	if msg.Positive {
		retargetTo(w, composition.ZwitscherInterfaceID)
		html = composition.GetZwitscherInterface(acc, parent, false, msg)
	} else {
		retargetToMessage(w)
		html = composition.GetMessage(msg)
	}

	if parent == "" {
		viewZwitscherRenderRequest(w, r, acc, html)
	} else {
		viewSingleZwitscherRenderRequest(w, r, acc, html)
	}
}

func GetSingleZwitscherViewService(w http.ResponseWriter, r *http.Request) {
	acc, isAdmin := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)

	html := composition.GetSingleZwitscher(acc, isAdmin, chi.URLParam(r, "uuid"))
	viewSingleZwitscherRenderRequest(w, r, acc, html)
}

func GetZwitscherViewService(w http.ResponseWriter, r *http.Request) {
	acc, isAdmin := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)

	extraInfo := &extraction.ZwitscherQueryInfo{
		IsAdmin: isAdmin,
	}
	extractURLFieldValues(extraInfo, r, composition.MinZwitscher, int64(standardAmount), composition.MaxZwitscher)

	html := composition.GetZwitschers(acc, extraInfo)
	viewZwitscherRenderRequest(w, r, acc, html)
}

var viewSingleZwitscherRenderRequest = genericRenderer(composition.ViewSingleZwitscher)
var viewZwitscherRenderRequest = genericRenderer(composition.ViewZwitscher)
