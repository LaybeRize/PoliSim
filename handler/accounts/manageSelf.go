package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
	"strconv"
	"time"
)

func GetMyProfile(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	setting := handler.ModifyPersonalSettings{FontScaling: acc.FontSize, TimeZone: acc.TimeZone.String()}

	handler.MakeFullPage(writer, acc, &handler.MyProfilePage{Settings: setting})
}

func PostUpdateMySettings(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.RedirectToErrorPage(writer)
		return
	}
	page := handler.ModifyPersonalSettings{IsError: true, FontScaling: acc.FontSize, TimeZone: acc.TimeZone.String()}

	err := request.ParseForm()
	if err != nil {
		page.Message = "Fehler beim parsen der Informationen"
		handler.MakeSpecialPagePart(writer, &page)
		return
	}

	newSize, err := strconv.Atoi(helper.GetFormEntry(request, "fontScaling"))
	if err != nil {
		page.Message = "Die Seitenskalierung ist keine valide Zahl"
		handler.MakeSpecialPagePart(writer, &page)
		return
	}
	if newSize < 10 {
		page.Message = "Die Seitenskalierung kann nicht auf eine Zahl kleiner 10 gesetzt werden"
		handler.MakeSpecialPagePart(writer, &page)
		return
	}
	newTimeZone, err := time.LoadLocation(helper.GetFormEntry(request, "timeZone"))
	if err != nil {
		page.Message = "Die ausgewählte Zeitzone ist nicht valide"
		handler.MakeSpecialPagePart(writer, &page)
		return
	}

	page.TimeZone = newTimeZone.String()
	page.FontScaling = int64(newSize)
	acc.TimeZone = newTimeZone
	acc.FontSize = page.FontScaling
	err = database.SetPersonalSettings(acc)
	if err != nil {
		page.Message = "Fehler beim speichern der persönlichen Informationen"
		handler.MakeSpecialPagePart(writer, &page)
		return
	}

	page.IsError = false
	page.Message = "Einstellungen erfolgreich gespeichert\nLaden sie die Seite neu, um den Effekt zu sehen"
	handler.MakeSpecialPagePart(writer, &page)
	return
}

func PostUpdateMyPassword(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	page := &handler.ChangePassword{IsError: true}
	if !loggedIn {
		handler.RedirectToErrorPage(writer)
		return
	}

	err := request.ParseForm()
	if err != nil {
		page.Message = "Fehler beim parsen der Informationen"
		handler.MakeSpecialPagePart(writer, page)
		return
	}

	page.OldPassword = helper.GetFormEntry(request, "oldPassword")
	page.NewPassword = request.Form.Get("newPassword")
	page.RepeatNewPassword = request.Form.Get("newPasswordRepeat")
	if !database.VerifyPassword(acc.Password, page.OldPassword) {
		page.Message = "Das alte Passwort ist falsch"
		handler.MakeSpecialPagePart(writer, page)
		return
	}
	if page.NewPassword != page.RepeatNewPassword {
		page.Message = "Die Wiederholung stimmt nicht mit dem neuen Passwort überein"
		handler.MakeSpecialPagePart(writer, page)
		return
	}
	if len(page.NewPassword) < 10 {
		page.Message = "Das neue Passwort ist kürzer als 10 Zeichen"
		handler.MakeSpecialPagePart(writer, page)
		return
	}
	newPassword, err := database.HashPassword(page.NewPassword)
	if err != nil {
		page.Message = "Fehler beim hashen des neuen Passworts"
		handler.MakeSpecialPagePart(writer, page)
		return
	}
	acc.Password = newPassword
	err = database.UpdatePassword(acc)
	if err != nil {
		page.Message = "Fehler beim speichern des neuen Passworts"
		handler.MakeSpecialPagePart(writer, page)
		return
	}

	handler.MakeSpecialPagePart(writer, &handler.ChangePassword{
		Message: "Passwort erfolgreich angepasst",
		IsError: false,
	})
}
