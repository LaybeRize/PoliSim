package database

import (
	loc "PoliSim/localisation"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"html/template"
	"strings"
	"time"
)

type LetterStatus int

const (
	Agreed LetterStatus = iota
	Declined
	NoDecision
	NoSignPossible

	letterCreation = `MATCH (aut:Account) WHERE aut.name = $Author
CREATE (l:Letter {id: $id, title: $title , author: $Author , flair: $Flair, 
written: $written , signable: $signable , body: $Body}) 
CREATE (aut)<-[:RECIPIENT {signature: $authorSign, viewed: true}]-(l) 
CREATE (aut)-[:WRITTEN]->(l);`
	letterLinkage = `MATCH (l:Letter), (a:Account) 
WHERE l.id = $id AND a.name IN $reader AND a.blocked = false AND a.name <> $Author
CREATE (a)<-[:RECIPIENT {signature: $signature, viewed: false}]-(l);`
)

type ReducedLetter struct {
	ID       string
	Title    string
	Author   string
	Flair    string
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

func (r *ReducedLetter) GetAuthor() string {
	if r.Flair == "" {
		return r.Author
	}
	return r.Author + "; " + r.Flair
}

type Letter struct {
	ID         string
	Title      string
	Author     string
	Flair      string
	Signable   bool
	Written    time.Time
	Body       template.HTML
	Recipient  string
	HasSigned  bool
	Reader     []string
	Agreed     []string
	Declined   []string
	NoDecision []string
}

func (l *Letter) GetReader() string {
	return fmt.Sprintf("Empfänger: %s", strings.Join(l.Reader, ", "))
}

func (l *Letter) GetAgreed() string {
	return fmt.Sprintf("Zugestimmt: %s", strings.Join(l.Agreed, ", "))
}

func (l *Letter) GetDeclined() string {
	if len(l.Declined) == 0 {
		return "Niemand hat abgelehnt"
	}
	return fmt.Sprintf("Abgelehnt: %s", strings.Join(l.Declined, ", "))
}

func (l *Letter) SomeoneHasntDecidedYet() bool {
	return len(l.NoDecision) != 0
}

func (l *Letter) GetNoDecision() string {
	return fmt.Sprintf("Keine Entscheidung: %s", strings.Join(l.NoDecision, ", "))
}

func (l *Letter) GetAuthor() string {
	if l.Flair == "" {
		return l.Author
	}
	return l.Author + "; " + l.Flair
}

func (l *Letter) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return l.Written.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return l.Written.Format("2006-01-02 15:04:05 MST")
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

func CreateLetter(letter *Letter) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `MATCH (acc:Account) WHERE acc.name = $Author AND acc.blocked = false 
RETURN acc;`,
		map[string]any{"Author": letter.Author})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	} else if result.Next(ctx); result.Record() == nil {
		return notAllowedError
	}

	result, err = tx.Run(ctx, `MATCH (a:Account) WHERE a.name IN $reader AND a.blocked = false 
RETURN a;`,
		map[string]any{"reader": letter.Reader})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	} else if result.Next(ctx); result.Record() == nil {
		return noRecipientFoundError
	}

	_, err = tx.Run(ctx, letterCreation, letter.GetCreationMap())
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Run(ctx, letterLinkage, letter.GetCreationMap())
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func GetLetterList(viewer []string, amount int, page int) ([]ReducedLetter, error) {
	result, err := makeRequest(`MATCH (a:Account)<-[r:RECIPIENT]-(l:Letter) 
WHERE a.name IN $viewer 
RETURN l.id, l.title, l.author, l.flair, l.written, a.name, r.viewed 
ORDER BY l.written DESC, a.name SKIP $skip LIMIT $amount;`,
		map[string]any{"viewer": viewer,
			"skip":   (page - 1) * amount,
			"amount": amount,
		})
	if err != nil {
		return nil, err
	}
	list := make([]ReducedLetter, len(result.Records))
	for i, record := range result.Records {
		list[i] = ReducedLetter{
			ID:       record.Values[0].(string),
			Title:    record.Values[1].(string),
			Author:   record.Values[2].(string),
			Flair:    record.Values[3].(string),
			Written:  record.Values[4].(time.Time),
			Recipent: record.Values[5].(string),
			Viewed:   record.Values[6].(bool),
		}
	}
	return list, nil
}

func GetLetterForReader(id string, reader string) (*Letter, error) {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return nil, err
	}
	var result neo4j.ResultWithContext

	if reader == loc.AdministrationName {
		result, err = tx.Run(ctx, `MATCH (l:Letter)
WHERE l.id = $id
RETURN l;`,
			map[string]any{"id": id})
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, err
		} else if result.Next(ctx); result.Record() == nil {
			_ = tx.Rollback(ctx)
			return nil, notFoundError
		}
	} else {
		result, err = tx.Run(ctx, `MATCH (a:Account)<-[r:RECIPIENT]-(l:Letter)
WHERE a.name = $reader AND l.id = $id 
SET r.viewed = true 
RETURN l, r.signature;`,
			map[string]any{"id": id, "reader": reader})
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, err
		} else if result.Next(ctx); result.Record() == nil {
			_ = tx.Rollback(ctx)
			return nil, notFoundError
		}
	}

	nodeTitle := result.Record().Values[0].(neo4j.Node)
	letter := &Letter{Recipient: reader}
	if reader == loc.AdministrationName {
		letter.HasSigned = true
	} else {
		letter.HasSigned = result.Record().Values[1].(int64) != int64(NoDecision)
	}
	letter.ID = id
	letter.Title = nodeTitle.Props["title"].(string)
	letter.Author = nodeTitle.Props["author"].(string)
	letter.Flair = nodeTitle.Props["flair"].(string)
	letter.Written = nodeTitle.Props["written"].(time.Time)
	letter.Signable = nodeTitle.Props["signable"].(bool)
	letter.Body = template.HTML(nodeTitle.Props["body"].(string))
	letter.Reader = make([]string, 0)
	if letter.Signable {
		letter.Agreed = make([]string, 0)
		letter.Declined = make([]string, 0)
		letter.NoDecision = make([]string, 0)
	}

	result, err = tx.Run(ctx, `MATCH (a:Account)<-[r:RECIPIENT]-(l:Letter) 
WHERE l.id = $id 
RETURN a.name, r.signature ORDER BY a.name;`,
		map[string]any{"id": id})
	if err != nil {
		_ = tx.Rollback(ctx)
		return letter, err
	}
	for result.Next(ctx) {
		name := result.Record().Values[0].(string)
		letter.Reader = append(letter.Reader, name)
		switch LetterStatus(result.Record().Values[1].(int64)) {
		case NoDecision:
			letter.NoDecision = append(letter.NoDecision, name)
		case Agreed:
			letter.Agreed = append(letter.Agreed, name)
		case Declined:
			letter.Declined = append(letter.Declined, name)
		default:
		}
	}

	err = tx.Commit(ctx)
	return letter, err
}

func UpdateSingatureStatus(id string, reader string, agree bool) error {
	newStatus := Declined
	if agree {
		newStatus = Agreed
	}
	_, err := makeRequest(`MATCH (a:Account)<-[r:RECIPIENT]-(l:Letter) 
WHERE l.id = $id AND a.name = $reader AND r.signature = $oldStatus 
SET r.signature = $status;`,
		map[string]any{"id": id, "reader": reader,
			"oldStatus": NoDecision,
			"status":    newStatus})
	return err
}
