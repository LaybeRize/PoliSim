package serving

import (
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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
	renderRequest(w, composition.GetCreateVotePage())
}

func PostCreateVoteInDatabaseService(w http.ResponseWriter, r *http.Request) {

}

func PatchGetVotePartial(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.ParseUint(chi.URLParam(r, "number"), 10, 64)
	if err != nil {
		i = 1337
	}
	renderRequest(w, composition.GetVotePartial(int64(i)))
}
