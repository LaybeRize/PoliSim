package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CreateArticle struct {
	Title          string `input:"title"`
	Subtitle       string `input:"subtitle"`
	Content        string `input:"content"`
	IsBreakingNews bool   `input:"breakingNews"`
	Account        string `input:"authorAccount"`
}

const (
	maxPressTitleLength    = 150
	maxPressContentLength  = 20_000
	maxPressSubtitleLength = 300
)

func (form *CreateArticle) CreateArticle(requestAccountID int64) (validate Message) {
	validate = Message{Positive: false}
	account, ok, err := IsAccountValidForUser(requestAccountID, form.Account)
	switch false {
	case isValidString(form.Title, maxPressTitleLength):
		// has no valid title
		validate.Message = fmt.Sprintf(builder.Translation["missingTitleForPress"], maxPressTitleLength)
		return
	case len([]rune(form.Subtitle)) < maxPressSubtitleLength:
		// has no valid subtitle
		validate.Message = fmt.Sprintf(builder.Translation["tooLongSubtitleForPress"], maxPressSubtitleLength)
		return
	case isValidString(form.Content, maxPressContentLength):
		// has no valid content
		validate.Message = fmt.Sprintf(builder.Translation["missingContentForPress"], maxPressContentLength)
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
	article := database.Article{
		UUID:        uuid.New().String(),
		Publication: database.EternatityPublicationName,
		Written:     time.Now(),
		Author:      account.DisplayName,
		Flair:       account.Flair,
		Headline:    form.Title,
		Subtitle:    sql.NullString{String: form.Subtitle, Valid: form.Subtitle != ""},
		Content:     form.Content,
		HTMLContent: helper.CreateHTML(form.Content),
	}
	if form.IsBreakingNews {
		err = createBreakingNewsPublication(&article)
		if err != nil {
			validate.Message = builder.Translation["databaseErrorArticleCreation"]
			return
		}
	}

	err = extraction.CreateArticle(&article)
	if err != nil {
		validate.Message = builder.Translation["databaseErrorArticleCreation"]
		return
	}

	form.Title = ""
	form.Subtitle = ""
	form.Content = ""
	return Message{
		Message:  builder.Translation["createdArticleSuccessfully"],
		Positive: true,
	}
}

func createBreakingNewsPublication(article *database.Article) error {
	pub := database.Publication{
		UUID:         uuid.New().String(),
		CreateTime:   time.Now(),
		Publicated:   false,
		BreakingNews: true,
	}
	err := extraction.CreatePublication(&pub)
	article.Publication = pub.UUID
	return err
}

var maxCharacterForRejection = 10_000

func RejectArticle(uuidStr string, content string) (validate Message) {
	validate = Message{Positive: false}
	if !isValidString(content, maxCharacterForRejection) {
		validate.Message = fmt.Sprintf(builder.Translation["missingRejectionMessage"], maxCharacterForRejection)
		return
	}
	article, err := extraction.FindHiddenArticle(uuidStr)
	if err != nil {
		validate.Message = builder.Translation["cantRejectArticle"]
		return
	}
	account, err := extraction.GetAccountByDisplayName(article.Author)
	if err != nil {
		validate.Message = builder.Translation["cantFindAuthorOfArticle"]
		return
	}
	params := []any{article.Content, content}
	msg := builder.Translation["rejectionLetterBody"]
	if article.Subtitle.Valid {
		msg = builder.Translation["rejectionLetterBodySubtitle"] + msg
		params = []any{article.Subtitle.String, article.Content, content}
	}
	letter := database.Letter{
		UUID:    uuid.New().String(),
		Written: time.Now(),
		Author:  builder.Translation["authorOfRejections"],
		Flair:   "",
		Title:   fmt.Sprintf(builder.Translation["rejectionLetterTitle"], article.Headline),
		Content: fmt.Sprintf(msg, params...),
		Info: database.LetterInfo{
			AllHaveToAgree:     false,
			NoSigning:          true,
			PeopleNotYetSigned: []string{},
			Signed:             []string{account.DisplayName},
			Rejected:           []string{},
		},
		Viewer:     []database.Account{*account},
		Removed:    false,
		ModMessage: true,
	}
	letter.HTMLContent = helper.CreateHTML(letter.Content)
	err = extraction.CreateLetter(&letter)
	if err != nil {
		validate.Message = builder.Translation["errorCreatingRejectionLetter"]
		return
	}
	err = extraction.DeleteArticle(article)
	if err != nil {
		validate.Message = builder.Translation["errorWhileDeletingArticle"]
		return
	}
	return Message{Positive: true}
}

func PublishNewspaper(uuidStr string) (validate Message, newUUID string) {
	validate = Message{Positive: false}
	pub, err := extraction.FindPublicationAndReturnIt(uuidStr, "false")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		validate.Message = builder.Translation["newspaperUUIDNotValid"]
		return
	} else if err != nil {
		validate.Message = builder.Translation["databaseErrorRetrievingNewspaper"]
		return
	}
	if pub.Publicated {
		validate.Message = builder.Translation["newspaperUUIDNotValid"]
		return
	}

	list, err := extraction.FindArticlesForPublicationUUID(pub.UUID)
	if err != nil {
		validate.Message = builder.Translation["newspaperUUIDNotValid"]
		return
	}

	if len(*list) == 0 {
		validate.Message = builder.Translation["noArticleToPublish"]
		return
	}

	if pub.BreakingNews {
		pub.PublishTime = time.Now()
		pub.Publicated = true
		err = extraction.ChangePublication(pub)
		newUUID = pub.UUID
	} else {
		newPub := database.Publication{
			UUID:         uuid.New().String(),
			CreateTime:   time.Now(),
			PublishTime:  time.Now(),
			Publicated:   true,
			BreakingNews: false,
		}
		err = extraction.CreatePublication(&newPub)
		if err != nil {
			validate.Message = builder.Translation["errorPublishingNewspaper"]
			return
		}
		err = extraction.UpdatePublication(pub.UUID, newPub.UUID)
		newUUID = newPub.UUID
	}

	if err != nil {
		validate.Message = builder.Translation["errorPublishingNewspaper"]
		return
	}

	validate = Message{Message: "", Positive: true}
	return
}
