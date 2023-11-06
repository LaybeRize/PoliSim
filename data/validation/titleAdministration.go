package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/html/builder"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type TitleModification struct {
	Name      string   `input:"name"`
	NewName   string   `input:"newName"`
	MainGroup string   `input:"mainGroup"`
	SubGroup  string   `input:"subGroup"`
	Flair     string   `input:"flair"`
	Holder    []string `input:"holder"`
}

// both below are also used in organisationAdminstration
const (
	maxTitleLength     = 200
	maxGroupNameLength = 200
	maxFlairLength     = 20
)

func (form *TitleModification) CreateTitle() (validate Message) {
	validate = Message{Positive: false}
	switch false {
	case isValidString(form.Name, maxTitleLength):
		// has no valid name
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongTitleName"], maxTitleLength)
		return
	case isValidString(form.MainGroup, maxGroupNameLength):
		// has no valid main group
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongMainGroupName"], maxGroupNameLength)
		return
	case isValidString(form.SubGroup, maxGroupNameLength):
		// has no valid subgroup
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongSubGroupName"], maxGroupNameLength)
		return
	case len([]rune(form.Flair)) <= maxFlairLength:
		// has no valid flair
		validate.Message = fmt.Sprintf(builder.Translation["TooLongTitleFlair"], maxFlairLength)
		return
	}
	accounts, ok, err := extraction.DoAccountsExist(form.Holder)
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}

	title := database.Title{
		Name:      form.Name,
		MainGroup: form.MainGroup,
		SubGroup:  form.SubGroup,
		Flair:     sql.NullString{String: form.Flair, Valid: form.Flair != ""},
		Holder:    *accounts,
	}

	err = extraction.CreateTitle(&title)
	if err != nil {
		validate.Message = builder.Translation["errorWhileCreatingTitle"]
		return
	}

	err = updateFlairs([]string{}, form.Holder, "", form.Flair)
	if err != nil {
		return Message{
			Message: builder.Translation["successfullyCreatedTitle"] + "\n" +
				builder.Translation["errorWithFlairUpdate"],
			Positive: true,
		}
	}
	return Message{
		Message:  builder.Translation["successfullyCreatedTitle"],
		Positive: true,
	}
}

func (form *TitleModification) SearchTitle() (validate Message) {
	validate = Message{Positive: false}
	title, err := extraction.GetTitle(form.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		validate.Message = builder.Translation["titleNotFound"]
		return
	} else if err != nil {
		validate.Message = builder.Translation["databaseErrorTitleSearch"]
		return
	}
	form.Name = title.Name
	form.NewName = title.Name
	form.MainGroup = title.MainGroup
	form.SubGroup = title.SubGroup
	form.Flair = title.Flair.String
	form.Holder = make([]string, len(title.Holder))
	for i, acc := range title.Holder {
		form.Holder[i] = acc.DisplayName
	}
	return Message{
		Message:  builder.Translation["successfullyFoundTitle"],
		Positive: true,
	}
}

func (form *TitleModification) ModifyTitle() (validate Message) {
	validate = Message{Positive: false}
	title, err := extraction.GetTitle(form.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		validate.Message = builder.Translation["titleNotFoundForModification"]
		return
	}
	switch false {
	case err == nil:
		//error with database access
		validate.Message = builder.Translation["errorWhileAccessingDatabaseForTitle"]
		return
	case isValidString(form.NewName, maxTitleLength):
		// has no valid name
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongTitleName"], maxNameLength)
		return
	case isValidString(form.MainGroup, maxGroupNameLength):
		// has no valid main group
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongMainGroupName"], maxGroupNameLength)
		return
	case isValidString(form.SubGroup, maxGroupNameLength):
		// has no valid subgroup
		validate.Message = fmt.Sprintf(builder.Translation["missingOrTooLongSubGroupName"], maxGroupNameLength)
		return
	case len([]rune(form.Flair)) <= maxFlairLength:
		// has no valid flair
		validate.Message = fmt.Sprintf(builder.Translation["TooLongTitleFlair"], maxFlairLength)
		return
	}
	accounts, ok, err := extraction.DoAccountsExist(form.Holder)
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	old := make([]string, len(title.Holder))
	oldFlair := title.Flair.String
	for i, acc := range title.Holder {
		old[i] = acc.DisplayName
	}
	title.Name = form.NewName
	title.MainGroup = form.MainGroup
	title.SubGroup = form.SubGroup
	title.Flair = sql.NullString{String: form.Flair, Valid: form.Flair != ""}
	title.Holder = *accounts
	err = extraction.UpdateTitle(title, form.Name)
	if err != nil {
		validate.Message = builder.Translation["errorWhileAccessingDatabaseForTitle"]
		return
	}
	err = updateFlairs(old, form.Holder, oldFlair, form.Flair)
	if err != nil {
		return Message{
			Message: builder.Translation["successfullyModifiedTitle"] + "\n" +
				builder.Translation["errorWithFlairUpdate"],
			Positive: true,
		}
	}
	return Message{
		Message:  builder.Translation["successfullyModifiedTitle"],
		Positive: true,
	}
}

func (form *TitleModification) DeleteTitle() (validate Message) {
	validate = Message{Positive: false}
	title, err := extraction.GetTitle(form.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		validate.Message = builder.Translation["titleNotFoundForDeletion"]
		return
	} else if err != nil {
		validate.Message = builder.Translation["databaseErrorTitleDeletion"]
		return
	}
	old := make([]string, len(title.Holder))
	oldFlair := title.Flair.String
	for i, acc := range title.Holder {
		old[i] = acc.DisplayName
	}
	err = extraction.DeleteTitle(title)
	if err != nil {
		validate.Message = builder.Translation["databaseErrorTitleDeletion"]
		return
	}
	err = updateFlairs(old, []string{}, oldFlair, "")
	if err != nil {
		return Message{
			Message: builder.Translation["successfullyDeletedTitle"] + "\n" +
				builder.Translation["errorWithFlairUpdate"],
			Positive: true,
		}
	}
	return Message{
		Message:  builder.Translation["successfullyDeletedTitle"],
		Positive: true,
	}
}
