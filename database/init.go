package database

import (
	loc "PoliSim/localisation"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var postgresDB *sql.DB

type DbError string

func (d DbError) Error() string {
	return string(d)
}

const (
	NotAllowedError              DbError = "action is for user not allowed"
	NoRecipientFoundError        DbError = "no recipient found for letter"
	AlreadyVoted                 DbError = "you already voted"
	DocumentHasInvalidVisibility DbError = "document has invalid visibility"
	DocumentHasNoAttachedVotes   DbError = "document has no attached votes"
	DoubleChatRoomEntry          DbError = "there already exists a room with these members"
)

var shutdown sync.Mutex

// HashPassword creates a hash of the given password, for later verification
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword checks if the provided password matches the stored hash
func VerifyPassword(storedHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

func init() {
	startPostgresDatabase()
	afterStartProcesses()
}

func startPostgresDatabase() {
	psqlInfo := fmt.Sprintf("host=%s port=5432 user=%s "+
		"password=%s dbname=%s sslmode=disable", os.Getenv("DB_ADDRESS"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	var err error
	postgresDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		closeErr := postgresDB.Close()
		log.Fatalf("Postgres DB connection error: %v | Driver close error: %v", err, closeErr)
	}
	err = postgresDB.Ping()
	if err != nil {
		closeErr := postgresDB.Close()
		log.Fatalf("Postgres DB Ping error: %v | Driver close error: %v", err, closeErr)
	}
	log.Println("Opened connection to the Postgres DB")
}

func afterStartProcesses() {
	migrate()

	loadColorPalettesFromDB()
	log.Println("Loading Cookies")
	loadCookiesFromDB()
	createRootAccount()
	createAdministrationAccount()
	log.Println("Starting Vote Cleanup Routine")
	generateResults()
	go resultRoutine()
}

func Shutdown() {
	shutdown.Lock()
	defer shutdown.Unlock()
	close(OwnerChangeOnAccountChannel)
	close(BlockedAccountChannel)
	saveColorPalettesToDB()
	saveCookiesToDB()
	err := postgresDB.Close()
	if err != nil {
		log.Printf("Postgres DB close error: %v\n", err)
	}
}

func createRootAccount() {
	acc, err := GetAccountByUsername(os.Getenv("USERNAME"))
	if err == nil && acc != nil {
		log.Println("Head Admin Account already exists")
		return
	} else if errors.Is(err, sql.ErrNoRows) && acc == nil {
		pass, hashError := HashPassword(os.Getenv("PASSWORD"))
		if hashError != nil {
			log.Fatalf("password hashing error for Head Admin Account: %v", hashError)
		}
		createError := CreateAccount(&Account{
			Name:     os.Getenv("NAME"),
			Username: os.Getenv("USERNAME"),
			Password: pass,
			Role:     RootAdmin,
			Blocked:  false,
		})
		if createError != nil {
			log.Fatalf("Head Admin Account creation error: %v", createError)
		}
		log.Println("Head Admin Account successfully created")
	} else {
		log.Fatalf("Head Admin Account search error: %v", err)
	}
}

func createAdministrationAccount() {
	acc, err := GetAccountByName(loc.AdministrationAccountName)
	if err == nil && acc != nil {
		log.Println("Administration Account already exists")
		return
	} else if errors.Is(err, sql.ErrNoRows) && acc == nil {
		createError := CreateAccount(&Account{
			Name:     loc.AdministrationAccountName,
			Username: loc.AdministrationAccountUsername,
			Password: loc.AdministrationAccountPassword,
			Role:     Special,
			Blocked:  false,
		})
		if createError != nil {
			log.Fatalf("Administration Account creation error: %v", createError)
		}
		log.Println("Administration Account successfully created")
	} else {
		log.Fatalf("Administration Account search error: %v", err)
	}
}

func rollback(tx *sql.Tx) {
	_ = tx.Rollback()
}

func closeRows(result *sql.Rows) {
	_ = result.Close()
}
