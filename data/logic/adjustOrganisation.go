package logic

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"sync"
)

var OrganisationModifyMutex = sync.Mutex{}

func UpdateOrganisationAccount(pressAccountID int64) (err error) {
	OrganisationModifyMutex.Lock()
	defer OrganisationModifyMutex.Unlock()
	var list *database.OrganisationList
	list, err = extraction.GetOrganisations(pressAccountID)
	if err != nil {
		return
	}
	for _, org := range *list {
		err = updateSingleOrg(&org)
		if err != nil {
			return
		}
	}
	return
}

func updateSingleOrg(org *database.Organisation) (err error) {
	list := make([]string, 0, 20)
	mapping := map[string]struct{}{}
	for _, acc := range org.Members {
		if _, ok := mapping[acc.DisplayName]; ok {
			continue
		}
		mapping[acc.DisplayName] = struct{}{}
		list = append(list, acc.DisplayName)
	}
	for _, acc := range org.Admins {
		if _, ok := mapping[acc.DisplayName]; ok {
			continue
		}
		mapping[acc.DisplayName] = struct{}{}
		list = append(list, acc.DisplayName)
	}
	var accounts *database.AccountList
	accounts, err = extraction.GetParentAccounts(list)
	if err != nil {
		return
	}
	org.Accounts = *accounts
	err = extraction.ModifiyOrganisation(org)
	return
}
