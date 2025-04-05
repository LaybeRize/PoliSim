package database

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var sessionStore = make(map[string]*SessionData)
var sessionMutex sync.RWMutex

type SessionData struct {
	Account   *Account
	ExpiresAt time.Time
	UpdateAt  time.Time
	Lock      sync.Mutex
	InUse     bool
}

func (s *SessionData) IsValid() bool {
	if s == nil {
		return false
	}
	s.Lock.Lock()
	if s.InUse {
		return true
	}
	s.Lock.Unlock()
	return false
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
		InUse:     true,
	}
	setSessionCookie(w, sessionID)
}

func RefreshSession(w http.ResponseWriter, r *http.Request) (*Account, bool) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, false
	}

	sessionID := cookie.Value

	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	data := sessionStore[sessionID]
	if !data.IsValid() {
		return nil, false
	}
	defer data.Lock.Unlock()

	if time.Now().After(data.ExpiresAt) || data.Account.Blocked {
		data.InUse = false
		return nil, false
	}

	if time.Now().After(data.UpdateAt) {
		data.InUse = false
		CreateSession(w, data.Account)
		return data.Account, true
	}

	setSessionCookie(w, sessionID)

	return data.Account, true
}

func EndSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}

	sessionID := cookie.Value

	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

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

var generator = rand.New(rand.NewSource(time.Now().UnixMilli()))

func generateSessionID() string {
	prefix := make([]byte, 8)
	suffix := make([]byte, 8)
	generator.Read(prefix)

	timeNano := time.Now().UnixNano()
	suffix[0] += byte(timeNano)
	suffix[1] += byte(timeNano >> 8)
	suffix[2] += byte(timeNano >> 16)
	suffix[3] += byte(timeNano >> 24)
	suffix[4] += byte(timeNano >> 32)
	suffix[5] += byte(timeNano >> 40)
	suffix[6] += byte(timeNano >> 48)
	suffix[7] += byte(timeNano >> 56)

	return fmt.Sprintf("TOKEN-%X-%X", prefix, suffix)
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
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	for _, sessionData := range sessionStore {
		sessionData.Lock.Lock()
		if sessionData.InUse && sessionData.Account.Name == account.Name {
			sessionData.Account = account
		}
		sessionData.Lock.Unlock()
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
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	for sessionID, sessionData := range sessionStore {
		if !sessionData.InUse ||
			time.Now().After(sessionData.ExpiresAt) ||
			sessionData.Account.Blocked {
			delete(sessionStore, sessionID)
		}
	}
}

func loadCookiesFromDB() {
	results, err := postgresDB.Query(`
SELECT account.username, account.password, account.role, account.blocked, account.font_size, account.time_zone,
    cookies.session_key, cookies.name, cookies.expires_at, cookies.update_at 
FROM cookies LEFT JOIN account ON cookies.name = account.name;`)
	if err != nil {
		log.Fatalf("Could not read postgres cookies tabel: %v", err)
	}

	for results.Next() {
		var key string
		var session = SessionData{Account: &Account{}, InUse: true}
		timeZoneStr := ""

		err = results.Scan(&session.Account.Username, &session.Account.Password, &session.Account.Role,
			&session.Account.Blocked, &session.Account.FontSize, &timeZoneStr, &key, &session.Account.Name,
			&session.ExpiresAt, &session.UpdateAt)
		if err != nil {
			slog.Error("could not scan entry correctly:", "err", err)
			continue
		}
		session.Account.TimeZone, err = time.LoadLocation(timeZoneStr)
		if err != nil {
			slog.Error("could not convert account time zone correctly:",
				"account_name", session.Account.Name, "err", err)
			continue
		}

		sessionStore[key] = &session
	}

	doCleanup()
}

func saveCookiesToDB() {
	queryStmt := `
        INSERT INTO cookies (session_key, name, expires_at, update_at)
        VALUES ($1, $2, $3, $4) ON CONFLICT (session_key) DO UPDATE SET name = $2, expires_at = $3, update_at = $4;
    `
	for key := range sessionStore {
		session := sessionStore[key]
		if !session.InUse {
			continue
		}
		_, err := postgresDB.Exec(queryStmt, &key, &session.Account.Name, &session.ExpiresAt, &session.UpdateAt)
		if err != nil {
			slog.Error("While saving colors encountered an error: ", "err", err)
		}
	}
}
