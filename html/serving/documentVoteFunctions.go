package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
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

func InstallVoteCreation() {
	composition.PageTitleMap[composition.CreateVoteDocument] = builder.Translation["voteCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateVoteDocument] = builder.Translation["voteCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateVoteDocument] = GetVoteCreationService
	composition.PostHTMXFunctions[composition.CreateVoteDocument] = PostCreateVoteInDatabaseService
	composition.PatchHTMXFunctions[composition.RequestVotePartial] = PatchGetVotePartial

	composition.PageTitleMap[composition.ViewVoteDocument] = builder.Translation["voteViewPageTitle"]
	composition.GetHTMXFunctions[composition.ViewVoteDocument] = GetVoteViewService
	composition.PatchHTMXFunctions[composition.MakeVote] = PatchMakeVoteService

	composition.GetHTMXFunctions[composition.VoteUpdateDocument] = GetUpdateVoteService
}

func GetUpdateVoteService(w http.ResponseWriter, r *http.Request) {
	acc, isAdmin := CheckUserPrivileges(r, database.Admin, database.HeadAdmin)
	renderRequest(w, composition.GetVoteViewPageUpdate(acc, chi.URLParam(r, "uuid"), isAdmin))
}

func PatchMakeVoteService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.Admin, database.HeadAdmin, database.MediaAdmin, database.User)
	isAdmin := CheckIfHasRole(acc, database.HeadAdmin, database.Admin)
	docUUID := chi.URLParam(r, "doc")
	voteUUID := chi.URLParam(r, "vote")
	voteType := database.VoteType(chi.URLParam(r, "type"))

	_, err := extraction.GetVoteForUser(docUUID, acc.ID, isAdmin)

	if !ok || err != nil {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{
		Message:  builder.Translation["invalidVoteType"],
		Positive: false,
	}
	if _, validType := database.VoteTranslation[voteType]; !validType {
		viewVoteDocumentOnlySwapMessage(w, r, msg, acc)
	}

	var create validation.CastVote
	switch voteType {
	case database.SingleVote:
		create = &validation.AddSingleVote{}
	case database.MultipleVotes:
		create = &validation.AddMultipleVote{}
	case database.RankedVotes:
		create = &validation.AddRankedVote{}
	case database.ThreeCategoryVoting:
		create = &validation.AddThreeChoice{}
	}

	err = json.NewDecoder(r.Body).Decode(create)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		viewVoteDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.CastVote(acc, docUUID, voteUUID, voteType)
	viewVoteDocumentOnlySwapMessage(w, r, msg, acc)
	/*if !msg.Positive {
		viewVoteDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	w.Header().Set("HX-Retarget", "#"+composition.MessageID)
	html := composition.GetMessage(msg)
	var swapInfo builder.Node
	switch vote.Info.VoteMethod {
	case database.SingleVote, database.MultipleVotes, database.ThreeCategoryVoting:
		swapInfo = composition.GetInfoStandardView(vote, true)
	case database.RankedVotes:
		swapInfo = composition.GetInfoRankedView(vote, true)
	}
	renderRequest(w, updateInformation(w, r, acc, composition.ViewVoteDocument), html, swapInfo)*/
}

func GetVoteViewService(w http.ResponseWriter, r *http.Request) {
	acc, admin := CheckUserPrivileges(r, database.Admin, database.HeadAdmin)

	html := composition.GetVoteViewPage(acc, chi.URLParam(r, "uuid"), admin,
		validation.Message{})
	viewVoteDocumentRenderRequest(w, r, acc, html)
}

var viewVoteDocumentRenderRequest = genericRenderer(composition.ViewVoteDocument)
var viewVoteDocumentOnlySwapMessage = genericMessageSwapper(composition.ViewVoteDocument)

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
	msg = create.CreateVote(acc.ID)
	if !msg.Positive {
		createVoteDocumentOnlySwapMessage(w, r, msg, acc)
		return
	}

	pushURL(w, "/"+string(composition.ViewVoteDocumentLink)+create.UUIDredirect)
	html := composition.GetVoteViewPage(acc, create.UUIDredirect,
		CheckIfHasRole(acc, database.HeadAdmin, database.Admin),
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
