package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
)

func GetCreatePressReleasePage(acc *extraction.AccountAuth, press *validation.CreateArticle, val validation.Message) Node {
	return getBasePageWrapper(
		getPageHeader(CreatePressRelease),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreatePressRelease), CLASS("w-[800px]"),
			getUserDropdown(acc, press.Account, Translation["accountPressRelease"]),
			getSimpleTextInput("title", "title", press.Title, Translation["pressTitle"]),
			getSimpleTextInput("subtitle", "subtitle", press.Subtitle, Translation["pressSubtitle"]),
			getCheckBox("breakingNews", press.IsBreakingNews, false, "true", "breakingNews", Translation["pressBreakingNews"], nil),
			getTextArea("content", "content", press.Content, Translation["pressContent"], true),
			getSubmitButton(Translation["createArticleButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}

func GetViewOfHiddenNewspaper() Node {
	list, err := extraction.GetHiddenPublications()
	if err != nil {
		return GetErrorPage(Translation["errorRetrievingPublication"])
	}
	nodes := make([]Node, len(*list))
	for i, item := range *list {
		link := string(ViewHiddenNewspaperList) + "/" + url.PathEscape(item.UUID)
		nodes[i] = getClickableLink("/"+APIPreRoute+link, "/"+link, Group(
			CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
			H1(CLASS("text-2xl"),
				IfElse(item.UUID == database.EternatityPublicationName, Text(Translation["hiddenStandardNewsPaper"]),
					Text(item.CreateTime.Format(Translation["hiddenBreakingNews"]))))))
	}
	return getBasePageWrapper(
		getPageHeader(ViewHiddenNewspaperList),
		Group(nodes...),
	)
}

func GetViewSingleHiddenNewspaper(uuid string) Node {
	ok, err := extraction.FindPublication(uuid, "false")
	if !ok || err != nil {
		return GetErrorPage(Translation["errorRetrievingSinglePublication"])
	}
	articleList, err := extraction.FindArticlesForPublicationUUID(uuid)
	nodes := make([]Node, len(*articleList))
	for i, item := range *articleList {
		//if i%2 == 1 {
		//	box = "box-f"
		//}
		link := string(rejectArticleLink) + "/" + item.UUID
		nodes[i] = renderSingleArticle(&item,
			getClickableLink("/"+APIPreRoute+link, "/"+link, Group(CLASS(buttonClassAttribute+" m-2"),
				STYLE("display: inline-block;"), Text(Translation["directToRejectArticleButton"]))))
	}
	if err != nil {
		return GetErrorPage(Translation["errorRetrievingArticles"])
	}
	return getBasePageWrapper(
		getCustomPageHeader(fmt.Sprintf(Translation["unpublishedNewsletterTitle"], uuid)),
		Group(nodes...),
	)
}

func GetRejectArticlePage(uuid string, content string, val validation.Message) Node {
	article, err := extraction.FindArticle(uuid, false)
	if err != nil {
		return GetErrorPage(Translation["errorFindingRejectableArticle"])
	}
	return getBasePageWrapper(
		getPageHeader(RejectArticle),
		renderSingleArticle(article, nil),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(rejectArticleLink)+"/"+url.PathEscape(uuid), CLASS("w-[800px]"),
			getTextArea("content", "content", content, Translation["rejectArticleMessage"], true),
			getSubmitButton(Translation["rejectArticleButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}

func renderSingleArticle(item *database.Article, specialNode Node) Node {
	return DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
		DIV(CLASS("w-full flex items-center flex-col"),
			H1(CLASS("text-3xl underline decoration-2 underline-offset-2"), Text(item.Headline)),
			If(item.Subtitle.Valid, H1(CLASS("text-2xl"), STYLE("font-style: italic;"), Text(item.Subtitle.String))),
		),
		P(CLASS("mx-2 mb-2"), I(Text(fmt.Sprintf(item.Written.Format(Translation["authorPressRelease"]), item.Author))),
			If(item.Flair != "", Group(I(Text("; ")), Text(item.Flair)))),
		Raw(item.HTMLContent),
		specialNode,
	)
}

func GetNewspaperListPage(extra *logic.ExtraInfo) Node {
	view, err := extra.GetNewspaper()
	if err != nil {
		return GetErrorPage(Translation["errorLoadingLetters"])
	}
	nodes := make([]Node, len(*view.PaperList))
	for i, item := range *view.PaperList {
		link := string(ViewNewspaperList) + "/" + url.PathEscape(item.UUID)
		nodes[i] = getClickableLink("/"+APIPreRoute+link, "/"+link, Group(
			CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
			H1(CLASS("text-2xl"), IfElse(item.BreakingNews, Text(item.PublishTime.Format(Translation["breakingNewsFormat"])),
				Text(item.PublishTime.Format(Translation["normalNewsFormat"]))))))
	}
	beforeLink, nextLink := generateNewspaperLink(extra, view)
	return getBasePageWrapper(
		getPageHeader(ViewNewspaperList),
		Group(nodes...),
		DIV(CLASS("w-[800px] flex justify-between flex-row"),
			DIV(If(view.BeforeUUID != "", getClickableLink("/"+APIPreRoute+beforeLink, "/"+beforeLink,
				P(CLASS("bg-slate-700 text-white p-2 mt-2"), Text(Translation["beforePage"])),
			))),
			DIV(If(view.NextUUID != "", getClickableLink("/"+APIPreRoute+nextLink, "/"+nextLink,
				P(CLASS("bg-slate-700 text-white p-2 mt-2"), Text(Translation["nextPage"])),
			))),
		),
	)
}

func generateNewspaperLink(extra *logic.ExtraInfo, view *logic.ViewNewspaper) (beforeLink string, nextLink string) {
	beforeLink = fmt.Sprintf("%s?uuid=%s&amount=%d&before=true", string(ViewNewspaperList),
		url.QueryEscape(view.BeforeUUID), extra.Amount)
	nextLink = fmt.Sprintf("%s?uuid=%s&amount=%d", string(ViewNewspaperList),
		url.QueryEscape(view.NextUUID), extra.Amount)
	return
}

func GetSingleNewspaperPage(uuid string) Node {
	ok, err := extraction.FindPublication(uuid, "false")
	if !ok || err != nil {
		return GetErrorPage(Translation["errorRetrievingSinglePublication"])
	}
	articleList, err := extraction.FindArticlesForPublicationUUID(uuid)
	nodes := make([]Node, len(*articleList))
	for i, item := range *articleList {
		nodes[i] = renderSingleArticle(&item, nil)
	}
	return getBasePageWrapper(
		getPageHeader(ViewNewspaper),
		Group(nodes...),
	)
}
