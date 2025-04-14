package database

import (
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
	"html/template"
	"strconv"
	"strings"
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

func (b *BlackboardNote) GetEmbed(references []string) *discordgo.MessageEmbed {
	base := &discordgo.MessageEmbed{
		URL:       helper.UrlPrefix + "/notes?loaded=" + b.ID,
		Type:      discordgo.EmbedTypeRich,
		Title:     b.Title,
		Timestamp: b.PostedAt.Format("2006-01-02T15:04:05Z"),
		Color:     0x334155,
		Footer:    nil,
		Image:     nil,
		Thumbnail: nil,
		Video:     nil,
		Provider:  nil,
		Author:    &discordgo.MessageEmbedAuthor{Name: b.GetAuthor()},
		Fields:    make([]*discordgo.MessageEmbedField, 0, 1),
	}
	if len(references) > 1 {
		base.Fields = append(base.Fields, &discordgo.MessageEmbedField{
			Name: loc.NotesInfoField,
			Value: fmt.Sprintf(loc.NotesMultipleReferences, len(references),
				helper.UrlPrefix+"/notes?loaded=", strings.Join(append([]string{b.ID}, references...), "&loaded=")),
		})
	} else if len(references) == 1 {
		base.Fields = append(base.Fields, &discordgo.MessageEmbedField{
			Name: loc.NotesInfoField,
			Value: fmt.Sprintf(loc.NotesSingleReference, helper.UrlPrefix+"/notes?loaded=",
				b.ID+"&loaded="+references[0]),
		})
	}
	return base
}

func CreateNote(note *BlackboardNote, references []string) error {
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}
	defer rollback(tx)
	note.PostedAt = time.Now().UTC()
	_, err = tx.Exec(`INSERT INTO blackboard_note (id, title, author, flair, posted, body, blocked) 
VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		&note.ID, &note.Title, &note.Author, &note.Flair, &note.PostedAt, &note.Body, &note.Removed)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO blackboard_references (base_note_id, reference_id)  
SELECT $1 AS base_note_id, id FROM blackboard_note
WHERE id = ANY($2);`, &note.ID, pq.Array(references))
	if err != nil {
		return err
	}
	err = tx.QueryRow(`SELECT ARRAY(SELECT reference_id FROM blackboard_references WHERE base_note_id = $1) 
FROM blackboard_references;`, note.ID).Scan(pq.Array(&references))
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err == nil {
		helper.SendDiscordEmbedMessage(helper.DiscordNotesChannelID, note.GetEmbed(references))
	}
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

type NoteSearch struct {
	Title            string
	ExactTitleMatch  bool
	Author           string
	ExactAuthorMatch bool
	ShowBlocked      bool
	values           []any
}

func (n *NoteSearch) GetQuery(acc *Account) string {
	var query string

	n.values = make([]any, 0)
	pos := 3
	if n.ShowBlocked && acc.IsAtLeastAdmin() {
		query += " true"
	} else {
		query += " blocked = false"
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

func (n *NoteSearch) GetValues(input []any) []any {
	return append(input, n.values...)
}

func SearchForNotesForwards(acc *Account, amount int, timeStamp time.Time, info *NoteSearch) ([]TruncatedBlackboardNotes, error) {
	result, err := postgresDB.Query(`SELECT id, title, author, flair, posted, blocked FROM blackboard_note
WHERE `+info.GetQuery(acc)+` AND posted <= $1 ORDER BY posted DESC LIMIT $2;`, info.GetValues([]any{timeStamp, amount + 1})...)
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

func SearchForNotesBackwards(acc *Account, amount int, timeStamp time.Time, info *NoteSearch) ([]TruncatedBlackboardNotes, error) {
	result, err := postgresDB.Query(`SELECT id, title, author, flair, posted, blocked FROM 
(SELECT id, title, author, flair, posted, blocked FROM blackboard_note
WHERE `+info.GetQuery(acc)+` AND posted >= $1 ORDER BY posted LIMIT $2) as note ORDER BY note.posted DESC;`, info.GetValues([]any{timeStamp, amount + 2})...)
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
