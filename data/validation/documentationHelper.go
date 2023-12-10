package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"fmt"
	"time"
)

type (
	BaseDocumentInfo struct {
		Account      string `input:"authorAccount" json:"authorAccount"`
		Organisation string `input:"organisation"  json:"organisation"`
		Title        string `input:"title"  json:"title"`
		Subtitle     string `input:"subtitle"  json:"subtitle"`
		Content      string `input:"content"  json:"content"`
		UUIDredirect string
	}
	PrivateDocumentInfo struct {
		EndTime               string   `input:"endTime" json:"endTime"`
		Private               bool     `input:"private" json:"private"`
		MembersCanParticipate bool     `input:"membersCanComment" json:"membersCanVote"`
		AnyoneCanParticipate  bool     `input:"anyoneCanComment" json:"anyoneCanVote"`
		Onlooker              []string `input:"reader" json:"attendents"`
		Participants          []string `input:"writer" json:"voter"`
	}
)

func (form *BaseDocumentInfo) validateBaseDocumentInformation(requestAccountID int64, account *database.Account, validate *Message) (result bool) {
	result = false
	temp, ok, err := IsAccountValidForUser(requestAccountID, form.Account)
	*account = *temp
	switch false {
	case isValidString(form.Title, maxDocumentTitleLength):
		// has no valid title
		validate.Message = fmt.Sprintf(builder.Translation["missingTitleForDocument"], maxDocumentTitleLength)
		return
	case len([]rune(form.Subtitle)) <= maxDocumentSubtitleLength:
		// has no valid title
		validate.Message = fmt.Sprintf(builder.Translation["tooLongSubtitleForDocument"], maxDocumentSubtitleLength)
		return
	case isValidString(form.Content, maxDocumentContentLength):
		// has no valid content
		validate.Message = fmt.Sprintf(builder.Translation["missingContentForDocument"], maxDocumentContentLength)
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
	return true
}

func (form *BaseDocumentInfo) validateOrganisation(account *database.Account, org *database.Organisation, isAdmin *bool, validate *Message) (result bool) {
	result = false
	temp, value, err := IsOrganisationValidForAccount(account.ID, form.Organisation)
	*org = *temp
	*isAdmin = value
	if err != nil {
		validate.Message = builder.Translation["databaseErrorWithOrganisationAccount"]
		return
	}
	return true
}

func (form *PrivateDocumentInfo) validateTime(endDiscussion *time.Time, validate *Message, errorString string) (result bool) {
	result = false
	var err error
	*endDiscussion, err = time.ParseInLocation("2006-01-02T15:04", form.EndTime, time.Local)
	if err != nil {
		validate.Message = errorString
		return
	}
	return true
}

func (form *PrivateDocumentInfo) validateAccounts(onlooker *database.AccountList, participants *database.AccountList, accounts *database.AccountList, validate *Message) (result bool) {
	result = false
	helper.RemoveEntriesFromList(&form.Onlooker, form.Participants)
	temp, ok, err := extraction.DoAccountsExist(form.Onlooker)
	*onlooker = *temp
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	temp, ok, err = extraction.DoAccountsExist(form.Participants)
	*participants = *temp
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	temp, err = extraction.GetParentAccounts(append(form.Onlooker, form.Participants...))
	*accounts = *temp
	if err != nil {
		validate.Message = builder.Translation["parentAccountError"]
		return
	}
	return true
}
