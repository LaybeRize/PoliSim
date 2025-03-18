package database

import (
	loc "PoliSim/localisation"
	"github.com/lib/pq"
	"html/template"
	"regexp"
	"time"
)

type TruncatedBlackboardNotes struct {
	ID       string
	Title    string
	Author   string
	Flair    string
	Removed  bool
	PostedAt time.Time
}

func (t *TruncatedBlackboardNotes) GetTimePostedAt(a *Account) string {
	if a.Exists() {
		return t.PostedAt.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return t.PostedAt.Format(loc.TimeFormatString)
}

func (t *TruncatedBlackboardNotes) IDinArray(arr []string) bool {
	for _, e := range arr {
		if e == t.ID {
			return true
		}
	}
	return false
}

func (t *TruncatedBlackboardNotes) GetAuthor() string {
	if t.Flair == "" {
		return t.Author
	}
	return t.Author + "; " + t.Flair
}

type BlackboardNote struct {
	ID       string
	Title    string
	Author   string
	Flair    string
	PostedAt time.Time
	Body     template.HTML
	Removed  bool
	Parents  []TruncatedBlackboardNotes
	Children []TruncatedBlackboardNotes
	Viewer   *Account
}

func (b *BlackboardNote) GetBody(acc *Account) template.HTML {
	if acc.IsAtLeastAdmin() || !b.Removed {
		return b.Body
	}
	return loc.NotesContentRemovedHTML
}

func (b *BlackboardNote) GetTitle(acc *Account) string {
	if acc.IsAtLeastAdmin() || !b.Removed {
		return b.Title
	}
	return loc.NotesRemovedTitelText
}

func (b *BlackboardNote) HasChildren() bool {
	return len(b.Children) != 0
}

func (b *BlackboardNote) HasParents() bool {
	return len(b.Parents) != 0
}

func (b *BlackboardNote) GetAuthor() string {
	if b.Flair == "" {
		return b.Author
	}
	return b.Author + "; " + b.Flair
}

func (b *BlackboardNote) GetTimePostedAt(a *Account) string {
	if a.Exists() {
		return b.PostedAt.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return b.PostedAt.Format(loc.TimeFormatString)
}

// Todo move this to migration
const tableDefinitionNotes = `
CREATE TABLE blackboard_note(
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
    author  TEXT NOT NULL,
    flair  TEXT NOT NULL,
    posted TIMESTAMP NOT NULL,
    body  TEXT NOT NULL,
	blocked BOOLEAN NOT NULL
);
CREATE TABLE blackboard_references(
	base_note_id TEXT NOT NULL,
	reference_id TEXT NOT NULL,
	CONSTRAINT fk_base_note_id
        FOREIGN KEY(base_note_id) REFERENCES blackboard_note(id),
    CONSTRAINT fk_reference_id
        FOREIGN KEY(reference_id) REFERENCES blackboard_note(id)
);
`

func CreateNote(note *BlackboardNote, references []string) error {
	_, err := postgresDB.Exec(`INSERT INTO blackboard_note (id, title, author, flair, posted, body, blocked) 
VALUES ($1, $2, $3, $4, $5, $6, $7);
INSERT INTO blackboard_references (base_note_id, reference_id)  
SELECT $1 AS base_note_id, id FROM blackboard_note
WHERE id = ANY($8);`,
		&note.ID, &note.Title, &note.Author, &note.Flair, time.Now().UTC(), &note.Body, &note.Removed,
		pq.Array(references))
	return err
}

func UpdateNoteRemovedStatus(note *BlackboardNote) error {
	_, err := postgresDB.Exec(`UPDATE blackboard_note SET blocked = $2 WHERE id = $1`, &note.ID, &note.Removed)
	return err
}

func GetNote(id string) (*BlackboardNote, error) {
	note := &BlackboardNote{}
	err := postgresDB.QueryRow(`SELECT id, title, author, flair, posted, body, blocked FROM blackboard_note
WHERE id = $1;`, &id).Scan(&note.ID, &note.Title, &note.Author, &note.Flair, &note.PostedAt, &note.Body, &note.Removed)
	if err != nil {
		return nil, err
	}
	result, err := postgresDB.Query(`SELECT id, title, author, flair, posted, blocked FROM blackboard_note 
    INNER JOIN blackboard_references br ON blackboard_note.id = br.reference_id WHERE base_note_id = $1`, &id)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	note.Parents = make([]TruncatedBlackboardNotes, 0)
	trunc := TruncatedBlackboardNotes{}
	for result.Next() {
		err = result.Scan(&trunc.ID, &trunc.Title, &trunc.Author, &trunc.Flair, &trunc.PostedAt, &trunc.Removed)
		if err != nil {
			return nil, err
		}
		note.Parents = append(note.Parents, trunc)
	}
	result, err = postgresDB.Query(`SELECT id, title, author, flair, posted, blocked FROM blackboard_note 
    INNER JOIN blackboard_references br ON blackboard_note.id = br.base_note_id WHERE reference_id = $1`, &id)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	note.Children = make([]TruncatedBlackboardNotes, 0)
	for result.Next() {
		err = result.Scan(&trunc.ID, &trunc.Title, &trunc.Author, &trunc.Flair, &trunc.PostedAt, &trunc.Removed)
		if err != nil {
			return nil, err
		}
		note.Children = append(note.Children, trunc)
	}
	return note, err
}

func SearchForNotes(acc *Account, amount int, page int, input string, showBlocked bool) ([]TruncatedBlackboardNotes, error) {
	parameter := []any{(page - 1) * amount, amount}
	query, parameter := queryAnalyzer(acc, parameter, input, showBlocked)
	result, err := postgresDB.Query(`SELECT id, title, author, flair, posted, blocked FROM blackboard_note
WHERE `+query+` ORDER BY posted DESC OFFSET $1 LIMIT $2;`, parameter...)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]TruncatedBlackboardNotes, 0)
	trunc := TruncatedBlackboardNotes{}
	for result.Next() {
		err = result.Scan(&trunc.ID, &trunc.Title, &trunc.Author, &trunc.Flair, &trunc.PostedAt, &trunc.Removed)
		if err != nil {
			return nil, err
		}
		arr = append(arr, trunc)
	}
	return arr, nil
}

var queryRegexNotes = regexp.MustCompile(`^\s*(.*?)\s*(\[|$)`)
var authorRegexNotes = regexp.MustCompile(`\[[bB][yY]:]\s*(.+?)\s*(\[|$)`)

func queryAnalyzer(acc *Account, parameter []any, input string, showBlocked bool) (string, []any) {
	query := ""
	if showBlocked && acc.IsAtLeastAdmin() {
		query += "true"
	} else {
		query += "blocked = false"
	}

	result := queryRegexNotes.FindStringSubmatch(input)
	if result != nil && result[1] != "" {
		parameter = append(parameter, result[1])
		query += " AND title LIKE '%' || $3 || '%'"
	}

	if result = authorRegexNotes.FindStringSubmatch(input); result != nil {
		parameter = append(parameter, result[1])
		query += " AND author LIKE '%' || $4 || '%'"
	}

	return query, parameter
}
