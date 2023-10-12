package dataExtraction

import "PoliSim/database"

var TitleGroupMap = make(map[string]map[string]struct{})

func GetAll() (*database.TitleList, error) {
	list := &database.TitleList{}
	err := database.DB.Find(list).Error
	return list, err
}

func GetAllDistinct() (*database.TitleList, error) {
	list := &database.TitleList{}
	err := database.DB.Distinct("main_group, sub_group").Find(list).Error
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
	for _, item := range *list {
		TitleGroupMap[item.MainGroup][item.SubGroup] = struct{}{}
	}
}
