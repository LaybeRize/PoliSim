package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"time"
)

type CreateDocument struct {
	Title        string `input:"title"`
	Subtitle     string `input:"subtitle"`
	Content      string `input:"content"`
	TagText      string `input:"tag"`
	TagColor     string `input:"color"`
	Account      string `input:"authorAccount"`
	Organisation string `input:"organisation"`
	UUIDredirect string
}

type CreateDiscussion struct {
	Title             string   `input:"title"`
	Subtitle          string   `input:"subtitle"`
	Content           string   `input:"content"`
	Private           bool     `input:"private"`
	MembersCanComment bool     `input:"membersCanComment"`
	AnyoneCanComment  bool     `input:"anyoneCanComment"`
	Account           string   `input:"authorAccount"`
	Organisation      string   `input:"organisation"`
	Reader            []string `input:"reader"`
	Writer            []string `input:"writer"`
	UUIDredirect      string
}

const (
	maxDocumentTitleLength    = 200
	maxDocumentSubtitleLength = 400
	maxDocumentContentLength  = 100_000
	maxDocumentInfoTagLength  = 200
)

func (form *CreateDocument) CreateDocument(requestAccountID int64) (validate Message) {
	validate = Message{Positive: false}
	account, ok, err := IsAccountValidForUser(requestAccountID, form.Account)
	var IsColor = regexp.MustCompile(`^#[a-fA-F0-9]{6}$`).MatchString
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
	case isValidString(form.TagText, maxDocumentInfoTagLength):
		// has no valid content
		validate.Message = fmt.Sprintf(builder.Translation["missingTagTextForDocument"], maxDocumentInfoTagLength)
		return
	case err == nil:
		// error with author account
		validate.Message = builder.Translation["databaseErrorWithAuthorAccount"]
		return
	case ok:
		// not allowed for author account
		validate.Message = builder.Translation["notAllowedToUseAccount"]
		return
	case IsColor(form.TagColor):
		//tag color doesn't fit the format anyway
		validate.Message = builder.Translation["invalidHexColor"]
		return
	}
	var org *database.Organisation
	org, ok, err = IsOrganisationValidForAccount(account.ID, form.Organisation)
	if err != nil {
		validate.Message = builder.Translation["databaseErrorWithOrganisationAccount"]
		return
	}
	if !ok {
		validate.Message = builder.Translation["notAnAdminOfOrganisation"]
		return
	}

	document := database.Document{
		UUID:         uuid.New().String(),
		Written:      time.Now(),
		Organisation: org.Name,
		Type:         database.LegislativeText,
		Author:       account.DisplayName,
		Flair:        account.Flair,
		Title:        form.Title,
		Subtitle: sql.NullString{
			String: form.Subtitle,
			Valid:  form.Subtitle != "",
		},
		HTMLContent:    helper.CreateHTML(form.Content),
		CurrentPostTag: form.TagText,
		Info: database.DocumentInfo{
			Finishing: time.Time{},
			Post: []database.Posts{{
				UUID:      uuid.New().String(),
				Submitted: time.Now(),
				Info:      form.TagText,
				Color:     form.TagColor,
			}},
		},
	}

	err = extraction.CreateDocument(&document)
	if err != nil {
		validate.Message = builder.Translation["errorCreatingDocument"]
		return
	}

	form.UUIDredirect = document.UUID
	return Message{Positive: true}
}

func (form *CreateDiscussion) CreateDiscussion(requestAccountID int64) (validate Message) {
	validate = Message{Positive: false}
	return
}
