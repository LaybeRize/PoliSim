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
	account, ok, err := IsAccountValidForUser(requestAccountID, form.Account)
	switch false {
	case isValidString(form.Title, maxLetterTitleLength):
		// has no valid title
		validate.Message = fmt.Sprintf(builder.Translation["missingTitleForLetter"], maxLetterTitleLength)
		return
	case isValidString(form.Content, maxLetterContentLength):
		// has no valid content
		validate.Message = fmt.Sprintf(builder.Translation["missingContentForLetter"], maxLetterContentLength)
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

const (
	maxAuthorNameLength         = 200
	maxModMailAuthorFlairLength = 200
)

func (form *CreateLetter) CreateModMail() (validate Message) {
	validate = Message{Positive: false}
	switch false {
	case isValidString(form.Title, maxLetterTitleLength):
		// has no valid title
		validate.Message = fmt.Sprintf(builder.Translation["missingTitleForLetter"], maxLetterTitleLength)
		return
	case isValidString(form.Content, maxLetterContentLength):
		// has no valid content
		validate.Message = fmt.Sprintf(builder.Translation["missingContentForLetter"], maxLetterContentLength)
		return
	case isValidString(form.Account, maxAuthorNameLength):
		// has no valid author
		validate.Message = fmt.Sprintf(builder.Translation["missingModMailAuthor"], maxAuthorNameLength)
		return
	case len([]rune(form.Flair)) <= maxModMailAuthorFlairLength:
		//has no valid flair
		validate.Message = fmt.Sprintf(builder.Translation["modMailFlairTooLong"], maxModMailAuthorFlairLength)
		return
	}
	reader, ok, err := extraction.DoAccountsExist(form.Reader)
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
		Author:      form.Account,
		Flair:       form.Flair,
		Title:       form.Title,
		Content:     form.Content,
		HTMLContent: helper.CreateHTML(form.Content),
		Info: database.LetterInfo{
			AllHaveToAgree:     form.AllHaveToAgree,
			NoSigning:          form.NoSigning,
			PeopleNotYetSigned: form.Reader,
			Signed:             []string{},
			Rejected:           []string{},
		},
		Viewer:     *reader,
		Removed:    false,
		ModMessage: true,
	}
	err = extraction.CreateLetter(&letter)
	if err != nil {
		validate.Message = builder.Translation["errorCreatingModMail"]
		return
	}

	form.Title = ""
	form.Content = ""
	return Message{
		Message:  builder.Translation["successfullyCreatedModMail"],
		Positive: true,
	}
}

func SignLetter(acc *database.AccountAuth, letterUUID string, accountSigningName string, action string) (account *extraction.AccountModification, validate Message) {
	validate = Message{Positive: false}
	var ok bool
	var err error
	account, ok, err = IsAccountValidForUser(acc.ID, accountSigningName)
	switch false {
	case err == nil:
		// error with author account
		validate.Message = builder.Translation["databaseErrorWithAuthorAccount"]
		return
	case ok:
		// not allowed for author account
		validate.Message = builder.Translation["notAllowedToUseAccount"]
		return
	case action == "reject" || action == "sign":
		// not allowed to do this action
		validate.Message = builder.Translation["actionNotAllowed"]
		return
	}
	letter, err := extraction.GetLetterByIDOnlyWithAccount(letterUUID, account.ID, false)
	if err != nil {
		validate.Message = builder.Translation["letterCouldNotBeFound"]
		return
	}
	sign := action == "sign"
	if sign {
		helper.RemoveFirstStringOccurrenceFromArray(&letter.Info.PeopleNotYetSigned, account.DisplayName)
		letter.Info.Signed = append(letter.Info.Signed, account.DisplayName)
	} else {
		helper.RemoveFirstStringOccurrenceFromArray(&letter.Info.PeopleNotYetSigned, account.DisplayName)
		letter.Info.Rejected = append(letter.Info.Rejected, account.DisplayName)
	}
	err = extraction.UpdateLetter(letter)
	if err != nil {
		validate.Message = builder.Translation["errorUpdatingLetter"]
		return
	}

	validate = Message{Positive: true}
	if sign {
		validate.Message = builder.Translation["successfullyDidSigningAction"]
	} else {
		validate.Message = builder.Translation["successfullyDidRejectionAction"]
	}
	return
}
