package extraction

import "PoliSim/data/database"

var (
	OrganisationNamesList  = []string{}
	OrganisationMainGroups = map[string]struct{}{}
	OrganisationSubGroups  = map[string]struct{}{}
)

func StartupUpdateOrganisation() (err error) {
	orgList := database.OrganisationList{}
	err = database.DB.Select("name, main_group, sub_group").Find(&orgList).Error
	if err != nil {
		return
	}
	OrganisationNamesList = make([]string, len(orgList))
	OrganisationMainGroups = make(map[string]struct{})
	OrganisationSubGroups = make(map[string]struct{})
	for i, orgs := range orgList {
		OrganisationNamesList[i] = orgs.Name
		OrganisationMainGroups[orgs.MainGroup] = struct{}{}
		OrganisationSubGroups[orgs.SubGroup] = struct{}{}
	}
	return
}

func CreateNewOrganisation(org *database.Organisation) (err error) {
	err = database.DB.Create(&org).Error
	if err != nil {
		return
	}
	updateNewOrganisation(org)
	return
}

func updateNewOrganisation(organisation *database.Organisation) {
	OrganisationNamesList = append(OrganisationNamesList, organisation.Name)
	OrganisationMainGroups[organisation.MainGroup] = struct{}{}
	OrganisationSubGroups[organisation.SubGroup] = struct{}{}
}

// updateOrganisationGroupings only updates the maps for main groups and subgroups.
func updateOrganisationGroupings() (err error) {
	orgList := database.OrganisationList{}
	err = database.DB.Select("main_group, sub_group").Distinct("main_group, sub_group").Find(&orgList).Error
	OrganisationMainGroups = make(map[string]struct{})
	OrganisationSubGroups = make(map[string]struct{})
	for _, orgs := range orgList {
		OrganisationMainGroups[orgs.MainGroup] = struct{}{}
		OrganisationSubGroups[orgs.SubGroup] = struct{}{}
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
