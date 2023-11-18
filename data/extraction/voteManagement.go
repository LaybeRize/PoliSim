package extraction

import "PoliSim/data/database"

func CreateVote(vote *database.Votes) error {
	return database.DB.Create(vote).Error
}

func UpdateVote(vote *database.Votes) error {
	return database.DB.Where("uuid = ?", vote.UUID).Updates(vote).Error
}

func GetSingleVote(voteUUID string) (*database.Votes, error) {
	vote := &database.Votes{}
	err := database.DB.Where("uuid = ?", voteUUID).First(vote).Error
	return vote, err
}

func GetVotesForDocument(docID string) (database.VotesList, error) {
	list := database.VotesList{}
	err := database.DB.Where("parent = ?", docID).Order("uuid").Find(&list).Error
	return list, err
}
