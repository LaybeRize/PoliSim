package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"log/slog"
	"net/http"
)

func GetSpecificPublicationPage(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	pubID := request.PathValue("id")
	err := database.GetPublicationForUser(pubID, acc.IsAtLeastPressAdmin())
	if err != nil {
		slog.Debug(err.Error())
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.ViewPublicationPage{}
	var pub *database.Publication
	pub, page.Articles, err = database.GetPublication(pubID)
	if page.QueryError = err != nil; err == nil {
		page.Publication = *pub
	} else {
		slog.Error(err.Error())
	}

	handler.MakeFullPage(writer, acc, page)
}

func PatchPublishPublication(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	pubID := request.PathValue("id")
	err := database.GetPublicationForUser(pubID, acc.IsAtLeastPressAdmin())
	if err != nil {
		slog.Debug(err.Error())
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	err = database.PublishPublication(pubID)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: loc.NewspaperErrorDuringPublication,
			IsError: true,
		})
	}

	page := &handler.ViewPublicationPage{}
	var pub *database.Publication
	pub, page.Articles, err = database.GetPublication(pubID)
	if page.QueryError = err != nil; err == nil {
		page.Publication = *pub
	} else {
		slog.Error(err.Error())
	}

	handler.MakePage(writer, acc, page)
}

func DeleteArticle(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	if !acc.IsAtLeastPressAdmin() {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: loc.RequestParseError,
			IsError: true,
		})
		return
	}
	rejectionText := values.GetTrimmedString("rejection")
	if rejectionText == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: loc.NewspaperRejectionMessageEmpty,
			IsError: true,
		})
		return
	}

	transaction, err := database.RejectableArticle(request.PathValue("id"))
	if err != nil {
		slog.Debug("Possible error while trying to locate article", "error", err)
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: loc.NewspaperErrorFindingArticleToReject,
			IsError: true,
		})
		return
	}
	err = transaction.DeleteArticle()
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: loc.NewspaperErrorDeletingArticle,
			IsError: true,
		})
		return
	}

	letter := &database.Letter{
		Title: fmt.Sprintf(loc.NewspaperFormatTitleForRejection,
			transaction.Article.Title, transaction.NewspaperName),
		Author:   loc.AdministrationName,
		Flair:    "",
		Signable: false,
		Body: handler.MakeMarkdown(fmt.Sprintf(loc.NewspaperFormatContentForRejection,
			rejectionText, transaction.Article.RawBody)),
		Reader: []string{transaction.Article.Author},
	}
	letter.ID = helper.GetUniqueID(letter.Author)

	err = transaction.CreateLetter(letter)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{
			Message: loc.NewspaperErrorCreatingRejectionLetter,
			IsError: true,
		})
		return
	}

	writer.Header().Add("HX-Redirect", "/check/newspapers")
	writer.WriteHeader(http.StatusFound)
}
