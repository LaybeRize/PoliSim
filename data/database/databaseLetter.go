package database

import (
	"time"
)

type (
	Letter struct {
		UUID        string `gorm:"primaryKey"`
		Written     time.Time
		Author      string
		Flair       string
		Title       string
		Content     string
		HTMLContent string
		Info        LetterInfo `gorm:"type:jsonb;serializer:json"`
		Viewer      []Account  `gorm:"many2many:letter_account;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
		Removed     bool
		ModMessage  bool `gorm:"column:mod_message"`
	}
	LetterList []Letter

	LetterInfo struct {
		AllHaveToAgree     bool     `json:"allAgree"`
		NoSigning          bool     `json:"noSigning"`
		PeopleNotYetSigned []string `json:"notSigned"`
		Signed             []string `json:"signed"`
		Rejected           []string `json:"rejected"`
	}
	LetterAccount struct {
		UUID string `gorm:"primaryKey"`
		ID   int64  `gorm:"primaryKey"`
		Read bool   `gorm:"default:false"`
	}
	ExtendedLetterList []ExtendedLetter
	ExtendedLetter     struct {
		Letter
		Read bool
	}
)

func (LetterAccount) TableName() string {
	return "letter_account"
}
