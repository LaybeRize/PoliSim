package database

import (
	"database/sql"
	"time"
)

type (
	Article struct {
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
	ArticleList []Article
)
