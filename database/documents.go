package database

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

const (
	DocTypePost       = "post"
	DocTypeDiscussion = "discussion"
	DocTypeVote       = "vote"
)

type Document struct {
	ID                  string
	Type                string
	Organisation        string
	Title               string
	Author              string
	Flair               string
	Written             time.Time
	Body                template.HTML
	Public              bool
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

func CreateDocument(document *Document) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	// Add validation logic for author/organisation etc.
	/*result, err := tx.Run(ctx, `MATCH (acc:Account) WHERE acc.name = $Author AND acc.blocked = false
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
		}*/

	_, err = tx.Run(ctx, `MATCH (a:Account) WHERE a.name = $author 
MATCH (o:Organisation) WHERE o.name = $organisation 
CREATE (d:Document {id: $id, title: $title, type: $type, author: $author, flair: $flair, body: $body,
end_time: $end_time, written: $written, public: $public, member_part: $member_part, admin_part: $admin_part}) 
MERGE (a)-[:WRITTEN]->(d)
MERGE (d)-[:IN]->(o);`, map[string]any{
		"organisation": document.Organisation,
		"id":           document.ID,
		"title":        document.Title,
		"type":         document.Type,
		"author":       document.Author,
		"end_time":     document.End,
		"written":      document.Written,
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
