package htmlComposition

import . "PoliSim/componentHelper"

func GetCreateVotePage() Node {
	return getBasePageWrapper(
		SCRIPT(SRC("/public/json-enc-custom.js")),
	)
}
