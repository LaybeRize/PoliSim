package serving

import (
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"net/http"
)

//TODO fill this with live

func InstallVoteCreation() {
	composition.PageTitleMap[composition.CreateVote] = builder.Translation["startPageTitle"]
	composition.SidebarTitleMap[composition.CreateVote] = builder.Translation["startSidebarText"]
	composition.GetHTMXFunctions[composition.CreateVote] = GetVoteCreationService
	composition.PostHTMXFunctions[composition.CreateVote] = PostCreateVoteInDatabaseService
	composition.PostHTMXFunctions[composition.RequestVotePartial] = PostGetVotePartial
}

func GetVoteCreationService(w http.ResponseWriter, r *http.Request) {

}

func PostCreateVoteInDatabaseService(w http.ResponseWriter, r *http.Request) {

}

func PostGetVotePartial(w http.ResponseWriter, r *http.Request) {

}
