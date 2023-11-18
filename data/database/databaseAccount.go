package database

import (
	"database/sql"
	"strconv"
)

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

// Roles has to keep PressAccount because of its special status as the first element of the array
var (
	Roles           = []RoleLevel{PressAccount, User, MediaAdmin, Admin, HeadAdmin}
	RoleTranslation = map[RoleLevel]string{
		PressAccount: "!placeholder!",
		User:         "!placeholder!",
		MediaAdmin:   "!placeholder!",
		Admin:        "!placeholder!",
		HeadAdmin:    "!placeholder!",
	}
)

func (t RoleLevel) String() string {
	return strconv.Itoa(int(t))
}
