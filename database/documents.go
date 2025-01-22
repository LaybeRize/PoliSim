package database

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"html/template"
	"strings"
	"time"
)

type DocumentType int

const (
	DocTypePost DocumentType = iota
	DocTypeDiscussion
	DocTypeVote
)

type (
	Document struct {
		ID                  string
		Type                DocumentType
		Organisation        string
		Title               string
		Author              string
		Flair               string
		Written             time.Time
		Body                template.HTML
		Public              bool
		Removed             bool
		MemberParticipation bool
		AdminParticipation  bool
		AllowedToAddTags    bool
		End                 time.Time
		Reader              []string
		Participants        []string
		Tags                []DocumentTag
		Links               []VoteInfo
		VoteIDs             []string
		Comments            []DocumentComment
	}
	SmallDocument struct {
		ID           string
		Type         DocumentType
		Organisation string
		Title        string
		Author       string
		Written      time.Time
		Removed      bool
	}
)

func (s *SmallDocument) IsPost() bool { return s.Type == DocTypePost }

func (s *SmallDocument) IsDiscussion() bool { return s.Type == DocTypeDiscussion }

func (s *SmallDocument) IsVote() bool { return s.Type == DocTypeVote }

func (s *SmallDocument) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return s.Written.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return s.Written.Format("2006-01-02 15:04:05 MST")
}

func (d *Document) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return d.Written.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return d.Written.Format("2006-01-02 15:04:05 MST")
}

func (d *Document) GetTimeEnd(a *Account) string {
	if a.Exists() {
		return d.End.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return d.End.Format("2006-01-02 15:04:05 MST")
}

func (d *Document) GetAuthor() string {
	if d.Flair == "" {
		return d.Author
	}
	return d.Author + "; " + d.Flair
}

func (d *Document) GetReader() string {
	if d.Public {
		return "Jeder kann dieses Dokument lesen."
	} else if len(d.Reader) == 0 {
		return "Leser: Alle Organisationsmitglieder"
	}
	return fmt.Sprintf("Leser: Alle Organisationsmitglieder plus %s", strings.Join(d.Reader, ", "))
}

func (d *Document) GetParticipants() string {
	if d.Type == DocTypePost {
		return "Nur Administratoren der Organisation dürfen Tags hinzufügen"
	} else if d.MemberParticipation && len(d.Participants) == 0 {
		return "Beteiligte: Alle Mitglieder der Organisation"
	} else if d.MemberParticipation {
		return fmt.Sprintf("Beteiligte: Alle Mitglieder der Organisation plus %s",
			strings.Join(d.Reader, ", "))
	} else if d.AdminParticipation {
		return "Beteiligte: Alle Administratoren der Organisation"
	}
	return fmt.Sprintf("Beteiligte: %s", strings.Join(d.Reader, ", "))
}

func (d *Document) IsPost() bool { return d.Type == DocTypePost }

func (d *Document) IsDiscussion() bool { return d.Type == DocTypeDiscussion }

func (d *Document) IsVote() bool { return d.Type == DocTypeVote }

func (d *Document) Ended() bool { return time.Now().After(d.End) }

type DocumentTag struct {
	ID              string
	Outgoing        bool
	Text            string
	Written         time.Time
	BackgroundColor string
	TextColor       string
	LinkColor       string
	Links           []string
	QueriedLinks    []any
}

func (t *DocumentTag) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return t.Written.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return t.Written.Format("2006-01-02 15:04:05 MST")
}

type DocumentComment struct {
	ID      string
	Author  string
	Flair   string
	Written time.Time
	Body    template.HTML
	Removed bool
}

func (c *DocumentComment) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return c.Written.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return c.Written.Format("2006-01-02 15:04:05 MST")
}

func (c *DocumentComment) GetAuthor() string {
	if c.Flair == "" {
		return c.Author
	}
	return c.Author + "; " + c.Flair
}

func CreateDocument(document *Document, acc *Account) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `MATCH (acc:Account)-[:ADMIN]->(o:Organisation) 
WHERE acc.name = $Author AND acc.blocked = false AND o.name = $organisation AND o.visibility <> $hidden
	RETURN o.visibility;`,
		map[string]any{"Author": document.Author,
			"organisation": document.Organisation,
			"hidden":       HIDDEN})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	} else if result.Next(ctx); result.Record() == nil {
		_ = tx.Rollback(ctx)
		return notAllowedError
	} else if result.Record().Values[0].(int64) == int64(SECRET) && document.Public {
		_ = tx.Rollback(ctx)
		return notAllowedError
	}

	_, err = tx.Run(ctx, `MATCH (a:Account) WHERE a.name = $author 
MATCH (o:Organisation) WHERE o.name = $organisation 
CREATE (d:Document {id: $id, title: $title, type: $type, author: $author, flair: $flair, body: $body, removed: false,
end_time: $end_time, written: $written, public: $public, member_part: $member_part, admin_part: $admin_part}) 
MERGE (a)-[:WRITTEN]->(d)
MERGE (d)-[:IN]->(o);`, map[string]any{
		"organisation": document.Organisation,
		"id":           document.ID,
		"title":        document.Title,
		"type":         document.Type,
		"author":       document.Author,
		"end_time":     document.End,
		"written":      time.Now().UTC(),
		"body":         document.Body,
		"flair":        document.Flair,
		"public":       document.Public,
		"member_part":  document.MemberParticipation,
		"admin_part":   document.AdminParticipation})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if !document.Public {
		_, err = tx.Run(ctx, `MATCH (a:Account), (d:Document) WHERE a.name IN $reader AND d.id = $id 
CREATE (a)<-[:READER]-(d);`, map[string]any{"reader": document.Reader, "id": document.ID})
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}

	if document.Type == DocTypeVote {
		result, err = tx.Run(ctx, `
MATCH (a:Account)-[r:MANAGES]->(v:Vote) WHERE a.name = $user AND v.id IN $ids 
MATCH (d:Document) WHERE d.id = $id 
DELETE r 
MERGE (d)-[:LINKS]->(v) 
RETURN v.id;`, map[string]any{
			"user": acc.Name,
			"ids":  document.VoteIDs,
			"id":   document.ID})
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		} else if result.Next(ctx); result.Record() == nil {
			_ = tx.Rollback(ctx)
			return notAllowedError
		}
	}

	_, err = tx.Run(ctx, `MATCH (a:Account), (d:Document) WHERE a.name IN $user AND d.id = $id 
CREATE (a)<-[:PARTICIPANT]-(d);`, map[string]any{"user": document.Participants, "id": document.ID})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func GetDocumentForUser(id string, acc *Account) (*Document, []string, error) {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return nil, nil, err
	}

	result, err := tx.Run(ctx, `MATCH (o:Organisation)<-[:IN]-(d:Document) WHERE d.id = $id 
RETURN d, o.name;`, map[string]any{
		"id": id})
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, nil, err
	} else if result.Next(ctx); result.Record() == nil {
		_ = tx.Rollback(ctx)
		return nil, nil, notAllowedError
	}
	props := result.Record().Values[0].(neo4j.Node).Props
	public := props["public"].(bool)
	allowedToAddTags := false

	if !public && !acc.Exists() {
		_ = tx.Rollback(ctx)
		return nil, nil, notAllowedError
	} else if !public {
		var userCheck neo4j.ResultWithContext
		userCheck, err = tx.Run(ctx, `
CALL { 
MATCH (a:Account)-[*..]->(o:Organisation)<-[:IN]-(d:Document) WHERE a.name = $name AND d.id = $id 
RETURN o, d 
UNION 
MATCH (o:Organisation)<-[:IN]-(d:Document)-[:READER|PARTICIPANT]->(a:Account) WHERE a.name = $name AND d.id = $id 
RETURN o,d 
} 
OPTIONAL MATCH (b:Account)-[:OWNER|ADMIN*..]->(o) WHERE b.name = $name 
RETURN d.id, b.name;`, map[string]any{
			"id":   id,
			"name": acc.Name})
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, nil, err
		} else if userCheck.Next(ctx); userCheck.Record() == nil && !acc.IsAtLeastAdmin() {
			_ = tx.Rollback(ctx)
			return nil, nil, notAllowedError
		} else if userCheck.Record() != nil {
			allowedToAddTags = userCheck.Record().Values[1] != nil
		}

	}

	doc := &Document{
		ID:                  id,
		Type:                DocumentType(props["type"].(int64)),
		Organisation:        result.Record().Values[1].(string),
		Title:               props["title"].(string),
		Author:              props["author"].(string),
		Flair:               props["flair"].(string),
		Written:             props["written"].(time.Time),
		Body:                template.HTML(props["body"].(string)),
		Public:              public,
		Removed:             props["removed"].(bool),
		MemberParticipation: props["member_part"].(bool),
		AdminParticipation:  props["admin_part"].(bool),
		End:                 props["end_time"].(time.Time),
		Tags:                make([]DocumentTag, 0),
		AllowedToAddTags:    allowedToAddTags,
	}
	var commentator []string

	if !doc.Public {
		doc.Reader = make([]string, 0)
		result, err = tx.Run(ctx, `MATCH (d:Document)-[:READER]->(a:Account)
WHERE d.id = $id RETURN a.name;`,
			map[string]any{"id": id})
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, nil, err
		}
		for result.Next(ctx) {
			doc.Reader = append(doc.Reader, result.Record().Values[0].(string))
		}
	}

	if !(doc.Type == DocTypePost) {
		doc.Participants = make([]string, 0)
		result, err = tx.Run(ctx, `MATCH (d:Document)-[:PARTICIPANT]->(a:Account)
WHERE d.id = $id RETURN a.name;`,
			map[string]any{"id": id})
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, nil, err
		}
		for result.Next(ctx) {
			doc.Participants = append(doc.Participants, result.Record().Values[0].(string))
		}
	}

	if doc.Type == DocTypeDiscussion {
		doc.Comments = make([]DocumentComment, 0)
		result, err = tx.Run(ctx, `MATCH (d:Document)<-[:ON]-(c:Comment)
WHERE d.id = $id RETURN c.id, c.author, c.flair, c.written, c.body, c.removed;`,
			map[string]any{"id": id})
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, nil, err
		}
		for result.Next(ctx) {
			doc.Comments = append(doc.Comments, DocumentComment{
				ID:      result.Record().Values[0].(string),
				Author:  result.Record().Values[1].(string),
				Flair:   result.Record().Values[2].(string),
				Written: result.Record().Values[3].(time.Time),
				Body:    template.HTML(result.Record().Values[4].(string)),
				Removed: result.Record().Values[5].(bool),
			})
		}
	}

	if doc.Type == DocTypeDiscussion && !doc.Ended() && acc.Exists() {
		commentator = make([]string, 0)
		result, err = tx.Run(ctx, `
MATCH (n:Account)-[:OWNER*0..]->(a:Account) WHERE n.name = $name 
RETURN a.name;`,
			map[string]any{"id": id, "name": acc.Name})
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, nil, err
		}
		for result.Next(ctx) {
			commentator = append(commentator, result.Record().Values[0].(string))
		}
	}

	if doc.Type == DocTypeVote {
		// Todo depending on if the vote has ended try also querying the results in addition
		doc.Links = make([]VoteInfo, 0)
		result, err = tx.Run(ctx, `MATCH (d:Document)-[:LINKS]->(v:Vote)
WHERE d.id = $id RETURN v.id, v.question;`,
			map[string]any{"id": id})
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, nil, err
		}
		for result.Next(ctx) {
			doc.Links = append(doc.Links, VoteInfo{
				ID:       result.Record().Values[0].(string),
				Question: result.Record().Values[1].(string),
			})
		}
	}

	result, err = tx.Run(ctx, `CALL { 
MATCH (d:Document)-[:LINKS]->(t:Tag)-[:LINKS]->(r:Document) 
WHERE d.id = $id 
RETURN t, collect(r.id) AS ids, true AS outgoing 
UNION 
MATCH (d:Document)<-[:LINKS]-(t:Tag)<-[:LINKS]-(r:Document) 
WHERE d.id = $id 
RETURN t, collect(r.id) AS ids, false AS outgoing 
} 
RETURN t, ids, outgoing ORDER BY t.written DESC;`, map[string]any{"id": id})
	if err != nil {
		return nil, nil, err
	}
	for result.Next(ctx) {
		values := result.Record().Values
		props = values[0].(neo4j.Node).Props
		doc.Tags = append(doc.Tags, DocumentTag{
			ID:              props["id"].(string),
			Outgoing:        values[2].(bool),
			Text:            props["text"].(string),
			Written:         props["written"].(time.Time),
			BackgroundColor: props["background"].(string),
			TextColor:       props["color"].(string),
			LinkColor:       props["link"].(string),
			QueriedLinks:    values[1].([]any),
		})
	}

	err = tx.Commit(ctx)
	return doc, commentator, err
}

func CreateTagForDocument(docID string, acc *Account, tag *DocumentTag) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `MATCH (a:Account)-[:ADMIN|OWNER*..]->(o:Organisation)<-[:IN]-(d:Document) 
WHERE a.name = $name AND d.id = $id 
RETURN a.name;`, map[string]any{"name": acc.Name, "id": docID})
	if err != nil {
		return err
	} else if result.Next(ctx); result.Record() == nil {
		return notAllowedError
	}

	_, err = tx.Run(ctx, `MATCH (d:Document) WHERE d.id = $id  
CREATE (t:Tag {id: $tagId, text: $text, written: $written, background: $background, color: $color, link: $link}) 
MERGE (d)-[:LINKS]->(t);`, map[string]any{
		"id":         docID,
		"tagId":      tag.ID,
		"links":      tag.Links,
		"text":       tag.Text,
		"written":    tag.Written,
		"background": tag.BackgroundColor,
		"color":      tag.TextColor,
		"link":       tag.LinkColor,
	})
	if err != nil {
		return err
	}

	result, err = tx.Run(ctx, `MATCH (target:Document) WHERE target.id <> $id AND target.id IN $links 
MATCH (t:Tag) WHERE t.id = $tagID 
MERGE (t)-[:LINKS]->(target);`, map[string]any{"tagID": tag.ID, "id": docID, "links": tag.Links})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func CreateDocumentComment(documentId string, comment *DocumentComment) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `CALL {
MATCH (a:Account)<-[:PARTICIPANT]-(d:Document) 
WHERE a.name = $author AND a.blocked = false AND d.id = $id AND $now > d.end_time
RETURN a 
UNION 
MATCH (a:Account)-[:USER|ADMIN]->(:Organisation)<-[:IN]-(d:Document) 
WHERE a.name = $author AND a.blocked = false AND d.id = $id AND d.member_part = true AND $now > d.end_time
RETURN a 
UNION 
MATCH (a:Account)-[:ADMIN]->(:Organisation)<-[:IN]-(d:Document) 
WHERE a.name = $author AND a.blocked = false AND d.id = $id AND d.admin_part = true AND $now > d.end_time
RETURN a 
} 
RETURN a;`,
		map[string]any{"id": documentId, "author": comment.Author, "now": time.Now().UTC()})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	} else if result.Next(ctx); result.Record() == nil {
		_ = tx.Rollback(ctx)
		return notAllowedError
	}

	_, err = tx.Run(ctx, `MATCH (a:Account) WHERE a.name = $author 
MATCH (d:Document) WHERE d.id = $doc_ID 
CREATE (c:Comment {id: $id, author: $author, flair: $flair, body: $body, removed: false, written: $written}) 
MERGE (a)-[:WRITTEN]->(c)
MERGE (c)-[:ON]->(d);`, map[string]any{
		"doc_ID":  documentId,
		"id":      comment.ID,
		"author":  comment.Author,
		"flair":   comment.Flair,
		"written": time.Now().UTC(),
		"body":    comment.Body})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func GetDocumentList(amount int, page int, acc *Account) ([]SmallDocument, error) {
	var query string
	if !acc.Exists() {
		query = `MATCH (o:Organisation)<-[:IN]-(d:Document) WHERE d.public = true `
	} else if !acc.IsAtLeastAdmin() {
		query = `CALL { MATCH (o:Organisation)<-[:IN]-(d:Document) WHERE d.public = true 
RETURN d, o
UNION
MATCH (a:Account)-[*..]->(o:Organisation)<-[:IN]-(d:Document) WHERE d.public = false AND a.name = $name 
RETURN d, o 
} `
	} else {
		query = `MATCH (o:Organisation)<-[:IN]-(d:Document) WHERE true `
	}

	result, err := makeRequest(query+`RETURN d.id, d.type, o.name, d.title, d.author, d.written, d.removed 
ORDER BY d.written DESC SKIP $skip LIMIT $amount;`,
		map[string]any{
			"amount": amount,
			"skip":   (page - 1) * amount,
			"name":   acc.GetName(),
		})
	if err != nil {
		return nil, err
	}
	arr := make([]SmallDocument, 0, len(result.Records))
	for _, record := range result.Records {
		arr = append(arr, SmallDocument{
			ID:           record.Values[0].(string),
			Type:         DocumentType(record.Values[1].(int64)),
			Organisation: record.Values[2].(string),
			Title:        record.Values[3].(string),
			Author:       record.Values[4].(string),
			Written:      record.Values[5].(time.Time),
			Removed:      record.Values[6].(bool),
		})
	}
	return arr, nil
}
