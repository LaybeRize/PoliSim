package account

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
	"strconv"
)

func GetMyProfile(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	setting := handler.ModifyPersonalSettings{}
	if acc.FontSize != nil {
		setting.FontScaling = *acc.FontSize
	} else {
		setting.FontScaling = 100
	}

	handler.MakeFullPage(writer, acc, &handler.MyProfilePage{Settings: setting})
}

func PostUpdateMySettings(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	page := handler.ModifyPersonalSettings{IsError: true}
	if !loggedIn {
		page.Message = "Du bist nicht eingelogt"
		handler.MakeSpecialPagePart(writer, handler.SETTING_CHANGE, page)
		return
	}

	err := request.ParseForm()
	if err != nil {
		page.Message = "Fehler beim parsen der Informationen"
		handler.MakeSpecialPagePart(writer, handler.SETTING_CHANGE, page)
		return
	}

	newSize, err := strconv.Atoi(request.Form.Get("fontScaling"))
	if acc.FontSize != nil {
		page.FontScaling = *acc.FontSize
	} else {
		page.FontScaling = 100
	}
	if err != nil {
		page.Message = "Die Seitenskalierung ist keine valide Zahl"
		handler.MakeSpecialPagePart(writer, handler.SETTING_CHANGE, page)
		return
	}
	if newSize < 10 {
		page.Message = "Die Seitenskalierung kann nicht auf eine Zahl kleiner 10 gesetzt werden"
		handler.MakeSpecialPagePart(writer, handler.SETTING_CHANGE, page)
		return
	}

	page.FontScaling = int64(newSize)
	acc.FontSize = &page.FontScaling
	err = database.SetPersonalSettings(acc)
	if err != nil {
		page.Message = "Fehler beim speichern der persönlichen"
		handler.MakeSpecialPagePart(writer, handler.SETTING_CHANGE, page)
		return
	}

	page.IsError = false
	page.Message = "Einstellungen erfolgreich gespeichert\nLaden sie die Seite neu, um den Effekt zu sehen"
	handler.MakeSpecialPagePart(writer, handler.SETTING_CHANGE, page)
	return
}

func PostUpdateMyPassword(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	page := handler.ChangePassword{IsError: true}
	if !loggedIn {
		page.Message = "Du bist nicht eingelogt"
		handler.MakeSpecialPagePart(writer, handler.PASSWORD_CHANGE, page)
		return
	}

	err := request.ParseForm()
	if err != nil {
		page.Message = "Fehler beim parsen der Informationen"
		handler.MakeSpecialPagePart(writer, handler.PASSWORD_CHANGE, page)
		return
	}

	page.OldPassword = request.Form.Get("oldPassword")
	page.NewPassword = request.Form.Get("newPassword")
	page.RepeatNewPassword = request.Form.Get("newPasswordRepeat")
	if !database.VerifyPassword(acc.Password, page.OldPassword) {
		page.Message = "Das alte Passwort ist falsch"
		handler.MakeSpecialPagePart(writer, handler.PASSWORD_CHANGE, page)
		return
	}
	if page.NewPassword != page.RepeatNewPassword {
		page.Message = "Die Wiederholung stimmt nicht mit dem neuen Passwort überein"
		handler.MakeSpecialPagePart(writer, handler.PASSWORD_CHANGE, page)
		return
	}
	if len(page.NewPassword) < 10 {
		page.Message = "Das neue Passwort ist kürzer als 10 Zeichen"
		handler.MakeSpecialPagePart(writer, handler.PASSWORD_CHANGE, page)
		return
	}
	newPassword, err := database.HashPassword(page.NewPassword)
	if err != nil {
		page.Message = "Fehler beim hashen des neuen Passworts"
		handler.MakeSpecialPagePart(writer, handler.PASSWORD_CHANGE, page)
		return
	}
	acc.Password = newPassword
	err = database.UpdatePassword(acc)
	if err != nil {
		page.Message = "Fehler beim speichern des neuen Passworts"
		handler.MakeSpecialPagePart(writer, handler.PASSWORD_CHANGE, page)
		return
	}

	handler.MakeSpecialPagePart(writer, handler.PASSWORD_CHANGE, &handler.ChangePassword{
		Message: "Passwort erfolgreich angepasst",
		IsError: false,
	})
}
