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

func (form *AccountModification) RequestAccount() (validate ValidationMessage) {
	validate = ValidationMessage{Positive: false}
	var err error
	acc := &dataExtraction.AccountModification{}
	if form.SearchByUsername {
		acc, err = dataExtraction.GetAccountModificationByUsername(form.Username)
	} else {
		acc, err = dataExtraction.GetAccountModificationByDisplayName(form.DisplayName)
	}
	if err != nil {
		validate.Message = componentHelper.Translation["accountDoesNotExists"]
		return
	}
	form.Username = acc.Username
	form.DisplayName = acc.DisplayName
	form.Flair = acc.Flair
	form.Suspended = acc.Suspended
	form.Role = int(acc.Role)
	form.Linked = acc.Linked.Int64

	validate.Message = componentHelper.Translation["accountFound"]
	validate.Positive = true
	return
}

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
		validate.Message = componentHelper.Translation["hashingPasswordError"]
		return
	}
	// empty password to make the account not login-able
	if form.Role == int(database.PressAccount) {
		pass = []byte("")
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
		validate.Message = componentHelper.Translation["creatingUserError"]
		return
	}
	form.Username, form.DisplayName, form.Password, form.Flair = "", "", "", ""

	validate.Positive = true
	validate.Message = componentHelper.Translation["userSuccessfullyCreated"]
	return
}

func (form *AccountModification) ValidateAccountModification(changer *dataExtraction.AccountAuth) (validate ValidationMessage) {
	var acc *dataExtraction.AccountModification
	var err error
	if form.SearchByUsername {
		acc, err = dataExtraction.GetAccountModificationByUsername(form.Username)
	} else {
		acc, err = dataExtraction.GetAccountModificationByDisplayName(form.DisplayName)
	}
	if err != nil {
		validate.Message = componentHelper.Translation["accountDoesNotExists"]
		return
	}

	switch true {
	case acc.ID == 1:
		return form.validateChangeRootAccount(acc, changer)
	case acc.Role == database.HeadAdmin:
		return form.validateChangeHeadAdmin(acc, changer)
	case acc.Role == database.PressAccount:
		return form.validateChangeToPressAccount(acc)
	default:
		return form.validateChangeToEveryoneElse(acc, changer)
	}
}

func (form *AccountModification) validateChangeRootAccount(acc *dataExtraction.AccountModification, changer *dataExtraction.AccountAuth) (validate ValidationMessage) {
	validate = ValidationMessage{Positive: false}
	if changer.ID != 1 {
		validate.Message = componentHelper.Translation["notAllowedToChangeAccount"]
		return
	}
	if !form.ChangeFlair {
		validate.Message = componentHelper.Translation["noChangesMade"]
		return
	}
	acc.Flair = form.Flair
	err := acc.OnlyUpdateFlair()
	if err != nil {
		validate.Message = componentHelper.Translation["changingUserError"]
		return
	}
	form.Username = acc.Username
	form.DisplayName = acc.DisplayName
	form.Suspended = acc.Suspended
	form.Role = int(acc.Role)
	validate.Message = componentHelper.Translation["changingUserSuccessfully"]
	validate.Positive = true
	return
}

func (form *AccountModification) validateChangeHeadAdmin(acc *dataExtraction.AccountModification, changer *dataExtraction.AccountAuth) (validate ValidationMessage) {
	var err error
	validate = ValidationMessage{Positive: false}
	if changer.ID != 1 && changer.ID != acc.ID {
		// Not allowed to change this account
		validate.Message = componentHelper.Translation["notAllowedToChangeAccount"]
		return
	}
	if changer.ID != 1 {
		// only allow changes to the flair for self
		if !form.ChangeFlair {
			validate.Message = componentHelper.Translation["noChangesMade"]
			return
		}
		acc.Flair = form.Flair
		err = acc.OnlyUpdateFlair()
		if err != nil {
			validate.Message = componentHelper.Translation["changingUserError"]
			return
		}
		form.Username = acc.Username
		form.DisplayName = acc.DisplayName
		form.Suspended = acc.Suspended
		form.Role = int(acc.Role)
		validate.Message = componentHelper.Translation["changingUserSuccessfully"]
		validate.Positive = true
		return
	}
	//validate Role
	if !isRoleValid(form.Role) || form.Role == int(database.PressAccount) {
		validate.Message = componentHelper.Translation["roleNotAllowed"]
		return
	}
	acc.Suspended = form.Suspended
	acc.Role = database.RoleLevel(form.Role)
	if form.ChangeFlair {
		acc.Flair = form.Flair
		err = acc.UpdateAllFields()
	} else {
		err = acc.UpdateEverythingExceptFlair()
	}
	if err != nil {
		validate.Message = componentHelper.Translation["changingUserError"]
		return
	}
	validate.Message = componentHelper.Translation["changingUserSuccessfully"]
	validate.Positive = true
	return
}

func (form *AccountModification) validateChangeToPressAccount(acc *dataExtraction.AccountModification) (validate ValidationMessage) {
	acc.Suspended = form.Suspended
	acc.Linked.Int64 = form.Linked
	acc.Linked.Valid = true
	form.Role = int(acc.Role)
	if form.Linked <= 0 {
		acc.Linked.Valid = false
		form.Linked = 0
	}
	var err error
	if form.ChangeFlair {
		acc.Flair = form.Flair
		err = acc.UpdateAllFields()
	} else {
		err = acc.UpdateEverythingExceptFlair()
	}
	// reset linked and try again, maybe the error was that the key was invalid
	if err != nil {
		acc.Linked.Valid = false
		form.Linked = 0
		if form.ChangeFlair {
			acc.Flair = form.Flair
			err = acc.UpdateAllFields()
		} else {
			err = acc.UpdateEverythingExceptFlair()
		}
	}

	// if there is still an error then we return an error message
	if err != nil {
		validate.Message = componentHelper.Translation["changingUserError"]
		return
	}
	validate.Message = componentHelper.Translation["changingUserSuccessfully"]
	validate.Positive = true
	return
}

func (form *AccountModification) validateChangeToEveryoneElse(acc *dataExtraction.AccountModification, changer *dataExtraction.AccountAuth) (validate ValidationMessage) {
	validate = ValidationMessage{Positive: false}
	if changer.ID != 1 && form.Role == int(database.HeadAdmin) {
		// Not allowed
		validate.Message = componentHelper.Translation["notAllowedToChangeAccount"]
		return
	}
	if !isRoleValid(form.Role) || form.Role == int(database.PressAccount) {
		validate.Message = componentHelper.Translation["roleNotAllowed"]
		return
	}

	acc.Suspended = form.Suspended
	acc.Role = database.RoleLevel(form.Role)

	var err error
	if form.ChangeFlair {
		acc.Flair = form.Flair
		err = acc.UpdateAllFields()
	} else {
		err = acc.UpdateEverythingExceptFlair()
	}
	if err != nil {
		validate.Message = componentHelper.Translation["changingUserError"]
		return
	}

	validate.Message = componentHelper.Translation["changingUserSuccessfully"]
	validate.Positive = true
	return
}
