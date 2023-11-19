package logic

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/html/builder"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		createCSVForVote(&item)
	}
}

const seperator = ","

func createCSVForVote(vote *database.Votes) {
	csvStr := builder.Translation["csvVoterColumnName"] + "," + builder.Translation["csvVoterMadeInvalidVote"]
	for _, opt := range vote.Info.Options {
		csvStr += seperator + "\"" + strings.ReplaceAll(opt, "\"", "\"\"") + "\""
	}
	csvStr += "\n"
	for _, person := range vote.Info.VoteOrder {
		addition := ""
		for _, opt := range vote.Info.Options {
			val := vote.Info.Results[person].Votes[opt]
			addition += seperator + strconv.FormatInt(val, 10)
		}
		if !vote.ShowNamesAfterVoting {
			person = builder.Translation["csvVoterNameObscure"]
		} else {
			person = "\"" + strings.ReplaceAll(person, "\"", "\"\"") + "\""
		}
		addition = person + seperator + "false" + addition + "\n"
		csvStr += addition
	}
	addition := strings.Repeat(seperator, len(vote.Info.Options))
	for _, person := range vote.Info.Summary.InvalidVotes {
		if !vote.ShowNamesAfterVoting {
			person = builder.Translation["csvVoterNameObscure"]
		}
		csvStr += person + seperator + "true" + addition + "\n"
	}
	csvStr += builder.Translation["csvSumVoterName"] + seperator
	addition = ""
	for _, opt := range vote.Info.Options {
		val := vote.Info.Summary.Sums[opt]
		addition += seperator + strconv.FormatInt(val, 10)
	}
	csvStr += addition
	vote.Info.Summary.CSV = csvStr
	vote.Finished = true
	err := extraction.UpdateVote(vote)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error updating vote with the UUID "+vote.UUID+": "+err.Error())
	}
}

const ResultExistsErrorText = "results for key already exists"

func AddNewResultToVote(uuid string, resultKey string, result database.Results) (*database.Votes, error) {
	voteMutex.Lock()
	defer voteMutex.Unlock()
	vote, err := extraction.GetSingleVote(uuid)

	if err != nil {
		return nil, err
	}
	if _, ok := vote.Info.Results[resultKey]; ok {
		return nil, errors.New(ResultExistsErrorText)
	}
	vote.Info.Results[resultKey] = result
	if result.InvalidVote {
		vote.Info.Summary.InvalidVotes = append(vote.Info.Summary.InvalidVotes, resultKey)
	} else {
		vote.Info.VoteOrder = append(vote.Info.VoteOrder, resultKey)
	}

	for option, value := range vote.Info.Results[resultKey].Votes {
		vote.Info.Summary.Sums[option] += value
	}

	err = extraction.UpdateVote(vote)
	return vote, err
}
