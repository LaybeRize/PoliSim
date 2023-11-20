package logic

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"sync"
)

var documentMutex = sync.Mutex{}

func BlockDocument(uuid string) (*database.Document, error) {
	documentMutex.Lock()
	defer documentMutex.Unlock()
	doc, err := extraction.GetDocument(uuid)
	if err != nil {
		return nil, err
	}
	doc.Blocked = !doc.Blocked
	err = extraction.UpdateBlock(doc)
	return doc, err
}
