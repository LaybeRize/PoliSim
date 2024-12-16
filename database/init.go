package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"os"
)

var ctx context.Context
var driver neo4j.DriverWithContext
var notFoundError = errors.New("item not found")
var multipleItemsError = errors.New("more then one item found")

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
		panic(err.Error() + "\n\nError while closing the driver: " + closeErr.Error())
	}

	_, _ = fmt.Fprintf(os.Stdout, "Opened connection to the DB\n")
}
