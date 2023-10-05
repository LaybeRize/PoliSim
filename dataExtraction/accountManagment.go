package dataExtraction

import (
	"PoliSim/database"
	"database/sql"
	"gorm.io/gorm"
)

type (
	AccountDisplayNameList []AccountDisplayName
	AccountDisplayName     struct {
		DisplayName string
	}
)

type AccountAuth struct {
	ID             int64
	DisplayName    string
	Suspended      bool
	RefreshToken   string
	ExpirationDate sql.NullTime
	Role           database.RoleLevel
}

type AccountLogin struct {
	ID             int64
	DisplayName    string
	Username       string
	Password       string
	Suspended      bool
	RefreshToken   string
	ExpirationDate sql.NullTime
	LoginTries     int
	NextLoginTime  sql.NullTime
	Role           database.RoleLevel
}

// RootAccountExists checks if the account with the ID 1 exists
// if so, it returns true and nil. If the DB gives back a gorm.ErrRecordNotFound
// it will return false and nil.
// If any other error occures it will return false and the error.
func RootAccountExists() (bool, error) {
	err := database.DB.Where("id=?", 1).First(&database.Account{}).Error
	if err == nil {
		return true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return false, err
}

// GetAccoutForLogin returns a filled *AccountLogin on sucessfully locating
// an account with that username. Otherwise, return an empty struct and the error.
func GetAccoutForLogin(username string) (*AccountLogin, error) {
	accoutLogin := &AccountLogin{}
	err := database.DB.Model(database.Account{}).Where("username=?", username).First(accoutLogin).Error
	return accoutLogin, err
}

// SaveBack updates the given account in the DB by id.
func (acc *AccountLogin) SaveBack() error {
	return database.DB.Model(database.Account{}).Where("id=?", acc.ID).Updates(acc).Error
}

// CreateMe creates a new account with all given rows. Every other row
// specified in database.Account will be filled with the standard struct value for that field.
func (acc *AccountLogin) CreateMe() error {
	return database.DB.Create(&database.Account{
		ID:             acc.ID,
		DisplayName:    acc.DisplayName,
		Username:       acc.Username,
		Password:       acc.Password,
		Suspended:      acc.Suspended,
		RefreshToken:   acc.RefreshToken,
		ExpirationDate: acc.ExpirationDate,
		LoginTries:     acc.LoginTries,
		NextLoginTime:  acc.NextLoginTime,
		Role:           acc.Role,
	}).Error
}

// GetAccountForAuth returns a filled *AccountAuth on sucessfully locating
// an account with that username. Otherwise, return an empty struct and the error.
func GetAccountForAuth(token string) (*AccountAuth, error) {
	accountAuth := &AccountAuth{}
	err := database.DB.Model(database.Account{}).Where("refresh_token=?", token).First(accountAuth).Error
	return accountAuth, err
}

// UpdateAuthToken updates the rows refresh_token and experation_date for the account with the given id.
func UpdateAuthToken(id int64, newToken string, newExperationDate sql.NullTime) error {
	return database.DB.Model(database.Account{}).Where("id=?", id).Update("refresh_token", newToken).Update("expiration_date", newExperationDate).Error
}

// GetAllChildrenDisplayNames returns a *AccountDisplayNameList on sucessfully finding any children.
// Will return an empty array and a gorm.ErrRecordNotFound on no children found and any other error if there was one.
func GetAllChildrenDisplayNames(parentID int64) (*AccountDisplayNameList, error) {
	array := &AccountDisplayNameList{}
	err := database.DB.Model(database.Account{}).Select("display_name").Where("linked=?", parentID).Order("display_name").Find(array).Error
	return array, err
}
