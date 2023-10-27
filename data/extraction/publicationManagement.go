package extraction

import "PoliSim/data/database"

func CreatePublication(pub *database.Publication) error {
	return database.DB.Create(&pub).Error
}
