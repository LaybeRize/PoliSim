package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/html/builder"
	"database/sql"
	"errors"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

// timeUntilTokenRunsOut defines the time in seconds until a token becomes invalid
const timeUntilTokenRunsOut = 60 * 60 * 24 * 7

var store = &sessions.CookieStore{}

// CreateStore sets up the cookie store for when the application starts
func CreateStore() {
	store = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_KEY")))
	store.MaxAge(timeUntilTokenRunsOut)
}

// ValidateToken returns a *extraction.AccountAuth with the Role database.NotLoggedIn if
// either no cookie exists on the request, or it is invalid. If the cookie is valid it renews it
// for the current session by writing it to the response.
func ValidateToken(r *http.Request) (returnAcc *extraction.AccountAuth) {
	session, err := store.Get(r, "session")
	temp, ok := session.Values["id"]
	id, okConvert := temp.(int64)

	returnAcc = &extraction.AccountAuth{Role: database.NotLoggedIn, Session: session}
	if err != nil || !ok || !okConvert {
		return
	}
	acc, err := extraction.GetAccountForAuth(id)
	if err != nil || acc.Suspended {
		return
	}

	returnAcc = acc
	returnAcc.Session = session
	return
}

func AddCookie(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	_ = store.Save(r, w, session)
}

// InvalidateAccountToken trys to invalidate the current cookie. ot returns a
// cookie to overwrite the current valid one.
func InvalidateAccountToken() (cookie *http.Cookie) {
	cookie = &http.Cookie{Name: "session", Value: "", Path: "/", HttpOnly: true, MaxAge: -1}
	return
}

type LoginForm struct {
	Username string `input:"username"`
	Password string `input:"password"`
}

// TryLogin always returns a ValidationMessage containing the error or success for the process.
// On success, it also returns a filled extraction.AccountLogin struct as well as writing
// a valid *http.Cookie to the response.
func (form LoginForm) TryLogin(w http.ResponseWriter, r *http.Request) (validate Message, acc *extraction.AccountLogin) {
	acc = &extraction.AccountLogin{}
	validate = Message{Positive: false}

	if form.Username == "" || form.Password == "" {
		validate.Message = builder.Translation["usernameOrPasswordMissing"]
		return
	}
	//check if user account exists
	var err error
	acc, err = extraction.GetAccountForLogin(form.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		validate.Message = builder.Translation["passwordOrUsernameWrong"]
		return
	}
	//if the database throws an error other than object not found, return an Internal Error
	if err != nil {
		validate.Message = builder.Translation["internalLoginAccountError"]
		return
	}
	//if the login block timer has not run out yet, return the time until it runs out
	if acc.NextLoginTime.Valid && !acc.NextLoginTime.Time.Before(time.Now()) {
		validate.Message = acc.NextLoginTime.Time.Format(builder.Translation["hasInternalLoginTimer"])
		return
	}
	//check password
	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(form.Password))
	if err != nil {
		//if the password is wrong update the login tries and return the correct error message
		if loginError, updateError := UpdateLoginTries(acc); updateError != nil {
			if loginError {
				validate.Message = acc.NextLoginTime.Time.Format(builder.Translation["hasInternalLoginTimer"])
				return
			} else {
				validate.Message = builder.Translation["internalLoginAccountError"]
				return
			}
		}
		validate.Message = builder.Translation["passwordOrUsernameWrong"]
		return
	}

	if acc.Suspended {
		validate.Message = builder.Translation["accountIsSuspended"]
		return
	}

	// reset account login tries and make the login timer invalid before returning the correct struct
	// and setting the cookie for identification
	acc.LoginTries = 0
	acc.NextLoginTime = sql.NullTime{}
	dbError := acc.SaveBack()
	session, _ := store.Get(r, "session")
	session.Values["id"] = acc.ID
	session.Values["role"] = -100
	err = session.Save(r, w)
	if err != nil || dbError != nil {
		//do the error handling
		validate.Message = builder.Translation["internalLoginAccountError"]
		return
	}

	validate.Positive = true
	validate.Message = builder.Translation["successFullLoggedIn"]
	return
}

// UpdateLoginTries increases the LoginTries by one and calculates the new NextLoginTime if needed
// then returns if the account is already timed out and if an error occurred on trying to save back the new
// data.
func UpdateLoginTries(acc *extraction.AccountLogin) (canNotBeLoggedIn bool, err error) {
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
