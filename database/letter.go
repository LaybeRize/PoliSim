package database

import (
	loc "PoliSim/localisation"
	"fmt"
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
	ID        string
	Title     string
	Author    string
	Flair     string
	Written   time.Time
	Recipient string
	Viewed    bool
}

func (r *ReducedLetter) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return r.Written.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return r.Written.Format(loc.TimeFormatString)
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
	return fmt.Sprintf(loc.LetterRecipientsFormatString, strings.Join(l.Reader, ", "))
}

func (l *Letter) GetAgreed() string {
	return fmt.Sprintf(loc.LetterAcceptedFormatString, strings.Join(l.Agreed, ", "))
}

func (l *Letter) GetDeclined() string {
	if len(l.Declined) == 0 {
		return loc.LetterNoOneDeclined
	}
	return fmt.Sprintf(loc.LetterDeclinedFormatString, strings.Join(l.Declined, ", "))
}

func (l *Letter) SomeoneHasNotDecidedYet() bool {
	return len(l.NoDecision) != 0
}

func (l *Letter) GetNoDecision() string {
	return fmt.Sprintf(loc.LetterNoDecisionFormatString, strings.Join(l.NoDecision, ", "))
}

func (l *Letter) GetAuthor() string {
	if l.Flair == "" {
		return l.Author
	}
	return l.Author + "; " + l.Flair
}

func (l *Letter) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return l.Written.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return l.Written.Format(loc.TimeFormatString)
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
	defer tx.Close()
	if err != nil {
		return err
	}

	result, err := tx.Run(`MATCH (acc:Account) WHERE acc.name = $Author AND acc.blocked = false 
RETURN acc;`,
		map[string]any{"Author": letter.Author})
	if err != nil {
		return err
	} else if !result.Peek() {
		return notAllowedError
	}

	result, err = tx.Run(`MATCH (a:Account) WHERE a.name IN $reader AND a.blocked = false 
RETURN a;`,
		map[string]any{"reader": letter.Reader})
	if err != nil {
		return err
	} else if !result.Peek() {
		return noRecipientFoundError
	}

	err = tx.RunWithoutResult(letterCreation, letter.GetCreationMap())
	if err != nil {
		return err
	}

	err = tx.RunWithoutResult(letterLinkage, letter.GetCreationMap())
	if err != nil {
		return err
	}

	err = tx.Commit()
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
	list := make([]ReducedLetter, len(result))
	for i, record := range result {
		list[i] = ReducedLetter{
			ID:        record.Values[0].(string),
			Title:     record.Values[1].(string),
			Author:    record.Values[2].(string),
			Flair:     record.Values[3].(string),
			Written:   record.Values[4].(time.Time),
			Recipient: record.Values[5].(string),
			Viewed:    record.Values[6].(bool),
		}
	}
	return list, nil
}

func GetLetterForReader(id string, reader string) (*Letter, error) {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return nil, err
	}
	var result *dbResult

	if reader == loc.AdministrationName {
		result, err = tx.Run(`MATCH (l:Letter)
WHERE l.id = $id
RETURN l;`,
			map[string]any{"id": id})
		if err != nil {
			return nil, err
		} else if !result.Next() {
			return nil, notFoundError
		}
	} else {
		result, err = tx.Run(`MATCH (a:Account)<-[r:RECIPIENT]-(l:Letter)
WHERE a.name = $reader AND l.id = $id 
SET r.viewed = true 
RETURN l, r.signature;`,
			map[string]any{"id": id, "reader": reader})
		if err != nil {
			return nil, err
		} else if !result.Next() {
			return nil, notFoundError
		}
	}

	props := GetPropsMapForRecordPosition(result.Record(), 0)
	letter := &Letter{Recipient: reader}
	if reader == loc.AdministrationName {
		letter.HasSigned = true
	} else {
		letter.HasSigned = result.Record().Values[1].(int64) != int64(NoDecision)
	}
	letter.ID = id
	letter.Title = props.GetString("title")
	letter.Author = props.GetString("author")
	letter.Flair = props.GetString("flair")
	letter.Written = props.GetTime("written")
	letter.Signable = props.GetBool("signable")
	letter.Body = template.HTML(props.GetString("body"))
	letter.Reader = make([]string, 0)
	if letter.Signable {
		letter.Agreed = make([]string, 0)
		letter.Declined = make([]string, 0)
		letter.NoDecision = make([]string, 0)
	}

	result, err = tx.Run(`MATCH (a:Account)<-[r:RECIPIENT]-(l:Letter) 
WHERE l.id = $id 
RETURN a.name, r.signature ORDER BY a.name;`,
		map[string]any{"id": id})
	if err != nil {
		return letter, err
	}
	for result.Next() {
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

	err = tx.Commit()
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
