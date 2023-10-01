package htmlComposition

import (
	. "PoliSim/componentHelper"
)

func GetStartPage(additionInfo string) Node {
	return El(DIV,
		El(P, Text("Das ist eine Startseite")),
		El(P, Text(additionInfo)),
	)
}
