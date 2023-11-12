package logic

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"github.com/google/uuid"
	"sync"
	"time"
)

var discussionMutex = sync.Mutex{}

func CloseDiscussionIfTimeIsUp(ending time.Time, uuidStr string) {
	if ending.After(time.Now()) {
		return
	}
	discussionMutex.Lock()
	defer discussionMutex.Unlock()

	doc, err := extraction.GetDocument(uuidStr)
	if err != nil {
		return
	}
	doc.Type = database.FinishedDiscussion
	_ = extraction.UpdateDocument(doc)
}

func AddComment(author string, flair string, comment string, uuidStr string) error {
	discussionMutex.Lock()
	defer discussionMutex.Unlock()
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
	discussionMutex.Lock()
	defer discussionMutex.Unlock()
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
