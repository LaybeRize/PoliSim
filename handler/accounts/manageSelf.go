package accounts

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	acc.FontSize = values.GetInt("fontScaling")
	if acc.FontSize < minFontSizeScale {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.AccountFontSizeMustBeBiggerThen, minFontSizeScale)})
		return
	}
	acc.TimeZone, err = time.LoadLocation(values.GetTrimmedString("timeZone"))
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountGivenTimezoneInvalid})
		return
	}

	err = database.SetPersonalSettings(acc)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountErrorSavingPersonalSettings})
		return
	}

	handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: false,
		Message: loc.AccountPersonalSettingsSavedSuccessfully})
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
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	oldPassword := values.GetString("oldPassword")
	newPassword := values.GetString("newPassword")
	repeatNewPassword := values.GetString("newPasswordRepeat")
	if !database.VerifyPassword(acc.Password, oldPassword) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountWrongOldPassword})
		return
	}
	if newPassword != repeatNewPassword {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountWrongRepeatPassword})
		return
	}
	if len(newPassword) < minLengthPassword {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.AccountNewPasswordMinimumLength, minLengthPassword)})
		return
	}
	newPassword, err = database.HashPassword(newPassword)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountNewPasswordTooLong})
		return
	} else if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountErrorHashingNewPassword})
		return
	}
	acc.Password = newPassword
	err = database.UpdatePassword(acc)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.AccountErrorSavingNewPassword})
		return
	}

	page := &handler.ChangePassword{}
	page.IsError = false
	page.Message = loc.AccountSuccessfullySavedNewPassword
	handler.MakeSpecialPagePart(writer, page)
}
