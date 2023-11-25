package database

import (
	"database/sql"
)

type (
	Organisation struct {
		Name      string `gorm:"primaryKey"`
		MainGroup string
		SubGroup  string
		Flair     sql.NullString
		Status    StatusString
		Members   []Account `gorm:"many2many:organisation_member;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
		Admins    []Account `gorm:"many2many:organisation_admins;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
		Accounts  []Account `gorm:"many2many:organisation_account;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
	}
	OrganisationList []Organisation

	StatusString string
)

const (
	Public  StatusString = "public"
	Private StatusString = "private"
	Secret  StatusString = "secret"
	Hidden  StatusString = "hidden"
)

var (
	StatusTranslation = map[StatusString]string{
		Public:  "!placeholder!",
		Private: "!placeholder!",
		Secret:  "!placeholder!",
		Hidden:  "!placeholder!",
	}
	Stati = []StatusString{Public, Private, Secret, Hidden}
)

func (t StatusString) String() string {
	return string(t)
}
