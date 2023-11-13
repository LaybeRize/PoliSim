package composition

import (
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"strconv"
)

const voteContainerDiv = "vote-container-div"

func GetCreateVotePage(val validation.Message) Node {
	return getBasePageWrapper(
		SCRIPT(SRC("/public/json-enc-custom.js")),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateVoteDocument), CLASS("w-[800px]"),
			HXEXTEND("json-enc-custom"),
			getSimpleTextInput("title", "title", "", Translation["titleTextDocument"]),
			getSimpleTextInput("subtitle", "subtitle", "", Translation["subtitleTextDocument"]),
			getTextArea("content", "content", "", Translation["contentTextDocument"],
				true, true),
			DIV(ID(voteContainerDiv),
				getPartialVote("1")),
			getPartialButton("2", false), BR(),
			getSubmitButton(Translation["createVoteDocument"]),
		),
		GetMessage(val),
		getPreviewElement(),
	)
}

const votePartialButtonID = "vote-partial-button"

func GetVotePartial(partialNumber int64) Node {
	partial := strconv.FormatInt(partialNumber, 10)
	return Group(getPartialVote(partial),
		getPartialButton(strconv.FormatInt(partialNumber+1, 10), true),
	)
}

type (
	ProofOfConcept struct {
		Title     string      `json:"title"`
		Subtitle  string      `json:"subtitle"`
		Content   string      `json:"content"`
		Questions []*Question `json:"question"`
	}
	Question struct {
		Text    string   `json:"questionText"`
		Answers []string `json:"answers"`
	}
)

func getPartialVote(number string) Node {
	return DIV(CLASS("w-[800px] box box-e p-2 mt-2"),
		STYLE("--clr-border: rgb(40 51 69);"),
		P(CLASS("text-xl"), Text(Translation["votePartialHeader"], number)),
		getSimpleTextInput("question-questionText-"+number, "question["+number+"][questionText]",
			"", Translation["voteQuestionText"]),
		getEditableList([]string{}, "question["+number+"][answers]", "",
			Translation["voteAddAnswersToQuestion"], "w-full"),
		BUTTON(CLASS("bg-slate-700 text-white p-2 mt-2 hover:bg-rose-800"),
			HYPERSCRIPT("on click tell me.parentElement remove yourself"), Text(Translation["deleteVote"])),
	)
}

func getPartialButton(number string, withSwap bool) Node {
	return BUTTON(If(withSwap, HXSWAPOOB("true")), ID(votePartialButtonID), TYPE("button"),
		HXTARGET("#"+voteContainerDiv), HXSWAP("beforeend"), HXPATCH("/"+APIPreRoute+string(requestVotePartialLink)+number),
		P(CLASS("bg-slate-700 text-white p-2 mt-2"), STYLE("text-align: center;"),
			Text(Translation["addnewVoteToPost"]),
		))
}
