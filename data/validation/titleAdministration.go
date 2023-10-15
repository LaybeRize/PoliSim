package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"database/sql"
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

var maxTitleLength = 200
var maxGroupNameLength = 200
var maxFlairLength = 20

func (form *TitleModification) CreateTitle() (validate Message) {
	validate = Message{Positive: false}
	switch false {
	case isValidString(form.Name, maxTitleLength):
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
	helper.ClearStringArray(&form.Holder)
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
	if err == gorm.ErrRecordNotFound {
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

func (form *TitleModification) ModifyTitle() Message {
	return Message{}
}

func (form *TitleModification) DeleteTitle() Message {
	return Message{}
}
