package logic

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

func CloseDiscussionIfTimeIsUp(ending time.Time, uuidStr string) {
	if ending.After(time.Now()) {
		return
	}
	documentMutex.Lock()
	defer documentMutex.Unlock()

	doc, err := extraction.GetDocument(uuidStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error finding discussion: "+err.Error())
		return
	}
	if doc.Type == database.FinishedDiscussion {
		return
	}
	doc.Type = database.FinishedDiscussion
	err = extraction.UpdateDocument(doc)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error updating discussion: "+err.Error())
	}
}

func AddComment(author string, flair string, comment string, uuidStr string) error {
	documentMutex.Lock()
	defer documentMutex.Unlock()
	doc, err := extraction.GetDocument(uuidStr)
	if err != nil {
		return err
	}
	doc.Info.Discussion = append(doc.Info.Discussion, database.Discussions{
		UUID:        uuid.New().String(),
		Hidden:      false,
		Written:     time.Now(),
		Author:      author,
		Flair:       flair,
		HTMLContent: comment,
	})
	err = extraction.UpdateDocument(doc)
	return err
}

func ChangeVisibiltyComment(commentUUID, docUUID string) (bool, error) {
	documentMutex.Lock()
	defer documentMutex.Unlock()
	doc, err := extraction.GetDocument(docUUID)
	if err != nil {
		return false, err
	}
	exists := false
	for i, disc := range doc.Info.Discussion {
		if disc.UUID == commentUUID {
			exists = true
			doc.Info.Discussion[i].Hidden = !disc.Hidden
		}
	}
	err = extraction.UpdateDocument(doc)
	return exists, err
}
