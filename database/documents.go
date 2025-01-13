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

func (d *Document) IsDiscussion() bool { return d.Type == DocTypeDiscussion }

func (d *Document) IsVote() bool { return d.Type == DocTypeVote }

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
