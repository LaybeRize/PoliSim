package logic

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var voteMutex = sync.Mutex{}

func CloseVoteIfTimeIsUp(ending time.Time, uuidStr string) {
	if ending.After(time.Now()) {
		return
	}
	voteMutex.Lock()
	defer voteMutex.Unlock()

	doc, err := extraction.GetDocument(uuidStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error finding vote document: "+err.Error())
		return
	}
	doc.Type = database.FinishedVote
	createSummaryForAllVotes(doc.UUID)
	err = extraction.UpdateDocument(doc)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error updating vote document: "+err.Error())
	}
}

func createSummaryForAllVotes(uuidStr string) {
	list, err := extraction.GetVotesForDocument(uuidStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error finding votes: "+err.Error())
		return
	}
	for _, item := range list {
		createSummaryForVote(&item)
	}
}

func createSummaryForVote(vote *database.Votes) {
	//TODO: add the logic
	vote.Finished = true
	err := extraction.UpdateVote(vote)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error updating vote with the UUID "+vote.UUID+": "+err.Error())
	}
}

func AddNewResultToVote(uuid string, resultKey string, result database.Results) error {
	voteMutex.Lock()
	defer voteMutex.Unlock()
	vote, err := extraction.GetSingleVote(uuid)

	if err != nil {
		return err
	}
	if _, ok := vote.Info.Results[resultKey]; ok {
		return errors.New("results for key already exists")
	}
	vote.Info.Results[resultKey] = result
	updateInfo(vote, resultKey)
	err = extraction.UpdateVote(vote)
	return err
}

func updateInfo(vote *database.Votes, key string) {
	//add values to summary
}
