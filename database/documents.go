package database

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

type DocumentType string

const (
	DocTypePost       DocumentType = "post"
	DocTypeDiscussion DocumentType = "discussion"
	DocTypeVote       DocumentType = "vote"
)

type Document struct {
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
	Reader              []string
	MemberParticipation bool
	AdminParticipation  bool
	End                 time.Time
	Participants        []string
	Tags                []DocumentTag
	Links               []string
	Comments            []DocumentComment
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
	Text            string
	Written         time.Time
	BackgroundColor string
	TextColor       string
	LinkColor       string
	Links           []string
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

type ColorPalette struct {
	Name       string
	Background string
	Text       string
	Link       string
}

func CreateDocument(document *Document) error {
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
	} else if result.Record().Values[0].(string) == string(SECRET) && document.Public {
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

	_, err = tx.Run(ctx, `MATCH (a:Account), (d:Document) WHERE a.name IN $user AND d.id = $id 
CREATE (a)<-[:PARTICIPANT]-(d);`, map[string]any{"user": document.Participants, "id": document.ID})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func CreateDocumentComment(documentId string, comment *DocumentComment) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `CALL {
MATCH (a:Account)<-[:PARTICIPANT]-(d:Document) 
WHERE a.name = $author AND a.blocked = false AND d.id = $id 
RETURN a 
UNION 
MATCH (a:Account)-[:USER|ADMIN]->(:Organisation)<-[:IN]-(d:Document) 
WHERE a.name = $author AND a.blocked = false AND d.id = $id AND d.member_part = true 
RETURN a 
UNION 
MATCH (a:Account)-[:ADMIN]->(:Organisation)<-[:IN]-(d:Document) 
WHERE a.name = $author AND a.blocked = false AND d.id = $id AND d.admin_part = true 
RETURN a 
} 
RETURN a;`,
		map[string]any{"id": documentId, "author": comment.Author})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	} else if result.Next(ctx); result.Record() == nil {
		_ = tx.Rollback(ctx)
		return notAllowedError
	}

	_, err = tx.Run(ctx, `MATCH (a:Account) WHERE a.name = $author 
MATCH (d:Document) WHERE d.id = $doc_ID 
CREATE (c:Document {id: $id, author: $author, flair: $flair, body: $body, removed: false, written: $written}) 
MERGE (a)-[:WRITTEN]->(c)
MERGE (c)-[:ON]->(d);`, map[string]any{
		"doc_ID":  documentId,
		"id":      comment.ID,
		"author":  comment.Author,
		"flair":   comment.Flair,
		"written": comment.Written,
		"body":    comment.Body})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}
