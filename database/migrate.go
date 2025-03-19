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

	_, err := postgresDB.Exec(`
-- Account --
CREATE TABLE account (
 	name TEXT PRIMARY KEY,
 	username TEXT UNIQUE NOT NULL,
 	password TEXT NOT NULL,
 	role INT NOT NULL,
 	blocked BOOLEAN NOT NULL,
 	font_size INT NOT NULL,
 	time_zone TEXT NOT NULL
);
CREATE INDEX account_is_blocked ON account USING hash (blocked);
CREATE TABLE ownership (
    account_name TEXT NOT NULL,
    owner_name TEXT NOT NULL,
    CONSTRAINT fk_account_name
        FOREIGN KEY(account_name) REFERENCES account(name),
    CONSTRAINT fk_owner_name
        FOREIGN KEY(owner_name) REFERENCES account(name)
);
CREATE INDEX ownership_account_name ON ownership USING hash (account_name);
CREATE INDEX ownership_owner_name ON ownership USING hash (owner_name);
-- Colors --
CREATE TABLE colors (
    name TEXT PRIMARY KEY,
    background TEXT,
    text TEXT,
    link TEXT
);
-- Cookies --
CREATE TABLE cookies (
    session_key TEXT PRIMARY KEY,
	name TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    update_at TIMESTAMP NOT NULL,
    CONSTRAINT account_name 
        FOREIGN KEY(name) REFERENCES account(name)
);
-- Letter --
CREATE TABLE letter(
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	author TEXT NOT NULL,
	flair TEXT NOT NULL,
	signable BOOLEAN NOT NULL,
	written TIMESTAMP NOT NULL,
	body TEXT NOT NULL
);
CREATE TABLE letter_to_account(
	letter_id TEXT NOT NULL,
	account_name TEXT NOT NULL,
	has_read BOOLEAN NOT NULL,
	sign_status INT NOT NULL,
	CONSTRAINT fk_letter_id
        FOREIGN KEY(letter_id) REFERENCES letter(id),
    CONSTRAINT fk_account_name
        FOREIGN KEY(account_name) REFERENCES account(name)
);
-- Newspaper --
CREATE TABLE newspaper (
    name TEXT PRIMARY KEY
);
CREATE TABLE newspaper_publication (
    id TEXT PRIMARY KEY,
    newspaper_name TEXT NOT NULL,
    special BOOLEAN NOT NULL,
    published BOOLEAN NOT NULL,
    publish_date TIMESTAMP NOT NULL,
    CONSTRAINT fk_newspaper_name
        FOREIGN KEY (newspaper_name) REFERENCES newspaper(name)
);
CREATE TABLE newspaper_article (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    subtitle TEXT NOT NULL,
    author TEXT NOT NULL,
    flair TEXT NOT NULL,
    html_body TEXT NOT NULL,
    raw_body TEXT NOT NULL,
    written TIMESTAMP NOT NULL,
    publication_id TEXT NOT NULL,
    CONSTRAINT fk_publication_id
        FOREIGN KEY (publication_id) REFERENCES newspaper_publication(id)
);
CREATE TABLE newspaper_to_account (
    newspaper_name TEXT NOT NULL,
    account_name TEXT NOT NULL,
    CONSTRAINT fk_newspaper_name
        FOREIGN KEY (newspaper_name) REFERENCES newspaper(name),
     CONSTRAINT fk_account_name
        FOREIGN KEY (account_name) REFERENCES account(name)
);
-- Notes --
CREATE TABLE blackboard_note(
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
    author  TEXT NOT NULL,
    flair  TEXT NOT NULL,
    posted TIMESTAMP NOT NULL,
    body  TEXT NOT NULL,
	blocked BOOLEAN NOT NULL
);
CREATE TABLE blackboard_references(
	base_note_id TEXT NOT NULL,
	reference_id TEXT NOT NULL,
	CONSTRAINT fk_base_note_id
        FOREIGN KEY(base_note_id) REFERENCES blackboard_note(id),
    CONSTRAINT fk_reference_id
        FOREIGN KEY(reference_id) REFERENCES blackboard_note(id)
);
-- Organisation --
CREATE TABLE organisation(
    name TEXT PRIMARY KEY,
    main_group TEXT NOT NULL,
    sub_group TEXT NOT NULL,
    visibility INT NOT NULL,
    flair TEXT NOT NULL,
    users TEXT[] NOT NULL,
    admins TEXT[] NOT NULL
);
CREATE TABLE organisation_to_account(
    organisation_name TEXT NOT NULL,
    account_name TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL,
    CONSTRAINT fk_organisation_name
        FOREIGN KEY(organisation_name) REFERENCES organisation(name) ON UPDATE CASCADE,
    CONSTRAINT fk_account_name
        FOREIGN KEY(account_name) REFERENCES account(name)
);
CREATE VIEW organisation_linked AS
    SELECT organisation.*, ota.account_name, ota.is_admin, ownership.owner_name FROM organisation
    LEFT JOIN organisation_to_account ota ON organisation.name = ota.organisation_name
    LEFT JOIN ownership ON ota.account_name = ownership.account_name;
-- Document --
CREATE TABLE document (
    id TEXT PRIMARY KEY,
    type INT NOT NULL,
    organisation TEXT NOT NULL,
    organisation_name TEXT NOT NULL,
    title TEXT NOT NULL,
    author  TEXT NOT NULL,
    flair  TEXT NOT NULL,
    body  TEXT NOT NULL,
    written TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    public  BOOLEAN NOT NULL,
    removed BOOLEAN NOT NULL,
    member_participation BOOLEAN NOT NULL,
    admin_participation BOOLEAN NOT NULL,
    extra_info jsonb NOT NULL,
    CONSTRAINT fk_organisation_name
        FOREIGN KEY (organisation_name) REFERENCES organisation(name) ON UPDATE CASCADE
);
CREATE TABLE document_to_account (
    document_id TEXT NOT NULL,
    account_name TEXT NOT NULL,
    participant BOOLEAN NOT NULL,
    CONSTRAINT fk_document_id
        FOREIGN KEY (document_id) REFERENCES document(id),
    CONSTRAINT fk_account_name
        FOREIGN KEY (account_name) REFERENCES account(name)
);
CREATE TABLE comment_to_document (
	comment_id TEXT PRIMARY KEY,
	document_id TEXT NOT NULL,
    author  TEXT NOT NULL,
    flair  TEXT NOT NULL,
    body  TEXT NOT NULL,
    written TIMESTAMP NOT NULL,
    removed BOOLEAN NOT NULL,
    CONSTRAINT fk_document_id
        FOREIGN KEY (document_id) REFERENCES document(id)	
);
CREATE VIEW document_linked AS
SELECT id, type, organisation, doc.organisation_name, title, author, flair, body, written, 
       end_time, public, removed, member_participation, admin_participation, extra_info, 
       dta.account_name as doc_account, ota.account_name as organisation_account, is_admin, 
       dta.participant as participant, owner_name FROM document doc
	LEFT JOIN document_to_account dta ON doc.id = dta.document_id
	LEFT JOIN organisation_to_account ota ON doc.organisation_name = ota.organisation_name
	LEFT JOIN ownership own ON ota.account_name = own.account_name OR dta.account_name = own.account_name;
-- Title --
CREATE TABLE title(
    name TEXT PRIMARY KEY,
    main_group TEXT NOT NULL,
    sub_group TEXT NOT NULL,
    flair TEXT NOT NULL
);
CREATE TABLE title_to_account(
    title_name TEXT NOT NULL,
    account_name TEXT NOT NULL,
    CONSTRAINT fk_organisation_name
        FOREIGN KEY(title_name) REFERENCES title(name) ON UPDATE CASCADE,
    CONSTRAINT fk_account_name
        FOREIGN KEY(account_name) REFERENCES account(name)
);
-- Votes --
CREATE TABLE personal_votes (
	number INT NOT NULL,
	account_name TEXT NOT NULL,
	id TEXT NOT NULL,
	question TEXT NOT NULL,
	answers TEXT[] NOT NULL,
	type INT NOT NULL,
	max_votes INT NOT NULL,
	show_votes BOOLEAN NOT NULL,
	anonymous BOOLEAN NOT NULL,
	end_date TIMESTAMP NOT NULL,
	vote_info jsonb NOT NULL,
	PRIMARY KEY (number, account_name)
);
CREATE TABLE document_to_vote (
	id TEXT PRIMARY KEY,
	document_id TEXT NOT NULL,
	question TEXT NOT NULL,
	answers TEXT[] NOT NULL,
	type INT NOT NULL,
	max_votes INT NOT NULL,
	show_votes BOOLEAN NOT NULL,
	anonymous BOOLEAN NOT NULL,
	end_date TIMESTAMP NOT NULL,
	vote_info jsonb NOT NULL,
    CONSTRAINT fk_document_id
        FOREIGN KEY(document_id) REFERENCES document(id)
);
CREATE TABLE has_voted (
	account_name TEXT NOT NULL,
	vote_id TEXT NOT NULL,
    CONSTRAINT fk_account_name
        FOREIGN KEY (account_name) REFERENCES account(name),
    CONSTRAINT fk_vote_id
        FOREIGN KEY(vote_id) REFERENCES document_to_vote(id)
);
`)
	if err != nil {
		log.Fatalf("Could not create tables for the current version %d: %v", currVersion, err)
	}

	_, err = postgresDB.Exec(`INSERT INTO version_management (version) VALUES ($1)`, currVersion)
	if err != nil {
		log.Fatalf("Could not save the information that the current version is now %d: %v", currVersion, err)
	}
}
