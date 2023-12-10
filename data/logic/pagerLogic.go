package logic

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
)

type QueryInfo struct {
	UUID            string `input:"uuid"`
	Before          bool   `input:"before"`
	Amount          int    `input:"amount"`
	ViewAccountID   int64
	ViewAccountName string
}

type ViewLetter struct {
	LetterList *database.ExtendedLetterList
	NextUUID   string
	BeforeUUID string
}

type ViewNewspaper struct {
	PaperList  *database.PublicationList
	NextUUID   string
	BeforeUUID string
}

type ViewDocuments struct {
	DocumentList *database.DocumentList
	NextUUID     string
	BeforeUUID   string
}

type ViewZwitscher struct {
	ZwitscherList *database.ZwitscherList
	NextUUID      string
	BeforeUUID    string
}

func (info *QueryInfo) GetLetter() (*ViewLetter, error) {
	viewInfo := &ViewLetter{NextUUID: "", BeforeUUID: ""}
	var exists bool
	var err error
	if info.Before {
		viewInfo.LetterList, exists, err = extraction.GetLettersBefore(info.UUID, info.Amount+1, info.ViewAccountID)
		if err != nil || len(*viewInfo.LetterList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.LetterList) == info.Amount+1 {
			viewInfo.BeforeUUID = (*viewInfo.LetterList)[1].UUID
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

func (info *QueryInfo) GetModMails() (*ViewLetter, error) {
	viewInfo := &ViewLetter{NextUUID: "", BeforeUUID: ""}
	var exists bool
	var err error
	if info.Before {
		viewInfo.LetterList, exists, err = extraction.GetModMailsBefore(info.UUID, info.Amount+1)
		if err != nil || len(*viewInfo.LetterList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.LetterList) == info.Amount+1 {
			viewInfo.BeforeUUID = (*viewInfo.LetterList)[1].UUID
			*viewInfo.LetterList = (*viewInfo.LetterList)[1:]
		}
		if exists {
			viewInfo.NextUUID = (*viewInfo.LetterList)[len(*viewInfo.LetterList)-1].UUID
		}
	} else {
		viewInfo.LetterList, exists, err = extraction.GetModMailsAfter(info.UUID, info.Amount+1)
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

func (info *QueryInfo) GetNewspaper() (*ViewNewspaper, error) {
	viewInfo := &ViewNewspaper{NextUUID: "", BeforeUUID: ""}
	var exists bool
	var err error
	if info.Before {
		viewInfo.PaperList, exists, err = extraction.GetPublicationBefore(info.UUID, info.Amount+1)
		if err != nil || len(*viewInfo.PaperList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.PaperList) == info.Amount+1 {
			viewInfo.BeforeUUID = (*viewInfo.PaperList)[1].UUID
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

func GetDocuments(info *extraction.DocumentQueryInfo) (*ViewDocuments, error) {
	viewInfo := &ViewDocuments{NextUUID: "", BeforeUUID: ""}
	var exists bool
	var err error
	if info.Before {
		viewInfo.DocumentList, exists, err = info.GetDocumentsBefore()
		if err != nil || len(*viewInfo.DocumentList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.DocumentList) == info.Amount+1 {
			viewInfo.BeforeUUID = (*viewInfo.DocumentList)[1].UUID
			*viewInfo.DocumentList = (*viewInfo.DocumentList)[1:]
		}
		if exists {
			viewInfo.NextUUID = (*viewInfo.DocumentList)[len(*viewInfo.DocumentList)-1].UUID
		}
	} else {
		viewInfo.DocumentList, exists, err = info.GetDocumentsAfter()
		if err != nil || len(*viewInfo.DocumentList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.DocumentList) == info.Amount+1 {
			viewInfo.NextUUID = (*viewInfo.DocumentList)[info.Amount-1].UUID
			*viewInfo.DocumentList = (*viewInfo.DocumentList)[:info.Amount]
		}
		if exists {
			viewInfo.BeforeUUID = (*viewInfo.DocumentList)[0].UUID
		}
	}
	return viewInfo, err
}

func GetZwitschers(info *extraction.ZwitscherQueryInfo) (*ViewZwitscher, error) {
	viewInfo := &ViewZwitscher{NextUUID: "", BeforeUUID: ""}
	var exists bool
	var err error
	if info.Before {
		viewInfo.ZwitscherList, exists, err = info.GetZwitscherBefore()
		if err != nil || len(*viewInfo.ZwitscherList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.ZwitscherList) == info.Amount+1 {
			viewInfo.BeforeUUID = (*viewInfo.ZwitscherList)[1].UUID
			*viewInfo.ZwitscherList = (*viewInfo.ZwitscherList)[1:]
		}
		if exists {
			viewInfo.NextUUID = (*viewInfo.ZwitscherList)[len(*viewInfo.ZwitscherList)-1].UUID
		}
	} else {
		viewInfo.ZwitscherList, exists, err = info.GetZwitscherAfter()
		if err != nil || len(*viewInfo.ZwitscherList) == 0 {
			return viewInfo, err
		}
		if len(*viewInfo.ZwitscherList) == info.Amount+1 {
			viewInfo.NextUUID = (*viewInfo.ZwitscherList)[info.Amount-1].UUID
			*viewInfo.ZwitscherList = (*viewInfo.ZwitscherList)[:info.Amount]
		}
		if exists {
			viewInfo.BeforeUUID = (*viewInfo.ZwitscherList)[0].UUID
		}
	}
	return viewInfo, err
}
