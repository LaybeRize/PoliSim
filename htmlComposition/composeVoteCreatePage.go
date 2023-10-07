package htmlComposition

import . "PoliSim/componentHelper"

func GetCreateVotePage() Node {
	return getBasePageWrapper(
		El(SCRIPT, Attr(SRC, "/public/json-enc-custom.js")),
	)
}
