package logic

import (
	"PoliSim/data/database"
	"sync"
)

type CommentUpdate struct {
	Change     bool
	Discussion database.Discussions
}

var updateComments = map[string]map[string]chan CommentUpdate{}
var updateVotes = map[string]map[string]chan database.Votes{}

var sseCommentMutex = sync.Mutex{}
var sseVoteMutex = sync.Mutex{}

func AddCommentChannel(uuid string, routineID string, NewChan chan CommentUpdate) {
	sseCommentMutex.Lock()
	defer sseCommentMutex.Unlock()
	if _, ok := updateComments[uuid]; !ok {
		updateComments[uuid] = make(map[string]chan CommentUpdate)
	}
	updateComments[uuid][routineID] = NewChan
}

func RemoveCommentChannel(uuid string, routineID string) {
	sseCommentMutex.Lock()
	defer sseCommentMutex.Unlock()
	delete(updateComments[uuid], routineID)
	if len(updateComments[uuid]) == 0 {
		delete(updateComments, uuid)
	}
}

func SendCommentToChannels(uuid string, disc CommentUpdate) {
	sseCommentMutex.Lock()
	defer sseCommentMutex.Unlock()
	for _, channel := range updateComments[uuid] {
		channel <- disc
	}
}

func AddVoteChannel(uuid string, routineID string, NewChan chan database.Votes) {
	sseVoteMutex.Lock()
	defer sseVoteMutex.Unlock()
	if _, ok := updateVotes[uuid]; !ok {
		updateVotes[uuid] = make(map[string]chan database.Votes)
	}
	updateVotes[uuid][routineID] = NewChan
}

func RemoveVoteChannel(uuid string, routineID string) {
	sseVoteMutex.Lock()
	defer sseVoteMutex.Unlock()
	delete(updateVotes[uuid], routineID)
	if len(updateVotes[uuid]) == 0 {
		delete(updateVotes, uuid)
	}
}

func SendVoteToChannels(uuid string, disc database.Votes) {
	sseVoteMutex.Lock()
	defer sseVoteMutex.Unlock()
	for _, channel := range updateVotes[uuid] {
		channel <- disc
	}
}
