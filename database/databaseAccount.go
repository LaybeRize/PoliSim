package database

import "database/sql"

type (
	RoleLevel   int
	AccountList []Account
)

type Account struct {
	ID            int64  `gorm:"primaryKey;autoIncrement:true"`
	DisplayName   string `gorm:"index:unique"`
	Flair         string
	Username      string `gorm:"index:unique"`
	Password      string
	Suspended     bool
	LoginTries    int
	NextLoginTime sql.NullTime
	Role          RoleLevel
	Linked        sql.NullInt64
	Parent        *Account  `gorm:"foreignKey:linked;joinReferences:id"`
	Children      []Account `gorm:"foreignKey:linked"`
}

const (
	PressAccount RoleLevel = iota - 1
	NotLoggedIn
	User
	MediaAdmin
	Admin
	HeadAdmin
)

var Roles = []RoleLevel{PressAccount, User, MediaAdmin, Admin, HeadAdmin}
var RoleTranslation = map[RoleLevel]string{
	PressAccount: "Presse-Account",
	User:         "Nutzer",
	MediaAdmin:   "Medien-Administrator",
	Admin:        "Administrator",
	HeadAdmin:    "Leitender Administrator",
}
