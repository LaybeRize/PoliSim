package database

import (
	loc "PoliSim/localisation"
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"sync"
)

const folderPath = "./data"

var ctx context.Context
var driver neo4j.DriverWithContext
var notFoundError = errors.New("item not found")
var notAllowedError = errors.New("action is for user not allowed")
var noRecipientFoundError = errors.New("no recipient found for letter")
var multipleItemsError = errors.New("more then one item found")
var AlreadyVoted = errors.New("you already voted")

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
	var err error
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUri := "bolt://" + os.Getenv("DB_ADDRESS")
	driver, err = neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
	ctx = context.Background()

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		closeErr := driver.Close(ctx)
		log.Fatalf("DB connection error: %v | Driver close error: %v", err, closeErr)
	}

	log.Println("Opened connection to the DB")
	log.Println("Loading Cookies")
	loadCookiesFromDisk()
	createConstraints()
	createRootAccount()
	createAdministrationAccount()
	migrate()
	log.Println("Starting Vote Cleanup Routine")
	generateResults()
	go resultRoutine()
}

func Shutdown() {
	shutdown.Lock()
	defer shutdown.Unlock()
	saveColorPalettesToDisk()
	saveCookiesToDisk()
	err := driver.Close(ctx)
	if err != nil {
		log.Fatalf("DB close error: %v", err)
	}
}

func createRootAccount() {
	acc, err := GetAccountByUsername(os.Getenv("USERNAME"))
	if err == nil && acc != nil {
		log.Println("Head Admin Account already exists")
		return
	} else if errors.Is(err, notFoundError) && acc == nil {
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
	} else if errors.Is(err, notFoundError) && acc == nil {
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

func createConstraints() {
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ""})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Fatalf("session close error: %v", err)
		}
	}(session, ctx)

	constraints := []string{
		"CREATE CONSTRAINT u_acc_name IF NOT exists FOR (acc:Account) REQUIRE acc.name IS UNIQUE;",
		"CREATE CONSTRAINT u_acc_username IF NOT exists FOR (acc:Account) REQUIRE acc.username IS UNIQUE;",
		//"CREATE CONSTRAINT r_acc_name IF NOT EXISTS FOR (acc:Account) REQUIRE acc.name IS NOT NULL;",
		//"CREATE CONSTRAINT r_acc_username IF NOT EXISTS FOR (acc:Account) REQUIRE acc.username IS NOT NULL;",
		"CREATE CONSTRAINT u_org_name IF NOT exists FOR (org:Organisation) REQUIRE org.name IS UNIQUE;",
		//"CREATE CONSTRAINT r_org_name IF NOT EXISTS FOR (org:Organisation) REQUIRE org.name IS NOT NULL;",
		"CREATE CONSTRAINT u_note_id IF NOT exists FOR (note:Note) REQUIRE note.id IS UNIQUE;",
		//"CREATE CONSTRAINT r_note_id IF NOT EXISTS FOR (note:Note) REQUIRE note.id IS NOT NULL;",
		"CREATE CONSTRAINT u_title_name IF NOT exists FOR (title:Title) REQUIRE title.name IS UNIQUE;",
		//"CREATE CONSTRAINT r_title_name IF NOT EXISTS FOR (title:Title) REQUIRE title.name IS NOT NULL;",
		"CREATE CONSTRAINT u_news_name IF NOT exists FOR (news:Newspaper) REQUIRE news.name IS UNIQUE;",
		//"CREATE CONSTRAINT r_news_name IF NOT EXISTS FOR (news:Newspaper) REQUIRE exists (news.name IS NOT NULL;",
		"CREATE CONSTRAINT u_pub_id IF NOT exists FOR (pub:Publication) REQUIRE pub.id IS UNIQUE;",
		//"CREATE CONSTRAINT r_pub_id IF NOT EXISTS FOR (pub:Publication) REQUIRE pub.id IS NOT NULL;",
		"CREATE CONSTRAINT u_art_id IF NOT exists FOR (art:Article) REQUIRE art.id IS UNIQUE;",
		//"CREATE CONSTRAINT r_art_id IF NOT EXISTS FOR (art:Article) REQUIRE art.id IS NOT NULL;",
		"CREATE CONSTRAINT u_letter_id IF NOT exists FOR (letter:Letter) REQUIRE letter.id IS UNIQUE;",
		//"CREATE CONSTRAINT r_letter_id IF NOT EXISTS FOR (letter:Letter) REQUIRE letter.id IS NOT NULL;",
		"CREATE CONSTRAINT u_document_id IF NOT exists FOR (doc:Document) REQUIRE doc.id IS UNIQUE;",
		//"CREATE CONSTRAINT r_document_id IF NOT EXISTS FOR (doc:Document) REQUIRE doc.id IS NOT NULL;",
		"CREATE CONSTRAINT u_vote_id IF NOT exists FOR (vote:Vote) REQUIRE vote.id IS UNIQUE;",
		//"CREATE CONSTRAINT r_vote_id IF NOT EXISTS FOR (vote:Vote) REQUIRE vote.id IS NOT NULL;",
		"CREATE CONSTRAINT u_comment_id IF NOT exists FOR (comment:Comment) REQUIRE comment.id IS UNIQUE;",
		//"CREATE CONSTRAINT r_comment_id IF NOT EXISTS FOR (comment:Comment) REQUIRE comment.id IS NOT NULL;",
	}

	for _, constraint := range constraints {
		_, err := session.Run(ctx, constraint, nil)
		if err != nil {
			log.Fatalf("constraint run error: %v", err)
		}
	}
}
