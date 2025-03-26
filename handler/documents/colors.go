package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"errors"
	"log/slog"
	"net/http"
)

func GetColorPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.EditColorPage{
		AllowedToCreate: database.HasPrivilegesForColorsAdd(acc),
		AllowedToDelete: database.HasPrivilegesForColorsDelete(acc),
		ColorPalettes:   database.ColorPaletteMap,
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateColor(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	color := &database.ColorPalette{
		Name:       values.GetTrimmedString("name"),
		Background: values.GetString("background-color"),
		Text:       values.GetString("text-color"),
		Link:       values.GetString("link-color"),
	}

	if color.Name == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentColorPaletteNameNotEmpty})
		return
	}

	if !helper.StringIsAColor(color.Background) || !helper.StringIsAColor(color.Text) ||
		!helper.StringIsAColor(color.Link) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentInvalidColor})
		return
	}

	err = database.AddColorPalette(color, acc)

	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentErrorCreatingColorPalette})
		return
	}

	page := &handler.EditColorPage{
		AllowedToCreate: true,
		AllowedToDelete: database.HasPrivilegesForColorsDelete(acc),
		ColorPalettes:   database.ColorPaletteMap,
		Color:           *color,
	}
	page.Message = loc.DocumentSuccessfullyCreatedChangedColorPalette
	page.IsError = false

	handler.MakePage(writer, acc, page)
}

func DeleteColor(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	name := values.GetTrimmedString("name")

	color, err := database.RemoveColorPalette(name, acc)
	if errors.Is(err, database.CanNotDeleteColor) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentStandardColorNotAllowedToBeDeleted})
		return
	}
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.DocumentErrorDeletingColorPalette})
		return
	}

	page := &handler.EditColorPage{
		AllowedToCreate: true,
		AllowedToDelete: true,
		ColorPalettes:   database.ColorPaletteMap,
		Color:           *color,
	}
	page.Message = loc.DocumentSuccessfullyDeletedColorPalette
	page.IsError = false

	handler.MakePage(writer, acc, page)
}
