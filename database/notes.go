package database

import (
	"html/template"
	"strings"
	"time"
)

type TruncatedBlackboardNotes struct {
	ID       string
	Title    string
	Author   string
	Flair    string
	PostedAt time.Time
}

func (t *TruncatedBlackboardNotes) GetTimePostedAt(a *Account) string {
	if a.Exists() {
		return t.PostedAt.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return t.PostedAt.Format("2006-01-02 15:04:05 MST")
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
		return b.PostedAt.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return b.PostedAt.Format("2006-01-02 15:04:05 MST")
}

func CreateNote(note *BlackboardNote, references []string) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(
		`MATCH (a:Account) WHERE a.name = $Author
CREATE (n:Note {id: $id , title: $title , author: $Author , flair: $Flair, 
posted_at: $PostedAt , body: $Body, removed: $Removed}) 
MERGE (a)-[:WRITTEN]->(n);`,
		map[string]any{"id": note.ID,
			"title":    note.Title,
			"Author":   note.Author,
			"Flair":    note.Flair,
			"PostedAt": note.PostedAt,
			"Body":     string(note.Body),
			"Removed":  note.Removed})
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(
		`MATCH (c:Note), (p:Note) WHERE c.id = $child AND p.id IN $parent MERGE (c)-[:LINKS]->(p);`,
		map[string]any{"parent": references,
			"child": note.ID})
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func UpdateNoteRemovedStatus(note *BlackboardNote) error {
	_, err := makeRequest(`MATCH (n:Note) WHERE n.id = $id 
SET n.removed = $removed
RETURN n;`, map[string]any{"id": note.ID, "removed": note.Removed})
	return err
}

func GetNote(id string) (*BlackboardNote, error) {
	idMap := map[string]any{"id": id}
	result, err := makeRequest(`MATCH (n:Note) WHERE n.id = $id RETURN n;`,
		idMap)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, notFoundError
	} else if len(result) > 1 {
		return nil, multipleItemsError
	}
	note := &BlackboardNote{}
	props := GetPropsMapForRecordPosition(result[0], 0)
	note.ID = props.GetString("id")
	note.Title = props.GetString("title")
	note.Author = props.GetString("author")
	note.Flair = props.GetString("flair")
	note.PostedAt = props.GetTime("posted_at")
	note.Body = template.HTML(props.GetString("body"))
	note.Removed = props.GetBool("removed")
	note.Parents, err = queryForRelations(`MATCH (n:Note) WHERE n.id = $id MATCH (r:Note)-[:LINKS]->(n) RETURN r;`, idMap)
	if err != nil {
		return nil, err
	}
	note.Children, err = queryForRelations(`MATCH (n:Note) WHERE n.id = $id MATCH (n)-[:LINKS]->(r:Note) RETURN r;`, idMap)
	return note, err
}

func queryForRelations(query string, idMap map[string]any) ([]TruncatedBlackboardNotes, error) {
	result, err := makeRequest(query, idMap)
	if err != nil {
		return nil, err
	}
	arr := make([]TruncatedBlackboardNotes, len(result))
	for i, record := range result {
		props := GetPropsMapForRecordPosition(record, 0)
		arr[i] = TruncatedBlackboardNotes{
			ID:     props.GetString("id"),
			Title:  props.GetString("title"),
			Author: props.GetString("author"),
		}
	}
	return arr, err
}

func SearchForNotes(amount int, page int, query string) ([]TruncatedBlackboardNotes, error) {
	title, author := queryAnalyzer(query)
	result, err := makeRequest(`MATCH (n:Note) 
WHERE n.removed = false AND n.title CONTAINS $title AND n.author CONTAINS $author 
RETURN n ORDER BY n.posted_at DESC SKIP $skip LIMIT $amount;`,
		map[string]any{
			"amount": amount,
			"skip":   (page - 1) * amount,
			"title":  title,
			"author": author,
		})
	if err != nil {
		return nil, err
	}
	arr := make([]TruncatedBlackboardNotes, len(result))
	for i, record := range result {
		props := GetPropsMapForRecordPosition(record, 0)
		arr[i] = TruncatedBlackboardNotes{
			ID:       props.GetString("id"),
			Title:    props.GetString("title"),
			Author:   props.GetString("author"),
			Flair:    props.GetString("flair"),
			PostedAt: props.GetTime("posted_at"),
		}
	}
	return arr, err
}

func queryAnalyzer(query string) (title string, author string) {
	if strings.Contains(query, "BY:") {
		res := strings.SplitN(query, "BY:", 2)
		title = strings.TrimSpace(res[0])
		author = strings.TrimSpace(res[1])
	} else {
		title = strings.TrimSpace(query)
	}
	return
}
