package database

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"html/template"
	"time"
)

const (
	Agreed         = "agreed"
	Declined       = "declined"
	NoDecision     = "no_decision"
	NoSignPossible = "no_sign"

	letterCreation = `MATCH (a:Account) WHERE a.name IN $reader 
MATCH (aut:Account) WHERE a.name = $Author 
CREATE (l:Letter {id: $id, title: $title , author: $Author , flair: $Flair, 
written: $written , signable: $signable , body: $Body}) 
MERGE (a)-[:RECIPIENT {signature: $signature, viewed: false}]->(l) 
MERGE (aut)-[:RECIPIENT {signature: $authorSign, viewed: true}]->(l) 
MERGE (aut)-[:WRITTEN]->(l);`
)

type ReducedLetter struct {
	ID       string
	Title    string
	Author   string
	Written  time.Time
	Recipent string
	Viewed   bool
}

func (r *ReducedLetter) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return r.Written.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return r.Written.Format("2006-01-02 15:04:05 MST")
}

type Letter struct {
	ID         string
	Title      string
	Author     string
	Flair      string
	Signable   bool
	Written    time.Time
	Body       template.HTML
	Recipent   string
	Reader     []string
	Agreed     []string
	Declined   []string
	NoDecision []string
}

func (l *Letter) GetCreationMap() map[string]any {
	signature := NoDecision
	authorSign := Agreed
	if !l.Signable {
		signature = NoSignPossible
		authorSign = NoSignPossible
	}
	return map[string]any{
		"id":         l.ID,
		"title":      l.Title,
		"Author":     l.Author,
		"Flair":      l.Flair,
		"written":    time.Now().UTC(),
		"signable":   l.Signable,
		"Body":       l.Body,
		"reader":     l.Reader,
		"signature":  signature,
		"authorSign": authorSign}
}

func (l *Letter) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return l.Written.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return l.Written.Format("2006-01-02 15:04:05 MST")
}

func CreateLetter(letter *Letter) error {
	_, err := makeRequest(letterCreation, letter.GetCreationMap())
	return err
}

func GetLetter(id string, viewer string) (*Letter, error) {
	result, err := makeRequest(`MATCH (a:Account)-[:RECIPIENT]->(l:Letter) 
WHERE l.id = $id AND a.name = $viewer
RETURN l;`, map[string]any{"id": id, "viewer": viewer})
	if err != nil {
		return nil, err
	}
	if len(result.Records) != 1 {
		return nil, notFoundError
	}

	node := result.Records[0].Values[0].(neo4j.Node).Props
	letter := &Letter{
		ID:         node["id"].(string),
		Title:      node["title"].(string),
		Author:     node["author"].(string),
		Flair:      node["flair"].(string),
		Signable:   node["signable"].(bool),
		Written:    node["written"].(time.Time),
		Body:       template.HTML(node["body"].(string)),
		Recipent:   viewer,
		Reader:     []string{},
		Agreed:     []string{},
		Declined:   []string{},
		NoDecision: []string{},
	}

	result, err = makeRequest(`MATCH (a:Account)-[r:RECIPIENT]->(l:Letter) 
WHERE l.id = $id
RETURN a.name, r.signature;`, map[string]any{"id": id})
	if err != nil {
		return nil, err
	}

	for _, record := range result.Records {
		name := record.Values[0].(string)
		letter.Reader = append(letter.Reader, name)
		switch record.Values[1].(string) {
		case Agreed:
			letter.Agreed = append(letter.Agreed, name)
		case Declined:
			letter.Declined = append(letter.Declined, name)
		case NoDecision:
			letter.NoDecision = append(letter.NoDecision, name)
		}
	}

	return letter, err
}

func GetLetterList(viewer []string) ([]ReducedLetter, error) {
	return nil, nil
}
