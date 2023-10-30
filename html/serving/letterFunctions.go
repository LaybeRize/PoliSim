package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"net/http"
)

func InstallLetter() {
	composition.PageTitleMap[composition.CreateLetter] = builder.Translation["letterCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateLetter] = builder.Translation["letterCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateLetter] = GetCreateLetterService
	composition.PostHTMXFunctions[composition.CreateLetter] = PostCreateletterService
}

func GetCreateLetterService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetCreateNormalLetterPage(acc, &validation.CreateLetter{}, validation.Message{})
	createLetterRenderRequest(w, r, acc, html)
}

func PostCreateletterService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	create := &validation.CreateLetter{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createLetterOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.CreateNormalLetter(acc.ID)
	if !msg.Positive {
		createLetterOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetCreateNormalLetterPage(acc, create, msg)
	createLetterRenderRequest(w, r, acc, html)
}

var createLetterRenderRequest = genericRenderer(composition.CreateLetter)
var createLetterOnlySwapMessage = genericMessageSwapper(composition.CreateLetter)
