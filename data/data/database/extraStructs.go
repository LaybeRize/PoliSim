package database

import "time"

type (
	// Extra Info for Documents
	DocumentInfo struct {
		Finishing  time.Time     `json:"time"`
		Post       []Posts       `json:"post"`
		Discussion []Discussions `json:"discussion"`
	}
	Posts struct {
		UUID      string    `json:"uuid"`
		Hidden    bool      `json:"hidden"`
		Submitted time.Time `json:"submitted"`
		Info      string    `json:"info"`
		Color     string    `json:"color"`
	}
	Discussions struct {
		UUID        string    `json:"uuid"`
		Hidden      bool      `json:"hidden"`
		Written     time.Time `json:"written"`
		Author      string    `json:"author"`
		Flair       string    `json:"flair"`
		HTMLContent string    `json:"htmlContent"`
	}

	// Extra Info for Letters
	LetterInfo struct {
		AllHaveToAgree     bool     `json:"allAgree"`
		NoSigning          bool     `json:"noSigning"`
		PeopleNotYetSigned []string `json:"notSigned"`
		Signed             []string `json:"signed"`
		Rejected           []string `json:"rejected"`
	}

	// Extra Info for Votes
	VoteInfo struct {
		VoteOrder  []string           `json:"voteOrder"` //order in which people voted
		Results    map[string]Results `json:"results"`   //the key is the voter name
		Summary    Summary            `json:"summary"`
		VoteMethod VoteType           `json:"voteMethod"`
		Options    []string           `json:"options"`
	}
	Results struct {
		InvalidVote bool             `json:"invalid"`
		Votes       map[string]int64 `json:"votes"`
	}
	Summary struct {
		Sums         map[string]int64 `json:"sums"`         //scores of every answer
		InvalidVotes []string         `json:"invalidVotes"` //list of all people that gave an invalid vote
		CSV          string           `json:"csv"`          //saves the data as a CSV for the ranked Map
	}
	VoteType string
)
