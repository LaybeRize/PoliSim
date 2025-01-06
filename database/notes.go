package database

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx,
		`CREATE (:Note {id: $id , title: $title , author: $Author , flair: $Flair, 
posted_at: $PostedAt , body: $Body, removed: $Removed});`,
		map[string]any{"id": note.ID,
			"title":    note.Title,
			"Author":   note.Author,
			"Flair":    note.Flair,
			"PostedAt": note.PostedAt,
			"Body":     string(note.Body),
			"Removed":  note.Removed})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx,
		`MATCH (c:Note), (p:Note) WHERE c.id = $child AND p.id IN $parent CREATE (c)-[:LINKS]->(p);`,
		map[string]any{"parent": references,
			"child": note.ID})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func GetNote(id string) (*BlackboardNote, error) {
	idMap := map[string]any{"id": id}
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (n:Note) WHERE n.id = $id RETURN n;`,
		idMap, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	if len(result.Records) == 0 {
		return nil, notFoundError
	} else if len(result.Records) > 1 {
		return nil, multipleItemsError
	}
	note := &BlackboardNote{}
	node := result.Records[0].Values[0].(neo4j.Node)
	note.ID = node.Props["id"].(string)
	note.Title = node.Props["title"].(string)
	note.Author = node.Props["author"].(string)
	note.Flair = node.Props["flair"].(string)
	note.PostedAt = node.Props["posted_at"].(time.Time)
	note.Body = template.HTML(node.Props["body"].(string))
	note.Removed = node.Props["removed"].(bool)
	note.Parents, err = queryForRelations(`MATCH (n:Note) WHERE n.id = $id MATCH (r:Note)-[:LINKS]->(n) RETURN r;`, idMap)
	if err != nil {
		return nil, err
	}
	note.Children, err = queryForRelations(`MATCH (n:Note) WHERE n.id = $id MATCH (n)-[:LINKS]->(r:Note) RETURN r;`, idMap)
	return note, err
}

func queryForRelations(query string, idMap map[string]any) ([]TruncatedBlackboardNotes, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, query,
		idMap, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	arr := make([]TruncatedBlackboardNotes, len(result.Records))
	for i, record := range result.Records {
		node := record.Values[0].(neo4j.Node)
		arr[i] = TruncatedBlackboardNotes{
			ID:     node.Props["id"].(string),
			Title:  node.Props["title"].(string),
			Author: node.Props["author"].(string),
		}
	}
	return arr, err
}

func SearchForNotes(amount int, page int, query string) ([]TruncatedBlackboardNotes, error) {
	title, author := queryAnalyzer(query)
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (n:Note) 
WHERE n.removed = false AND n.title CONTAINS $title AND n.author CONTAINS $author 
RETURN n ORDER BY n.posted_at DESC SKIP $skip LIMIT $amount;`,
		map[string]any{
			"amount": amount,
			"skip":   (page - 1) * amount,
			"title":  title,
			"author": author,
		}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	arr := make([]TruncatedBlackboardNotes, len(result.Records))
	for i, record := range result.Records {
		node := record.Values[0].(neo4j.Node)
		arr[i] = TruncatedBlackboardNotes{
			ID:       node.Props["id"].(string),
			Title:    node.Props["title"].(string),
			Author:   node.Props["author"].(string),
			Flair:    node.Props["flair"].(string),
			PostedAt: node.Props["posted_at"].(time.Time),
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
