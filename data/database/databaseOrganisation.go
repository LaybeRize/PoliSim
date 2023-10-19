package database

import (
	"database/sql"
)

type (
	StatusString     string
	OrganisationList []Organisation
	Organisation     struct {
		Name      string `gorm:"primaryKey"`
		MainGroup string
		SubGroup  string
		Flair     sql.NullString
		Status    StatusString
		Members   []Account `gorm:"many2many:organisation_member;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
		Admins    []Account `gorm:"many2many:organisation_admins;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
		Accounts  []Account `gorm:"many2many:organisation_account;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
	}
)

const (
	Public  StatusString = "public"
	Private StatusString = "private"
	Secret  StatusString = "secret"
	Hidden  StatusString = "hidden"
)

// TODO: add translation to .json
var (
	StatusTranslation = map[StatusString]string{
		Public:  "Öffentlich",
		Private: "Privat",
		Secret:  "Geheim",
		Hidden:  "Versteckt",
	}
	Stati = []StatusString{Public, Private, Secret, Hidden}
)

func (t StatusString) String() string {
	return string(t)
}
