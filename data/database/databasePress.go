package database

import (
	"database/sql"
	"time"
)

var EternatityPublicationName = "theFirstOfThemAll"

type (
	PublicationList []Publication
	Publication     struct {
		UUID         string    `gorm:"primaryKey"`
		CreateTime   time.Time `gorm:"column:creation_time"`
		PublishTime  time.Time `gorm:"column:publication_time"`
		Publicated   bool
		BreakingNews bool `gorm:"column:hast"`
	}
	ArticleList []Article
	Article     struct {
		UUID        string `gorm:"primaryKey"`
		Publication string
		Written     time.Time
		Author      string
		Flair       string
		Headline    string
		Subtitle    sql.NullString
		Content     string
		HTMLContent string `gorm:"column:html_content"`
	}
)
