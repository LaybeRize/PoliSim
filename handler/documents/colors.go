package documents

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
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

	err := request.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen"})
		return
	}

	color := &database.ColorPalette{
		Name:       helper.GetFormEntry(request, "name"),
		Background: helper.GetFormEntry(request, "background-color"),
		Text:       helper.GetFormEntry(request, "text-color"),
		Link:       helper.GetFormEntry(request, "link-color"),
	}

	if color.Name == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Name der Farbpalette darf nicht leer sein"})
		return
	}

	if !helper.StringIsAColor(color.Background) || !helper.StringIsAColor(color.Text) ||
		!helper.StringIsAColor(color.Link) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Einer der übergebene Farben ist kein valider Hexcode"})
		return
	}

	err = database.AddColorPalette(color, acc)

	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Erstellen der Farbpalette"})
		return
	}

	page := &handler.EditColorPage{
		AllowedToCreate: true,
		AllowedToDelete: database.HasPrivilegesForColorsDelete(acc),
		ColorPalettes:   database.ColorPaletteMap,
		Color:           *color,
	}

	handler.MakePage(writer, acc, page)
}

func DeleteColor(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	err := request.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen"})
		return
	}

	color, err := database.RemoveColorPalette(helper.GetFormEntry(request, "name"), acc)

	if err != nil {
		slog.Debug(err.Error())
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Löschen der Farbpalette"})
		return
	}

	page := &handler.EditColorPage{
		AllowedToCreate: true,
		AllowedToDelete: true,
		ColorPalettes:   database.ColorPaletteMap,
		Color:           *color,
	}

	handler.MakePage(writer, acc, page)
}
