package extraction

import (
	"PoliSim/data/database"
	"gorm.io/gorm"
	"sync"
)

var (
	OrganisationNamesList  = make([]string, 0, 20)
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
	err = database.DB.Select("name, main_group, sub_group").Distinct("main_group, sub_group, name").
		Where("status = 'hidden'").Find(&data).Error
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

func GetAllOrganisationsInSubGroup(accountID int64, mainGroup string, subGroup string) (*database.OrganisationList, error) {
	list := &database.OrganisationList{}
	err := database.DB.Joins("LEFT JOIN organisation_account ON organisations.name = organisation_account.name").
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, display_name")
		}).Preload("Admins", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, display_name")
	}).Where("main_group = ? AND sub_group = ?", mainGroup, subGroup).
		Where("organisation_account.id = ? OR status = 'public' OR status = 'private'", accountID).Select("DISTINCT organisations.name, main_group, sub_group, flair, status").Order("organisations.name").Find(list).Error
	return list, err
}

func GetAllOrganisationsInSubGroupForAdmins(mainGroup string, subGroup string) (*database.OrganisationList, error) {
	list := &database.OrganisationList{}
	err := database.DB.
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, display_name")
		}).Preload("Admins", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, display_name")
	}).Where("main_group = ? AND sub_group = ?", mainGroup, subGroup).
		Where("NOT status = 'hidden'").Order("name").Find(list).Error
	return list, err
}

// GetOrganisationGroupings returns a 2d array. The first array (array[0]) is filled with all the names of the subgroups.
// The corresponding subgroups to a main group are listed in the array that is one offset of the position in the first array.
// Meaning that all subgroups for the main group listed at array[0][0] are found in array[1] and all subgroups for array[0][12] are found in array[13].
// If the query for the subrgoups throws an error it is returned too.
func GetOrganisationGroupings(accountID int64) (*[][]string, error) {
	list := &database.OrganisationList{}
	err := database.DB.Joins("LEFT JOIN organisation_account ON organisations.name = organisation_account.name").
		Where("organisation_account.id = ? OR status = 'public' OR status = 'private'", accountID).
		Distinct("main_group, sub_group").Order("main_group, sub_group").Find(list).Error
	return transformList(list), err
}

func GetOrganisationGroupingsForAdmins() (*[][]string, error) {
	list := &database.OrganisationList{}
	err := database.DB.
		Where("NOT status = 'hidden'").
		Distinct("main_group, sub_group").Order("main_group, sub_group").Find(list).Error
	return transformList(list), err
}

func transformList(list *database.OrganisationList) *[][]string {
	var array = make([][]string, 1, 21)
	array[0] = make([]string, 0, 20)
	if len(*list) != 0 {
		array[0] = append(array[0], (*list)[0].MainGroup)
		currentMainGroup := (*list)[0].MainGroup
		array = append(array, make([]string, 0, 20))
		pos := 1
		array[pos] = append(array[pos], (*list)[0].SubGroup)

		for i := 1; i < len(*list); i++ {
			if (*list)[i].MainGroup != currentMainGroup {
				currentMainGroup = (*list)[i].MainGroup
				array[0] = append(array[0], (*list)[i].MainGroup)
				array = append(array, make([]string, 0, 20))
				pos++
				array[pos] = append(array[pos], (*list)[i].SubGroup)
				continue
			}
			array[pos] = append(array[pos], (*list)[i].SubGroup)
		}
	}
	return &array
}

func GetOrganisationForWithUserInIt(userID int64, organisationName string) (*database.Organisation, error) {
	org := &database.Organisation{}
	err := database.DB.Joins("LEFT JOIN organisation_member ON organisations.name = organisation_member.name").
		Joins("LEFT JOIN organisation_admins ON organisations.name = organisation_admins.name").
		Preload("Admins", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, display_name")
		}).
		Where("(organisation_admins.id = ? OR organisation_member.id = ?) AND organisations.name = ?", userID, userID, organisationName).
		Select("organisations.name, main_group, sub_group, flair, status").First(org).Error
	return org, err
}

func GetOrganisationsForWithUserInIt(userID int64, isAdmin bool) (*database.OrganisationList, error) {
	org := &database.OrganisationList{}
	var err error
	if isAdmin {
		err = database.DB.Joins("LEFT JOIN organisation_admins ON organisations.name = organisation_admins.name").
			Where("organisation_admins.id = ?", userID).
			Select("DISTINCT organisations.name").Order("organisations.name").Find(org).Error
	} else {
		err = database.DB.Joins("LEFT JOIN organisation_member ON organisations.name = organisation_member.name").
			Joins("LEFT JOIN organisation_admins ON organisations.name = organisation_admins.name").
			Where("organisation_admins.id = ? OR organisation_member.id = ?", userID, userID).
			Select("DISTINCT organisations.name").Order("organisations.name").Find(org).Error
	}
	return org, err
}
