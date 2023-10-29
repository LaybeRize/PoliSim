package extraction

import (
	"PoliSim/data/database"
	"errors"
	"gorm.io/gorm"
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
