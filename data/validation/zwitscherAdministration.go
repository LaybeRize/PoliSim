package validation

import "PoliSim/data/extraction"

type CreateZwitscher struct {
	Account         string `input:"authorAccount"`
	Content         string `input:"content"`
	ParentZwitscher string
}

func (form *CreateZwitscher) CreateZwitscher(acc *extraction.AccountAuth) (validation Message) {
	return
}
