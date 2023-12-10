package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
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
	composition.PatchHTMXFunctions[composition.MarkAllLetterAccount] = PatchMarkLetterService

	composition.PageTitleMap[composition.ViewSingleLetter] = builder.Translation["letterViewSinglePageTitle"]
	composition.GetHTMXFunctions[composition.ViewSingleLetter] = GetViewSingleLetterService
	composition.PatchHTMXFunctions[composition.UpdateLetter] = PatchSigningLetterService

	composition.PageTitleMap[composition.ViewModMails] = builder.Translation["modMailListViewPageTitle"]
	composition.SidebarTitleMap[composition.ViewModMails] = builder.Translation["modMailListViewSidebarText"]
	composition.GetHTMXFunctions[composition.ViewModMails] = GetViewModMailListService

	composition.PageTitleMap[composition.CreateModmail] = builder.Translation["modMailCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreateModmail] = builder.Translation["modMailCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreateModmail] = GetCreateModMailService
	composition.PostHTMXFunctions[composition.CreateModmail] = PostCreateModMailService
}

func PostCreateModMailService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	create := &validation.CreateLetter{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createModMailOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.CreateModMail()
	if !msg.Positive {
		createModMailOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetCreateModMailPage(create, msg)
	createModMailRenderRequest(w, r, acc, html)
}

func GetCreateModMailService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetCreateModMailPage(&validation.CreateLetter{}, validation.Message{})
	createModMailRenderRequest(w, r, acc, html)
}

var createModMailRenderRequest = genericRenderer(composition.CreateModmail)
var createModMailOnlySwapMessage = genericMessageSwapper(composition.CreateModmail)

var standardAmount = 10

func GetViewModMailListService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}
	extraInfo := &logic.QueryInfo{}
	extractURLFieldValues(extraInfo, r, 5, int64(standardAmount), 50)

	html := composition.GetViewModmailList(acc, extraInfo)
	viewModMailListRenderRequest(w, r, acc, html)
}

var viewModMailListRenderRequest = genericRenderer(composition.ViewModMails)

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

func PatchMarkLetterService(w http.ResponseWriter, r *http.Request) {
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
	var account *database.Account
	account, ok, err = validation.IsAccountValidForUser(acc.ID, name)
	if !ok || err != nil {
		ShowErrorPage(w, r, acc, builder.Translation["letterAccountError"])
		return
	}
	err = extraction.SetAllLetterAsReadForAccount(account.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		viewLetterSwapMessage(w, r, validation.Message{Message: builder.Translation["errorMarkingAllAsRead"]}, acc)
	}

	extraInfo := &logic.QueryInfo{
		UUID:            "",
		Before:          false,
		Amount:          standardAmount,
		ViewAccountID:   account.ID,
		ViewAccountName: account.DisplayName,
	}
	html := composition.GetLetterViewPersonalLetters(acc, extraInfo)
	viewLetterRenderRequest(w, r, acc, html)
}

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
	var account *database.Account
	account, ok, err = validation.IsAccountValidForUser(acc.ID, name)
	if !ok || err != nil {
		ShowErrorPage(w, r, acc, builder.Translation["letterAccountError"])
		return
	}
	pushURL(w, "/"+string(composition.ViewLetterLink)+
		url.PathEscape(account.DisplayName))

	extraInfo := &logic.QueryInfo{
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
	extraInfo := &logic.QueryInfo{
		ViewAccountID:   account.ID,
		ViewAccountName: account.DisplayName,
	}
	extractURLFieldValues(extraInfo, r, 5, int64(standardAmount), 50)

	html := composition.GetLetterViewPersonalLetters(acc, extraInfo)
	viewLetterRenderRequest(w, r, acc, html)
}

var viewLetterRenderRequest = genericRenderer(composition.ViewLetter)
var viewLetterSwapMessage = genericMessageSwapper(composition.ViewLetter)

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
