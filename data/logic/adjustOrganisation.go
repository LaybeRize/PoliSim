package logic

import "PoliSim/data/database"

// TODO: this is shit but somehow we need to update all the organistations a press account is part of, if their linked value changes
// rework in to a coherent function that can be called by the press account adjusting function

func UpdateOrganisationAccount(oldAccountID int64, newAccountID int64) (err error) {
	err = database.DB.Raw("UPDATE organisation_account SET id = ? WHERE id = ?;", newAccountID, oldAccountID).Error
	return
}
