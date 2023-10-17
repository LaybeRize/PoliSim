package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var DB *gorm.DB

// ConnectDatabase establishes a connection stored in DB
// made with the Env Parameter: DB_USER, DB_PASSWORD, DB_ADDRESS and DB_NAME.
// It exists the program on error.
func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_NAME"))), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while connecting to postgres DB:\n"+err.Error()+"\n")
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Connection to DB established\n")

	err = DB.AutoMigrate(Account{}, Votes{}, Document{}, Title{}, Organisation{})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while migrating data model to postgres DB:\n"+err.Error()+"\n")
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Data models migrated\n")
}
