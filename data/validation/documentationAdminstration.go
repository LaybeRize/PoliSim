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
	EndTime           string   `input:"endTime"`
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
	account, ok, err := IsAccountValidForUser(requestAccountID, form.Account)
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

	org, isAdmin, err := IsOrganisationValidForAccount(account.ID, form.Organisation)
	if err != nil {
		validate.Message = builder.Translation["databaseErrorWithOrganisationAccount"]
		return
	}

	//TODO add needed cases for admins and other stuff
	if org.Status == database.Secret {
		if ((len(form.Reader) != 0 || len(form.Writer) != 0) && !isAdmin) || form.AnyoneCanComment {
			validate.Message = builder.Translation["noExternalReaderOrWriterAllowed"]
			return
		}
		if !form.Private {
			validate.Message = builder.Translation["errorBecauseNotPrivate"]
			return
		}
	} else if org.Status == database.Private {
		if form.Private && form.AnyoneCanComment {
			validate.Message = builder.Translation["mutuallyExlusiveSelection"]
			return
		}
	} else if org.Status == database.Public {
		if form.Private {
			validate.Message = builder.Translation["notAllowedToBePrivate"]
			return
		}
	}

	var endDiscussion time.Time
	endDiscussion, err = time.ParseInLocation("2006-01-02T15:04", form.EndTime, time.Local)
	if err != nil {
		validate.Message = builder.Translation["timeIsInvalidString"]
		return
	}

	helper.RemoveEntriesFromList(&form.Reader, form.Writer)
	reader, ok, err := extraction.DoAccountsExist(form.Reader)
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	writer, ok, err := extraction.DoAccountsExist(form.Writer)
	if !ok {
		validate.Message = fmt.Sprintf(builder.Translation["nameCouldNotBeFound"], err.Error())
		return
	}
	accounts, err := extraction.GetParentAccounts(append(form.Reader, form.Writer...))
	if err != nil {
		validate.Message = builder.Translation["parentAccountError"]
		return
	}

	written := time.Now()
	discType := database.RunningDiscussion
	if len(form.Reader) == 0 && len(form.Writer) == 0 &&
		!form.AnyoneCanComment && !form.MembersCanComment {
		endDiscussion = written
		discType = database.FinishedDiscussion
	}

	document := database.Document{
		UUID:         uuid.New().String(),
		Written:      written,
		Organisation: org.Name,
		Type:         discType,
		Author:       account.DisplayName,
		Flair:        account.Flair,
		Title:        form.Title,
		Subtitle: sql.NullString{
			String: form.Subtitle,
			Valid:  form.Subtitle != "",
		},
		HTMLContent:               helper.CreateHTML(form.Content),
		AnyPosterAllowed:          form.AnyoneCanComment,
		OrganisationPosterAllowed: form.MembersCanComment,
		Info: database.DocumentInfo{
			Finishing:  endDiscussion,
			Discussion: []database.Discussions{},
		},
		Viewer:  *accounts,
		Poster:  *writer,
		Allowed: *reader,
	}

	err = extraction.CreateDocument(&document)
	if err != nil {
		validate.Message = builder.Translation["errorCreatingDocument"]
		return
	}

	form.UUIDredirect = document.UUID
	return Message{Positive: true}
}
