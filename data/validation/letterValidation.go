package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type CreateLetter struct {
	Title          string   `input:"title"`
	Content        string   `input:"content"`
	Account        string   `input:"authorAccount"`
	Flair          string   `input:"flair"`
	AllHaveToAgree bool     `input:"allHaveToAgree"`
	NoSigning      bool     `input:"noSigning"`
	Reader         []string `input:"reader"`
}

const (
	maxLetterTitleLength   = 100
	maxLetterContentLength = 10_000
)

func (form *CreateLetter) CreateNormalLetter(requestAccountID int64) (validate Message) {
	validate = Message{Positive: false}
	account, ok, err := isAccountValidForUser(requestAccountID, form.Account)
	switch false {
	case isValidString(form.Title, maxLetterTitleLength):
		// has no valid title
		validate.Message = fmt.Sprintf(builder.Translation["missingTitleForLetter"], maxPressTitleLength)
		return
	case isValidString(form.Content, maxLetterContentLength):
		// has no valid content
		validate.Message = fmt.Sprintf(builder.Translation["missingContentForLetter"], MaxPressContentLength)
		return
	case err == nil:
		// error with author account
		validate.Message = builder.Translation["databaseErrorWithAuthorAccount"]
		return
	case ok:
		// not allowed for author account
		validate.Message = builder.Translation["notAllowedToUseAccount"]
		return
	}
	helper.RemoveFirstStringOccurrenceFromArray(&form.Reader, account.DisplayName)
	var reader *database.AccountList
	reader, ok, err = extraction.DoAccountsExist(form.Reader)
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	if len(*reader) == 0 {
		validate.Message = builder.Translation["noReaderForLetter"]
		return
	}
	letter := database.Letter{
		UUID:        uuid.New().String(),
		Written:     time.Now(),
		Author:      account.DisplayName,
		Flair:       account.Flair,
		Title:       form.Title,
		Content:     form.Content,
		HTMLContent: helper.CreateHTML(form.Content),
		Info: database.LetterInfo{
			AllHaveToAgree:     form.AllHaveToAgree,
			NoSigning:          form.NoSigning,
			PeopleNotYetSigned: form.Reader,
			Signed:             []string{account.DisplayName},
			Rejected:           []string{},
		},
		Viewer:     append(*reader, database.Account{ID: account.ID, DisplayName: account.DisplayName}),
		Removed:    false,
		ModMessage: false,
	}
	err = extraction.CreateLetter(&letter)
	if err != nil {
		validate.Message = builder.Translation["errorCreatingLetter"]
		return
	}

	form.Title = ""
	form.Content = ""
	return Message{
		Message:  builder.Translation["successfullyCreatedLetter"],
		Positive: true,
	}
}
