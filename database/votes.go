package database

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log/slog"
	"strconv"
	"strings"
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
		MaxVotes              int64
		ShowVotesDuringVoting bool
		Anonymous             bool
		EndDate               time.Time
	}

	AccountVotes struct {
		Anonymous    bool
		Type         VoteType
		AnswerAmount int
		IllegalVotes []string
		Illegal      []any
		List         map[string][]any
		BaseMap      map[string]any
	}
)

func (a *AccountVotes) GetIllegalVotes() string {
	if a.IllegalVotes == nil {
		if a.Anonymous {
			return strconv.Itoa(len(a.Illegal))
		}

		a.IllegalVotes = make([]string, len(a.Illegal))
		for i, value := range a.Illegal {
			a.IllegalVotes[i] = value.(string)
		}
	} else if a.Anonymous {
		return strconv.Itoa(len(a.IllegalVotes))
	}
	return strings.Join(a.IllegalVotes, ", ")
}

func (a *AccountVotes) VoteIterator() func(func(string, []string) bool) {
	if a.List == nil {
		return func(yield func(string, []string) bool) {
			pos := 0
			for name, list := range a.BaseMap {
				newList := make([]string, len(list.([]any)))
				for i, val := range list.([]any) {
					newList[i] = strconv.Itoa(int(val.(int64)))
					if a.Type.IsRankedVoting() && newList[i] == "-1" {
						newList[i] = ""
					}
				}
				if a.Anonymous {
					pos += 1
					name = strconv.Itoa(pos)
				}
				if !yield(name, newList) {
					return
				}
			}
		}
	}
	return func(yield func(string, []string) bool) {
		pos := 0
		for name, list := range a.List {
			newList := make([]string, len(list))
			for i, val := range list {
				newList[i] = strconv.Itoa(int(val.(int64)))
				if a.Type.IsRankedVoting() && newList[i] == "-1" {
					newList[i] = ""
				}
			}
			if a.Anonymous {
				pos += 1
				name = strconv.Itoa(pos)
			}
			if !yield(name, newList) {
				return
			}
		}
	}
}

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
		for i, str := range v.IterableAnswers {
			if !yield(i, str.(string)) {
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

func (v *VoteInstance) GetTimeEnd(a *Account) string {
	if a.Exists() {
		return v.EndDate.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return v.EndDate.Format("2006-01-02 15:04:05 MST")
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
		Type:                  VoteType(record.Props["type"].(int64)),
		MaxVotes:              record.Props["max_votes"].(int64),
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

func CastVoteWithAccount(name string, id string, votes []int) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `CALL {
MATCH (a:Account)<-[:PARTICIPANT]-(d:Document)-[:LINKS]->(v:Vote)
WHERE a.name = $author AND a.blocked = false AND v.id = $id AND $now > d.end_time
RETURN a  
UNION 
MATCH (a:Account)-[:USER|ADMIN]->(:Organisation)<-[:IN]-(d:Document)-[:LINKS]->(v:Vote) 
WHERE a.name = $author AND a.blocked = false AND v.id = $id AND d.member_part = true AND $now > d.end_time
RETURN a 
UNION 
MATCH (a:Account)-[:ADMIN]->(:Organisation)<-[:IN]-(d:Document)-[:LINKS]->(v:Vote) 
WHERE a.name = $author AND a.blocked = false AND v.id = $id AND d.admin_part = true AND $now > d.end_time
RETURN a  
} 
RETURN a;`,
		map[string]any{"id": id, "author": name, "now": time.Now().UTC()})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	} else if result.Next(ctx); result.Record() == nil {
		_ = tx.Rollback(ctx)
		return notAllowedError
	}

	result, err = tx.Run(ctx, `MATCH (a:Account)-[:VOTED]->(v:Vote) WHERE a.name = $author AND v.id = $id 
RETURN a;`,
		map[string]any{"id": id, "author": name})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	} else if result.Next(ctx); result.Record() != nil {
		_ = tx.Rollback(ctx)
		return notAllowedError
	}

	mapIsNil := votes == nil
	_, err = tx.Run(ctx, `MATCH (a:Account), (v:Vote) WHERE a.name = $author AND v.id = $id 
MERGE (a)-[:VOTED {written: $now, illegal: $illegal, vote: $voteMap}]->(v);`,
		map[string]any{"id": id, "author": name, "illegal": mapIsNil, "voteMap": votes, "now": time.Now().UTC()})
	if err == nil {
		return tx.Commit(ctx)
	}
	return err
}

func GetVoteForUser(id string, acc *Account) (*VoteInstance, *AccountVotes, error) {
	extra := `CALL { MATCH (o:Organisation)<-[:IN]-(d:Document)-[:LINKS]->(v:Vote) WHERE d.public = true 
RETURN o
UNION
MATCH (a:Account)-[*..]->(o:Organisation)<-[:IN]-(d:Document) WHERE d.public = false AND a.name = $name 
RETURN o }`
	if acc.IsAtLeastAdmin() {
		extra = ""
	}

	result, err := makeRequest(`MATCH (d:Document)-[:LINKS]->(v:Vote) WHERE v.id = $id 
`+extra+` RETURN d.id, d.end_time, v;`, map[string]any{"id": id, "name": acc.Name})
	if err != nil {
		return nil, nil, err
	}
	Props := result.Records[0].Values[2].(neo4j.Node).Props
	vote := &VoteInstance{
		DocumentID:            result.Records[0].Values[0].(string),
		ID:                    Props["id"].(string),
		Question:              Props["question"].(string),
		IterableAnswers:       Props["answers"].([]any),
		Type:                  VoteType(Props["type"].(int64)),
		MaxVotes:              Props["max_votes"].(int64),
		ShowVotesDuringVoting: Props["show_during"].(bool),
		Anonymous:             Props["anonymous"].(bool),
		EndDate:               result.Records[0].Values[1].(time.Time),
	}
	var voteList *AccountVotes
	if vote.ShowVotesDuringVoting {
		voteList = &AccountVotes{
			Type:         vote.Type,
			Anonymous:    vote.Anonymous,
			AnswerAmount: len(vote.IterableAnswers),
			IllegalVotes: []string{},
			List:         make(map[string][]any),
		}
		result, err = makeRequest(`MATCH (a:Account)-[r:VOTED]->(v:Vote) WHERE v.id = $id RETURN r, a.name 
ORDER BY r.written;`,
			map[string]any{"id": id})
		if err != nil {
			return nil, nil, err
		}
		for _, record := range result.Records {
			props := record.Values[0].(neo4j.Node).Props
			if props["illegal"].(bool) {
				voteList.IllegalVotes = append(voteList.IllegalVotes, record.Values[1].(string))
				continue
			}
			voteList.List[record.Values[1].(string)] = props["vote"].([]any)
		}
	}

	return vote, voteList, err
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

	result, err := makeRequest(`MATCH (a:Account)-[r:VOTED]->(v:Vote)<-[:LINKS]-(d:Document) WHERE $now < d.end_time 
RETURN v, d.id AS docID,a.name AS accName, r ORDER BY r.written 
RETURN v, docID, collect(accName), collect(r);`, map[string]any{"now": time.Now().UTC()})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	votes, docIds, voteIds := transformVotesForResults(result)

	for i, vote := range votes {
		_, err = makeRequest(`MATCH (v:Vote) WHERE v.id = $vote 
MATCH (d:Document) WHERE d.id = $id 
CREATE (r:Result {type: $Type, anonymous: $Anonymous, amount: $AnswerAmount, illegal: $IllegalVotes, 
list: $List}) 
MERGE (d)-[:VOTED]->(r) 
DETACH DELETE v;`,
			map[string]any{
				"vote":         voteIds[i],
				"id":           docIds[i],
				"Type":         vote.Type,
				"Anonymous":    vote.Anonymous,
				"AnswerAmount": vote.AnswerAmount,
				"IllegalVotes": vote.IllegalVotes,
				"List":         vote.List,
			})
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// returns first the docID then the VoteID array
func transformVotesForResults(result *neo4j.EagerResult) ([]AccountVotes, []string, []string) {
	votes := make([]AccountVotes, len(result.Records))
	docIDs := make([]string, len(result.Records))
	voteIDs := make([]string, len(result.Records))

	for i, record := range result.Records {
		docIDs[i] = record.Values[1].(string)
		voteProps := record.Values[0].(neo4j.Node).Props
		voteIDs[i] = voteProps["id"].(string)

		votes[i] = AccountVotes{
			Anonymous:    voteProps["anonymous"].(bool),
			Type:         VoteType(voteProps["type"].(int64)),
			AnswerAmount: len(voteProps["answers"].([]any)),
			IllegalVotes: []string{},
			List:         make(map[string][]any),
		}

		nameList := record.Values[2].([]any)
		voteList := record.Values[3].([]any)

		for j, name := range nameList {
			props := voteList[j].(neo4j.Node).Props
			if props["illegal"].(bool) {
				votes[i].IllegalVotes = append(votes[i].IllegalVotes, name.(string))
				continue
			}
			votes[i].List[record.Values[1].(string)] = props["vote"].([]any)
		}
	}

	return votes, docIDs, voteIDs
}
