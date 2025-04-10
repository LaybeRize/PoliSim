package database

import (
	loc "PoliSim/localisation"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type (
	VoteType int

	VoteInfo struct {
		ID       string `json:"id"`
		Question string `json:"question"`
	}

	VoteInstance struct {
		ID                    string
		DocumentID            string
		Question              string
		Answers               []string
		Type                  VoteType
		MaxVotes              int
		ShowVotesDuringVoting bool
		Anonymous             bool
		EndDate               time.Time
	}

	AccountVotes struct {
		Question     string             `json:"question"`
		Answers      []string           `json:"answers"`
		Anonymous    bool               `json:"anonymous"`
		Type         VoteType           `json:"type"`
		AnswerAmount int                `json:"answer_amount"`
		IllegalVotes []string           `json:"illegal_votes"`
		List         []SingleCastedVote `json:"list"`
		CSV          string             `json:"CSV"`
	}

	SingleCastedVote struct {
		Voter         string `json:"voter"`
		BallotNumbers []int  `json:"ballot"`
	}
)

func (a *AccountVotes) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AccountVotes) Scan(src interface{}) error {
	switch src.(type) {
	case []byte:
		return json.Unmarshal(src.([]byte), a)
	case string:
		return json.Unmarshal([]byte(src.(string)), a)
	default:
		return errors.New("value can not be unmarshalled into account vote")
	}
}

func (a *AccountVotes) GetEscapeCSV() string {
	return strings.ReplaceAll(url.QueryEscape(a.CSV), "+", "%20")
}

func (a *AccountVotes) GetHeaderWidth() int {
	return a.AnswerAmount + 1
}

func (a *AccountVotes) GetIllegalVotes() string {
	if a.Anonymous {
		return strconv.Itoa(len(a.IllegalVotes))
	}
	if len(a.IllegalVotes) == 0 {
		return loc.VoteNoIllegalVotesCasted
	}
	return strings.Join(a.IllegalVotes, ", ")
}

func (a *AccountVotes) NoVotes() bool {
	return len(a.List) == 0
}

func (a *AccountVotes) VoteIterator() func(func(string, []string) bool) {
	return func(yield func(string, []string) bool) {
		var name string
		for pos, votes := range a.List {
			newList := make([]string, a.AnswerAmount)
			for i, val := range votes.BallotNumbers {
				if val > 0 {
					newList[i] = strconv.Itoa(val)
				}
			}
			if a.Anonymous {
				name = strconv.Itoa(pos + 1)
			} else {
				name = votes.Voter
			}
			if !yield(name, newList) {
				return
			}
		}
	}
}

func (t VoteType) IsSingleVote() bool { return t == SingleVote }

func (t VoteType) IsMultipleVotes() bool { return t == MultipleVotes }

func (t VoteType) IsRankedVoting() bool { return t == RankedVoting }

func (t VoteType) IsVoteSharing() bool { return t == VoteShares }

func (v *VoteInstance) Ended() bool { return time.Now().After(v.EndDate) }

func (v *VoteInstance) HasValidType() bool { return v.Type >= SingleVote && v.Type <= VoteShares }

func (v *VoteInstance) AnswerLength() int { return len(v.Answers) }

func (v *VoteInstance) AnswerIterator() func(func(int, string) bool) {
	return func(yield func(int, string) bool) {
		for i, str := range v.Answers {
			if !yield(i+1, str) {
				return
			}
		}
	}
}

func (v *VoteInstance) GetAnswerAsList() string {
	return strings.Join(v.Answers, ", ")
}

func (v *VoteInstance) GetTimeEnd(a *Account) string {
	if a.Exists() {
		return v.EndDate.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return v.EndDate.Format(loc.TimeFormatString)
}

const (
	SingleVote VoteType = iota
	MultipleVotes
	RankedVoting
	VoteShares
)

func CreateOrUpdateVote(instance *VoteInstance, acc *Account, number int) error {
	accompaniment := &AccountVotes{
		Question:     instance.Question,
		Answers:      instance.Answers,
		Anonymous:    instance.Anonymous,
		Type:         instance.Type,
		AnswerAmount: instance.AnswerLength(),
		IllegalVotes: []string{},
		List:         []SingleCastedVote{},
		CSV:          "",
	}
	_, err := postgresDB.Exec(`INSERT INTO personal_votes (number, account_name, id, question, answers, type, max_votes, show_votes, anonymous, end_date, vote_info) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (number, account_name) DO UPDATE 
SET question = $4, answers = $5, type = $6, max_votes = $7, show_votes = $8, anonymous = $9, end_date = $10, vote_info = $11;`,
		number, acc.GetName(), instance.ID, instance.Question, pq.Array(instance.Answers), instance.Type,
		instance.MaxVotes, instance.ShowVotesDuringVoting, instance.Anonymous, instance.EndDate, accompaniment)
	return err
}

func GetVote(acc *Account, number int) (*VoteInstance, error) {
	props := &VoteInstance{}
	err := postgresDB.QueryRow(`SELECT id, question, answers, type, max_votes, show_votes, anonymous, end_date 
FROM personal_votes WHERE number = $1 AND account_name = $2;`,
		number, acc.GetName()).Scan(&props.ID, &props.Question, pq.Array(&props.Answers), &props.Type, &props.MaxVotes,
		&props.ShowVotesDuringVoting, &props.Anonymous, &props.EndDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return props, nil
}

func GetVoteInfoList(acc *Account) ([]VoteInfo, error) {
	result, err := postgresDB.Query(`SELECT id, question FROM personal_votes WHERE account_name = $1 ORDER BY number;`, acc.GetName())
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	list := make([]VoteInfo, 0)
	vote := VoteInfo{}
	for result.Next() {
		err = result.Scan(&vote.ID, &vote.Question)
		if err != nil {
			return nil, err
		}
		list = append(list, vote)
	}
	return list, err
}

func CastVoteWithAccount(name string, id string, votes []int) error {
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}
	defer rollback(tx)
	var docID string
	err = tx.QueryRow(`SELECT document_id FROM document_to_vote WHERE id = $1`, id).Scan(&docID)
	if errors.Is(err, sql.ErrNoRows) {
		return NotAllowedError
	} else if err != nil {
		return err
	}

	err = tx.QueryRow(`SELECT id FROM document_linked WHERE id = $1 AND end_time > $2 AND removed = false AND
                                     ((doc_account = $3 AND participant = true) OR 
                                      (organisation_account = $3 AND member_participation = true) OR 
                                      (organisation_account = $3 AND is_admin = true AND admin_participation = true)) LIMIT 1;`,
		docID, time.Now().UTC(), name).Scan(&docID)
	if errors.Is(err, sql.ErrNoRows) {
		return NotAllowedError
	} else if err != nil {
		return err
	}

	var emptyString string
	err = tx.QueryRow(`SELECT vote_id FROM has_voted WHERE account_name = $1 AND vote_id = $2`, name, id).Scan(&emptyString)
	switch true {
	case errors.Is(err, sql.ErrNoRows):
	case err != nil:
		return err
	default:
		return AlreadyVoted
	}

	var byteStr []byte
	if votes == nil {
		byteStr, err = json.Marshal(name)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`UPDATE document_to_vote 
SET vote_info = jsonb_insert(vote_info, '{illegal_votes, -1}', $2, true) WHERE id = $1;`,
			id, string(byteStr))
	} else {
		byteStr, err = json.Marshal(&SingleCastedVote{Voter: name, BallotNumbers: votes})
		if err != nil {
			return err
		}
		_, err = tx.Exec(`UPDATE document_to_vote SET vote_info = jsonb_insert(vote_info, '{list, -1}', $2, true) WHERE id = $1;`,
			id, string(byteStr))
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO has_voted (account_name, vote_id) VALUES ($1, $2);`, name, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetAnswersAndTypeForVote(id string, acc *Account) ([]string, VoteType, int, error) {
	var docID string
	var answers []string
	var vType VoteType
	var maxVotes int
	err := postgresDB.QueryRow(`SELECT document_id, answers, type, max_votes FROM document_to_vote WHERE id = $1;`, id).Scan(&docID, pq.Array(&answers), &vType, &maxVotes)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, -1, -1, NotAllowedError
	} else if err != nil {
		return nil, -1, -1, err
	}

	err = postgresDB.QueryRow(`SELECT id FROM document_linked 
          WHERE end_time > $2 AND id = $1 AND (public = true OR owner_name = $3) LIMIT 1;`,
		docID, time.Now().UTC(), acc.GetName()).Scan(&docID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, -1, -1, NotAllowedError
	} else if err != nil {
		return nil, -1, -1, err
	}
	return answers, vType, maxVotes, nil
}

func GetVoteForUser(id string, acc *Account) (*VoteInstance, *AccountVotes, error) {
	voteInstance := &VoteInstance{}
	accountVotes := &AccountVotes{}
	err := postgresDB.QueryRow(`SELECT id, document_id, question, answers, type, max_votes, 
       show_votes, anonymous, end_date, vote_info  FROM document_to_vote WHERE id = $1;`, id).Scan(
		&voteInstance.ID, &voteInstance.DocumentID, &voteInstance.Question, pq.Array(&voteInstance.Answers), &voteInstance.Type,
		&voteInstance.MaxVotes, &voteInstance.ShowVotesDuringVoting, &voteInstance.Anonymous, &voteInstance.EndDate,
		accountVotes)
	if err != nil {
		return nil, nil, err
	}
	err = postgresDB.QueryRow(`SELECT id FROM document_linked 
          WHERE id = $1 AND (public = true OR owner_name = $2 OR $3 = true) LIMIT 1;`,
		voteInstance.DocumentID, acc.GetName(), acc.IsAtLeastAdmin()).Scan(&voteInstance.DocumentID)
	if err != nil {
		return nil, nil, err
	}
	return voteInstance, accountVotes, nil
}

func resultRoutine() {
	curr := time.Now().UTC()
	next := time.Date(curr.Year(), curr.Month(), curr.Day()+1, 0, 0, 0, 0, time.UTC)
	ticker := time.NewTicker(next.Sub(curr) + time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		generateResults()
		curr = time.Now().UTC()
		next = time.Date(curr.Year(), curr.Month(), curr.Day()+1, 0, 0, 0, 0, time.UTC)
		ticker.Reset(next.Sub(curr) + time.Second)
	}
}

func generateResults() {
	shutdown.Lock()
	defer shutdown.Unlock()

	log.Println("Cleaning Vote Results")
	result, err := postgresDB.Query(`SELECT id FROM document WHERE end_time < $1 AND type = $2;`, time.Now().UTC(), DocTypeVote)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer closeRows(result)
	for result.Next() {
		var id string
		err = result.Scan(&id)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		transactionForSingleDocument(id)
	}
}

func transactionForSingleDocument(documentID string) {
	tx, err := postgresDB.Begin()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer rollback(tx)

	result, err := tx.Query(`SELECT id, vote_info FROM document_to_vote WHERE document_id = $1 ORDER BY id;`, documentID)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer closeRows(result)

	accountVotesArray := make([]AccountVotes, 0)
	voteIDs := make([]string, 0)
	accVote := AccountVotes{}
	voteID := ""
	for result.Next() {
		err = result.Scan(&voteID, &accVote)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		generateCSV(&accVote)
		accountVotesArray = append(accountVotesArray, accVote)
		voteIDs = append(voteIDs, voteID)
	}

	_, err = tx.Exec(`DELETE FROM has_voted WHERE vote_id = ANY($1);`, pq.Array(voteIDs))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	_, err = tx.Exec(`DELETE FROM document_to_vote WHERE document_id = $1;`, documentID)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	byteArr, err := json.Marshal(accountVotesArray)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	_, err = tx.Exec(`UPDATE document SET extra_info = jsonb_set(extra_info, '{result}', $2, false) WHERE id = $1;`,
		documentID, string(byteArr))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	_, err = tx.Exec(`UPDATE document SET extra_info = jsonb_set(extra_info, '{links}', '[]', false) WHERE id = $1;`,
		documentID)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	_ = tx.Commit()
	return
}

func generateCSV(vote *AccountVotes) {
	(*vote).CSV = "\"" + strings.ReplaceAll(vote.Question, "\"", "\"\"") + "\""
	for _, answer := range vote.Answers {
		(*vote).CSV += ",\"" + strings.ReplaceAll(answer, "\"", "\"\"") + "\""
	}
	for name, votes := range vote.VoteIterator() {
		if vote.Anonymous {
			(*vote).CSV += "\n,"
		} else {
			(*vote).CSV += "\n\"" + strings.ReplaceAll(name, "\"", "\"\"") + "\","
		}

		if vote.Type != RankedVoting {
			for i, count := range votes {
				if count == "" {
					votes[i] = "0"
				}
			}
		}
		(*vote).CSV += strings.Join(votes, ",")
	}
	(*vote).CSV += fmt.Sprintf("\nInvalid,%d", len(vote.IllegalVotes)) + strings.Repeat(",", len(vote.Answers)-1)
}
