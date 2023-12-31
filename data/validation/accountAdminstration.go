package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"sync"
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

const maxNameLength = 100

func (form *AccountModification) RequestAccount() (validate Message) {
	validate = Message{Positive: false}
	var err error
	acc := &database.Account{}
	if form.SearchByUsername {
		acc, err = extraction.GetAccountByUsername(form.Username)
	} else {
		acc, err = extraction.GetAccountByDisplayName(form.DisplayName)
	}
	if err != nil {
		validate.Message = builder.Translation["accountDoesNotExists"]
		return
	}
	form.Username = acc.Username
	form.DisplayName = acc.DisplayName
	form.Flair = acc.Flair
	form.Suspended = acc.Suspended
	form.Role = int(acc.Role)
	form.Linked = acc.Linked.Int64

	validate.Message = builder.Translation["accountFound"]
	validate.Positive = true
	return
}

func (form *AccountModification) ValidateAccountCreation(changer *database.AccountAuth) (validate Message) {
	validate = Message{Positive: false}
	switch false {
	case isValidString(form.DisplayName, maxNameLength):
		// has no valid display name
		validate.Message = fmt.Sprintf(builder.Translation["missingDisplayName"], maxNameLength)
		return
	case isRoleValid(form.Role):
		// has no valid Role
		validate.Message = builder.Translation["roleNotAllowed"]
		return
	case changer.ID == 1 || form.Role != int(database.HeadAdmin):
		// non-root account tries to create head admin
		validate.Message = fmt.Sprintf(builder.Translation["cantCreateHeadAdmin"], database.RoleTranslation[database.HeadAdmin])
		return
	case form.Role == int(database.PressAccount) || isValidString(form.Username, maxNameLength):
		// non-pressaccount misses username
		validate.Message = fmt.Sprintf(builder.Translation["missingUsernameForNonPressAccounts"], maxNameLength)
		return
	case form.Role == int(database.PressAccount) || isValidString(form.Password, -1):
		// non-pressaccount is missing password
		validate.Message = fmt.Sprintf(builder.Translation["missingPasswordForNonPressAccounts"])
		return
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		validate.Message = builder.Translation["hashingPasswordError"]
		return
	}
	// empty password to make the account not login-able
	if form.Role == int(database.PressAccount) {
		pass = []byte("")
	}
	acc := database.Account{
		DisplayName: form.DisplayName,
		Suspended:   false,
		Role:        database.RoleLevel(form.Role),
		Username:    form.Username,
		Password:    string(pass),
		Flair:       form.Flair,
		Linked:      sql.NullInt64{Int64: form.Linked},
	}
	if acc.Role == database.PressAccount {
		acc.Linked.Valid = true
		acc.Username = acc.DisplayName
	}
	err = extraction.CreateAccount(&acc)
	if err != nil {
		validate.Message = builder.Translation["creatingUserError"]
		return
	}
	form.Username, form.DisplayName, form.Password, form.Flair = "", "", "", ""

	validate.Positive = true
	validate.Message = builder.Translation["userSuccessfullyCreated"]
	return
}

func (form *AccountModification) ValidateAccountModification(changer *database.AccountAuth) (validate Message) {
	var acc *database.Account
	var err error
	if form.SearchByUsername {
		acc, err = extraction.GetAccountByUsername(form.Username)
	} else {
		acc, err = extraction.GetAccountByDisplayName(form.DisplayName)
	}
	if err != nil {
		validate.Message = builder.Translation["accountDoesNotExists"]
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

func (form *AccountModification) validateChangeRootAccount(acc *database.Account, changer *database.AccountAuth) (validate Message) {
	validate = Message{Positive: false}
	if changer.ID != 1 {
		validate.Message = builder.Translation["notAllowedToChangeAccount"]
		return
	}
	if !form.ChangeFlair {
		validate.Message = builder.Translation["noChangesMade"]
		return
	}
	acc.Flair = form.Flair
	err := extraction.OnlyUpdateFlair(acc)
	if err != nil {
		validate.Message = builder.Translation["changingUserError"]
		return
	}
	form.Username = acc.Username
	form.DisplayName = acc.DisplayName
	form.Suspended = acc.Suspended
	form.Role = int(acc.Role)
	validate.Message = builder.Translation["changingUserSuccessfully"]
	validate.Positive = true
	return
}

func (form *AccountModification) validateChangeHeadAdmin(acc *database.Account, changer *database.AccountAuth) (validate Message) {
	var err error
	validate = Message{Positive: false}
	if changer.ID != 1 && changer.ID != acc.ID {
		// Not allowed to change this account
		validate.Message = builder.Translation["notAllowedToChangeAccount"]
		return
	}
	if changer.ID != 1 {
		// only allow changes to the flair for self
		if !form.ChangeFlair {
			validate.Message = builder.Translation["noChangesMade"]
			return
		}
		acc.Flair = form.Flair
		err = extraction.OnlyUpdateFlair(acc)
		if err != nil {
			validate.Message = builder.Translation["changingUserError"]
			return
		}
		form.Username = acc.Username
		form.DisplayName = acc.DisplayName
		form.Suspended = acc.Suspended
		form.Role = int(acc.Role)
		validate.Message = builder.Translation["changingUserSuccessfully"]
		validate.Positive = true
		return
	}
	//validate Role
	if !isRoleValid(form.Role) || form.Role == int(database.PressAccount) {
		validate.Message = builder.Translation["roleNotAllowed"]
		return
	}
	acc.Suspended = form.Suspended
	acc.Role = database.RoleLevel(form.Role)
	if form.ChangeFlair {
		acc.Flair = form.Flair
		err = extraction.UpdateAllFields(acc)
	} else {
		err = extraction.UpdateEverythingExceptFlair(acc)
	}
	if err != nil {
		validate.Message = builder.Translation["changingUserError"]
		return
	}
	validate.Message = builder.Translation["changingUserSuccessfully"]
	validate.Positive = true
	return
}

func (form *AccountModification) validateChangeToPressAccount(acc *database.Account) (validate Message) {
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
		err = extraction.UpdateAllFields(acc)
	} else {
		err = extraction.UpdateEverythingExceptFlair(acc)
	}
	// reset linked and try again, maybe the error was that the key was invalid
	if err != nil {
		acc.Linked.Valid = false
		form.Linked = 0
		if form.ChangeFlair {
			acc.Flair = form.Flair
			err = extraction.UpdateAllFields(acc)
		} else {
			err = extraction.UpdateEverythingExceptFlair(acc)
		}
	}

	// if there is still an error then we return an error message
	if err != nil {
		validate.Message = builder.Translation["changingUserError"]
		return
	}
	validate.Message = builder.Translation["changingUserSuccessfully"]
	err = logic.UpdateOrganisationAccount(acc.ID)
	if err != nil {
		validate.Message += "\n" + builder.Translation["errorForChangesToOrganisationWithPressAccount"]
	}
	validate.Positive = true
	return
}

func (form *AccountModification) validateChangeToEveryoneElse(acc *database.Account, changer *database.AccountAuth) (validate Message) {
	validate = Message{Positive: false}
	if changer.ID != 1 && form.Role == int(database.HeadAdmin) {
		// Not allowed
		validate.Message = builder.Translation["notAllowedToChangeAccount"]
		return
	}
	if !isRoleValid(form.Role) || form.Role == int(database.PressAccount) {
		validate.Message = builder.Translation["roleNotAllowed"]
		return
	}

	acc.Suspended = form.Suspended
	acc.Role = database.RoleLevel(form.Role)

	var err error
	if form.ChangeFlair {
		acc.Flair = form.Flair
		err = extraction.UpdateAllFields(acc)
	} else {
		err = extraction.UpdateEverythingExceptFlair(acc)
	}
	if err != nil {
		validate.Message = builder.Translation["changingUserError"]
		return
	}

	validate.Message = builder.Translation["changingUserSuccessfully"]
	validate.Positive = true
	return
}

var updateFlair = sync.Mutex{}

func updateFlairs(old []string, new []string, oldFlair string, newFlair string) error {
	updateFlair.Lock()
	defer updateFlair.Unlock()
	switch true {
	case oldFlair == "" && newFlair == "":
	//do nothing
	case oldFlair != "" && oldFlair == newFlair:
		return onlyUpdateUser(old, new, oldFlair)
	case oldFlair != "" && newFlair != "":
		return removeFromAllAndAddToNew(old, new, oldFlair, newFlair)
	case oldFlair != "" && newFlair == "":
		return justRemoveFlair(old, oldFlair)
	case oldFlair == "" && newFlair != "":
		return justAddFlair(new, newFlair)
	}
	return nil
}

func justAddFlair(new []string, flair string) error {
	toAdd, err := extraction.GetFlairAccountList(new)
	if err != nil {
		return err
	}
	addFlair(toAdd, flair)
	return extraction.UpdateFlairs(toAdd)
}

func justRemoveFlair(old []string, flair string) error {
	toRemove, err := extraction.GetFlairAccountList(old)
	if err != nil {
		return err
	}
	removeFlair(toRemove, flair)
	return extraction.UpdateFlairs(toRemove)
}

func removeFromAllAndAddToNew(old []string, new []string, oldFlair string, newFlair string) error {
	remove, change, add, err := extraction.GetDifferentAccountGroups(old, new)
	if err != nil {
		return err
	}
	addFlair(add, newFlair)
	err = extraction.UpdateFlairs(add)
	if err != nil {
		return err
	}
	removeFlair(remove, oldFlair)
	err = extraction.UpdateFlairs(remove)
	if err != nil {
		return err
	}
	removeFlair(change, oldFlair)
	addFlair(change, newFlair)
	return extraction.UpdateFlairs(change)
}

func onlyUpdateUser(old []string, new []string, flair string) error {
	remove, _, add, err := extraction.GetDifferentAccountGroups(old, new)
	if err != nil {
		return err
	}
	addFlair(add, flair)
	err = extraction.UpdateFlairs(add)
	if err != nil {
		return err
	}
	removeFlair(remove, flair)
	return extraction.UpdateFlairs(remove)
}

const (
	flairSeperator     = ","
	longFlairSeperator = flairSeperator + " "
)

func addFlair(update *database.AccountList, flair string) {
	for i, acc := range *update {
		if acc.Flair == "" {
			(*update)[i].Flair = flair
		} else {
			(*update)[i].Flair += longFlairSeperator + flair
		}
	}
}

func removeFlair(update *database.AccountList, flair string) {
	for i, acc := range *update {
		switch true {
		case acc.Flair == flair:
			(*update)[i].Flair = ""
		case strings.HasPrefix(acc.Flair, flair+flairSeperator):
			(*update)[i].Flair = helper.TrimPrefix(acc.Flair, flair+longFlairSeperator)
		case strings.Contains(acc.Flair, longFlairSeperator+flair+flairSeperator):
			var re = regexp.MustCompile(`(?m)` + regexp.QuoteMeta(longFlairSeperator+flair+flairSeperator))
			(*update)[i].Flair = re.ReplaceAllString(acc.Flair, flairSeperator)
		case strings.HasSuffix(acc.Flair, longFlairSeperator+flair):
			(*update)[i].Flair = helper.TrimSuffix(acc.Flair, longFlairSeperator+flair)
		}
	}
}

type ChangePassword struct {
	OldPassword      string `input:"ordPassword"`
	NewPassword      string `input:"newPassword"`
	NewPasswordAgain string `input:"newPasswordAgain"`
}

const (
	maxPasswordLength = 50
	minPasswordLength = 10
)

func (form *ChangePassword) ChangePassword(acc *database.AccountAuth) (validate Message) {
	validate.Positive = false
	account, err := extraction.GetAccountByID(acc.ID)
	if err != nil {
		validate.Message = builder.Translation["couldNotFindAccount"]
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(form.OldPassword))
	if err != nil {
		validate.Message = builder.Translation["oldPasswordIsWrong"]
		return
	}
	length := len([]rune(form.NewPassword))
	if length > maxPasswordLength || length < minPasswordLength {
		validate.Message = fmt.Sprintf(builder.Translation["newPasswordWrongLength"], minPasswordLength, maxPasswordLength)
		return
	}
	if form.NewPassword != form.NewPasswordAgain {
		validate.Message = builder.Translation["newPasswordAgainNotTheSame"]
		return
	}
	var pass []byte
	pass, err = bcrypt.GenerateFromPassword([]byte(form.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		validate.Message = builder.Translation["errorGeneratingNewPassword"]
		return
	}
	err = extraction.ChangePassword(account.ID, string(pass))
	if err != nil {
		validate.Message = builder.Translation["errorSavingNewPassword"]
		return
	}
	return Message{Positive: true, Message: builder.Translation["successfullyChangePassword"]}
}
