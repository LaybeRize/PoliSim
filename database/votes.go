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
		MaxVotes              int
		ShowVotesDuringVoting bool
		Anonymous             bool
		EndDate               time.Time
	}

	AccountVotes struct {
		Question        string
		IterableAnswers []any
		Anonymous       bool
		Type            VoteType
		AnswerAmount    int
		IllegalVotes    []string
		Illegal         []any
		List            map[string][]any
		Voter           []any
		Votes           []any
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
	if len(a.IllegalVotes) == 0 {
		return "Keine"
	}
	return strings.Join(a.IllegalVotes, ", ")
}

func (a *AccountVotes) NoVotes() bool {
	if a.List == nil {
		return len(a.Voter) == 0
	}
	return len(a.List) == 0
}

func (a *AccountVotes) VoteIterator() func(func(string, []string) bool) {
	if a.List == nil {
		return func(yield func(string, []string) bool) {
			pos := 0
			for _, name := range a.Voter {
				newList := make([]string, a.AnswerAmount)
				for i, val := range a.Votes[a.AnswerAmount*pos : (pos+1)*a.AnswerAmount] {
					newVal := int(val.(int64))
					if newVal > 0 {
						newList[i] = strconv.Itoa(newVal)
					}
				}
				pos += 1
				if a.Anonymous {
					name = strconv.Itoa(pos)
				}
				if !yield(name.(string), newList) {
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

func (t VoteType) IsSingleVote() bool { return t == SingleVote }

func (t VoteType) IsMultipleVotes() bool { return t == MultipleVotes }

func (t VoteType) IsRankedVoting() bool { return t == RankedVoting }

func (t VoteType) IsVoteSharing() bool { return t == VoteShares }

func (v *VoteInstance) Ended() bool { return time.Now().After(v.EndDate) }

func (v *VoteInstance) HasValidType() bool { return v.Type >= SingleVote && v.Type <= VoteShares }

func (v *VoteInstance) AnswerLength() int { return len(v.IterableAnswers) }

func (v *VoteInstance) AnswerIterator() func(func(int, string) bool) {
	return func(yield func(int, string) bool) {
		for i, str := range v.IterableAnswers {
			if !yield(i+1, str.(string)) {
				return
			}
		}
	}
}

func (v *VoteInstance) ConvertToAnswer() {
	v.Answers = make([]string, len(v.IterableAnswers))
	for i, str := range v.AnswerIterator() {
		v.Answers[i-1] = str
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
	defer tx.Close()
	if err != nil {
		return err
	}

	result, err := tx.Run(`MATCH (acc:Account)-[r:MANAGES]->(v:Vote) 
WHERE acc.name = $Manager AND r.position = $position 
RETURN v.id;`, map[string]any{
		"Manager":  acc.Name,
		"position": number})
	if err != nil {
		return err
	} else if !result.Next() {
		err = tx.RunWithoutResult(`MATCH (a:Account) WHERE a.name = $Manager 
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

		err = tx.RunWithoutResult(`
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
		return err
	}

	return tx.Commit()
}

func GetVote(acc *Account, number int) (*VoteInstance, error) {
	result, err := makeRequest(`MATCH (a:Account)-[r:MANAGES]->(v:Vote) 
WHERE a.name = $name AND r.position = $position
RETURN v;`, map[string]any{"name": acc.Name, "position": number})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	props := GetPropsMapForRecordPosition(result[0], 0)
	return &VoteInstance{
		ID:                    props.GetString("id"),
		Question:              props.GetString("question"),
		IterableAnswers:       props.GetArray("answers"),
		Type:                  VoteType(props.GetInt("type")),
		MaxVotes:              props.GetInt("max_votes"),
		ShowVotesDuringVoting: props.GetBool("show_during"),
		Anonymous:             props.GetBool("anonymous"),
	}, nil
}

func GetVoteInfoList(acc *Account) ([]VoteInfo, error) {
	result, err := makeRequest(`MATCH (a:Account)-[:MANAGES]->(v:Vote) 
WHERE a.name = $name RETURN v.id, v.question;`, map[string]any{"name": acc.Name})
	if err != nil {
		return nil, err
	}
	list := make([]VoteInfo, len(result))
	for i, record := range result {
		list[i] = VoteInfo{
			ID:       record.Values[0].(string),
			Question: record.Values[1].(string),
		}
	}
	return list, err
}

func CastVoteWithAccount(name string, id string, votes []int) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}

	result, err := tx.Run(`CALL {
MATCH (a:Account)<-[:PARTICIPANT]-(d:Document)-[:LINKS]->(v:Vote)
WHERE a.name = $author AND a.blocked = false AND v.id = $id AND datetime($now) < datetime(d.end_time) 
RETURN a  
UNION 
MATCH (a:Account)-[:USER|ADMIN]->(:Organisation)<-[:IN]-(d:Document)-[:LINKS]->(v:Vote) 
WHERE a.name = $author AND a.blocked = false AND v.id = $id AND d.member_part = true AND datetime($now) < datetime(d.end_time) 
RETURN a 
UNION 
MATCH (a:Account)-[:ADMIN]->(:Organisation)<-[:IN]-(d:Document)-[:LINKS]->(v:Vote) 
WHERE a.name = $author AND a.blocked = false AND v.id = $id AND d.admin_part = true AND datetime($now) < datetime(d.end_time) 
RETURN a  
} 
RETURN a;`,
		map[string]any{"id": id, "author": name, "now": time.Now().UTC()})
	if err != nil {
		return err
	} else if !result.Next() {
		return notAllowedError
	}

	result, err = tx.Run(`MATCH (a:Account)-[:VOTED]->(v:Vote) WHERE a.name = $author AND v.id = $id 
RETURN a;`,
		map[string]any{"id": id, "author": name})
	if err != nil {
		return err
	} else if result.Next() {
		return AlreadyVoted
	}

	mapIsNil := votes == nil
	_, err = tx.Run(`MATCH (a:Account), (v:Vote) WHERE a.name = $author AND v.id = $id 
MERGE (a)-[:VOTED {written: $now, illegal: $illegal, vote: $voteMap}]->(v);`,
		map[string]any{"id": id, "author": name, "illegal": mapIsNil, "voteMap": votes, "now": time.Now().UTC()})
	if err == nil {
		return tx.Commit()
	}
	return err
}

func GetAnswersAndTypeForVote(id string, acc *Account) ([]any, VoteType, int, error) {
	extra := `CALL { MATCH (o:Organisation)<-[:IN]-(d)-[:LINKS]->(v:Vote) WHERE d.public = true 
RETURN o
UNION
MATCH (a:Account)-[*..]->(o:Organisation)<-[:IN]-(d) WHERE d.public = false AND a.name = $name 
RETURN o }`
	if acc.IsAtLeastAdmin() {
		extra = ""
	}

	result, err := makeRequest(`MATCH (d:Document)-[:LINKS]->(v:Vote) WHERE v.id = $id `+extra+
		` RETURN v.answers, v.type, v.max_votes;`, map[string]any{"id": id, "name": acc.GetName()})
	if err != nil {
		return nil, -1, -1, err
	}
	if len(result) == 0 {
		return nil, -1, -1, notAllowedError
	}
	return result[0].Values[0].([]any), VoteType(result[0].Values[1].(int64)), int(result[0].Values[2].(int64)), nil
}

func GetVoteForUser(id string, acc *Account) (*VoteInstance, *AccountVotes, error) {
	extra := `CALL { MATCH (o:Organisation)<-[:IN]-(d)-[:LINKS]->(v:Vote) WHERE d.public = true AND d.removed = false
RETURN o
UNION
MATCH (a:Account)-[*..]->(o:Organisation)<-[:IN]-(d) WHERE d.public = false AND d.removed = false AND a.name = $name 
RETURN o }`
	if acc.IsAtLeastAdmin() {
		extra = ""
	}

	result, err := makeRequest(`MATCH (d:Document)-[:LINKS]->(v:Vote) WHERE v.id = $id `+extra+
		` RETURN d.id, d.end_time, v;`, map[string]any{"id": id, "name": acc.GetName()})
	if err != nil {
		return nil, nil, err
	}
	if len(result) == 0 {
		return nil, nil, notAllowedError
	}
	props := GetPropsMapForRecordPosition(result[0], 2)
	vote := &VoteInstance{
		DocumentID:            result[0].Values[0].(string),
		ID:                    props.GetString("id"),
		Question:              props.GetString("question"),
		IterableAnswers:       props.GetArray("answers"),
		Type:                  VoteType(props.GetInt("type")),
		MaxVotes:              props.GetInt("max_votes"),
		ShowVotesDuringVoting: props.GetBool("show_during"),
		Anonymous:             props.GetBool("anonymous"),
		EndDate:               result[0].Values[1].(time.Time),
	}
	vote.ShowVotesDuringVoting = vote.ShowVotesDuringVoting || vote.Ended()

	var voteList *AccountVotes
	if vote.ShowVotesDuringVoting {
		voteList = &AccountVotes{
			Question:        vote.Question,
			IterableAnswers: vote.IterableAnswers,
			Type:            vote.Type,
			Anonymous:       vote.Anonymous,
			AnswerAmount:    len(vote.IterableAnswers),
			IllegalVotes:    []string{},
			List:            make(map[string][]any),
		}
		result, err = makeRequest(`MATCH (a:Account)-[r:VOTED]->(v:Vote) WHERE v.id = $id RETURN r, a.name 
ORDER BY r.written;`,
			map[string]any{"id": id})
		if err != nil {
			return nil, nil, err
		}
		for _, record := range result {
			props = GetPropsMapForRecordPosition(record, 0)
			if props.GetBool("illegal") {
				voteList.IllegalVotes = append(voteList.IllegalVotes, record.Values[1].(string))
				continue
			}
			voteList.List[record.Values[1].(string)] = props.GetArray("vote")
		}
	}

	return vote, voteList, err
}

func resultRoutine() {
	generateResults()
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

	result, err := makeRequest(`MATCH (a:Account)-[r:VOTED]->(v:Vote)<-[:LINKS]-(d:Document) 
WHERE datetime($now) > datetime(d.end_time) 
WITH v, d.id AS docID,a.name AS accName, r ORDER BY r.written 
RETURN v, docID, collect(accName), collect(r);`, map[string]any{"now": time.Now().UTC()})
	if err != nil {
		slog.Error(err.Error())
		return
	}
	votes, docIds, voteIds := transformVotesForResults(result)

	tx, err := openTransaction()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer tx.Close()
	for i, vote := range votes {
		err = tx.RunWithoutResult("MATCH (v:Vote) WHERE v.id = $vote DETACH DELETE v;",
			map[string]any{"vote": voteIds[i]})
		if err != nil {
			slog.Error(err.Error())
			return
		}
		err = tx.RunWithoutResult(`MATCH (d:Document) WHERE d.id = $id 
CREATE (r:Result {type: $Type, anonymous: $Anonymous, amount: $AnswerAmount, illegal: $IllegalVotes, 
question: $question, answers: $answers, voter: $voter, votes: $votes}) 
MERGE (d)-[:VOTED]->(r);`,
			map[string]any{
				"id":           docIds[i],
				"Type":         vote.Type,
				"Anonymous":    vote.Anonymous,
				"AnswerAmount": vote.AnswerAmount,
				"IllegalVotes": vote.IllegalVotes,
				"question":     vote.Question,
				"answers":      vote.IterableAnswers,
				"voter":        vote.Voter,
				"votes":        vote.Votes,
			})
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.Debug(err.Error())
	}
}

// returns first the docID then the VoteID array
func transformVotesForResults(result []*neo4j.Record) ([]AccountVotes, []string, []string) {
	votes := make([]AccountVotes, len(result))
	docIDs := make([]string, len(result))
	voteIDs := make([]string, len(result))

	for i, record := range result {
		docIDs[i] = record.Values[1].(string)
		voteProps := GetPropsMapForRecordPosition(record, 0)
		voteIDs[i] = voteProps.GetString("id")
		slog.Debug("Working on: ", "Document", docIDs[i], "Vote", voteIDs[i])

		votes[i] = AccountVotes{
			Question:        voteProps.GetString("question"),
			IterableAnswers: voteProps.GetArray("answers"),
			Anonymous:       voteProps.GetBool("anonymous"),
			Type:            VoteType(voteProps.GetInt("type")),
			IllegalVotes:    []string{},
			Voter:           make([]any, 0),
			Votes:           make([]any, 0),
		}
		votes[i].AnswerAmount = len(votes[i].IterableAnswers)

		nameList := record.Values[2].([]any)
		voteList := record.Values[3].([]any)

		for j, name := range nameList {
			props := PropsMap(voteList[j].(neo4j.Relationship).Props)
			if props.GetBool("illegal") {
				votes[i].IllegalVotes = append(votes[i].IllegalVotes, name.(string))
				continue
			}
			votes[i].Voter = append(votes[i].Voter, name)
			votes[i].Votes = append(votes[i].Votes, props.GetArray("vote")...)
		}
	}

	return votes, docIDs, voteIDs
}
