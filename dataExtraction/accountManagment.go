package dataExtraction

import (
	"PoliSim/database"
	"database/sql"
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
	Role           database.RoleString
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
	Role           database.RoleString
}

func GetAccoutForLogin(username string) (*AccountLogin, error) {
	accoutLogin := &AccountLogin{}
	err := database.DB.Model(database.Account{}).Where("username=?", username).First(accoutLogin).Error
	return accoutLogin, err
}

func (acc *AccountLogin) SaveBack() error {
	return database.DB.Model(database.Account{}).Save(acc).Error
}

func GetAccountForAuth(token string) (*AccountAuth, error) {
	accountAuth := &AccountAuth{}
	err := database.DB.Model(database.Account{}).Where("refresh_token=?", token).First(accountAuth).Error
	return accountAuth, err
}

func UpdateAuthToken(id int64, newToken string, newExperationDate sql.NullTime) error {
	return database.DB.Model(database.Account{}).Where("id=?", id).Update("refresh_token", newToken).Update("expiration_date", newExperationDate).Error
}

func GetAllChildrenDisplayNames(parentID int64) (*AccountDisplayNameList, error) {
	array := &AccountDisplayNameList{}
	err := database.DB.Model(database.Account{}).Select("display_name").Where("linked=?", parentID).Order("display_name").Find(array).Error
	return array, err
}
