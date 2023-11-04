package database

import (
	"database/sql"
	"time"
)

type (
	DocumentType     string
	DocumentTypeList []DocumentType
	DocumentList     []Document
	Document         struct {
		UUID                      string `gorm:"primaryKey"`
		Written                   time.Time
		Organisation              string
		Type                      DocumentType
		Author                    string
		Flair                     string
		Title                     string
		Subtitle                  sql.NullString
		HTMLContent               string
		Private                   bool
		Blocked                   bool
		CurrentPostTag            string       `gorm:"column:current_tag_info"`
		AnyPosterAllowed          bool         `gorm:"column:allowed_any"`
		OrganisationPosterAllowed bool         `gorm:"column:allowed_members"`
		Info                      DocumentInfo `gorm:"type:jsonb;serializer:json"`
		Viewer                    []Account    `gorm:"many2many:doc_viewer;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
		Poster                    []Account    `gorm:"many2many:doc_poster;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
		Allowed                   []Account    `gorm:"many2many:doc_allowed;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
	}
	DocumentInfo struct {
		Finishing  time.Time     `json:"time"`
		Post       []Posts       `json:"post"`
		Discussion []Discussions `json:"discussion"`
		Votes      []string      `json:"vote"`
	}
	Posts struct {
		UUID      string    `json:"uuid"`
		Hidden    bool      `json:"hidden"`
		Submitted time.Time `json:"submitted"`
		Info      string    `json:"info"`
		Color     string    `json:"color"`
	}
	Discussions struct {
		UUID        string    `json:"uuid"`
		Hidden      bool      `json:"hidden"`
		Written     time.Time `json:"written"`
		Author      string    `json:"author"`
		Flair       string    `json:"flair"`
		HTMLContent string    `json:"htmlContent"`
	}
)

const (
	LegislativeText    DocumentType = "legislative_text"
	RunningDiscussion  DocumentType = "running_discussion"
	FinishedDiscussion DocumentType = "finished_discussion"
	RunningVote        DocumentType = "running_vote"
	FinishedVote       DocumentType = "finished_vote"
)
