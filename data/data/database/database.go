package database

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"os"
)

var DB *sqlx.DB

var schema = `
CREATE TABLE IF NOT EXISTS account (
   id BIGSERIAL PRIMARY KEY,
   display_name TEXT,
   flair TEXT,
   username TEXT,
   password TEXT,
   suspended BOOLEAN,
   login_tries BOOLEAN,
   next_login_time TIMESTAMP WITH TIME ZONE,
   role SMALLINT,
   linked BIGINT,
   has_letters BOOLEAN,
   parent BIGINT
);

CREATE TABLE IF NOT EXISTS document (
   uuid TEXT PRIMARY KEY,
   written TIMESTAMP WITH TIME ZONE,
   organisation TEXT,
   type TEXT,
   author TEXT,
   flair TEXT,
   title TEXT,
   subtitle TEXT,
   html_content TEXT,
   private BOOLEAN,
   blocked BOOLEAN,
   current_tag TEXT,
   any_p_allowed BOOLEAN,
   org_p_allowed BOOLEAN,
   info JSONB
);

CREATE TABLE IF NOT EXISTS doc_allowed (
   id BIGINT NOT NULL,
   uuid TEXT NOT NULL,
   PRIMARY KEY (id, uuid),
   FOREIGN KEY (id) REFERENCES account (id),
   FOREIGN KEY (uuid) REFERENCES document (uuid)
);

`

// ConnectDatabase establishes a connection stored in DB
// made with the Env Parameter: DB_USER, DB_PASSWORD, DB_ADDRESS and DB_NAME.
// It exists the program on error.
func ConnectDatabase() {
	var err error
	DB, err = sqlx.Connect("pgx", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_NAME")))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while connecting to postgres DB:\n"+err.Error()+"\n")
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Connection to DB established\n")

	_, err = DB.Exec(schema)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to create database schema:\n"+err.Error())
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Database schemas created\n")

	setupBasicData()
}

func setupBasicData() {
	createRootAccountIfNotExist()
	createEternatityPublicationIfNotExist()
}

func createRootAccountIfNotExist() {

}

func createEternatityPublicationIfNotExist() {

}

/*
ORDER BY column1 DESC, column2

This sorts everything by column1 (descending) first, and then by column2 (ascending, which is the default)
whenever the column1 fields for two or more rows are equal.
*/
