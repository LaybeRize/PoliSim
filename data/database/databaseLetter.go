package database

import (
	"time"
)

type (
	LetterList []Letter
	Letter     struct {
		UUID        string `gorm:"primaryKey"`
		Written     time.Time
		Author      string
		Flair       string
		Title       string
		Content     string
		HTMLContent string     `gorm:"column:html_content"`
		Info        LetterInfo `gorm:"type:jsonb;serializer:json"`
		Viewer      []Account  `gorm:"many2many:letter_account;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
		Removed     bool
		ModMessage  bool `gorm:"column:mod_message"`
	}
	LetterInfo struct {
		AllHaveToAgree     bool     `json:"allAgree"`
		NoSigning          bool     `json:"noSigning"`
		PeopleNotYetSigned []string `json:"notSigned"`
		Signed             []string `json:"signed"`
		Rejected           []string `json:"rejected"`
	}
)
