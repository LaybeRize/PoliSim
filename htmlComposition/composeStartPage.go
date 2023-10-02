package htmlComposition

import (
	. "PoliSim/componentHelper"
	"PoliSim/dataExtraction"
)

func GetStartPage(acc *dataExtraction.AccountAuth) Node {
	return El(DIV,
		El(P, Text("Das ist eine Startseite")),
		El(P, Text(acc.DisplayName)),
	)
}
