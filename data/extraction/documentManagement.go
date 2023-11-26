package extraction

import (
	"PoliSim/data/database"
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

func CreateDocument(doc *database.Document) error {
	return database.DB.Create(&doc).Error
}

func GetDocumentIfNotPrivate(docType database.DocumentType, uuid string, admin bool) (*database.Document, error) {
	doc := &database.Document{}
	err := database.DB.Where("type = ? AND uuid = ? AND private = false AND (blocked = false OR true = ?)", string(docType), uuid, admin).First(doc).Error
	return doc, err
}

func GetVoteForUser(uuid string, userID int64, isAdmin bool) (*database.Document, error) {
	doc, err := GetDocumentForUser(uuid, userID, isAdmin)
	if doc.Type == database.FinishedVote || doc.Type == database.RunningVote {
		return doc, err
	}
	return nil, errors.New("is not a vote")
}

func GetDiscussionForUser(uuid string, userID int64, isAdmin bool) (*database.Document, error) {
	doc, err := GetDocumentForUser(uuid, userID, isAdmin)
	if doc.Type == database.FinishedDiscussion || doc.Type == database.RunningDiscussion {
		return doc, err
	}
	return nil, errors.New("is not a vote")
}

func GetDocumentForUser(uuid string, userID int64, isAdmin bool) (*database.Document, error) {
	doc := &database.Document{}
	err := database.DB.Joins("LEFT JOIN organisation_account ON documents.organisation = organisation_account.name").
		Joins("LEFT JOIN doc_allowed ON doc_allowed.uuid = documents.uuid").
		Preload("Viewer", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, display_name")
		}).Preload("Poster", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, display_name")
	}).Select("documents.uuid, written, documents.organisation, type, author, flair, title, subtitle"+
		", html_content, private, blocked, info, allowed_any, allowed_members").
		Where("documents.uuid = ? AND (blocked = false Or true = ?) AND (private = false OR organisation_account.id = ? OR doc_allowed.id=? OR true = ?)",
			uuid, isAdmin, userID, userID, isAdmin).First(doc).Error
	return doc, err
}

func GetDocument(uuid string) (*database.Document, error) {
	doc := &database.Document{}
	err := database.DB.Where("uuid = ?", uuid).First(doc).Error
	return doc, err
}

func UpdateDocument(document *database.Document) error {
	return database.DB.Updates(document).Error
}

func UpdateBlock(document *database.Document) error {
	return database.DB.Model(&database.Document{UUID: document.UUID}).Updates(map[string]interface{}{
		"blocked": document.Blocked,
	}).Error
}

func GetDocumentIfCanParticipate(uuid string, accountId int64) (*database.Document, error) {
	doc := &database.Document{}
	err := database.DB.Joins("LEFT JOIN organisation_admins ON documents.organisation = organisation_admins.name").
		Joins("LEFT JOIN organisation_member ON documents.organisation = organisation_member.name").
		Joins("LEFT JOIN doc_poster ON doc_poster.uuid = documents.uuid").
		Preload("Viewer", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, display_name")
		}).Preload("Poster", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, display_name")
	}).Select("documents.uuid, type").Where("documents.uuid = ? AND blocked = false AND "+
		"(doc_poster.id = ? OR (allowed_members = true AND (organisation_admins.id = ? OR organisation_member.id = ?)) OR allowed_any = true)",
		uuid, accountId, accountId, accountId).First(doc).Error
	return doc, err
}

func GetFirstDocumentBeforeTime(t time.Time, userID int64, isAdmin bool) (*database.Document, error) {
	doc := &database.Document{}
	err := database.DB.Joins("LEFT JOIN organisation_account ON documents.organisation = organisation_account.name").
		Joins("LEFT JOIN doc_allowed ON doc_allowed.uuid = documents.uuid").
		Select("documents.uuid, written, documents.organisation, type, author, flair, title, subtitle"+
			", html_content, private, blocked, info, allowed_any, allowed_members").
		Where("(blocked = false Or true = ?) AND (private = false OR organisation_account.id = ? OR doc_allowed.id=? OR true = ?) AND written > ?",
			isAdmin, userID, userID, isAdmin, t).Order("written").First(doc).Error
	return doc, err
}

func (extra *ExtraInfo) GetDocumentsAfter() (documentList *database.DocumentList, exists bool, err error) {
	documentList = &database.DocumentList{}
	exists = true
	var doc *database.Document
	t, err := time.Parse("2006-01-02", extra.Written)
	if err == nil {
		doc, err = GetFirstDocumentBeforeTime(t, extra.ViewAccountID, extra.IsAdmin)
	} else {
		doc, err = GetDocumentForUser(extra.UUID, extra.ViewAccountID, extra.IsAdmin)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		exists = false
		doc.Written = time.Now()
	} else if err != nil {
		return
	}
	err = extra.getBasicDocumentQuery().Where("written < ?", doc.Written).Order("written desc").Find(documentList).Error
	return
}

func (extra *ExtraInfo) GetDocumentsBefore() (documentList *database.DocumentList, exists bool, err error) {
	documentList = &database.DocumentList{}
	exists = true
	var doc *database.Document
	doc, err = GetDocumentForUser(extra.UUID, extra.ViewAccountID, extra.IsAdmin)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		exists = false
		doc.Written = time.Now()
	} else if err != nil {
		return
	}
	err = database.DB.Select("*").Table("(?) as X", extra.getBasicDocumentQuery().Order("written").Where("written > ?", doc.Written)).Order("X.written desc").Find(documentList).Error
	return
}

type ExtraInfo struct {
	UUID          string `input:"uuid"`
	Before        bool   `input:"before"`
	Amount        int    `input:"amount"`
	HideBlock     bool   `input:"hideblock"`
	Text          bool   `input:"text"`
	Discussion    bool   `input:"discussion"`
	Votes         bool   `input:"votes"`
	Written       string `input:"written"`
	Organisation  string `input:"organisation"`
	Author        string `input:"author"`
	Title         string `input:"title"`
	IsAdmin       bool
	ViewAccountID int64
}

func (extra *ExtraInfo) getBasicDocumentQuery() *gorm.DB {
	query := "documents.uuid != ? AND (private = false OR organisation_account.id = ? OR doc_allowed.id=? OR true = ?) "
	params := []any{extra.UUID, extra.ViewAccountID, extra.ViewAccountID, extra.IsAdmin}
	if extra.HideBlock {
		query += "AND blocked = false "
	} else {
		query += "AND (blocked = false Or true = ?) "
		params = append(params, extra.IsAdmin)
	}
	if extra.Text || extra.Discussion || extra.Votes {
		types := make([]string, 0, 5)
		if extra.Text {
			types = append(types, "type='"+string(database.LegislativeText)+"'")
		}
		if extra.Discussion {
			types = append(types, "type='"+string(database.RunningDiscussion)+"'",
				"type='"+string(database.FinishedDiscussion)+"'")
		}
		if extra.Votes {
			types = append(types, "type='"+string(database.RunningVote)+"'",
				"type='"+string(database.FinishedVote)+"'")
		}
		query += "AND (" + strings.Join(types, " OR ") + ") "
	}
	if extra.Organisation != "" {
		query += "AND organisation LIKE ? "
		params = append(params, "%"+extra.Organisation+"%")
	}
	if extra.Author != "" {
		query += "AND author LIKE ? "
		params = append(params, "%"+extra.Author+"%")
	}
	if extra.Title != "" {
		query += "AND title LIKE ? "
		params = append(params, "%"+extra.Title+"%")
	}
	return database.DB.Joins("LEFT JOIN organisation_account ON documents.organisation = organisation_account.name").
		Joins("LEFT JOIN doc_allowed ON doc_allowed.uuid = documents.uuid").
		Select("DISTINCT documents.uuid, title, type, author, organisation, written").
		Where(query, params...).Limit(extra.Amount + 1).Table("documents")
}
