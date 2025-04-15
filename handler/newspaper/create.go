package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func GetCreateArticlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.CreateArticlePage{}
	page.IsError = true

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		slog.Debug(err.Error())
		page.Message = loc.CouldNotFindAllAuthors
		arr = make([]string, 0)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = acc.Name
	page.PossibleAuthors = arr
	page.PossibleNewspaper, err = database.GetNewspaperNameListForAccount(acc.Name)
	if err != nil {
		slog.Debug(err.Error())
		page.Message = "\n" + loc.NewspaperCouldNotLoadAllNewspaperForAccount
		page.Message = strings.TrimSpace(page.Message)
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateArticlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	article := &database.NewspaperArticle{
		Title:    values.GetTrimmedString("title"),
		Subtitle: values.GetTrimmedString("subtitle"),
		Author:   values.GetTrimmedString("author"),
		RawBody:  values.GetTrimmedString("markdown"),
	}
	isSpecial := values.GetBool("special")
	newspaper := values.GetTrimmedString("newspaper")

	if article.Title == "" || article.RawBody == "" {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ContentOrBodyAreEmpty})
		return
	}

	const maxTitleLength = 400
	if len([]rune(article.Title)) > maxTitleLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ErrorTitleTooLong, maxTitleLength)})
		return
	}

	const maxSubtitleLength = 600
	if len([]rune(article.Subtitle)) > maxSubtitleLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.NewspaperSubtitleTooLong, maxSubtitleLength)})
		return
	}

	allowed, err := database.CheckIfUserAllowedInNewspaper(acc, article.Author, newspaper)
	if !allowed || err != nil {
		if err != nil {
			slog.Debug(err.Error())
		}
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.NewspaperMissingPermissionForNewspaper})
		return
	}

	article.Flair, err = database.GetAccountFlairs(&database.Account{Name: article.Author})
	article.Body = helper.MakeMarkdown(article.RawBody)
	if err != nil {
		slog.Debug(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ErrorLoadingFlairInfoForAccount})
		return
	}

	err = database.CreateArticle(article, isSpecial, newspaper)
	if err != nil {
		slog.Error(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.NewspaperErrorWhileCreatingArticle})
		return
	}

	page := &handler.CreateArticlePage{}
	page.IsError = false
	page.Message = loc.NewspaperSuccessfullyCreatedArticle

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		slog.Debug(err.Error())
		page.Message = "\n" + loc.CouldNotFindAllAuthors
		arr = make([]string, 0)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = article.Author
	page.PossibleAuthors = arr
	page.PossibleNewspaper, err = database.GetNewspaperNameListForAccount(acc.Name)
	if err != nil {
		page.Message = "\n" + loc.NewspaperCouldNotLoadAllNewspaperForAccount
	}

	handler.MakePage(writer, acc, page)

}

func GetFindNewspaperForAccountPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		slog.Debug(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.CreateArticlePage{}
	page.Author = values.GetTrimmedString("author")
	allowed, err := database.IsAccountAllowedToPostWith(acc, page.Author)
	if !allowed || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissionForAccountInfo})
		return
	}

	page.PossibleNewspaper, err = database.GetNewspaperNameListForAccount(page.Author)
	if err != nil {
		slog.Debug(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.NewspaperCouldNotLoadAllNewspaperForAccount})
		return
	}
	handler.MakeSpecialPagePart(writer, page)
}
