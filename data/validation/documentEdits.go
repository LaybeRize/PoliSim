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

func (form *AddComment) AddComment(uuidStr string, acc *database.AccountAuth) (validate Message) {
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

	var doc *database.Document
	doc, err = extraction.GetDocumentIfCanParticipate(uuidStr, account.ID)
	if err != nil || doc.Type == database.FinishedDiscussion {
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
		CastVote(acc *database.AccountAuth, docUUID string, voteUUID string, voteType database.VoteType) Message
	}
	GeneralVote struct {
		InvalidateVote bool   `json:"invalidateVote"`
		Account        string `json:"authorAccount"`
	}
	AddSingleVote struct {
		GeneralVote
		Answer int64 `json:"answerSingle"`
	}
	AddMultipleVote struct {
		GeneralVote
		Answer []bool `json:"answerMultiple"`
	}
	AddRankedVote struct {
		GeneralVote
		Answers []int64 `json:"answerRanked"`
	}
	AddThreeChoice struct {
		GeneralVote
		Answers []int64 `json:"answerThree"`
	}
)

func (form *AddSingleVote) CastVote(acc *database.AccountAuth, docUUID string, voteUUID string, voteType database.VoteType) Message {
	return generelizeCast(acc, docUUID, voteUUID, voteType, form.Account, func(vote *database.Votes) (Message, database.Results) {
		result := database.Results{
			InvalidVote: form.InvalidateVote,
			Votes:       map[string]int64{},
		}
		if result.InvalidVote {
			return Message{Positive: true}, result
		}

		if form.Answer < 0 || form.Answer >= int64(len(vote.Info.Options)) {
			return Message{Message: builder.Translation["voteInvalidWrongParameter"]}, result
		}
		result.Votes[vote.Info.Options[form.Answer]] = 1
		return Message{Positive: true}, result
	})
}

func (form *AddMultipleVote) CastVote(acc *database.AccountAuth, docUUID string, voteUUID string, voteType database.VoteType) Message {
	return generelizeCast(acc, docUUID, voteUUID, voteType, form.Account, func(vote *database.Votes) (Message, database.Results) {
		result := database.Results{
			InvalidVote: form.InvalidateVote,
			Votes:       map[string]int64{},
		}
		if result.InvalidVote {
			return Message{Positive: true}, result
		}

		if len(form.Answer) > len(vote.Info.Options) {
			return Message{Message: builder.Translation["voteInvalidWrongParameter"]}, result
		}
		for i, b := range form.Answer {
			if b {
				result.Votes[vote.Info.Options[i]] = 1
			}
		}
		return Message{Positive: true}, result
	})
}

func (form *AddRankedVote) CastVote(acc *database.AccountAuth, docUUID string, voteUUID string, voteType database.VoteType) Message {
	return generelizeCast(acc, docUUID, voteUUID, voteType, form.Account, func(vote *database.Votes) (Message, database.Results) {
		result := database.Results{
			InvalidVote: form.InvalidateVote,
			Votes:       map[string]int64{},
		}
		if result.InvalidVote {
			return Message{Positive: true}, result
		}

		maxOptions := len(vote.Info.Options)
		if len(form.Answers) > maxOptions {
			return Message{Message: builder.Translation["voteInvalidWrongParameter"]}, result
		}
		alreadyRanked := map[int64]struct{}{}
		for i, b := range form.Answers {
			if b > int64(maxOptions) || b < 0 {
				return Message{Message: builder.Translation["voteInvalidWrongParameter"]}, result
			}
			if _, ok := alreadyRanked[b]; ok && b != 0 {
				return Message{
					Message: fmt.Sprintf(builder.Translation["voteAlreadyAssignedRank"], b),
				}, result
			}
			alreadyRanked[b] = struct{}{}
			result.Votes[vote.Info.Options[i]] = b
		}
		return Message{Positive: true}, result
	})
}

func (form *AddThreeChoice) CastVote(acc *database.AccountAuth, docUUID string, voteUUID string, voteType database.VoteType) Message {
	return generelizeCast(acc, docUUID, voteUUID, voteType, form.Account, func(vote *database.Votes) (Message, database.Results) {
		result := database.Results{
			InvalidVote: form.InvalidateVote,
			Votes:       map[string]int64{},
		}
		if result.InvalidVote {
			return Message{Positive: true}, result
		}

		if len(form.Answers) > len(vote.Info.Options) {
			return Message{Message: builder.Translation["voteInvalidWrongParameter"]}, result
		}
		for i, b := range form.Answers {
			if b > 1 || b < -1 {
				return Message{Message: builder.Translation["voteInvalidWrongParameter"]}, result
			}
			result.Votes[vote.Info.Options[i]] = b
		}
		return Message{Positive: true}, result
	})
}

func generelizeCast(acc *database.AccountAuth, docUUID string, voteUUID string, voteType database.VoteType, accountName string, f func(vote *database.Votes) (Message, database.Results)) (validate Message) {
	validate = Message{Positive: false}
	account, ok, err := IsAccountValidForUser(acc.ID, accountName)
	switch false {
	case err == nil:
		// error with author account
		validate.Message = builder.Translation["databaseErrorWithAuthorAccount"]
		return
	case ok:
		// not allowed for author account
		validate.Message = builder.Translation["notAllowedToUseAccount"]
		return
	}
	var doc *database.Document
	doc, err = extraction.GetDocumentIfCanParticipate(docUUID, account.ID)
	if err != nil || doc.Type == database.FinishedDiscussion {
		validate.Message = builder.Translation["notAllowedToVote"]
		return
	}

	var oldVote *database.Votes
	oldVote, err = extraction.GetSingleVote(voteUUID)
	if err != nil || oldVote.Info.VoteMethod != voteType || oldVote.Parent != docUUID {
		validate.Message = builder.Translation["voteDoesNotExistsInThatForm"]
		return
	}

	var result database.Results
	validate, result = f(oldVote)
	if !validate.Positive {
		return
	}
	validate.Positive = false

	err = logic.AddNewResultToVote(voteUUID, account.DisplayName, result)
	if err != nil {
		validate.Message = builder.Translation["notAllowedToVote"]
		return
	}

	validate = Message{
		Message:  builder.Translation["successfullyVoted"],
		Positive: true,
	}
	return
}
