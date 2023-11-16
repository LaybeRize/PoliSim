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

type (
	CreateDocument struct {
		BaseDocumentInfo
		TagText  string `input:"tag"`
		TagColor string `input:"color"`
	}
	CreateDiscussion struct {
		BaseDocumentInfo
		PrivateDocumentInfo
	}
	CreateVote struct {
		BaseDocumentInfo
		PrivateDocumentInfo
		Questions []*Question `json:"question"`
	}
	Question struct {
		Text                   string   `json:"questionText"`
		Answers                []string `json:"answers"`
		QuestionType           string   `json:"type"`
		ViewCountsWhileRunning bool     `json:"viewCountsWhileRunning"`
		ViewNamesWhileRunning  bool     `json:"viewNamesWhileRunning"`
		ViewNamesAfterFinished bool     `json:"viewNamesAfterFinished"`
	}
)

const (
	maxDocumentTitleLength    = 200
	maxDocumentSubtitleLength = 400
	maxDocumentContentLength  = 100_000
	maxDocumentInfoTagLength  = 200
	minDays                   = 1
	maxDays                   = 14
	maxQuestions              = 10
	maxVoteQuestionLength     = 200
	maxVoteAnswers            = 50
)

func (form *CreateDocument) CreateDocument(requestAccountID int64) (validate Message) {
	validate = Message{Positive: false}
	var account *extraction.AccountModification
	var org *database.Organisation
	var isAdmin bool
	var IsColor = regexp.MustCompile(`^#[a-fA-F0-9]{6}$`).MatchString
	switch false {
	case form.BaseDocumentInfo.validateBaseDocumentInformation(requestAccountID, account, &validate):
		return
	case isValidString(form.TagText, maxDocumentInfoTagLength):
		// has no valid content
		validate.Message = fmt.Sprintf(builder.Translation["missingTagTextForDocument"], maxDocumentInfoTagLength)
		return
	case IsColor(form.TagColor):
		//tag color doesn't fit the format anyway
		validate.Message = builder.Translation["invalidHexColor"]
		return
	case form.BaseDocumentInfo.validateOrganisation(account, org, &isAdmin, &validate):
		return
	case isAdmin:
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

	err := extraction.CreateDocument(&document)
	if err != nil {
		validate.Message = builder.Translation["errorCreatingDocument"]
		return
	}

	form.UUIDredirect = document.UUID
	return Message{Positive: true}
}

func (form *CreateDiscussion) CreateDiscussion(requestAccountID int64) (validate Message) {
	validate = Message{Positive: false}
	var account *extraction.AccountModification
	var org *database.Organisation
	var isAdmin bool
	var endDiscussion time.Time
	var reader, writer, accounts *database.AccountList
	switch false {
	case form.BaseDocumentInfo.validateBaseDocumentInformation(requestAccountID, account, &validate):
		return
	case form.BaseDocumentInfo.validateOrganisation(account, org, &isAdmin, &validate):
		return
	case form.PrivateDocumentInfo.validateTime(&endDiscussion, &validate, builder.Translation["timeIsInvalidString"]):
		return
	case !endDiscussion.Before(time.Now().Add(24 * time.Hour * minDays)):
		validate.Message = fmt.Sprintf(builder.Translation["timeUnderMinAmountDays"], minDays)
		return
	case !endDiscussion.After(time.Now().Add(24 * time.Hour * maxDays)):
		validate.Message = fmt.Sprintf(builder.Translation["timeOverMaxAmountDays"], maxDays)
		return
	case form.PrivateDocumentInfo.validateAccounts(reader, writer, accounts, &validate):
		return
	}

	//TODO add needed cases for admins and other stuff
	switch org.Status {
	case database.Secret:
		if ((len(form.Onlooker) != 0 || len(form.Participants) != 0) && !isAdmin) || form.AnyoneCanParticipate {
			validate.Message = builder.Translation["noExternalReaderOrWriterAllowed"]
			return
		} else if !form.Private {
			validate.Message = builder.Translation["errorBecauseNotPrivate"]
			return
		}
	case database.Private:
		if form.Private && form.AnyoneCanParticipate {
			validate.Message = builder.Translation["mutuallyExlusiveSelection"]
			return
		}
	case database.Public:
		if form.Private {
			validate.Message = builder.Translation["notAllowedToBePrivate"]
			return
		}
	}
	if form.AnyoneCanParticipate && !isAdmin {
		validate.Message = builder.Translation["needsToBeAdminForAnyone"]
		return
	}

	written := time.Now()
	discType := database.RunningDiscussion
	if len(form.Participants) == 0 && !form.AnyoneCanParticipate && !form.MembersCanParticipate {
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
		Private:                   form.Private,
		AnyPosterAllowed:          form.AnyoneCanParticipate,
		OrganisationPosterAllowed: form.MembersCanParticipate,
		Info: database.DocumentInfo{
			Finishing:  endDiscussion,
			Discussion: []database.Discussions{},
		},
		Viewer:  *reader,
		Poster:  *writer,
		Allowed: *accounts,
	}

	err := extraction.CreateDocument(&document)
	if err != nil {
		validate.Message = builder.Translation["errorCreatingDocument"]
		return
	}

	form.UUIDredirect = document.UUID
	return Message{Positive: true}
}

func (form *CreateVote) CreateVote(requestAccountID int64) (validate Message) {
	validate = Message{Positive: false}
	var account *extraction.AccountModification
	var org *database.Organisation
	var isAdmin bool
	var endVote time.Time
	var spectator, voter, accounts *database.AccountList
	switch false {
	case form.BaseDocumentInfo.validateBaseDocumentInformation(requestAccountID, account, &validate):
		return
	case form.BaseDocumentInfo.validateOrganisation(account, org, &isAdmin, &validate):
		return
	case form.PrivateDocumentInfo.validateTime(&endVote, &validate, builder.Translation["timeVoteEndIsInvalidString"]):
		return
	case !endVote.Before(time.Now().Add(24 * time.Hour * minDays)):
		validate.Message = fmt.Sprintf(builder.Translation["timeForVoteUnderMinAmountDays"], minDays)
		return
	case !endVote.After(time.Now().Add(24 * time.Hour * maxDays)):
		validate.Message = fmt.Sprintf(builder.Translation["timeForVoteOverMaxAmountDays"], maxDays)
		return
	case form.PrivateDocumentInfo.validateAccounts(spectator, voter, accounts, &validate):
		return
	case checkQuestions(&form.Questions, &validate):
		return
	}

	switch org.Status {
	case database.Secret:
		if ((len(form.Onlooker) != 0 || len(form.Participants) != 0) && !isAdmin) || form.AnyoneCanParticipate {
			validate.Message = builder.Translation["noExternalReaderOrWriterAllowedVote"]
			return
		} else if !form.Private {
			validate.Message = builder.Translation["errorBecauseNotPrivateVote"]
			return
		}
	case database.Private:
		if form.Private && form.AnyoneCanParticipate {
			validate.Message = builder.Translation["mutuallyExlusiveSelectionVote"]
			return
		}
	case database.Public:
		if form.Private {
			validate.Message = builder.Translation["notAllowedToBePrivateVote"]
			return
		}
	}
	if form.AnyoneCanParticipate && !isAdmin {
		validate.Message = builder.Translation["needsToBeAdminForAnyoneVote"]
		return
	}

	if !form.AnyoneCanParticipate && !form.MembersCanParticipate && len(form.Participants) == 0 {
		validate.Message = builder.Translation["voteNeedsParticipants"]
		return
	}

	parentID := uuid.New().String()
	for _, item := range form.Questions {
		newVote := &database.Votes{
			UUID:                   uuid.New().String(),
			Parent:                 parentID,
			Question:               item.Text,
			ShowNumbersWhileVoting: item.ViewCountsWhileRunning,
			ShowNamesWhileVoting:   item.ViewNamesWhileRunning,
			ShowNamesAfterVoting:   item.ViewNamesAfterFinished,
			Finished:               false,
			Info: database.VoteInfo{
				Results:    map[string]database.Results{},
				Summary:    database.Summary{},
				VoteMethod: database.VoteType(item.QuestionType),
				Options:    item.Answers,
			},
		}
		err := extraction.CreateVote(newVote)
		if err != nil {
			validate.Message = builder.Translation["errorCreatingDocument"]
			return
		}
	}

	document := database.Document{
		UUID:         parentID,
		Written:      time.Now(),
		Organisation: org.Name,
		Type:         database.RunningVote,
		Author:       account.DisplayName,
		Flair:        account.Flair,
		Title:        form.Title,
		Subtitle: sql.NullString{
			String: form.Subtitle,
			Valid:  form.Subtitle != "",
		},
		HTMLContent:               helper.CreateHTML(form.Content),
		Private:                   form.Private,
		AnyPosterAllowed:          form.AnyoneCanParticipate,
		OrganisationPosterAllowed: form.MembersCanParticipate,
		Info: database.DocumentInfo{
			Finishing: endVote,
		},
		Viewer:  *spectator,
		Poster:  *voter,
		Allowed: *accounts,
	}

	err := extraction.CreateDocument(&document)
	if err != nil {
		validate.Message = builder.Translation["errorCreatingDocument"]
		return
	}

	form.UUIDredirect = parentID
	return Message{Positive: true}
}

func checkIfQuestionIsValid(pos int, item *Question, validate *Message) (result bool) {
	result = false
	switch false {
	case isValidString(item.Text, maxVoteQuestionLength):
		// has no valid question
		validate.Message = fmt.Sprintf(builder.Translation["missingVoteQuestion"], pos, maxVoteQuestionLength)
		return
	case database.VoteTranslation[database.VoteType(item.QuestionType)] != "":
		validate.Message = fmt.Sprintf(builder.Translation["wrongVoteQuestionType"], pos)
		return
	case len(item.Answers) == 0:
		validate.Message = fmt.Sprintf(builder.Translation["needsAnswersForVote"], pos)
		return
	case len(item.Answers) > maxVoteAnswers:
		validate.Message = fmt.Sprintf(builder.Translation["tooManyAnswersForVote"], pos, maxVoteAnswers)
		return
	}
	return true
}

func checkQuestions(questions *[]*Question, validate *Message) (result bool) {
	result = false
	newQuestions := make([]*Question, 0, len(*questions))
	counter := 0
	for pos, item := range *questions {
		if item == nil {
			continue
		}

		if item.Answers == nil {
			validate.Message = fmt.Sprintf(builder.Translation["needsAnswersForVote"], pos)
			return
		}
		helper.ClearStringArray(&item.Answers)
		if !checkIfQuestionIsValid(pos, item, validate) {
			return
		}

		counter++
		newQuestions = append(newQuestions, item)
	}
	if counter == 0 || counter > maxQuestions {
		validate.Message = fmt.Sprintf(builder.Translation["questionLimit"], maxQuestions)
		return
	}
	questions = &newQuestions
	return true
}
