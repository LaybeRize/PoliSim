package extraction

import (
	"PoliSim/data/database"
	"gorm.io/gorm"
)

type (
	SpecificDocument struct {
		UUID   string
		Title  string
		Viewer []AccountDisplayName `gorm:"many2many:doc_viewer;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
	}
	QueryInterface interface {
		GetSelection() string
		GetSelf() any
	}
	AdvancedQueryInterface interface {
		QueryInterface
		GetPreload(DB *gorm.DB) *gorm.DB
	}
)

func (*SpecificDocument) GetSelection() string {
	return "documents.uuid, documents.title"
}

func (*SpecificDocument) GetPreload(DB *gorm.DB) *gorm.DB {
	return DB.Preload("Viewer", func(db *gorm.DB) *gorm.DB {
		return db.Model(database.Account{}).Select("id, display_name").Order("id")
	})
}

func (d *SpecificDocument) GetSelf() any { return d }

func FindDocumentInterfaceByUUID(uuid string, queryInterface AdvancedQueryInterface) error {
	return queryInterface.GetPreload(database.DB.Model(&database.Document{})).Select(queryInterface.GetSelection()).
		Where("documents.uuid = ?", uuid).First(queryInterface.GetSelf()).Error
}
