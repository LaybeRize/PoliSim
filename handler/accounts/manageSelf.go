package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
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

const minFontSizeScale = 10

func PostUpdateMySettings(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.RedirectToErrorPage(writer)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	acc.FontSize = values.GetInt("fontScaling")
	if acc.FontSize < minFontSizeScale {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.AccountFontSizeMustBeBiggerThen, minFontSizeScale)})
		return
	}
	acc.TimeZone, err = time.LoadLocation(values.GetTrimmedString("timeZone"))
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountGivenTimezoneInvalid})
		return
	}

	err = database.SetPersonalSettings(acc)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountErrorSavingPersonalSettings})
		return
	}

	page := &handler.ModifyPersonalSettings{
		FontScaling: acc.FontSize,
		TimeZone:    acc.TimeZone.String(),
		MessageUpdate: handler.MessageUpdate{
			Message: loc.AccountPersonalSettingsSavedSuccessfully,
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

	values, err := helper.GetAdvancedFormValuesWithoutDebugLogger(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError, ElementID: messageIdPassword})
		return
	}

	oldPassword := values.GetString("oldPassword")
	newPassword := values.GetString("newPassword")
	repeatNewPassword := values.GetString("newPasswordRepeat")
	if !database.VerifyPassword(acc.Password, oldPassword) {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountWrongOldPassword, ElementID: messageIdPassword})
		return
	}
	if newPassword != repeatNewPassword {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountWrongRepeatPassword, ElementID: messageIdPassword})
		return
	}
	if len(newPassword) < minLengthPassword {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message:   fmt.Sprintf(loc.AccountNewPasswordMinimumLength, minLengthPassword),
			ElementID: messageIdPassword})
		return
	}
	newPassword, err = database.HashPassword(newPassword)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountErrorHashingNewPassword, ElementID: messageIdPassword})
		return
	}
	acc.Password = newPassword
	err = database.UpdatePassword(acc)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountErrorSavingNewPassword, ElementID: messageIdPassword})
		return
	}

	page := &handler.ChangePassword{}
	page.ElementID = messageIdPassword
	page.IsError = false
	page.Message = loc.AccountSuccessfullySavedNewPassword
	handler.MakeSpecialPagePart(writer, page)
}
