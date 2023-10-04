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

func GetAccoutForLogin(username string) (*AccountLogin, error) {
	accoutLogin := &AccountLogin{}
	err := database.DB.Model(database.Account{}).Where("username=?", username).First(accoutLogin).Error
	return accoutLogin, err
}

func (acc *AccountLogin) SaveBack() error {
	return database.DB.Model(database.Account{}).Where("id=?", acc.ID).Updates(acc).Error
}

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
