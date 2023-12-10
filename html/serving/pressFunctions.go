package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InstallPress() {
	composition.PageTitleMap[composition.CreatePressRelease] = builder.Translation["pressReleaseCreatePageTitle"]
	composition.SidebarTitleMap[composition.CreatePressRelease] = builder.Translation["pressReleaseCreateSidebarText"]
	composition.GetHTMXFunctions[composition.CreatePressRelease] = GetCreatePressCreateService
	composition.PostHTMXFunctions[composition.CreatePressRelease] = PostCreatePressCreateService

	composition.PageTitleMap[composition.ViewHiddenNewspaperList] = builder.Translation["viewHiddenNewspaperPageTitle"]
	composition.SidebarTitleMap[composition.ViewHiddenNewspaperList] = builder.Translation["viewHiddenNewspaperSidebarText"]
	composition.GetHTMXFunctions[composition.ViewHiddenNewspaperList] = GetHiddenNewsPaperListService
	composition.PageTitleMap[composition.ViewHiddenNewspaper] = builder.Translation["viewSingleHiddenNewspaperPageTitle"]
	composition.GetHTMXFunctions[composition.ViewHiddenNewspaper] = GetHiddenNewsPaperService
	composition.PageTitleMap[composition.RejectArticle] = builder.Translation["rejectArticlePageTitle"]
	composition.GetHTMXFunctions[composition.RejectArticle] = GetRejectArticleService
	composition.PostHTMXFunctions[composition.RejectArticle] = PostRejectArticleService
	composition.PatchHTMXFunctions[composition.PublishNewspaper] = PatchPublishNewspaperService

	composition.PageTitleMap[composition.ViewNewspaperList] = builder.Translation["viewNewspaperPageTitle"]
	composition.SidebarTitleMap[composition.ViewNewspaperList] = builder.Translation["viewNewspaperSidebarText"]
	composition.GetHTMXFunctions[composition.ViewNewspaperList] = GetNewsPaperListService
	composition.PageTitleMap[composition.ViewNewspaper] = builder.Translation["viewSingleNewspaperPageTitle"]
	composition.GetHTMXFunctions[composition.ViewNewspaper] = GetNewsPaperService
}

func PatchPublishNewspaperService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}
	msg, uuid := validation.PublishNewspaper(chi.URLParam(r, "uuid"))
	if !msg.Positive {
		viewSingleHiddenNewspaperOnlySwapMessage(w, r, msg, acc)
		return
	}

	pushURL(w, "/"+string(composition.ViewNewspaperList)+"/"+uuid)
	html := composition.GetSingleNewspaperPage(uuid)
	newspaperRenderRequest(w, r, acc, html)
}

var viewSingleHiddenNewspaperOnlySwapMessage = genericMessageSwapper(composition.ViewHiddenNewspaper)

func GetNewsPaperService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(r)

	html := composition.GetSingleNewspaperPage(chi.URLParam(r, "uuid"))
	newspaperRenderRequest(w, r, acc, html)
}

var newspaperRenderRequest = genericRenderer(composition.ViewNewspaper)

func GetNewsPaperListService(w http.ResponseWriter, r *http.Request) {
	acc, _ := CheckUserPrivileges(r)

	extraInfo := &logic.QueryInfo{}
	extractURLFieldValues(extraInfo, r, 5, 10, 50)
	html := composition.GetNewspaperListPage(extraInfo)
	newspaperListRenderRequest(w, r, acc, html)
}

var newspaperListRenderRequest = genericRenderer(composition.ViewNewspaperList)

func GetCreatePressCreateService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetCreatePressReleasePage(acc, &validation.CreateArticle{}, validation.Message{})
	createPressReleaseRenderRequest(w, r, acc, html)
}

func PostCreatePressCreateService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	create := &validation.CreateArticle{}
	err := extractFormValuesForFields(create, r, 0)
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		createPressReleaseOnlySwapMessage(w, r, msg, acc)
		return
	}

	msg = create.CreateArticle(acc.ID)
	if !msg.Positive {
		createPressReleaseOnlySwapMessage(w, r, msg, acc)
		return
	}

	html := composition.GetCreatePressReleasePage(acc, create, msg)
	createPressReleaseRenderRequest(w, r, acc, html)
}

var createPressReleaseRenderRequest = genericRenderer(composition.CreatePressRelease)
var createPressReleaseOnlySwapMessage = genericMessageSwapper(composition.CreatePressRelease)

func GetHiddenNewsPaperListService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	html := composition.GetViewOfHiddenNewspaper()
	viewHiddenNewspaperRenderRequest(w, r, acc, html)
}

var viewHiddenNewspaperRenderRequest = genericRenderer(composition.ViewHiddenNewspaperList)

func GetHiddenNewsPaperService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	uuid := chi.URLParam(r, "uuid")
	html := composition.GetViewSingleHiddenNewspaper(uuid)
	viewSingleHiddenNewspaperRenderRequest(w, r, acc, html)
}

var viewSingleHiddenNewspaperRenderRequest = genericRenderer(composition.ViewHiddenNewspaper)

func PostRejectArticleService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	msg := validation.Message{Positive: false}

	err := r.ParseForm()
	content := ""
	if err == nil {
		content = r.PostFormValue("content")
	}
	if err != nil {
		msg.Message = builder.Translation["extractionError"]
		rejectArticleOnlySwapMessage(w, r, msg, acc)
		return
	}
	uuid := chi.URLParam(r, "uuid")

	msg = validation.RejectArticle(uuid, content)
	if !msg.Positive {
		rejectArticleOnlySwapMessage(w, r, msg, acc)
		return
	}
	//return to the hidden newspapers
	pushURL(w, "/"+string(composition.ViewHiddenNewspaperList))
	html := composition.GetViewOfHiddenNewspaper()
	viewHiddenNewspaperRenderRequest(w, r, acc, html)
}

func GetRejectArticleService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	uuid := chi.URLParam(r, "uuid")

	html := composition.GetRejectArticlePage(uuid, "", validation.Message{})
	rejectArticleRenderRequest(w, r, acc, html)
}

var rejectArticleRenderRequest = genericRenderer(composition.RejectArticle)
var rejectArticleOnlySwapMessage = genericMessageSwapper(composition.RejectArticle)
