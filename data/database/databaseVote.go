package database

type (
	VoteType  string
	VotesList []Votes
	Votes     struct {
		UUID                   string `gorm:"primaryKey"`
		Parent                 string
		Question               string
		ShowNumbersWhileVoting bool
		ShowNamesWhileVoting   bool
		ShowNamesAfterVoting   bool
		Finished               bool
		Info                   VoteInfo `gorm:"type:jsonb;serializer:json"`
	}
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
)

const (
	SingleVote          VoteType = "single_vote"
	MultipleVotes       VoteType = "multiple_votes"
	RankedVotes         VoteType = "ranked_votes"
	ThreeCategoryVoting VoteType = "three_category_voting" // for/against/neutral
)

var (
	VoteTypes       = []VoteType{SingleVote, MultipleVotes, RankedVotes, ThreeCategoryVoting}
	VoteTranslation = map[VoteType]string{
		SingleVote:          "!placeholder!",
		MultipleVotes:       "!placeholder!",
		RankedVotes:         "!placeholder!",
		ThreeCategoryVoting: "!placeholder!",
	}
)

func (vType VoteType) String() string {
	return string(vType)
}
