package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

//TODO fill this with live

func InstallVoteCreation() {
	composition.PageTitleMap[composition.CreateVoteDocument] = builder.Translation["voteCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateVoteDocument] = builder.Translation["voteCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateVoteDocument] = GetVoteCreationService
	composition.PostHTMXFunctions[composition.CreateVoteDocument] = PostCreateVoteInDatabaseService
	composition.PatchHTMXFunctions[composition.RequestVotePartial] = PatchGetVotePartial
	composition.PageTitleMap[composition.ViewVoteDocument] = builder.Translation["voteViewPageTitle"]
	composition.GetHTMXFunctions[composition.ViewVoteDocument] = GetVoteViewService
}

func GetVoteViewService(w http.ResponseWriter, r *http.Request) {
	acc, admin := CheckUserPrivileges(r, database.Admin, database.HeadAdmin)

	html := composition.GetVoteViewPage(acc, chi.URLParam(r, "uuid"), admin,
		validation.Message{})
	viewVoteDocumentRenderRequest(w, r, acc, html)
}

var viewVoteDocumentRenderRequest = genericRenderer(composition.ViewVoteDocument)

func GetVoteCreationService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetCreateVotePage(acc, &validation.CreateVote{PrivateDocumentInfo: validation.PrivateDocumentInfo{
		EndTime: time.Now().Add(time.Hour * 25).Format("2006-01-02T15:04"),
	}}, validation.Message{})
	createVoteDocumentRenderRequest(w, r, acc, html)
}

func PostCreateVoteInDatabaseService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}
	create := validation.CreateVote{}
	err := json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createVoteDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}
	helper.ClearStringArray(&create.Onlooker)
	helper.ClearStringArray(&create.Participants)
	if !msg.Positive {
		createVoteDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	//TODO change this to the appropirate way
	w.Header().Set("HX-Push-Url", "/"+string(composition.ViewVoteDocumentLink)+create.UUIDredirect)
	html := composition.GetVoteViewPage(acc, create.UUIDredirect,
		acc.Role == database.Admin || acc.Role == database.HeadAdmin,
		validation.Message{})
	viewVoteDocumentRenderRequest(w, r, acc, html)
}

func PatchGetVotePartial(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.ParseUint(chi.URLParam(r, "number"), 10, 64)
	if err != nil {
		acc, _ := CheckUserPrivileges(r)
		createVoteDocumentOnlySwapMessage(w, r, validation.Message{Message: builder.Translation["invalidNumber"]}, acc)
	}
	renderRequest(w, composition.GetVotePartial(int64(i)))
}

var createVoteDocumentRenderRequest = genericRenderer(composition.CreateVoteDocument)
var createVoteDocumentOnlySwapMessage = genericMessageSwapper(composition.CreateVoteDocument)
