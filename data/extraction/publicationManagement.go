package extraction

import (
	"PoliSim/data/database"
	"errors"
	"gorm.io/gorm"
	"time"
)

func CreatePublication(pub *database.Publication) error {
	return database.DB.Create(&pub).Error
}

func GetHiddenPublications() (*database.PublicationList, error) {
	list := &database.PublicationList{}
	err := database.DB.Where("publicated = false").Order("creation_time").Find(list).Error
	return list, err
}

func FindPublication(uuid string, visible string) (bool, error) {
	pub := &database.Publication{}
	err := database.DB.Where("publicated = ? AND uuid = ?", visible, uuid).First(pub).Error
	if err == nil {
		return true, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func FindPublicationAndReturnIt(uuid string, visible string) (*database.Publication, error) {
	pub := &database.Publication{}
	err := database.DB.Where("publicated = ? AND uuid = ?", visible, uuid).First(pub).Error
	return pub, err
}

func GetPublicationAfter(publicationUUID string, amount int) (publicationList *database.PublicationList, exists bool, err error) {
	publicationList = &database.PublicationList{}
	exists = true
	var pub *database.Publication
	pub, err = FindPublicationAndReturnIt(publicationUUID, "true")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		exists = false
		pub.PublishTime = time.Now()
	} else if err != nil {
		return
	}
	err = getBasicNewsQuery(publicationUUID, amount).Where("publication_time < ?", pub.PublishTime).Order("publication_time desc").Find(publicationList).Error
	return
}

func GetPublicationBefore(publicationUUID string, amount int) (publicationList *database.PublicationList, exists bool, err error) {
	exists = true
	var pub *database.Publication
	pub, err = FindPublicationAndReturnIt(publicationUUID, "true")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		exists = false
		pub.PublishTime = time.Now()
	} else if err != nil {
		return
	}
	err = database.DB.Select("*").Table("(?) as X", getBasicNewsQuery(publicationUUID, amount).Order("publication_time").Where("publication_time > ?", pub.PublishTime)).Order("X.publication_time desc").Find(publicationList).Error
	return
}

func getBasicNewsQuery(uuid string, amount int) *gorm.DB {
	return database.DB.Select("*").Where("uuid != ? AND publicated = true", uuid).Limit(amount).Table("publications")
}
