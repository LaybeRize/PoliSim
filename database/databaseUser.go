package database

import "database/sql"

type (
	RoleString  string
	AccountList []Account
)

type Account struct {
	ID             int64  `gorm:"primaryKey;autoIncrement:true"`
	DisplayName    string `gorm:"index:unique"`
	Flair          string
	Username       string `gorm:"index:unique"`
	Password       string
	Suspended      bool
	RefreshToken   sql.NullString
	ExpirationDate sql.NullTime
	LoginTries     int
	NextLoginTime  sql.NullTime
	Role           RoleString
	Linked         sql.NullInt64
	Parent         *Account  `gorm:"foreignKey:linked;joinReferences:id"`
	Children       []Account `gorm:"foreignKey:linked"`
}

const (
	User         RoleString = "user"
	MediaAdmin   RoleString = "media_admin"
	Admin        RoleString = "admin"
	HeadAdmin    RoleString = "head_admin"
	PressAccount RoleString = "press_account"
	NotLoggedIn  RoleString = "notLoggedIn"
)

var Roles = []RoleString{PressAccount, User, MediaAdmin, Admin, HeadAdmin}
var RoleTranslation = map[RoleString]string{
	PressAccount: "Presse-Account",
	User:         "Nutzer",
	MediaAdmin:   "Medien-Administrator",
	Admin:        "Administrator",
	HeadAdmin:    "Leitender Administrator",
}
