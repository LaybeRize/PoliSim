package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"log/slog"
	"net/http"
)

func PostCreateComment(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentGeneralFunctionNotAvailable})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	docId := request.PathValue("id")
	comment := &database.DocumentComment{
		Author: values.GetTrimmedString("author"),
		Body:   helper.MakeMarkdown(values.GetTrimmedString("markdown")),
	}
	comment.ID = helper.GetUniqueID(comment.Author)

	if comment.Body == "" {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ContentIsEmpty})
		return
	}

	allowed, err := database.IsAccountAllowedToPostWith(acc, comment.Author)
	if !allowed || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentMissingPermissionForComment})
		return
	}

	comment.Flair, err = database.GetAccountFlairs(&database.Account{Name: comment.Author})
	if err != nil {
		slog.Info(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ErrorLoadingFlairInfoForAccount})
		return
	}

	err = database.CreateDocumentComment(docId, comment)
	if err != nil {
		slog.Info(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentErrorWhileSavingComment})
		return
	}

	if obj := getDocumentPageObject(acc, request); obj != nil {
		handler.MakePage(writer, acc, obj)
	} else {
		handler.PartialGetNotFoundPage(writer, request)
	}
}
