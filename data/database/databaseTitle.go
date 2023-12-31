package database

import "database/sql"

type (
	Title struct {
		Name      string `gorm:"primaryKey"`
		MainGroup string
		SubGroup  string
		Flair     sql.NullString `gorm:"index:unique"`
		Holder    []Account      `gorm:"many2many:title_account;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
	}
	TitleList []Title
)
