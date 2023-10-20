package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"database/sql"
	"fmt"
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

var maxOrganisationLength = 200

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
		//handel error
	}

	return Message{}
}

func (form *OrganisationModification) SearchOrganisation() (validate Message) {
	return Message{}
}

func (form *OrganisationModification) ModifyOrganisation() (validate Message) {
	return Message{}
}
