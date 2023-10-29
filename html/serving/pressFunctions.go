package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"github.com/go-chi/chi"
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
	composition.GetHTMXFunctions[composition.ViewHiddenNewspaper] = GetHiddenNewsPaperService
	composition.PageTitleMap[composition.RejectArticle] = builder.Translation["rejectArticlePageTitle"]
	composition.GetHTMXFunctions[composition.RejectArticle] = GetRejectArticleService
	composition.PostHTMXFunctions[composition.RejectArticle] = PostRejectArticleService
}

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

func GetHiddenNewsPaperService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	uuid := chi.URLParam(r, "uuid")
	html := composition.GetViewSingleHiddenNewspaper(uuid)
	viewHiddenNewspaperRenderRequest(w, r, acc, html)
}

var viewHiddenNewspaperRenderRequest = genericRenderer(composition.ViewHiddenNewspaperList)

func PostRejectArticleService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	//uuid := chi.URLParam(r, "uuid")
	rejectArticleRenderRequest(w, r, acc, nil)
}

func GetRejectArticleService(w http.ResponseWriter, r *http.Request) {
	acc, ok := CheckUserPrivileges(r, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !ok {
		ShowErrorPage(w, r, acc, builder.Translation["notAllowedToViewThisPage"])
		return
	}

	//uuid := chi.URLParam(r, "uuid")
	rejectArticleRenderRequest(w, r, acc, nil)
}

var rejectArticleRenderRequest = genericRenderer(composition.RejectArticle)
