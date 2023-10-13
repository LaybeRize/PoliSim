package database

import "database/sql"

type (
	TitleList []Title
	Title     struct {
		Name      string `gorm:"primaryKey"`
		MainGroup string
		SubGroup  string
		Flair     sql.NullString
		Holder    []Account `gorm:"many2many:title_account;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
	}
)
