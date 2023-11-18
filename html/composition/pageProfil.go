package composition

import (
	"PoliSim/data/extraction"
	. "PoliSim/html/builder"
)

func GetPersonalProfil(acc *extraction.AccountAuth) Node {
	return getBasePageWrapper()
}
