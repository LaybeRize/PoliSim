package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type (
	DocumentType     string
	DocumentTypeList []DocumentType
	DocumentList     []Document
	Document         struct {
		UUID         string `gorm:"primaryKey"`
		Written      time.Time
		Organisation string
		Type         DocumentType
		Author       string
		Flair        string
		Title        string
		Subtitle     sql.NullString
		HTMLContent  string
		Private      bool
		Blocked      bool
		Info         DocumentInfo `gorm:"type:jsonb"`
		Viewer       []Account    `gorm:"many2many:doc_viewer;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
		Poster       []Account    `gorm:"many2many:doc_poster;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
		Allowed      []Account    `gorm:"many2many:doc_allowed;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
	}
	DocumentInfo struct {
		AnyPosterAllowed          bool          `json:"anyPosterAllowed"`
		OrganisationPosterAllowed bool          `json:"organisationPosterAllowed"`
		Finishing                 time.Time     `json:"time"`
		Post                      []Posts       `json:"post"`
		Discussion                []Discussions `json:"discussion"`
		Votes                     []string      `json:"vote"`
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

func (docI *DocumentInfo) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &docI)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &docI)
		return err
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (docI *DocumentInfo) Value() driver.Value {
	l, _ := json.Marshal(&docI)
	return l
}

func (docIValue DocumentTypeList) Value() (arr []string) {
	for _, val := range docIValue {
		arr = append(arr, string(val))
	}
	return
}

const (
	LegislativeText    DocumentType = "legislative_text"
	RunningDiscussion  DocumentType = "running_discussion"
	FinishedDiscussion DocumentType = "finished_discussion"
	RunningVote        DocumentType = "running_vote"
	FinishedVote       DocumentType = "finished_vote"
)
