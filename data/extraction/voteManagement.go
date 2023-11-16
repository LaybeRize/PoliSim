package extraction

import "PoliSim/data/database"

func CreateVote(vote *database.Votes) error {
	return database.DB.Create(vote).Error
}
