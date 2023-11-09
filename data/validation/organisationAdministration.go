package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type OrganisationModification struct {
	Name      string   `input:"name"`
	MainGroup string   `input:"mainGroup"`
	SubGroup  string   `input:"subGroup"`
	Status    string   `input:"status"`
	Flair     string   `input:"flair"`
	User      []string `input:"user"`
	Admins    []string `input:"admins"`
}

const maxOrganisationLength = 200

func (form *OrganisationModification) CreateOrganisation() (validate Message) {
	validate = Message{Positive: false}
	switch false {
	case isValidString(form.Name, maxOrganisationLength):
		// has no valid name
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongOrganisationName"], maxOrganisationLength)
		return
	case isValidString(form.MainGroup, maxGroupNameLength):
		// has no valid main group
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongMainGroupName"], maxGroupNameLength)
		return
	case isValidString(form.SubGroup, maxGroupNameLength):
		// has no valid subgroup
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongSubGroupName"], maxGroupNameLength)
		return
	case isOrgStatusValid(form.Status):
		// has no valid status
		validate.Message = builder.Translation["invalidOrganisationStatus"]
		return
	case len([]rune(form.Flair)) <= maxFlairLength:
		// has no valid flair
		validate.Message = fmt.Sprintf(builder.Translation["tooLongOrganisationFlair"], maxFlairLength)
		return
	}
	helper.RemoveEntriesFromList(&form.Admins, form.User)
	user, ok, err := extraction.DoAccountsExist(form.User)
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	var admins *database.AccountList
	admins, ok, err = extraction.DoAccountsExist(form.Admins)
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	accounts, err := extraction.GetParentAccounts(append(form.User, form.Admins...))
	if err != nil {
		validate.Message = fmt.Sprintf(builder.Translation["parentAccountError"], err.Error())
		return
	}

	org := database.Organisation{
		Name:      form.Name,
		MainGroup: form.MainGroup,
		SubGroup:  form.SubGroup,
		Flair:     sql.NullString{String: form.Flair, Valid: form.Flair != ""},
		Status:    database.StatusString(form.Status),
		Members:   *user,
		Admins:    *admins,
		Accounts:  *accounts,
	}
	err = extraction.CreateNewOrganisation(&org)
	if err != nil {
		validate.Message = builder.Translation["databaseErrorOrganisationCreation"]
		return
	}

	err = updateFlairs([]string{}, append(form.User, form.Admins...), "", org.Flair.String)
	if err != nil {
		return Message{
			Message: builder.Translation["successfullyCreatedOrganisation"] + "\n" +
				builder.Translation["errorWithFlairUpdate"],
			Positive: true,
		}
	}

	return Message{Message: builder.Translation["successfullyCreatedOrganisation"],
		Positive: true}
}

func (form *OrganisationModification) SearchOrganisation() (validate Message) {
	validate = Message{Positive: false}
	org, err := extraction.GetOrganisation(form.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		validate.Message = builder.Translation["organisationNotFound"]
		return
	} else if err != nil {
		validate.Message = builder.Translation["databaseErrorOrganisationSearch"]
		return
	}

	form.Flair = org.Flair.String
	form.Status = string(org.Status)
	form.MainGroup = org.MainGroup
	form.SubGroup = org.SubGroup
	form.User = make([]string, len(org.Members))
	for i, acc := range org.Members {
		form.User[i] = acc.DisplayName
	}
	form.Admins = make([]string, len(org.Admins))
	for i, acc := range org.Admins {
		form.Admins[i] = acc.DisplayName
	}

	return Message{
		Message:  builder.Translation["successfullyFoundOrganisation"],
		Positive: true,
	}
}

func (form *OrganisationModification) ModifyOrganisation() (validate Message) {
	validate = Message{Positive: false}
	org, err := extraction.GetOrganisation(form.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		validate.Message = builder.Translation["organisationDoesNotExist"]
		return
	} else if err != nil {
		validate.Message = builder.Translation["databaseErrorOrganisationModification"]
		return
	}
	switch false {
	case isValidString(form.MainGroup, maxGroupNameLength):
		// has no valid main group
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongMainGroupName"], maxGroupNameLength)
		return
	case isValidString(form.SubGroup, maxGroupNameLength):
		// has no valid subgroup
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongSubGroupName"], maxGroupNameLength)
		return
	case isOrgStatusValid(form.Status):
		// has no valid status
		validate.Message = builder.Translation["invalidOrganisationStatus"]
		return
	case len([]rune(form.Flair)) <= maxFlairLength:
		// has no valid flair
		validate.Message = fmt.Sprintf(builder.Translation["tooLongOrganisationFlair"], maxFlairLength)
		return
	}
	//clear hidden organisation
	if form.Status == string(database.Hidden) {
		form.Admins = []string{}
		form.User = []string{}
	}

	helper.RemoveEntriesFromList(&form.Admins, form.User)
	user, ok, err := extraction.DoAccountsExist(form.User)
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	var admins *database.AccountList
	admins, ok, err = extraction.DoAccountsExist(form.Admins)
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	accounts, err := extraction.GetParentAccounts(append(form.User, form.Admins...))
	if err != nil {
		validate.Message = builder.Translation["parentAccountError"]
		return
	}

	oldNames := make([]string, len(org.Members)+len(org.Admins))
	for i, acc := range org.Members {
		oldNames[i] = acc.DisplayName
	}
	for i, acc := range org.Admins {
		oldNames[i+len(org.Members)] = acc.DisplayName
	}

	org.MainGroup = form.MainGroup
	org.SubGroup = form.SubGroup
	oldFlair := org.Flair.String
	org.Flair = sql.NullString{String: form.Flair, Valid: form.Flair != ""}
	org.Status = database.StatusString(form.Status)
	org.Members = *user
	org.Admins = *admins
	org.Accounts = *accounts

	err = extraction.ModifiyOrganisation(org)
	if err != nil {
		validate.Message = builder.Translation["databaseErrorOrganisationModification"]
		return
	}

	err = updateFlairs(oldNames, append(form.User, form.Admins...), oldFlair, org.Flair.String)
	if err != nil {
		return Message{
			Message: builder.Translation["successfullyModifiedOrganisation"] + "\n" +
				builder.Translation["errorWithFlairUpdate"],
			Positive: true,
		}
	}

	return Message{Message: builder.Translation["successfullyModifiedOrganisation"],
		Positive: true}
}
