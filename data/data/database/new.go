package database

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
)

var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
)`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

func main() {
	// this connects & tries a simple 'SELECT 1', panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("pgx", "postgres://pgx_md5:secret@localhost:5432/pgx_test")
	if err != nil {
		log.Fatalln(err)
	}
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	tx.Commit()

	// Selects Mr. Smith from the database
	rows, err := db.NamedQuery(`SELECT * FROM person WHERE first_name=:fn`, map[string]interface{}{"fn": "John"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rows)
}

/*
create table public.doc_allowed (
  uuid text not null,
  id bigint not null default nextval('doc_allowed_id_seq'::regclass),
  primary key (uuid, id),
  foreign key (id) references public.accounts (id)
  match simple on update no action on delete no action,
  foreign key (uuid) references public.documents (uuid)
  match simple on update no action on delete no action
);

ORDER BY column1 DESC, column2

This sorts everything by column1 (descending) first, and then by column2 (ascending, which is the default)
whenever the column1 fields for two or more rows are equal.
*/
