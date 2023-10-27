package database

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
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

func createEternatityPublicationIfNotExist() {
	pub := Publication{}
	err := DB.Where("uuid = ?", EternatityPublicationName).First(&pub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		pub.UUID = EternatityPublicationName
		err = DB.Create(&pub).Error
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while finding/creating rolling publication:\n"+err.Error()+"\n")
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Rolling publication created/found\n")
}
