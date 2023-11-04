package composition

import (
	. "PoliSim/html/builder"
)

func GetCreateVotePage() Node {
	return getBasePageWrapper(
		SCRIPT(SRC("/public/json-enc-custom.js")),
	)
}
