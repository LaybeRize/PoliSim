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

type AccountNames struct {
	DisplayName string
	Username    string
}

type AccountAuth struct {
	ID          int64
	DisplayName string
	Suspended   bool
	Role        database.RoleLevel
}

type AccountLogin struct {
	ID            int64
	DisplayName   string
	Username      string
	Password      string
	Suspended     bool
	LoginTries    int
	NextLoginTime sql.NullTime
	Role          database.RoleLevel
}

type AccountModification struct {
	ID          int64
	DisplayName string
	Username    string
	Password    string
	Flair       string
	Suspended   bool
	Role        database.RoleLevel
	Linked      sql.NullInt64
}

// RootAccountExists checks if the account with the ID 1 exists
// if so, it returns true and nil. If the DB gives back a gorm.ErrRecordNotFound
// it will return false and nil.
// If any other error occurs it will return false and the error.
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

// GetAccountForLogin returns a filled *AccountLogin on successfully locating
// an account with that username. Otherwise, return an empty struct and the error.
func GetAccountForLogin(username string) (*AccountLogin, error) {
	accountLogin := &AccountLogin{}
	err := database.DB.Model(database.Account{}).Where("username=?", username).First(accountLogin).Error
	return accountLogin, err
}

// SaveBack updates the given account in the DB by id.
func (acc *AccountLogin) SaveBack() error {
	return database.DB.Model(database.Account{}).Where("id=?", acc.ID).Updates(acc).Error
}

// CreateMe creates a new account with all given rows. Every other row
// specified in database.Account will be filled with the standard struct value for that field.
func (acc *AccountLogin) CreateMe() error {
	return database.DB.Create(&database.Account{
		ID:            acc.ID,
		DisplayName:   acc.DisplayName,
		Username:      acc.Username,
		Password:      acc.Password,
		Suspended:     acc.Suspended,
		LoginTries:    acc.LoginTries,
		NextLoginTime: acc.NextLoginTime,
		Role:          acc.Role,
	}).Error
}

// GetAccountForAuth returns a filled *AccountAuth on successfully locating
// an account with that id. Otherwise, return an empty struct and the error.
func GetAccountForAuth(id int64) (*AccountAuth, error) {
	accountAuth := &AccountAuth{}
	err := database.DB.Model(database.Account{}).Where("id=?", id).First(accountAuth).Error
	return accountAuth, err
}

// GetAllChildrenDisplayNames returns a *AccountDisplayNameList on successfully finding any children.
// Will return an empty array and a gorm.ErrRecordNotFound on no children found and any other error if there was one.
func GetAllChildrenDisplayNames(parentID int64) (*AccountDisplayNameList, error) {
	array := &AccountDisplayNameList{}
	err := database.DB.Model(database.Account{}).Select("display_name").Where("linked=?", parentID).Order("display_name").Find(array).Error
	return array, err
}

func (acc *AccountModification) CreateMe() error {
	return database.DB.Create(&database.Account{
		DisplayName: acc.DisplayName,
		Username:    acc.Username,
		Password:    acc.Password,
		Flair:       acc.Flair,
		Suspended:   acc.Suspended,
		Role:        acc.Role,
		Linked:      acc.Linked,
	}).Error
}

func GetAccountModificationByUsername(username string) (*AccountModification, error) {
	acc := &AccountModification{}
	err := database.DB.Model(database.Account{}).Where("username=?", username).First(acc).Error
	return acc, err
}

func GetAccountModificationByDisplayName(displayName string) (*AccountModification, error) {
	acc := &AccountModification{}
	err := database.DB.Model(database.Account{}).Where("display_name=?", displayName).First(acc).Error
	return acc, err
}

func (acc *AccountModification) OnlyUpdateFlair() error {
	return database.DB.Model(&database.Account{ID: acc.ID}).Update("flair", acc.Flair).Error
}

// UpdateAllFields updates all allowed fields (flair, suspended, role, linked)
func (acc *AccountModification) UpdateAllFields() error {
	return database.DB.Model(&database.Account{ID: acc.ID}).Updates(map[string]interface{}{
		"flair":     acc.Flair,
		"suspended": acc.Suspended,
		"role":      acc.Role,
		"linked":    acc.Linked,
	}).Error
}

func (acc *AccountModification) UpdateEverythingExceptFlair() error {
	return database.DB.Model(&database.Account{ID: acc.ID}).Updates(map[string]interface{}{
		"suspended": acc.Suspended,
		"role":      acc.Role,
		"linked":    acc.Linked,
	}).Error
}

// ReturnNames returns as the first argument the DisplayNames and as the second Argument the Usernames
func ReturnNames() ([]string, []string, error) {
	rows, err := database.DB.Model(&database.Account{}).Rows()
	defer rows.Close()
	var names = make([]string, 0, 20)
	var users = make([]string, 0, 20)
	if err != nil {
		return names, users, err
	}
	for rows.Next() {
		var user AccountNames
		err = database.DB.ScanRows(rows, &user)
		if err != nil {
			return names, users, err
		}

		names = append(names, user.DisplayName)
		users = append(users, user.Username)
	}
	return names, users, nil
}
