package database

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type (
	VoteType int

	VoteInfo struct {
		ID       string
		Question string
	}

	VoteInstance struct {
		ID                    string
		DocumentID            string
		Question              string
		Answers               []string
		IterableAnswers       []any
		Type                  VoteType
		MaxVotes              int
		ShowVotesDuringVoting bool
		Anonymous             bool
		EndDate               time.Time
	}
)

func (t VoteType) IsSingleVote() bool      { return t == SingleVote }
func (v *VoteInstance) IsSingleVote() bool { return v.Type.IsSingleVote() }

func (t VoteType) IsMultipleVotes() bool      { return t == MultipleVotes }
func (v *VoteInstance) IsMultipleVotes() bool { return v.Type.IsMultipleVotes() }

func (t VoteType) IsRankedVoting() bool      { return t == RankedVoting }
func (v *VoteInstance) IsRankedVoting() bool { return v.Type.IsRankedVoting() }

func (t VoteType) IsVoteSharing() bool      { return t == VoteShares }
func (v *VoteInstance) IsVoteSharing() bool { return v.Type.IsVoteSharing() }

func (v *VoteInstance) Ended() bool { return time.Now().After(v.EndDate) }

func (v *VoteInstance) HasValidType() bool { return v.Type >= SingleVote && v.Type <= VoteShares }

func (v *VoteInstance) AnswerIterator() func(func(int, string) bool) {
	return func(yield func(int, string) bool) {
		for i, v := range v.IterableAnswers {
			if !yield(i, v.(string)) {
				return
			}
		}
	}
}

func (v *VoteInstance) ConvertToAnswer() {
	v.Answers = make([]string, len(v.IterableAnswers))
	for i, str := range v.AnswerIterator() {
		v.Answers[i] = str
	}
}

const (
	SingleVote VoteType = iota
	MultipleVotes
	RankedVoting
	VoteShares
)

func CreateOrUpdateVote(instance *VoteInstance, acc *Account, number int) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `MATCH (acc:Account)-[r:MANAGES]->(v:Vote) 
WHERE acc.name = $Manager AND r.position = $position 
RETURN v.id;`, map[string]any{
		"Manager":  acc.Name,
		"position": number})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	} else if result.Next(ctx); result.Record() == nil {
		_, err = tx.Run(ctx, `MATCH (a:Account) WHERE a.name = $Manager 
CREATE (v:Vote {id: $id, type: $type, question: $question, answers: $answers, max_votes: $max_votes, 
show_during: $show_during, anonymous: $anonymous}) 
MERGE (a)-[:MANAGES {position: $position}]->(v);`, map[string]any{
			"Manager":     acc.Name,
			"position":    number,
			"id":          instance.ID,
			"type":        instance.Type,
			"question":    instance.Question,
			"answers":     instance.Answers,
			"max_votes":   instance.MaxVotes,
			"show_during": instance.ShowVotesDuringVoting,
			"anonymous":   instance.Anonymous})
	} else {
		instance.ID = result.Record().Values[0].(string)

		_, err = tx.Run(ctx, `
MATCH (v:Vote) WHERE v.id = $id 
SET v.type = $type, v.question = $question, v.answers = $answers, v.max_votes = $max_votes, 
v.show_during = $show_during, v.anonymous = $anonymous;`, map[string]any{
			"id":          instance.ID,
			"type":        instance.Type,
			"question":    instance.Question,
			"answers":     instance.Answers,
			"max_votes":   instance.MaxVotes,
			"show_during": instance.ShowVotesDuringVoting,
			"anonymous":   instance.Anonymous})
	}

	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func GetVote(acc *Account, number int) (*VoteInstance, error) {
	result, err := makeRequest(`MATCH (a:Account)-[r:MANAGES]->(v:Vote) 
WHERE a.name = $name AND r.position = $position
RETURN v;`, map[string]any{"name": acc.Name, "position": number})
	if err != nil {
		return nil, err
	}
	if len(result.Records) == 0 {
		return nil, nil
	}
	record := result.Records[0].Values[0].(neo4j.Node)
	return &VoteInstance{
		ID:                    record.Props["id"].(string),
		Question:              record.Props["question"].(string),
		IterableAnswers:       record.Props["answers"].([]any),
		Type:                  VoteType(record.Props["type"].(int)),
		MaxVotes:              record.Props["max_votes"].(int),
		ShowVotesDuringVoting: record.Props["show_during"].(bool),
		Anonymous:             record.Props["anonymous"].(bool),
	}, nil
}

func GetVoteInfoList(acc *Account) ([]VoteInfo, error) {
	result, err := makeRequest(`MATCH (a:Account)-[:MANAGES]->(v:Vote) 
WHERE a.name = $name RETURN v.id, v.question;`, map[string]any{"name": acc.Name})
	if err != nil {
		return nil, err
	}
	list := make([]VoteInfo, len(result.Records))
	for i, record := range result.Records {
		list[i] = VoteInfo{
			ID:       record.Values[0].(string),
			Question: record.Values[1].(string),
		}
	}
	return list, err
}
