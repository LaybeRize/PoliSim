package logic

import "PoliSim/data/extraction"

func UpdateLetterNotification(accountId int64) {
	count, err := extraction.GetCountOfLetters(accountId)
	if err != nil {
		return
	}
	_ = extraction.UpdateHasLettersFlag(accountId, count != 0)
}

func SetLetterAsRead(uuid string, accountID int64) {
	extraction.SetLetterAsRead(uuid, accountID)
	UpdateLetterNotification(accountID)
}
