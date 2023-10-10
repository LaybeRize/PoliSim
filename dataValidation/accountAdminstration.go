package dataValidation

import (
	"PoliSim/componentHelper"
	"PoliSim/dataExtraction"
	"PoliSim/database"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type AccountModification struct {
	SearchByUsername bool   `input:"searchByUsername"`
	DisplayName      string `input:"displayName"`
	Username         string `input:"username"`
	Password         string `input:"password"`
	ChangeFlair      bool   `input:"changeFlair"`
	Flair            string `input:"flair"`
	Suspended        bool   `input:"suspended"`
	Role             int    `input:"role"`
	Linked           int64  `input:"linked"`
}

var maxNameLength = 100
var maxPasswordLength = 50

func (form *AccountModification) ValidateAccountCreation(changer *dataExtraction.AccountAuth) (validate ValidationMessage) {
	validate = ValidationMessage{Positive: false}
	switch false {
	case !isEmptyOrNotInRange(form.DisplayName, maxNameLength):
		// has no display name
		validate.Message = fmt.Sprintf(componentHelper.Translation["missingDisplayName"], maxNameLength)
		return
	case isRoleValid(form.Role):
		// has no valid Role
		validate.Message = componentHelper.Translation["roleNotAllowed"]
		return
	case changer.ID == 1 || form.Role != int(database.HeadAdmin):
		// non-root account tries to create head admin
		validate.Message = fmt.Sprintf(componentHelper.Translation["cantCreateHeadAdmin"], database.RoleTranslation[database.HeadAdmin])
		return
	case form.Role == int(database.PressAccount) || !isEmptyOrNotInRange(form.Username, maxNameLength):
		// non-pressaccount misses username
		validate.Message = fmt.Sprintf(componentHelper.Translation["missingUsernameForNonPressAccounts"], maxNameLength)
		return
	case form.Role == int(database.PressAccount) || !isEmptyOrNotInRange(form.Password, maxPasswordLength):
		// non-pressaccount is missing password
		validate.Message = fmt.Sprintf(componentHelper.Translation["missingPasswordForNonPressAccounts"], maxPasswordLength)
		return
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		// handel error
		validate.Message = componentHelper.Translation["hashingPasswordError"]
		return
	}
	acc := dataExtraction.AccountModification{
		DisplayName: form.DisplayName,
		Username:    form.Username,
		Password:    string(pass),
		Flair:       form.Flair,
		Suspended:   false,
		Role:        database.RoleLevel(form.Role),
		Linked:      sql.NullInt64{Int64: form.Linked},
	}
	if acc.Role == database.PressAccount {
		acc.Linked.Valid = true
		acc.Username = acc.DisplayName
	}
	err = acc.CreateMe()
	if err != nil {
		// handel error
		validate.Message = componentHelper.Translation["creatingUserError"]
		return
	}
	form.Username, form.DisplayName, form.Password, form.Flair = "", "", "", ""

	validate.Positive = true
	validate.Message = componentHelper.Translation["userSuccessfullyCreated"]
	return
}

func (form *AccountModification) ValidateAccountModification(changer *dataExtraction.AccountAuth) {
	var acc *dataExtraction.AccountModification
	var err error
	if form.SearchByUsername {
		acc, err = dataExtraction.GetAccountModificationByUsername(form.Username)
	} else {
		acc, err = dataExtraction.GetAccountModificationByDisplayName(form.DisplayName)
	}
	if err != nil {
		// do error handling
	}

	switch true {
	case acc.ID == 1:
		form.validateChangeRootAccount(acc, changer)
	case acc.Role == database.HeadAdmin:
		form.validateChangeHeadAdmin(acc, changer)
	case acc.Role == database.PressAccount:
		form.validateChangeToPressAccount(acc)
	default:
		form.validateChangeToEveryoneElse(acc, changer)
	}
}

func (form *AccountModification) validateChangeRootAccount(acc *dataExtraction.AccountModification, changer *dataExtraction.AccountAuth) {
	if changer.ID != 1 {
		// not allowed to change root account
	}
	acc.Flair = form.Flair
	err := acc.OnlyUpdateFlair()
	//handel error
	_ = err
}

func (form *AccountModification) validateChangeHeadAdmin(acc *dataExtraction.AccountModification, changer *dataExtraction.AccountAuth) {
	if changer.ID != 1 && changer.ID != acc.ID {
		// Not allowed to change this account
	}
	if changer.ID != 1 {
		//only allowed to change own flair
		acc.Flair = form.Flair
		err := acc.OnlyUpdateFlair()
		//handel error
		_ = err
	}
	//validate Role
}

func (form *AccountModification) validateChangeToPressAccount(acc *dataExtraction.AccountModification) {

}

func (form *AccountModification) validateChangeToEveryoneElse(acc *dataExtraction.AccountModification, changer *dataExtraction.AccountAuth) {
	if changer.ID != 1 && form.Role == int(database.HeadAdmin) {
		//Not allowed
	}
}

func isRoleValid(level int) bool {
	return level >= int(database.PressAccount) && level != int(database.NotLoggedIn) && level <= int(database.HeadAdmin)
}

func isEmptyOrNotInRange(str string, length int) bool {
	return str == "" || len([]rune(str)) > length
}
