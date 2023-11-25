package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"time"
)

var EternatityPublicationName = "rollingNewsletter"

type (
	Publication struct {
		UUID         string    `gorm:"primaryKey"`
		CreateTime   time.Time `gorm:"column:creation_time"`
		PublishTime  time.Time `gorm:"column:publication_time"`
		Publicated   bool
		BreakingNews bool `gorm:"column:hast"`
	}
	PublicationList []Publication
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
