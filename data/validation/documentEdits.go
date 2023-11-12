package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
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

func (form *AddComment) AddComment(uuidStr string, acc *extraction.AccountAuth) (validate Message) {
	validate = Message{Positive: false}

	return Message{Positive: true}
}
