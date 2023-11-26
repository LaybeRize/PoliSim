package database

import (
	"database/sql"
	"time"
)

type (
	Zwitscher struct {
		UUID         string `gorm:"primaryKey"`
		Written      time.Time
		Author       string
		Flair        string
		HTMLContent  string
		Blocked      bool
		Linked       sql.NullString
		ReadByParent bool
		Parent       *Zwitscher  `gorm:"foreignKey:linked;joinReferences:uuid"`
		Children     []Zwitscher `gorm:"foreignKey:linked"`
	}
	ZwitscherList []Zwitscher
)
