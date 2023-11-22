package logic

import (
	"PoliSim/data/database"
	"sync"
)

type CommentUpdate struct {
	Change     bool
	Discussion database.Discussions
}

var updateComments = map[string]map[int64]chan CommentUpdate{}
var updateVotes = map[string]map[int64]chan database.Votes{}

var sseCommentMutex = sync.Mutex{}
var sseVoteMutex = sync.Mutex{}

func AddCommentChannel(uuid string, accountID int64, NewChan chan CommentUpdate) {
	sseCommentMutex.Lock()
	defer sseCommentMutex.Unlock()
	if _, ok := updateComments[uuid]; !ok {
		updateComments[uuid] = make(map[int64]chan CommentUpdate)
	}
	updateComments[uuid][accountID] = NewChan
}

func RemoveCommentChannel(uuid string, accountID int64) {
	sseCommentMutex.Lock()
	defer sseCommentMutex.Unlock()
	delete(updateComments[uuid], accountID)
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

func AddVoteChannel(uuid string, accountID int64, NewChan chan database.Votes) {
	sseVoteMutex.Lock()
	defer sseVoteMutex.Unlock()
	if _, ok := updateVotes[uuid]; !ok {
		updateVotes[uuid] = make(map[int64]chan database.Votes)
	}
	updateVotes[uuid][accountID] = NewChan
}

func RemoveVoteChannel(uuid string, accountID int64) {
	sseVoteMutex.Lock()
	defer sseVoteMutex.Unlock()
	delete(updateVotes[uuid], accountID)
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
