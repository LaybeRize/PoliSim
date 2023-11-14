package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
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
}

func GetVoteCreationService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetCreateVotePage(acc, &validation.CreateVote{
		EndTime: time.Now().Add(time.Hour * 25).Format("2006-01-02T15:04"),
	}, validation.Message{})
	createVoteDocumentRenderRequest(w, r, acc, html)
}

func PostCreateVoteInDatabaseService(w http.ResponseWriter, r *http.Request) {
	var p validation.CreateVote
	err := json.NewDecoder(r.Body).Decode(&p)
	msg := "success"
	if err != nil {
		msg = "failur"
	}
	acc, _ := CheckUserPrivileges(r)
	createVoteDocumentOnlySwapMessage(w, r, validation.Message{Message: msg}, acc)
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
