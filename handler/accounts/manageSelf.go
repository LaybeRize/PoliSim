package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
	"time"
)

const messageIdPassword = "message-div-password"

func GetMyProfile(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.MyProfilePage{
		Settings: handler.ModifyPersonalSettings{
			FontScaling: acc.FontSize,
			TimeZone:    acc.TimeZone.String(),
		},
		Password: handler.ChangePassword{
			MessageUpdate: handler.MessageUpdate{
				ElementID: messageIdPassword,
			},
		},
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostUpdateMySettings(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.RedirectToErrorPage(writer)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen"})
		return
	}

	acc.FontSize = values.GetInt("fontScaling")
	if acc.FontSize < 10 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die Seitenskalierung kann nicht auf eine Zahl kleiner 10 gesetzt werden"})
		return
	}
	acc.TimeZone, err = time.LoadLocation(values.GetTrimmedString("timeZone"))
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die ausgewählte Zeitzone ist nicht valide"})
		return
	}

	err = database.SetPersonalSettings(acc)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim speichern der persönlichen Informationen"})
		return
	}

	page := &handler.ModifyPersonalSettings{
		FontScaling: acc.FontSize,
		TimeZone:    acc.TimeZone.String(),
		MessageUpdate: handler.MessageUpdate{
			Message: "Einstellungen erfolgreich gespeichert\nLaden sie die Seite neu, um den Effekt zu sehen",
			IsError: false,
		},
	}

	handler.MakeSpecialPagePart(writer, page)
	return
}

func PostUpdateMyPassword(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)

	if !loggedIn {
		handler.RedirectToErrorPage(writer)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim parsen der Informationen", ElementID: messageIdPassword})
		return
	}

	oldPassword := values.GetString("oldPassword")
	newPassword := values.GetString("newPassword")
	repeatNewPassword := values.GetString("newPasswordRepeat")
	if !database.VerifyPassword(acc.Password, oldPassword) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Das alte Passwort ist falsch", ElementID: messageIdPassword})
		return
	}
	if newPassword != repeatNewPassword {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Die Wiederholung stimmt nicht mit dem neuen Passwort überein", ElementID: messageIdPassword})
		return
	}
	if len(newPassword) < 10 {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Das neue Passwort ist kürzer als 10 Zeichen", ElementID: messageIdPassword})
		return
	}
	newPassword, err = database.HashPassword(newPassword)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim hashen des neuen Passworts", ElementID: messageIdPassword})
		return
	}
	acc.Password = newPassword
	err = database.UpdatePassword(acc)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim speichern des neuen Passworts", ElementID: messageIdPassword})
		return
	}

	page := &handler.ChangePassword{}
	page.ElementID = messageIdPassword
	page.IsError = false
	page.Message = "Passwort erfolgreich angepasst"
	handler.MakeSpecialPagePart(writer, page)
}
