package composition

import (
	. "PoliSim/html/builder"
	"strconv"
)

const voteContainerDiv = "vote-container-div"

func GetCreateVotePage() Node {
	return getBasePageWrapper(
		SCRIPT(SRC("/public/json-enc-custom.js")),
		DIV(ID(voteContainerDiv)),
		getPartialButton("1", false),
	)
}

const votePartialButtonID = "vote-partial-button"

func GetVotePartial(partialNumber int64) Node {
	partial := strconv.FormatInt(partialNumber, 10)
	return Group(DIV(
		Text(partial),
	),
		getPartialButton(strconv.FormatInt(partialNumber+1, 10), true),
	)
}

func getPartialButton(number string, withSwap bool) Node {
	return BUTTON(If(withSwap, HXSWAPOOB("true")), ID(votePartialButtonID), TYPE("button"),
		HXTARGET("#"+voteContainerDiv), HXSWAP("beforeend"), HXPATCH("/"+APIPreRoute+string(requestVotePartialLink)+number),
		P(CLASS("bg-slate-700 text-white p-2 mt-2 ml-2"), STYLE("text-align: center;"),
			Text(Translation["addnewVoteToPost"]),
		))
}
