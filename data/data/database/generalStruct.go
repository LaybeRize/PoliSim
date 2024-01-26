package database

import "time"

type Account struct {
	ID            int64     `db:"id"`
	DisplayName   string    `db:"display_name"`
	Flair         string    `db:"flair"`
	Username      string    `db:"username"`
	Password      string    `db:"password"`
	Suspended     bool      `db:"suspended"`
	LoginTries    bool      `db:"login_tries"`
	NextLoginTime time.Time `db:"next_login_time"`
	Role          int8      `db:"role"`
	Linked        int64     `db:"linked"`
	HasLetters    bool      `db:"has_letters"`
	Parent        int64     `db:"parent"`
}

type Document struct {
	UUID                      string       `db:"uuid"`
	Date                      time.Time    `db:"written"`
	Organisation              string       `db:"organisation"`
	Type                      string       `db:"type"`
	Author                    string       `db:"author"`
	Flair                     string       `db:"flair"`
	Title                     string       `db:"title"`
	Subtitle                  string       `db:"subtitle"`
	HTML                      string       `db:"html_content"`
	Private                   bool         `db:"private"`
	Blocked                   bool         `db:"blocked"`
	CurrentPostTag            string       `db:"current_tag"`
	AnyPosterAllowed          bool         `db:"any_p_allowed"`
	OrganisationPosterAllowed bool         `db:"org_p_allowed"`
	DocumentInfo              DocumentInfo `db:"info"`
}
