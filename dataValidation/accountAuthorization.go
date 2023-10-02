package dataValidation

import (
	"PoliSim/dataExtraction"
	"PoliSim/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// timeUntilTokenRunsOut defines the time in seconds until a token becomes invalid
var timeUntilTokenRunsOut = 60 * 60 * 24 * 7

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
