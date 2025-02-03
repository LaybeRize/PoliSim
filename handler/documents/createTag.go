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

const elementID = "tag-message"

func PostNewDocumentTagPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentGeneralFunctionNotAvailable, ElementID: elementID})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError, ElementID: elementID})
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
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagTextEmpty, ElementID: elementID})
		return
	}

	const maxLength = 400
	if len([]rune(tag.Text)) > maxLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.DocumentTagTextTooLong, maxLength), ElementID: elementID})
		return
	}

	if !helper.StringIsAColor(tag.BackgroundColor) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagColorInvalidBackground, ElementID: elementID})
		return
	}

	if !helper.StringIsAColor(tag.TextColor) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagColorInvalidText, ElementID: elementID})
		return
	}

	if !helper.StringIsAColor(tag.LinkColor) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagColorInvalidLink, ElementID: elementID})
		return
	}

	err = database.CreateTagForDocument(request.PathValue("id"), acc, tag)
	if err != nil {
		slog.Error(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentTagCreationError, ElementID: elementID})
		return
	}

	if obj := getDocumentPageObject(acc, request); obj != nil {
		handler.MakePage(writer, acc, obj)
	} else {
		handler.GetNotFoundPage(writer, request)
	}
}
