package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/crypto/bcrypt"
	"os"
)

var ctx context.Context
var driver neo4j.DriverWithContext
var notFoundError = errors.New("item not found")
var multipleItemsError = errors.New("more then one item found")

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
		addition := ""
		if closeErr != nil {
			addition = "\n\nError while closing the driver: " + closeErr.Error()
		}
		panic(err.Error() + addition)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Opened connection to the DB\n")
	createConstraints()
	createRootAccount()

}

func createRootAccount() {
	acc, err := GetAccountByUsername(os.Getenv("USERNAME"))
	if err == nil && acc != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Head Admin Account already exists\n")
		return
	} else if errors.Is(err, notFoundError) && acc == nil {
		pass, err := HashPassword(os.Getenv("PASSWORD"))
		if err != nil {
			panic(err)
		}
		err = CreateAccount(&Account{
			Name:     os.Getenv("NAME"),
			Username: os.Getenv("USERNAME"),
			Password: pass,
			Role:     RootAdmin,
			Blocked:  false,
		})
		if err != nil {
			panic(err)
		}
		_, _ = fmt.Fprintf(os.Stdout, "Head Admin Account successfully created\n")
	} else {
		panic(err)
	}
}

func createConstraints() {
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ""})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(session, ctx)

	constraints := []string{
		"CREATE CONSTRAINT ON (acc:Account) ASSERT acc.name IS UNIQUE;",
		"CREATE CONSTRAINT ON (acc:Account) ASSERT acc.username IS UNIQUE;",
		"CREATE CONSTRAINT ON (acc:Account) ASSERT exists (acc.name);",
		"CREATE CONSTRAINT ON (acc:Account) ASSERT exists (acc.username);",
		"CREATE CONSTRAINT ON (org:Organisation) ASSERT org.name IS UNIQUE;",
		"CREATE CONSTRAINT ON (org:Organisation) ASSERT exists (org.name);",
		"CREATE CONSTRAINT ON (note:Note) ASSERT note.id IS UNIQUE;",
		"CREATE CONSTRAINT ON (note:Note) ASSERT exists (note.id);",
	}

	for _, constraint := range constraints {
		_, err := session.Run(ctx, constraint, nil)
		if err != nil {
			panic(err)
		}
	}
}
