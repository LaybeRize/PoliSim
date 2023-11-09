package extraction

import "PoliSim/data/database"

func CreateDocument(doc *database.Document) error {
	return database.DB.Create(&doc).Error
}

func GetDocument(docType database.DocumentType, uuid string) (*database.Document, error) {
	doc := &database.Document{}
	err := database.DB.Where("type = ? AND uuid = ?", string(docType), uuid).First(doc).Error
	return doc, err
}
