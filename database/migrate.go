package database

import (
	"log"
)

var version = 0

func migrate() {
	log.Println("Looking for migration status")
	_, err := postgresDB.Exec(`CREATE TABLE IF NOT EXISTS version_management (
    version INTEGER PRIMARY KEY
    )`)
	if err != nil {
		log.Fatalf("Could not create version tabel to migrate automatically: %v", err)
	}

	results, err := postgresDB.Query("SELECT version FROM version_management ORDER BY version DESC LIMIT 1;")
	if err != nil {
		log.Fatalf("Could not read on the version management table: %v", err)
	}
	for results.Next() {
		if results.Err() != nil {
			log.Fatalf("Could not read on the version management table row: %v", err)
		}
		err = results.Scan(&version)
		if err != nil {
			log.Fatalf("Could not scan the version management row: %v", err)
		}
	}

	switch version {
	case 0:
		migrateToCurrentVersion()
	}

}

func migrateToCurrentVersion() {
	const currVersion = 1
	log.Println("Setting up the database for current version ", currVersion)
	var err error

	// Todo create the chat tables

	version = currVersion
	_, err = postgresDB.Exec(`INSERT INTO version_management (version) VALUES ($1)`, &version)
	if err != nil {
		log.Fatalf("Could not save the information that the current version is now %d: %v", currVersion, err)
	}
}
