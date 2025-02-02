package database

import (
	loc "PoliSim/localisation"
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
	props := GetPropsMapForRecordPosition(result[0], 0)
	note := &BlackboardNote{
		ID:       props.GetString("id"),
		Title:    props.GetString("title"),
		Author:   props.GetString("author"),
		Flair:    props.GetString("flair"),
		PostedAt: props.GetTime("posted_at"),
		Body:     template.HTML(props.GetString("body")),
		Removed:  props.GetBool("removed"),
	}

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
			ID:      props.GetString("id"),
			Title:   props.GetString("title"),
			Author:  props.GetString("author"),
			Removed: props.GetBool("removed"),
		}
	}
	return arr, err
}

func SearchForNotes(acc *Account, amount int, page int, input string, showBlocked bool) ([]TruncatedBlackboardNotes, error) {
	query, parameter := queryAnalyzer(acc, input, showBlocked)
	parameter["amount"] = amount
	parameter["skip"] = (page - 1) * amount
	result, err := makeRequest(`MATCH (n:Note) 
WHERE `+query+` 
RETURN n ORDER BY n.posted_at DESC SKIP $skip LIMIT $amount;`,
		parameter)
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

var queryRegexNotes = regexp.MustCompile(`^\s*(.*?)\s*(\[|$)`)
var authorRegexNotes = regexp.MustCompile(`\[[bB][yY]:\]\s*(.+?)\s*(\[|$)`)

func queryAnalyzer(acc *Account, input string, showBlocked bool) (query string, parameter map[string]any) {
	parameter = make(map[string]any)
	query = ""
	if showBlocked && acc.IsAtLeastAdmin() {
		query += "true"
	} else {
		query += "n.removed = false"
	}

	result := queryRegexNotes.FindStringSubmatch(input)
	if result != nil && result[1] != "" {
		parameter["title"] = result[1]
		query += " AND n.title CONTAINS $title"
	}

	if result = authorRegexNotes.FindStringSubmatch(input); result != nil {
		parameter["author"] = result[1]
		query += " AND n.author CONTAINS $author"
	}

	return
}
