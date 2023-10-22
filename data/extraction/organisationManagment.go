package extraction

import (
	"PoliSim/data/database"
	"sync"
)

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

var organisationMutex = sync.Mutex{}

func CreateNewOrganisation(org *database.Organisation) (err error) {
	organisationMutex.Lock()
	defer organisationMutex.Unlock()
	err = database.DB.Create(&org).Error
	if err != nil {
		return
	}
	updateNewOrganisation(org)
	return
}

func GetOrganisation(name string) (org *database.Organisation, err error) {
	organisationMutex.Lock()
	defer organisationMutex.Unlock()
	err = database.DB.Preload("Members").Preload("Admins").Where("name = ?", name).First(&org).Error
	return
}

func ModifiyOrganisation(org *database.Organisation) (err error) {
	organisationMutex.Lock()
	defer organisationMutex.Unlock()
	err = database.DB.Model(database.Title{}).Where("name = ?", org.Name).Updates(&org).Error
	if err == nil {
		err = database.DB.Model(&org).Association("Members").Replace(&org.Members)
	}
	if err == nil {
		err = database.DB.Model(&org).Association("Admins").Replace(&org.Admins)
	}
	if err == nil {
		err = database.DB.Model(&org).Association("Accounts").Replace(&org.Accounts)
	}
	if err != nil {
		_ = updateOrganisationGroupings()
	}
	return
}
