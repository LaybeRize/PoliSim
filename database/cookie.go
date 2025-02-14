package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"
	"net/http"
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

func loadCookiesFromDB() {
	_, err := postgresDB.Exec(`CREATE TABLE IF NOT EXISTS cookies (
    session_key TEXT PRIMARY KEY,
	name TEXT,
    expiresAt TIMESTAMP,
    updateAt TIMESTAMP
)`)
	if err != nil {
		log.Fatalf("Could not create postgres cookies tabel: %v", err)
	}

	results, err := postgresDB.Query("SELECT session_key, name, expiresAt, updateAt FROM cookies;")
	if err != nil {
		log.Fatalf("Could not read postgres cookies tabel: %v", err)
	}

	for results.Next() {
		var key string
		var session = SessionData{Account: &Account{}}

		err = results.Scan(&key, &session.Account.Name, &session.ExpiresAt, &session.UpdateAt)
		if err != nil {
			slog.Error("could not scan entry correctly:", "err", err)
		}

		sessionStore[key] = &session
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

func saveCookiesToDB() {
	queryStmt := `
        INSERT INTO cookies (session_key, name, expiresAt, updateAt)
        VALUES ($1, $2, $3, $4);
    `
	for key := range sessionStore {
		session := sessionStore[key]
		_, err := postgresDB.Exec(queryStmt, &key, &session.Account.Name, &session.ExpiresAt, &session.UpdateAt)
		if err != nil {
			slog.Error("While saving colors encountered an error: ", "err", err)
		}
	}
}
