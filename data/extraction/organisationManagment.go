package extraction

import "PoliSim/data/database"

func GetAllOrganisationInfo() (orgNames []string, mainGroups map[string]struct{}, subGroups map[string]struct{}, err error) {
	orgList := database.OrganisationList{}
	err = database.DB.Select("name, main_group, sub_group").Find(&orgList).Error
	orgNames = make([]string, len(orgList))
	mainGroups = make(map[string]struct{})
	subGroups = make(map[string]struct{})
	for i, orgs := range orgList {
		orgNames[i] = orgs.Name
		mainGroups[orgs.MainGroup] = struct{}{}
		subGroups[orgs.SubGroup] = struct{}{}
	}
	return
}

func GetMainAndSubOrganisationGroups() (mainGroups map[string]struct{}, subGroups map[string]struct{}, err error) {
	orgList := database.OrganisationList{}
	err = database.DB.Select("main_group, sub_group").Distinct("main_group, sub_group").Find(&orgList).Error
	mainGroups = make(map[string]struct{})
	subGroups = make(map[string]struct{})
	for _, orgs := range orgList {
		mainGroups[orgs.MainGroup] = struct{}{}
		subGroups[orgs.SubGroup] = struct{}{}
	}
	return
}

func GetHiddenOrganistaions() (data *database.OrganisationList, err error) {
	*data = database.OrganisationList{}
	err = database.DB.Select("name, main_group, sub_group").Distinct("main_group, sub_group, name").
		Where("status = 'hidden'").Find(data).Error
	return
}

// TODO: this is shit but somehow we need to update all the organistations a press account is part of, if their linked value changes
func UpdateOrganisationAccount(oldAccountID int64, newAccountID int64) (err error) {
	err = database.DB.Raw("UPDATE organisation_account SET id = ? WHERE id = ?;", newAccountID, oldAccountID).Error
	return
}
