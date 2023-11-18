package extraction

import (
	"PoliSim/data/database"
	"PoliSim/helper"
	"database/sql"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type (
	AccountDisplayNameList []AccountDisplayName
	AccountDisplayName     struct {
		ID          int64
		DisplayName string
	}

	AccountFlairUpdateList []AccountFlairUpdate
	AccountFlairUpdate     struct {
		ID    int64
		Flair string
	}
)

type AccountNames struct {
	DisplayName string
	Username    string
}

type (
	AccountList        []AccountListElement
	AccountListElement struct {
		ID          int64
		DisplayName string
		Username    string
		Flair       string
		Suspended   bool
		Role        database.RoleLevel
		Linked      sql.NullInt64
	}
)

type AccountAuth struct {
	ID          int64
	DisplayName string
	Suspended   bool
	Role        database.RoleLevel
	Session     *sessions.Session `gorm:"-"`
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
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

// GetAccountForPasswordChange returns a filled *AccountLogin on successfully locating
// an account with that id. Otherwise, return an empty struct and the error.
func GetAccountForPasswordChange(id int64) (*AccountLogin, error) {
	accountLogin := &AccountLogin{}
	err := database.DB.Model(database.Account{}).Where("id=?", id).First(accountLogin).Error
	return accountLogin, err
}

func ChangePassword(id int64, newPassword string) error {
	return database.DB.Model(&database.Account{ID: id}).Updates(map[string]interface{}{
		"password": newPassword,
	}).Error
}

// SaveBack updates the given account in the DB by id.
func (acc *AccountLogin) SaveBack() error {
	return database.DB.Model(database.Account{}).Where("id=?", acc.ID).Updates(acc).Error
}

// CreateMe creates a new account with all given rows. Every other row
// specified in database.Account will be filled with the standard struct value for that field.
// Only use this for the root account creation. Otherwise, use the ID agnostic creation with AccountModification.
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
	err := database.DB.Model(database.Account{}).Select("display_name").Where("linked=? AND suspended = false", parentID).Order("display_name").Find(array).Error
	return array, err
}

// CreateMe creates the account in the database based on the display name, username, password, flair, role and linked value
// theoretically also sets the suspended value, but creating a suspended accounts seems dumb.
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

// GetAccountModificationByUsername returns an AccountModification pointer for the given username
// or an error.
func GetAccountModificationByUsername(username string) (*AccountModification, error) {
	acc := &AccountModification{}
	err := database.DB.Model(database.Account{}).Where("username=?", username).First(acc).Error
	return acc, err
}

// GetAccountModificationByDisplayName returns an AccountModification Pointer to the user with the displayname
// or an error.
func GetAccountModificationByDisplayName(displayName string) (*AccountModification, error) {
	acc := &AccountModification{}
	err := database.DB.Model(database.Account{}).Where("display_name=?", displayName).First(acc).Error
	return acc, err
}

// OnlyUpdateFlair updates the given account flair by id
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

// UpdateEverythingExceptFlair updates suspended, role and linked
func (acc *AccountModification) UpdateEverythingExceptFlair() error {
	return database.DB.Model(&database.Account{ID: acc.ID}).Updates(map[string]interface{}{
		"suspended": acc.Suspended,
		"role":      acc.Role,
		"linked":    acc.Linked,
	}).Error
}

// ReturnNames returns as the first argument the DisplayNames and as the second Argument the Usernames
func ReturnNames() ([]string, []string, error) {
	rows, err := database.DB.Model(&database.Account{}).Select("display_name, username").Rows()
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

// ReturnListOfDisplayNames returns all not suspended accounts display names
func ReturnListOfDisplayNames() ([]string, error) {
	rows, err := database.DB.Model(&database.Account{}).Select("display_name").Where("suspended = false").Rows()
	defer rows.Close()
	var names = make([]string, 0, 20)
	if err != nil {
		return names, err
	}
	for rows.Next() {
		var user AccountNames
		err = database.DB.ScanRows(rows, &user)
		if err != nil {
			return names, err
		}

		names = append(names, user.DisplayName)
	}
	return names, nil
}

// ReturnAccountList returns itself and all their press accounts
func ReturnAccountList(id int64) (AccountList, error) {
	array := AccountList{}
	err := database.DB.Model(database.Account{}).Where("id=? OR linked=? OR ?=0", id, id, id).Order("id").Find(&array).Error
	return array, err
}

func DoAccountsExist(displayNames []string) (accountList *database.AccountList, b bool, err error) {
	accountList = &database.AccountList{}

	err = database.DB.Model(database.Account{}).Select("id, display_name").Where("display_name = ANY($1) AND suspended = false", pq.StringArray(displayNames)).Order("display_name").Find(&accountList).Error
	if len(displayNames) == len(*accountList) {
		b = true
		return
	}

	for _, item := range *accountList {
		helper.RemoveFirstStringOccurrenceFromArray(&displayNames, item.DisplayName)
	}

	accountList = &database.AccountList{}
	b = false
	err = errors.New(displayNames[0])
	return
}

// GetDifferentAccountGroups returns three arrays. The first containing only the old accounts, the second containing all accounts in both groups, and the last containing only the new accounts.
// If a query error arises, it gets returned too.
func GetDifferentAccountGroups(old []string, new []string) (onlyOld *database.AccountList, onBoth *database.AccountList, onlyNew *database.AccountList, err error) {
	err = database.DB.Select("id, display_name, flair").Where("display_name = ANY($1) AND NOT (display_name = ANY($2))", pq.StringArray(old), pq.StringArray(new)).Order("display_name").Find(&onlyOld).Error
	if err != nil {
		return
	}
	err = database.DB.Select("id, display_name, flair").Where("display_name = ANY($1) AND display_name = ANY($2)", pq.StringArray(old), pq.StringArray(new)).Order("display_name").Find(&onBoth).Error
	if err != nil {
		return
	}
	err = database.DB.Select("id, display_name, flair").Where("NOT (display_name = ANY($1)) AND display_name = ANY($2)", pq.StringArray(old), pq.StringArray(new)).Order("display_name").Find(&onlyNew).Error
	return
}

// GetFlairAccountList returns the flair account list, for the queried display names or an error if one occurs.
func GetFlairAccountList(accounts []string) (accList *database.AccountList, err error) {
	err = database.DB.Model(database.Account{}).Select("id, display_name, flair").Where("display_name = ANY($1)", pq.StringArray(accounts)).Order("display_name").Find(&accList).Error
	return
}

func UpdateFlairs(acc *database.AccountList) error {
	for _, singleAcc := range *acc {
		err := database.DB.Model(&singleAcc).Update("flair", singleAcc.Flair).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func GetParentAccounts(names []string) (accounts *database.AccountList, err error) {
	err = database.DB.Select("id, linked").Where("display_name = ANY($1)", pq.StringArray(names)).Order("display_name").Find(&accounts).Error
	if err != nil {
		return
	}
	for i, acc := range *accounts {
		if acc.Linked.Valid {
			(*accounts)[i] = database.Account{ID: acc.Linked.Int64}
		}
	}
	return
}
