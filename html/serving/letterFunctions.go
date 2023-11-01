package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi"
	"net/http"
	"net/url"
)

func InstallLetter() {
	composition.PageTitleMap[composition.CreateLetter] = builder.Translation["letterCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateLetter] = builder.Translation["letterCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateLetter] = GetCreateLetterService
	composition.PostHTMXFunctions[composition.CreateLetter] = PostCreateletterService
	composition.PageTitleMap[composition.ViewLetter] = builder.Translation["letterViewPageTitle"]
	composition.SidebarTitleMap[composition.ViewLetter] = builder.Translation["letterViewSidebarText"]
	composition.GetHTMXFunctions[composition.ViewLetter] = GetViewLetterService
	composition.PatchHTMXFunctions[composition.ChangeViewLetterAccount] = PatchViewLetterService
	composition.PageTitleMap[composition.ViewSingleLetter] = builder.Translation["letterViewSinglePageTitle"]
	composition.GetHTMXFunctions[composition.ViewSingleLetter] = GetViewSingleLetterService
	composition.PatchHTMXFunctions[composition.UpdateLetter] = PatchSigningLetterService
}

func PatchSigningLetterService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	uuid := chi.URLParam(r, "uuid")
	account, msg := validation.SignLetter(acc,
		uuid, chi.URLParam(r, "account"),
		chi.URLParam(r, "action"))
	if !msg.Positive {
		viewSingleLetterOnlySwapMessage(w, r, msg, acc)
	}

	html := composition.GetSingLetterView(account, uuid,
		account.Role != database.User, msg)
	viewSingleLetterRenderRequest(w, r, acc, html)
}

func GetViewSingleLetterService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	name := chi.URLParam(r, "account")
	uuid := chi.URLParam(r, "uuid")
	account, ok, err := validation.IsAccountValidForUser(acc.ID, name)
	if !ok || err != nil {
		ShowErrorPage(w, r, acc, builder.Translation["letterAccountError"])
		return
	}

	html := composition.GetSingLetterView(account, uuid,
		account.Role != database.User, validation.Message{})
	viewSingleLetterRenderRequest(w, r, acc, html)
}

var viewSingleLetterRenderRequest = genericRenderer(composition.ViewSingleLetter)
var viewSingleLetterOnlySwapMessage = genericMessageSwapper(composition.ViewSingleLetter)

var standardAmount = 10

func PatchViewLetterService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}
	err := r.ParseForm()
	name := ""
	if err == nil {
		name = r.PostFormValue("reader")
	}
	var account *extraction.AccountModification
	account, ok, err = validation.IsAccountValidForUser(acc.ID, name)
	if !ok || err != nil {
		ShowErrorPage(w, r, acc, builder.Translation["letterAccountError"])
		return
	}
	w.Header().Set("HX-Push-Url", "/"+string(composition.ViewLetterLink)+"/"+
		url.PathEscape(account.DisplayName))

	extraInfo := &logic.ExtraInfo{
		UUID:            "",
		Before:          false,
		Amount:          standardAmount,
		ViewAccountID:   account.ID,
		ViewAccountName: account.DisplayName,
	}
	html := composition.GetLetterViewPersonalLetters(acc, extraInfo)
	viewLetterRenderRequest(w, r, acc, html)
}

func GetViewLetterService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}
	account, ok, err := validation.IsAccountValidForUser(acc.ID, chi.URLParam(r, "account"))
	if !ok || err != nil {
		ShowErrorPage(w, r, acc, builder.Translation["letterAccountError"])
		return
	}
	extraInfo := &logic.ExtraInfo{
		ViewAccountID:   account.ID,
		ViewAccountName: account.DisplayName,
	}
	extractURLFieldValues(extraInfo, r, 5, int64(standardAmount), 50)

	html := composition.GetLetterViewPersonalLetters(acc, extraInfo)
	viewLetterRenderRequest(w, r, acc, html)
}

var viewLetterRenderRequest = genericRenderer(composition.ViewLetter)

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
