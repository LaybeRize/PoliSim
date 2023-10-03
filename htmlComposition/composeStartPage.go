package htmlComposition

import (
	. "PoliSim/componentHelper"
	"PoliSim/dataExtraction"
)

func GetStartPage(acc *dataExtraction.AccountAuth) Node {
	return getBasePageWrapper(
		getCustomPageHeader(Translation["welcomMessage"]),
		El(P, Text("Das ist eine Startseite")),
		El(P, Text(acc.DisplayName)),
		Raw(RawStartPageContent),
	)
}
