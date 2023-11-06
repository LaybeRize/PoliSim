package extraction

import "PoliSim/data/database"

func CreateDocument(doc *database.Document) error {
	return database.DB.Create(&doc).Error
}
