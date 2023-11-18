package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/html/builder"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"time"
)

type AddTag struct {
	TagText  string `input:"tag"`
	TagColor string `input:"color"`
}

func (form *AddTag) AddTagToDocument(doc *database.Document) (validate Message) {
	validate = Message{Positive: false}
	var IsColor = regexp.MustCompile(`^#[a-fA-F0-9]{6}$`).MatchString
	switch false {
	case isValidString(form.TagText, maxDocumentInfoTagLength):
		// has no valid content
		validate.Message = fmt.Sprintf(builder.Translation["missingTagTextForDocument"], maxDocumentInfoTagLength)
		return
	case IsColor(form.TagColor):
		//tag color doesn't fit the format anyway
		validate.Message = builder.Translation["invalidHexColor"]
		return
	}

	doc.Info.Post = append([]database.Posts{{
		UUID:      uuid.New().String(),
		Hidden:    false,
		Submitted: time.Now(),
		Info:      form.TagText,
		Color:     form.TagColor,
	}}, doc.Info.Post...)

	err := extraction.UpdateDocument(doc)
	if err != nil {
		validate.Message = builder.Translation["errorAddingTag"]
		return
	}

	return Message{Positive: true}
}

func FlipTagHidden(tagUUID string, doc *database.Document) (validate Message) {
	validate = Message{Positive: false}

	exists := false
	for i, post := range doc.Info.Post {
		if post.UUID == tagUUID {
			doc.Info.Post[i].Hidden = !doc.Info.Post[i].Hidden
			exists = true
		}
	}
	if !exists {
		validate.Message = builder.Translation["tagDoesNotExist"]
		return
	}

	err := extraction.UpdateDocument(doc)
	if err != nil {
		validate.Message = builder.Translation["errorChangingTag"]
		return
	}

	return Message{Positive: true}
}

type AddComment struct {
	Content string `input:"content"`
	Account string `input:"authorAccount"`
}

const (
	maxDocumentCommentLength = 5_000
)

func (form *AddComment) AddComment(uuidStr string, acc *extraction.AccountAuth) (validate Message) {
	validate = Message{Positive: false}
	account, ok, err := IsAccountValidForUser(acc.ID, form.Account)
	switch false {
	case isValidString(form.Content, maxDocumentCommentLength):
		// has no valid title
		validate.Message = fmt.Sprintf(builder.Translation["missingContentForDiscussionComment"], maxDocumentCommentLength)
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

	err = extraction.GetDocumentIfCanParticipate(uuidStr, account.ID)
	if err != nil {
		validate.Message = builder.Translation["notAllowedToComment"]
		return
	}

	err = logic.AddComment(account.DisplayName, account.Flair, form.Content, uuidStr)
	if err != nil {
		validate.Message = builder.Translation["errorWhenMakingComment"]
		return
	}

	return Message{
		Message:  builder.Translation["addedComment"],
		Positive: true,
	}
}

type (
	CastVote interface {
		CastVote(acc *extraction.AccountAuth, isAdmin bool, docUUID string, voteUUID string, voteType database.VoteType) (validate Message)
	}

	AddSingleVote struct {
		InvalidateVote bool   `json:"invalidateVote"`
		Answer         string `json:"answer"`
		Account        string `json:"authorAccount"`
	}
	AddMultipleVote struct {
		InvalidateVote bool     `json:"invalidateVote"`
		Answer         []string `json:"answer"`
		Account        string   `json:"authorAccount"`
	}
	AddRankedVote struct {
		InvalidateVote bool     `json:"invalidateVote"`
		Answers        []string `json:"answer"`
		Account        string   `json:"authorAccount"`
	}
	AddThreeChoice struct {
		InvalidateVote bool     `json:"invalidateVote"`
		Answers        []string `json:"answer"`
		Account        string   `json:"authorAccount"`
	}
)

func (form *AddSingleVote) CastVote(acc *extraction.AccountAuth, isAdmin bool, docUUID string, voteUUID string, voteType database.VoteType) (validate Message) {
	return Message{}
}
func (form *AddMultipleVote) CastVote(acc *extraction.AccountAuth, isAdmin bool, docUUID string, voteUUID string, voteType database.VoteType) (validate Message) {
	return Message{}
}
func (form *AddRankedVote) CastVote(acc *extraction.AccountAuth, isAdmin bool, docUUID string, voteUUID string, voteType database.VoteType) (validate Message) {
	return Message{}
}
func (form *AddThreeChoice) CastVote(acc *extraction.AccountAuth, isAdmin bool, docUUID string, voteUUID string, voteType database.VoteType) (validate Message) {
	return Message{}
}
