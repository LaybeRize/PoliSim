package htmlServer

import (
	"PoliSim/componentHelper"
	"PoliSim/htmlComposition"
	"net/http"
)

//TODO fill this with live

func InstallVoteCreation() {
	htmlComposition.PageTitleMap[htmlComposition.CreateVote] = componentHelper.Translation["startPageTitle"]
	htmlComposition.SidebarTitleMap[htmlComposition.CreateVote] = componentHelper.Translation["startSidebarText"]
	htmlComposition.GetHTMXFunctions[htmlComposition.CreateVote] = GetVoteCreationService
	htmlComposition.PostHTMXFunctions[htmlComposition.CreateVote] = PostCreateVoteInDatabaseService
	htmlComposition.PostHTMXFunctions[htmlComposition.RequestVotePartial] = PostGetVotePartial
}

func GetVoteCreationService(w http.ResponseWriter, r *http.Request) {

}

func PostCreateVoteInDatabaseService(w http.ResponseWriter, r *http.Request) {

}

func PostGetVotePartial(w http.ResponseWriter, r *http.Request) {

}
