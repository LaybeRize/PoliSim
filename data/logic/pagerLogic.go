package logic

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
)

type ExtraInfo struct {
	UUID            string `input:"uuid"`
	Before          bool   `input:"before"`
	Amount          int    `input:"amount"`
	ViewAccountID   int64
	ViewAccountName string
}
type ViewLetter struct {
	LetterList *database.LetterList
	NextUUID   string
	BeforeUUID string
}

type ViewNewspaper struct {
	PaperList  *database.PublicationList
	NextUUID   string
	BeforeUUID string
}

func (info *ExtraInfo) GetLetter() (*ViewLetter, error) {
	viewInfo := &ViewLetter{NextUUID: "", BeforeUUID: ""}
	var exists bool
	var err error
	if info.Before {
		viewInfo.LetterList, exists, err = extraction.GetLettersBefore(info.UUID, info.Amount+1, info.ViewAccountID)
		if err != nil || len(*viewInfo.LetterList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.LetterList) == info.Amount+1 {
			viewInfo.BeforeUUID = (*viewInfo.LetterList)[0].UUID
			*viewInfo.LetterList = (*viewInfo.LetterList)[1:]
		}
		if exists {
			viewInfo.NextUUID = (*viewInfo.LetterList)[len(*viewInfo.LetterList)-1].UUID
		}
	} else {
		viewInfo.LetterList, exists, err = extraction.GetLettersAfter(info.UUID, info.Amount+1, info.ViewAccountID)
		if err != nil || len(*viewInfo.LetterList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.LetterList) == info.Amount+1 {
			viewInfo.NextUUID = (*viewInfo.LetterList)[info.Amount-1].UUID
			*viewInfo.LetterList = (*viewInfo.LetterList)[:info.Amount]
		}
		if exists {
			viewInfo.BeforeUUID = (*viewInfo.LetterList)[0].UUID
		}
	}
	return viewInfo, err
}

func (info *ExtraInfo) GetNewspaper() (*ViewNewspaper, error) {
	viewInfo := &ViewNewspaper{NextUUID: "", BeforeUUID: ""}
	var exists bool
	var err error
	if info.Before {
		viewInfo.PaperList, exists, err = extraction.GetPublicationBefore(info.UUID, info.Amount+1)
		if err != nil || len(*viewInfo.PaperList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.PaperList) == info.Amount+1 {
			viewInfo.BeforeUUID = (*viewInfo.PaperList)[0].UUID
			*viewInfo.PaperList = (*viewInfo.PaperList)[1:]
		}
		if exists {
			viewInfo.NextUUID = (*viewInfo.PaperList)[len(*viewInfo.PaperList)-1].UUID
		}
	} else {
		viewInfo.PaperList, exists, err = extraction.GetPublicationAfter(info.UUID, info.Amount+1)
		if err != nil || len(*viewInfo.PaperList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.PaperList) == info.Amount+1 {
			viewInfo.NextUUID = (*viewInfo.PaperList)[info.Amount-1].UUID
			*viewInfo.PaperList = (*viewInfo.PaperList)[:info.Amount]
		}
		if exists {
			viewInfo.BeforeUUID = (*viewInfo.PaperList)[0].UUID
		}
	}
	return viewInfo, err
}
