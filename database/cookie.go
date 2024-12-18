package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var sessionStore = make(map[string]*SessionData)
var mu sync.Mutex

type SessionData struct {
	Account   *Account
	ExpiresAt time.Time
	UpdateAt  time.Time
}

const cleanupInterval = 5 * time.Hour
const expirationTime = 7 * 24 * time.Hour
const updateTime = 30 * time.Minute
const cookieName = "poli_sim_cookie"

func init() {
	_, _ = fmt.Fprintf(os.Stdout, "Starting Cookie Cleanup Routine\n")
	go startCleanup()
}

func CreateSession(w http.ResponseWriter, account *Account) {
	sessionID := generateSessionID()

	sessionStore[sessionID] = &SessionData{
		Account:   account,
		ExpiresAt: time.Now().Add(expirationTime),
		UpdateAt:  time.Now().Add(updateTime),
	}
	setSessionCookie(w, sessionID)
}

func RefreshSession(w http.ResponseWriter, r *http.Request) (*Account, bool) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, false
	}

	sessionID := cookie.Value

	mu.Lock()
	defer mu.Unlock()

	data, exists := sessionStore[sessionID]
	if !exists || time.Now().After(data.ExpiresAt) || data.Account.Blocked {
		delete(sessionStore, sessionID)
		return nil, false
	}

	if time.Now().After(data.UpdateAt) {
		delete(sessionStore, sessionID)
		CreateSession(w, data.Account)
		return data.Account, true
	}

	sessionStore[sessionID].ExpiresAt = time.Now().Add(expirationTime)
	setSessionCookie(w, sessionID)

	return data.Account, true
}

func EndSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}

	sessionID := cookie.Value

	mu.Lock()
	defer mu.Unlock()

	delete(sessionStore, sessionID)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
}

func generateSessionID() string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		randomBytes = []byte(fmt.Sprintf("%x", time.Now().String()))
	}

	combined := append(randomBytes, []byte(fmt.Sprintf("%d", time.Now().UnixNano()))...)

	hash := sha256.Sum256(combined)
	return hex.EncodeToString(hash[:])
}

func setSessionCookie(w http.ResponseWriter, sessionID string) {
	expiration := time.Now().Add(expirationTime)
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func updateAccount(account *Account) {
	mu.Lock()
	defer mu.Unlock()

	for sessionID, sessionData := range sessionStore {
		if sessionData.Account.Name == account.Name {
			sessionStore[sessionID].Account = account
		}
	}
}

func startCleanup() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		doCleanup()
	}
}

func doCleanup() {
	mu.Lock()
	defer mu.Unlock()

	for sessionID, sessionData := range sessionStore {
		if time.Now().After(sessionData.ExpiresAt) {
			delete(sessionStore, sessionID)
		}
	}
}
