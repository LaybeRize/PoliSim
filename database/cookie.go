package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

var sessionStore = make(map[string]*SessionData)
var mu sync.Mutex

type SessionData struct {
	Account   *Account  `json:"Account,omitempty"`
	ExpiresAt time.Time `json:"ExpireDate,omitempty"`
	UpdateAt  time.Time `json:"UpdateDate,omitempty"`
}

const cleanupInterval = 5 * time.Hour
const expirationTime = 7 * 24 * time.Hour
const updateTime = 30 * time.Minute
const cookieName = "poli_sim_cookie"

func init() {
	log.Println("Starting Cookie Cleanup Routine")
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

const cookiesFilePath = folderPath + "/cookies.json"

func loadCookiesFromDisk() {
	if _, err := os.Stat(cookiesFilePath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Directioary can not be created: %v", err)
		}
		sessionStore = make(map[string]*SessionData)
		return
	}
	file, err := os.Open(cookiesFilePath)
	if err != nil {
		log.Fatalf("Cookie file not found: %v", err)
	}
	err = json.NewDecoder(file).Decode(&sessionStore)
	if err != nil {
		log.Fatalf("Cookie file not correctly decoded: %v", err)
	}
	doCleanup()
	for key, session := range sessionStore {
		session.Account, err = GetAccountByName(session.Account.Name)
		if err != nil {
			slog.Error("Could not retrieve Account for Cookie:", "error", err.Error())
			delete(sessionStore, key)
			continue
		}
	}

}

func saveCookiesToDisk() {
	file, err := os.Create(cookiesFilePath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	err = file.Truncate(0)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	err = json.NewEncoder(file).Encode(&sessionStore)
	if err != nil {
		slog.Error(err.Error())
	}
}
