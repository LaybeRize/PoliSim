package extraction

import (
	"PoliSim/data/database"
	"errors"
	"gorm.io/gorm"
)

func CreateDocument(doc *database.Document) error {
	return database.DB.Create(&doc).Error
}

func GetDocumentIfNotPrivate(docType database.DocumentType, uuid string) (*database.Document, error) {
	doc := &database.Document{}
	err := database.DB.Where("type = ? AND uuid = ? AND private = false AND blocked = false", string(docType), uuid).First(doc).Error
	return doc, err
}

func GetDocumentForUser(uuid string, userID int64, isAdmin bool, docType ...database.DocumentType) (*database.Document, error) {
	doc := &database.Document{}
	err := database.DB.Joins("LEFT JOIN organisation_account ON documents.organisation = organisation_account.name").
		Joins("LEFT JOIN doc_allowed ON doc_allowed.uuid = documents.uuid").
		Preload("Viewer", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, display_name")
		}).Preload("Poster", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, display_name")
	}).Select("documents.uuid, written, documents.organisation, type, author, flair, title, subtitle"+
		", html_content, private, blocked, info, allowed_any, allowed_members").
		Where("documents.uuid = ? AND blocked = false AND (private = false OR organisation_account.id = ? OR doc_allowed.id=? OR true = ?)",
			uuid, userID, userID, isAdmin).First(doc).Error
	for _, singleType := range docType {
		if doc.Type == singleType {
			return doc, err
		}
	}
	return doc, errors.New("is Not One Of the specified Types")
}

func GetDocument(uuid string) (*database.Document, error) {
	doc := &database.Document{}
	err := database.DB.Where("uuid = ?", uuid).First(doc).Error
	return doc, err
}

func UpdateDocument(document *database.Document) error {
	return database.DB.Updates(document).Error
}
