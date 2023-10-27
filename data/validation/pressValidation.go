package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type CreateArticle struct {
	Title          string `input:"title"`
	Subtitle       string `input:"subtitle"`
	Content        string `input:"content"`
	IsBreakingNews bool   `input:"breakingNews"`
	Account        string `input:"authorAccount"`
}

var maxPressTitleLength = 150
var maxPressSubtitleLength = 300
var MaxPressContentLength = 20_000

func (form *CreateArticle) CreateArticle(requestAccountID int64) (validate Message) {
	validate = Message{Positive: false}
	account, ok, err := isAccountValidForUser(requestAccountID, form.Account)
	switch false {
	case isValidString(form.Title, maxPressTitleLength):
		// has no valid title
		validate.Message = fmt.Sprintf(builder.Translation["missingTitleForPress"], maxPressTitleLength)
		return
	case len([]rune(form.Subtitle)) > maxPressSubtitleLength:
		// has no valid subtitle
		validate.Message = fmt.Sprintf(builder.Translation["tooLongSubtitleForPress"], maxPressSubtitleLength)
		return
	case isValidString(form.Content, MaxPressContentLength):
		// has no valid content
		validate.Message = fmt.Sprintf(builder.Translation["missingContentForPress"], MaxPressContentLength)
		return
	case err != nil:
		// error with author account
		validate.Message = builder.Translation["databaseErrorWithAuthorAccount"]
		return
	case !ok:
		// not allowed for author account
		validate.Message = builder.Translation["notAllowedToUseAccount"]
		return
	}
	article := database.Article{
		UUID:        uuid.New().String(),
		Publication: "",
		Written:     time.Now(),
		Author:      account.DisplayName,
		Flair:       account.Flair,
		Headline:    form.Title,
		Subtitle:    sql.NullString{String: form.Subtitle, Valid: form.Subtitle != ""},
		Content:     form.Content,
		HTMLContent: helper.CreateHTML(form.Content),
	}
	err = extraction.CreateArticle(&article)
	return
}
