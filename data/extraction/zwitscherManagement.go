package extraction

import (
	"PoliSim/data/database"
	"errors"
	"gorm.io/gorm"
	"time"
)

func CreateZwitscher(zwitscher *database.Zwitscher) error {
	return database.DB.Create(zwitscher).Error
}

func GetZwitscher(uuid string) (*database.Zwitscher, error) {
	zwt := &database.Zwitscher{}
	err := database.DB.Preload("Parent").Preload("Children").
		Where("uuid = ?", uuid).First(zwt).Error
	return zwt, err
}

func ChangeBlockStatus(uuid string, blockstatus bool) error {
	return database.DB.Model(&database.Zwitscher{UUID: uuid}).Updates(map[string]interface{}{
		"blocked": blockstatus,
	}).Error
}

func GetAllZwitscherFromAuthor(author string) (*database.ZwitscherList, error) {
	list := &database.ZwitscherList{}
	err := database.DB.Where("author = ?", author).Find(list).Error
	return list, err
}

type ZwitscherQueryInfo struct {
	UUID              string `input:"uuid"`
	Before            bool   `input:"before"`
	HideBlock         bool   `input:"hideblock"`
	Amount            int    `input:"amount"`
	Author            string `input:"author"`
	ShowOnlyReplies   bool   `input:"onlyreplies"`
	ShowOnlyZwitscher bool   `input:"onlyzwitscher"`
	IsAdmin           bool
}

func (info *ZwitscherQueryInfo) GetZwitscherAfter() (zwitscherList *database.ZwitscherList, exists bool, err error) {
	zwitscherList = &database.ZwitscherList{}
	exists = true
	var zwt *database.Zwitscher
	zwt, err = GetZwitscher(info.UUID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		exists = false
		zwt.Written = time.Now()
	} else if err != nil {
		return
	}
	err = info.getBasicZwitscherQuery().Where("written < ?", zwt.Written).Order("written desc").Find(zwitscherList).Error
	return
}

func (info *ZwitscherQueryInfo) GetZwitscherBefore() (zwitscherList *database.ZwitscherList, exists bool, err error) {
	zwitscherList = &database.ZwitscherList{}
	exists = true
	var zwt *database.Zwitscher
	zwt, err = GetZwitscher(info.UUID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		exists = false
		zwt.Written = time.Now()
	} else if err != nil {
		return
	}
	err = database.DB.Select("*").Table("(?) as X", info.getBasicZwitscherQuery().Order("written").Where("written > ?", zwt.Written)).Order("X.written desc").Find(zwitscherList).Error
	return
}

func (info *ZwitscherQueryInfo) getBasicZwitscherQuery() *gorm.DB {
	query := "uuid != ? AND blocked = false"
	params := make([]any, 1, 2)
	params[0] = info.UUID
	if info.IsAdmin && !info.HideBlock {
		query = "uuid != ?"
	}
	if info.Author != "" {
		query += " AND author = ?"
		params = append(params, info.Author)
		if info.ShowOnlyZwitscher {
			query += " AND linked IS NULL"
		} else if info.ShowOnlyReplies {
			query += " AND linked IS NOT NULL"
		}
	} else {
		query += " AND linked IS NULL"
	}
	return database.DB.Select("*").Where(query, params...).Limit(info.Amount).Table("zwitschers")
}
