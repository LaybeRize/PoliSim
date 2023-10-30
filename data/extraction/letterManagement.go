package extraction

import "PoliSim/data/database"

func CreateLetter(letter *database.Letter) error {
	return database.DB.Create(&letter).Error
}
