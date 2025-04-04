package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"log/slog"
	"net/http"
)

func PostNewDocumentTagPage(writer http.ResponseWriter, request *http.Request) {
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

	tag := &database.DocumentTag{
		ID:              helper.GetUniqueID(acc.Name),
		Text:            values.GetTrimmedString("text"),
		BackgroundColor: values.GetTrimmedString("background-color"),
		TextColor:       values.GetTrimmedString("text-color"),
		LinkColor:       values.GetTrimmedString("link-color"),
		Links:           values.GetCommaSeperatedArray("links"),
	}

	if tag.Text == "" {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagTextEmpty})
		return
	}

	const maxLength = 400
	if len([]rune(tag.Text)) > maxLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.DocumentTagTextTooLong, maxLength)})
		return
	}

	if !helper.StringIsAColor(tag.BackgroundColor) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagColorInvalidBackground})
		return
	}

	if !helper.StringIsAColor(tag.TextColor) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagColorInvalidText})
		return
	}

	if !helper.StringIsAColor(tag.LinkColor) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagColorInvalidLink})
		return
	}

	err = database.CreateTagForDocument(request.PathValue("id"), acc, tag)
	if err != nil {
		slog.Error(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagCreationError})
		return
	}

	if obj := getDocumentPageObject(acc, request); obj != nil {
		handler.MakePage(writer, acc, obj)
	} else {
		handler.GetNotFoundPage(writer, request)
	}
}
