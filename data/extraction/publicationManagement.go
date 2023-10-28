package extraction

import "PoliSim/data/database"

func CreatePublication(pub *database.Publication) error {
	return database.DB.Create(&pub).Error
}

func GetHiddenPublication() (*database.PublicationList, error) {
	list := &database.PublicationList{}
	err := database.DB.Where("publicated = false").Order("creation_time").Find(list).Error
	return list, err
}
