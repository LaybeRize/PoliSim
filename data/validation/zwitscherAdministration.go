package validation

import (
	"PoliSim/data/database"
)

type CreateZwitscher struct {
	Account         string `input:"authorAccount"`
	Content         string `input:"content"`
	ParentZwitscher string
}

func (form *CreateZwitscher) CreateZwitscher(acc *database.AccountAuth) (validation Message) {
	return
}
