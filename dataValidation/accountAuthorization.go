package dataValidation

import (
	"PoliSim/componentHelper"
	"PoliSim/dataExtraction"
	"PoliSim/database"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// timeUntilTokenRunsOut defines the time in seconds until a token becomes invalid
var timeUntilTokenRunsOut = 60 * 60 * 24 * 7

var store = &sessions.CookieStore{}

// CreateStore sets up the cookie store for when the application starts
func CreateStore() {
	store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
	store.MaxAge(timeUntilTokenRunsOut)
}

// ValidateToken returns a *dataExtraction.AccountAuth with the Role database.NotLoggedIn if
// either no token exists, or the token is invalid. If the token is valid it renews the cookie
// for the current session by writing it to the response.
func ValidateToken(w http.ResponseWriter, r *http.Request) (returnAcc *dataExtraction.AccountAuth) {
	session, err := store.Get(r, "session")
	temp, ok := session.Values["id"]
	id, okConvert := temp.(int64)

	returnAcc = &dataExtraction.AccountAuth{Role: database.NotLoggedIn}
	if err != nil || !ok || !okConvert {
		return
	}
	acc, err := dataExtraction.GetAccountForAuth(id)
	if err != nil || acc.Suspended {
		return
	}

	_ = sessions.Save(r, w)
	returnAcc = acc
	return
}

// InvalidateAccountToken trys to invalidate the current cookie. on success, it retuns a nil error
// and a cookie to overwrite the current valid one. On failure, it returns a nil cookie and the error.
func InvalidateAccountToken() (cookie *http.Cookie) {
	cookie = &http.Cookie{Name: "session", Value: "", Path: "/", HttpOnly: true, MaxAge: -1}
	return
}

type LoginForm struct {
	Username string `input:"username"`
	Password string `input:"password"`
}

// TryLogin always returns a ValidationMessage containg the error or sucess for the process.
// On sucess it also returns a filled dataExtraction.AccountLogin struct as well as a new
// valid *http.Cookie.
func (form LoginForm) TryLogin(w http.ResponseWriter, r *http.Request) (validate ValidationMessage, acc *dataExtraction.AccountLogin) {
	acc = &dataExtraction.AccountLogin{}
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
	session, getError := store.Get(r, "session")
	session.Values["id"] = acc.ID
	err = session.Save(r, w)
	if err != nil || getError != nil {
		//do the error handling
		validate.Message = componentHelper.Translation["internalLoginAccountError"]
		return
	}
	/*acc.LoginTries = 0
	acc.NextLoginTime = sql.NullTime{}
	acc.ExpirationDate.Time = time.Now().UTC().Add(time.Second * time.Duration(timeUntilTokenRunsOut))
	acc.ExpirationDate.Valid = true
	acc.RefreshToken = uuid.New().String()
	err = acc.SaveBack()
	if err != nil {
		validate.Message = componentHelper.Translation["internalLoginAccountError"]
		return
	}*/

	validate.Positive = true
	validate.Message = componentHelper.Translation["successFullLoggedIn"]
	//cookie = &http.Cookie{Name: "token", Value: acc.RefreshToken, Path: "/", MaxAge: timeUntilTokenRunsOut}
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
