package extraction

import (
	"PoliSim/data/database"
	"gorm.io/gorm"
	"sync"
)

var TitleGroupMap = make(map[string][]string)
var TitleMainGroupList = make([]string, 0)
var TitleSubGroupNameMap = make(map[string]struct{})

type (
	TitleNameList []TitleName
	TitleName     struct {
		Name string
	}
)

func getAllNames() (*TitleNameList, error) {
	list := &TitleNameList{}
	err := database.DB.Model(&database.Title{}).Find(list).Error
	return list, err
}

func GetAllTitleNames() ([]string, error) {
	list, err := getAllNames()
	if err != nil {
		return []string{}, err
	}
	strList := make([]string, len(*list))
	for i, title := range *list {
		strList[i] = title.Name
	}
	return strList, nil
}

func GetAllDistinct() (*database.TitleList, error) {
	list := &database.TitleList{}
	err := database.DB.Distinct("main_group, sub_group").Order("main_group, sub_group").Find(list).Error
	return list, err
}

func GetAllTitlesInSubGroup(mainGroup string, subGroup string) (*database.TitleList, error) {
	list := &database.TitleList{}
	err := database.DB.Preload("Holder", func(db *gorm.DB) *gorm.DB {
		return db.Select("display_name")
	}).Where("main_group = ? AND sub_group = ?", mainGroup, subGroup).Order("name").Find(list).Error
	return list, err
}

func UpdateTitleGroupMap() {
	list, err := GetAllDistinct()
	if err != nil {
		return
	}
	TitleGroupMap = make(map[string][]string)
	TitleSubGroupNameMap = make(map[string]struct{})
	TitleMainGroupList = make([]string, 0, 20)
	index := -1
	for _, item := range *list {
		if _, ok := TitleGroupMap[item.MainGroup]; !ok {
			TitleGroupMap[item.MainGroup] = make([]string, 0)
		}
		TitleGroupMap[item.MainGroup] = append(TitleGroupMap[item.MainGroup], item.SubGroup)
		if index == -1 || TitleMainGroupList[index] != item.MainGroup {
			TitleMainGroupList = append(TitleMainGroupList, item.MainGroup)
			index++
		}
		TitleSubGroupNameMap[item.SubGroup] = struct{}{}
	}
}

var titleMutex = sync.Mutex{}

func GetTitle(name string) (title *database.Title, err error) {
	titleMutex.Lock()
	defer titleMutex.Unlock()
	err = database.DB.Preload("Holder").Where("name = ?", name).First(&title).Error
	return
}

func CreateTitle(title *database.Title) error {
	titleMutex.Lock()
	defer titleMutex.Unlock()
	err := database.DB.Create(title).Error
	UpdateTitleGroupMap()
	return err
}

func UpdateTitle(title *database.Title, oldTitleName string) error {
	titleMutex.Lock()
	defer titleMutex.Unlock()
	err := database.DB.Model(database.Title{}).Where("name = ?", oldTitleName).Updates(&title).Error
	UpdateTitleGroupMap()
	if err == nil {
		err = database.DB.Model(&title).Association("Holder").Replace(&title.Holder)
	}
	return err
}

func DeleteTitle(title *database.Title) error {
	titleMutex.Lock()
	defer titleMutex.Unlock()
	err := database.DB.Model(&title).Association("Holder").Replace(&[]database.Account{})
	if err == nil {
		err = database.DB.Delete(&title).Error
	}
	UpdateTitleGroupMap()
	return err
}

func GetTitlesForUser(userID int64) (titleList *database.TitleList, err error) {
	titleList = &database.TitleList{}
	err = database.DB.Joins("LEFT JOIN title_account ON titles.name = title_account.name").
		Select("DISTINCT titles.name").
		Where("title_account.id = ?", userID).Find(&titleList).Error
	return
}
