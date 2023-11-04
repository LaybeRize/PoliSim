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
		Results     map[string]Results `json:"results"`
		Summary     Summary            `json:"summary"`
		VoteMethod  VoteType           `json:"voteMethod"`
		MaxPosition int                `json:"maxPosition"`
		Options     []string           `json:"options"`
	}
	Results struct {
		Voter       string         `json:"voter"`
		InvalidVote bool           `json:"invalid"`
		Votes       map[string]int `json:"votes"`
	}
	Summary struct {
		Sums         map[string]int            `json:"sums"`
		RankedMap    map[string]map[string]int `json:"rankedMap"`
		Person       map[string]string         `json:"person"` //the option the person voted for
		InvalidVotes []string                  `json:"invalidVotes"`
		CSV          string                    `json:"csv"` //saves the data as a CSV for the ranked Map
	}
)

const (
	SingleVote          VoteType = "single_vote"
	MultipleVotes       VoteType = "multiple_votes"
	VoteRanking         VoteType = "vote_ranking"
	ThreeCategoryVoting VoteType = "three_category_voting" //for against neutral
)

// TODO: add translation to .json
var (
	VoteTypes       = []VoteType{SingleVote, MultipleVotes, VoteRanking, ThreeCategoryVoting}
	VoteTranslation = map[VoteType]string{
		SingleVote:          "Einzelstimmenwahl",
		MultipleVotes:       "Mehrstimmenwahl",
		VoteRanking:         "Gewichtete Wahl",
		ThreeCategoryVoting: "Dafür-Dagegen-Enthaltung-Wahl",
	}
)
