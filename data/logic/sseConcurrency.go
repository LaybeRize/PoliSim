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
var watchComments = map[int64]int{}
var watchVotes = map[int64]int{}

var sseCommentMutex = sync.Mutex{}
var sseVoteMutex = sync.Mutex{}

func AddCommentChannel(uuid string, userID string, accountID int64, NewChan chan CommentUpdate) {
	sseCommentMutex.Lock()
	defer sseCommentMutex.Unlock()
	if _, ok := watchComments[accountID]; !ok {
		watchComments[accountID] = 1
	} else {
		watchComments[accountID] += 1
	}
	if _, ok := updateComments[uuid]; !ok {
		updateComments[uuid] = make(map[string]chan CommentUpdate)
	}
	updateComments[uuid][userID] = NewChan
}

func RemoveCommentChannel(uuid string, userID string, accountID int64) {
	sseCommentMutex.Lock()
	defer sseCommentMutex.Unlock()
	watchComments[accountID] -= 1
	delete(updateComments[uuid], userID)
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

func AddVoteChannel(uuid string, userID string, accountID int64, NewChan chan database.Votes) {
	sseVoteMutex.Lock()
	defer sseVoteMutex.Unlock()
	if _, ok := watchVotes[accountID]; !ok {
		watchVotes[accountID] = 1
	} else {
		watchVotes[accountID] += 1
	}
	if _, ok := updateVotes[uuid]; !ok {
		updateVotes[uuid] = make(map[string]chan database.Votes)
	}
	updateVotes[uuid][userID] = NewChan
}

func RemoveVoteChannel(uuid string, userID string, accountID int64) {
	sseVoteMutex.Lock()
	defer sseVoteMutex.Unlock()
	watchVotes[accountID] -= 1
	delete(updateVotes[uuid], userID)
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
