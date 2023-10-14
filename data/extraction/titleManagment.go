package extraction

import (
	"PoliSim/data/database"
	"sync"
)

var TitleGroupMap = make(map[string]map[string]struct{})
var TitleSubGroupNameMap = make(map[string]struct{})

func GetAll() (*database.TitleList, error) {
	list := &database.TitleList{}
	err := database.DB.Find(list).Error
	return list, err
}

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

func GetAllInSubGroup(mainGroup string, subGroup string) (*database.TitleList, error) {
	list := &database.TitleList{}
	err := database.DB.Preload("Holder").Where("main_group = ? AND sub_group = ?", mainGroup, subGroup).Order("name").Find(list).Error
	return list, err
}

func UpdateTitleGroupMap() {
	list, err := GetAllDistinct()
	if err != nil {
		return
	}
	TitleGroupMap = make(map[string]map[string]struct{})
	TitleSubGroupNameMap = make(map[string]struct{})
	for _, item := range *list {
		if _, ok := TitleGroupMap[item.MainGroup]; !ok {
			TitleGroupMap[item.MainGroup] = make(map[string]struct{})
		}
		TitleGroupMap[item.MainGroup][item.SubGroup] = struct{}{}
		TitleSubGroupNameMap[item.SubGroup] = struct{}{}
	}
}

var titleMutex = sync.Mutex{}

func GetTitle(name string) (title *database.Title, err error) {
	titleMutex.Lock()
	defer titleMutex.Unlock()
	err = database.DB.Where("name = ?", name).First(title).Error
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
	err := database.DB.Model(database.Title{}).Where("name = ?", oldTitleName).Updates(title).Error
	UpdateTitleGroupMap()
	if err != nil {
		err = database.DB.Model(database.Title{}).Association("Holder").Replace(title.Holder)
	}
	return err
}

func DeleteTitle(title *database.Title) error {
	titleMutex.Lock()
	defer titleMutex.Unlock()
	err := database.DB.Delete(title).Error
	UpdateTitleGroupMap()
	return err
}
