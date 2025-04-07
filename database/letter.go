package database

import (
	loc "PoliSim/localisation"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"html/template"
	"strconv"
	"strings"
	"time"
)

type LetterStatus int

const (
	Agreed LetterStatus = iota
	Declined
	NoDecision
	NoSignPossible
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

func CreateLetter(letter *Letter) error {
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}
	name := ""
	err = tx.QueryRow(`SELECT name FROM account WHERE name = $1 AND blocked = false`, &letter.Author).Scan(&name)
	if err != nil {
		_ = tx.Rollback()
		return NotAllowedError
	}
	err = tx.QueryRow(`SELECT name FROM account WHERE name = ANY($1) AND blocked = false LIMIT 1`,
		pq.Array(letter.Reader)).Scan(&name)
	if err != nil {
		_ = tx.Rollback()
		return NoRecipientFoundError
	}

	err = createLetter(tx, letter)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func createLetter(tx *sql.Tx, letter *Letter) error {
	_, err := tx.Exec(`INSERT INTO letter (id, title, author, flair, signable, written, body) VALUES 
($1, $2, $3, $4, $5, $6, $7);`, &letter.ID, &letter.Title, &letter.Author, &letter.Flair, &letter.Signable, time.Now().UTC(), &letter.Body)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	signature := NoDecision
	authorSign := Agreed
	if !letter.Signable {
		signature = NoSignPossible
		authorSign = NoSignPossible
	}

	if letter.Author != loc.AdministrationName {
		_, err = tx.Exec(`INSERT INTO letter_to_account (letter_id, account_name, has_read, sign_status) 
VALUES ($1, $2, true, $3);`, &letter.ID, &letter.Author, &authorSign)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	result, err := tx.Query(`
INSERT INTO letter_to_account (letter_id, account_name, has_read, sign_status)
SELECT $1 AS letter_id, name, false AS has_read, $3 AS sign_status FROM account
WHERE name = ANY($2) AND blocked = false RETURNING account_name;`, &letter.ID, pq.Array(letter.Reader), &signature)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	val := newEntry{ID: letter.ID, Accounts: make([]string, 0)}
	for result.Next() {
		var name string
		var tempErr error
		tempErr = result.Scan(&name)
		if tempErr != nil {
			continue
		}
		val.Accounts = append(val.Accounts, name)
	}
	newLetter <- val

	return err
}

type LetterSearch struct {
	Title            string
	ExactTitleMatch  bool
	Author           string
	ExactAuthorMatch bool
	ShowOnlyUnread   bool
	values           []any
}

func (n *LetterSearch) GetQuery() string {
	var query string

	n.values = make([]any, 0)
	pos := 5
	if n.ShowOnlyUnread {
		query += " AND has_read = false "
	}

	if n.Title != "" {
		if n.ExactTitleMatch {
			query += " AND title = $" + strconv.Itoa(pos) + " "
		} else {
			query += " AND title LIKE '%' || $" + strconv.Itoa(pos) + " || '%' "
		}
		pos += 1
		n.values = append(n.values, n.Title)
	}

	if n.Author != "" {
		if n.ExactAuthorMatch {
			query += " AND author = $" + strconv.Itoa(pos) + " "
		} else {
			query += " AND author LIKE '%' || $" + strconv.Itoa(pos) + " || '%' "
		}
		pos += 1
		n.values = append(n.values, n.Author)
	}

	return query
}

func (n *LetterSearch) GetValues(input []any) []any {
	return append(input, n.values...)
}

func GetLetterListForwards(viewer []string, amount int, timeStamp time.Time, recName string, info *LetterSearch) ([]ReducedLetter, error) {
	result, err := postgresDB.Query(`SELECT id, title, author, flair, written, account_name, has_read FROM letter
 INNER JOIN letter_to_account lta on letter.id = lta.letter_id WHERE account_name = ANY($1) `+info.GetQuery()+` AND (written, account_name) <= ($2, $3)
 ORDER BY (written, account_name) DESC LIMIT $4;`, info.GetValues([]any{pq.Array(viewer), timeStamp, recName, amount + 1})...)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	list := make([]ReducedLetter, 0)
	item := ReducedLetter{}
	for result.Next() {
		err = result.Scan(&item.ID, &item.Title, &item.Author, &item.Flair, &item.Written, &item.Recipient, &item.Viewed)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, err
}

func GetLetterListBackwards(viewer []string, amount int, timeStamp time.Time, recName string, info *LetterSearch) ([]ReducedLetter, error) {
	result, err := postgresDB.Query(`SELECT id, title, author, flair, written, account_name, has_read FROM (
SELECT id, title, author, flair, written, account_name, has_read FROM letter
 INNER JOIN letter_to_account lta on letter.id = lta.letter_id WHERE account_name = ANY($1) `+info.GetQuery()+` AND (written, account_name) >= ($2, $3)
 ORDER BY (written, account_name) LIMIT $4) as let ORDER BY (let.written, let.account_name) DESC;`,
		info.GetValues([]any{pq.Array(viewer), timeStamp, recName, amount + 2})...)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	list := make([]ReducedLetter, 0)
	item := ReducedLetter{}
	for result.Next() {
		err = result.Scan(&item.ID, &item.Title, &item.Author, &item.Flair, &item.Written, &item.Recipient, &item.Viewed)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, err
}

func GetLetterForReader(id string, reader string) (*Letter, error) {
	var row *sql.Row

	if reader == loc.AdministrationName {
		row = postgresDB.QueryRow(`SELECT title, author, flair, written, signable, body FROM letter 
                                                     WHERE id = $1`, &id)
	} else {
		row = postgresDB.QueryRow(`UPDATE letter_to_account SET has_read = true FROM letter 
WHERE letter_to_account.letter_id = letter.id AND account_name = $1 AND letter_id = $2
RETURNING letter.title, letter.author, letter.flair, letter.written, letter.signable, letter.body;`, &reader, &id)
	}
	letter := &Letter{ID: id, Recipient: reader}
	err := row.Scan(&letter.Title, &letter.Author, &letter.Flair, &letter.Written, &letter.Signable, &letter.Body)
	if err != nil {
		return nil, err
	}
	markLetterAsRead <- []string{reader, id}
	result, err := postgresDB.Query(`SELECT account_name, sign_status FROM letter_to_account WHERE letter_id = $1`, &id)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	letter.Reader = make([]string, 0)
	if letter.Signable {
		letter.Agreed = make([]string, 0)
		letter.Declined = make([]string, 0)
		letter.NoDecision = make([]string, 0)
	}
	letter.HasSigned = true
	accountName := ""
	status := LetterStatus(-10)
	for result.Next() {
		err = result.Scan(&accountName, &status)
		if err != nil {
			return nil, err
		}
		letter.Reader = append(letter.Reader, accountName)
		switch status {
		case NoDecision:
			if accountName == reader {
				letter.HasSigned = false
			}
			letter.NoDecision = append(letter.NoDecision, accountName)
		case Agreed:
			letter.Agreed = append(letter.Agreed, accountName)
		case Declined:
			letter.Declined = append(letter.Declined, accountName)
		default:
		}
	}
	return letter, err
}

func UpdateSignatureStatus(id string, reader string, agree bool) error {
	newStatus := Declined
	if agree {
		newStatus = Agreed
	}
	_, err := postgresDB.Exec(`UPDATE letter_to_account SET sign_status = $1 
                         WHERE sign_status = $2 AND letter_id = $3 AND account_name = $4`,
		newStatus, NoDecision, id, reader)
	return err
}
