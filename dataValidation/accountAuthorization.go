package dataValidation

import (
	"PoliSim/componentHelper"
	"PoliSim/dataExtraction"
	"PoliSim/database"
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// timeUntilTokenRunsOut defines the time in seconds until a token becomes invalid
var timeUntilTokenRunsOut = 60 * 60 * 24 * 7

// ValidateToken returns a *dataExtraction.AccountAuth with the Role database.NotLoggedIn if
// either no token exists, or the token is invalid. If the token is valid it writes a new token
// to the database for the account and returns a new valid cookie for the account.
func ValidateToken(token string) (returnAcc *dataExtraction.AccountAuth, cookie *http.Cookie) {
	cookie = nil
	returnAcc = &dataExtraction.AccountAuth{Role: database.NotLoggedIn}
	acc, err := dataExtraction.GetAccountForAuth(token)
	if err != nil {
		return
	}

	if !acc.ExpirationDate.Valid || acc.ExpirationDate.Time.Before(time.Now()) || acc.Suspended {
		return
	}

	acc.ExpirationDate.Time = time.Now().UTC().Add(time.Second * time.Duration(timeUntilTokenRunsOut))
	acc.ExpirationDate.Valid = true
	acc.RefreshToken = uuid.New().String()
	err = dataExtraction.UpdateAuthToken(acc.ID, acc.RefreshToken, acc.ExpirationDate)
	if err != nil {
		return
	}

	cookie = &http.Cookie{Name: "token", Value: acc.RefreshToken, Path: "/", MaxAge: timeUntilTokenRunsOut}
	returnAcc = acc
	return
}

// InvalidateAccountToken trys to invalidate the current cookie. on success, it retuns a nil error
// and a cookie to overwrite the current valid one. On failure, it returns a nil cookie and the error.
func InvalidateAccountToken(acc *dataExtraction.AccountAuth) (cookie *http.Cookie, err error) {
	cookie = nil
	err = dataExtraction.UpdateAuthToken(acc.ID, acc.RefreshToken, sql.NullTime{})
	if err != nil {
		return
	}
	cookie = &http.Cookie{Name: "token", Value: "", Path: "/", MaxAge: 0}
	return
}

type LoginForm struct {
	Username string `input:"username"`
	Password string `input:"password"`
}

// TryLogin always returns a ValidationMessage containg the error or sucess for the process.
// On sucess it also returns a filled dataExtraction.AccountLogin struct as well as a new
// valid *http.Cookie.
func (form LoginForm) TryLogin() (validate ValidationMessage, acc *dataExtraction.AccountLogin, cookie *http.Cookie) {
	acc = &dataExtraction.AccountLogin{}
	cookie = nil
	validate = ValidationMessage{Positive: false}

	if form.Username == "" || form.Password == "" {
		validate.Message = componentHelper.Translation["usernameOrPasswordMissing"]
		return
	}
	//check if user account exists
	var err error
	acc, err = dataExtraction.GetAccoutForLogin(form.Username)
	if err == gorm.ErrRecordNotFound {
		validate.Message = componentHelper.Translation["passwordOrUsernameWrong"]
		return
	}
	//if the database throws an error other than object not found, return an Internal Error
	if err != nil {
		validate.Message = componentHelper.Translation["internalLoginAccountError"]
		return
	}
	//if the login block timer has not run out yet, return the time until it runs out
	if acc.NextLoginTime.Valid && !acc.NextLoginTime.Time.Before(time.Now()) {
		validate.Message = acc.NextLoginTime.Time.Format(componentHelper.Translation["hasInternalLoginTimer"])
		return
	}
	//check password
	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(form.Password))
	if err != nil {
		//if the password is wrong update the login tries and return the correct error message
		if loginError, updateError := UpdateLoginTries(acc); updateError != nil {
			if loginError {
				validate.Message = acc.NextLoginTime.Time.Format(componentHelper.Translation["hasInternalLoginTimer"])
				return
			} else {
				validate.Message = componentHelper.Translation["internalLoginAccountError"]
				return
			}
		}
		validate.Message = componentHelper.Translation["passwordOrUsernameWrong"]
		return
	}

	if acc.Suspended {
		validate.Message = componentHelper.Translation["accountIsSupended"]
		return
	}
	//reset account login tries and make the login timer invalid before returning the correct struct
	acc.LoginTries = 0
	acc.NextLoginTime = sql.NullTime{}
	acc.ExpirationDate.Time = time.Now().UTC().Add(time.Second * time.Duration(timeUntilTokenRunsOut))
	acc.ExpirationDate.Valid = true
	acc.RefreshToken = uuid.New().String()
	err = acc.SaveBack()
	if err != nil {
		validate.Message = componentHelper.Translation["internalLoginAccountError"]
		return
	}

	validate.Positive = true
	validate.Message = componentHelper.Translation["successFullLoggedIn"]
	cookie = &http.Cookie{Name: "token", Value: acc.RefreshToken, Path: "/", MaxAge: timeUntilTokenRunsOut}
	return
}

// UpdateLoginTries increases the LoginTries by one and calculates the new NextLoginTime if needed
// then returns if the account is already timed out and if an error occured on trying to save back the new
// data.
func UpdateLoginTries(acc *dataExtraction.AccountLogin) (canNotBeLoggedIn bool, err error) {
	canNotBeLoggedIn = false
	acc.LoginTries += 1
	//set the timer appropriate for the tries
	switch acc.LoginTries {
	case 1, 2, 3:
	case 4, 5:
		acc.NextLoginTime.Time = time.Now().UTC().Add(time.Second * 5)
	case 6, 7:
		acc.NextLoginTime.Time = time.Now().UTC().Add(time.Minute)
	case 8, 9:
		acc.NextLoginTime.Time = time.Now().UTC().Add(time.Minute * 5)
	default:
		minimum := acc.LoginTries * acc.LoginTries * 10
		acc.NextLoginTime.Time = time.Now().UTC().Add(time.Second * time.Duration(minimum))
	}
	//make it valid if it had been set
	if acc.LoginTries > 3 {
		acc.NextLoginTime.Valid = true
	}
	err = acc.SaveBack()
	//check if the timer was saved correctly
	if err == nil && acc.LoginTries > 3 {
		canNotBeLoggedIn = true
	}
	acc.NextLoginTime.Time = acc.NextLoginTime.Time.In(time.Local)
	return
}
