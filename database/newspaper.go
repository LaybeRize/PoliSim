package database

import (
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"html/template"
	"time"
)

type NewspaperArticle struct {
	ID       string
	Title    string
	Subtitle string
	Author   string
	Flair    string
	Written  time.Time
	RawBody  string
	Body     template.HTML
}

func (n *NewspaperArticle) GetAuthor() string {
	if n.Flair == "" {
		return n.Author
	}
	return n.Author + "; " + n.Flair
}

func (n *NewspaperArticle) HasSubtitle() bool {
	return n.Subtitle != ""
}

func (n *NewspaperArticle) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return n.Written.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return n.Written.Format(loc.TimeFormatString)
}

type Newspaper struct {
	Name    string
	Authors []string
}

type Publication struct {
	ID            string
	NewspaperName string
	Special       bool
	Published     bool
	PublishedDate time.Time
}

func (n *Publication) GetPublishedDate(a *Account) string {
	if a.Exists() {
		return n.PublishedDate.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return n.PublishedDate.Format(loc.TimeFormatString)
}

// Todo move to migration
const newspaperCreatTableDefinition = `
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
);`

func CreateNewspaper(newspaper *Newspaper) error {
	_, err := postgresDB.Exec(`INSERT INTO newspaper (name) VALUES ($1);
INSERT INTO newspaper_publication (id, newspaper_name, special, published, publish_date) 
                            VALUES ($2, $1, false, false, $3);`, &newspaper.Name, helper.GetUniqueID(newspaper.Name), time.Now().UTC())

	return err
}

func GetFullNewspaperInfo(name string) (*Newspaper, error) {
	newspaper := &Newspaper{}
	err := postgresDB.QueryRow(`SELECT name FROM newspaper WHERE name = $1`, &name).Scan(&newspaper.Name)
	if err != nil {
		return nil, err
	}
	newspaper.Authors = make([]string, 0)
	result, err := postgresDB.Query(`SELECT account_name FROM newspaper_to_account WHERE newspaper_name = $1`, &name)
	if err != nil {
		return nil, err
	}
	accName := ""
	for result.Next() {
		err = result.Scan(&accName)
		if err != nil {
			return nil, err
		}
		newspaper.Authors = append(newspaper.Authors, accName)
	}

	return newspaper, err
}

func GetNewspaperNameList() ([]string, error) {
	result, err := postgresDB.Query(`SELECT name FROM newspaper ORDER BY name`)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	name := ""
	for result.Next() {
		err = result.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, err
}

func GetNewspaperNameListForAccount(name string) ([]string, error) {
	result, err := postgresDB.Query(`SELECT newspaper_name FROM newspaper_to_account WHERE account_name = $1`, &name)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0)
	newspaperName := ""
	for result.Next() {
		err = result.Scan(&newspaperName)
		if err != nil {
			return nil, err
		}
		names = append(names, newspaperName)
	}
	return names, err
}

func RemoveAccountsFromNewspaper(newspaper *Newspaper) error {
	_, err := postgresDB.Exec(`DELETE FROM newspaper_to_account WHERE newspaper_name = $1`, newspaper.Name)
	return err
}

func UpdateNewspaper(newspaper *Newspaper) error {
	_, err := postgresDB.Exec(`INSERT INTO newspaper_to_account (newspaper_name, account_name) 
SELECT $1 AS newspaper_name, name FROM account
WHERE name = ANY($2);`,
		&newspaper.Name, pq.Array(newspaper.Authors))
	return err
}

func CheckIfUserAllowedInNewspaper(acc *Account, author string, newspaper string) (bool, error) {
	err := postgresDB.QueryRow(`SELECT newspaper_name FROM newspaper_to_account
 INNER JOIN ownership o on newspaper_to_account.account_name = o.account_name 
 WHERE newspaper_name = $1 AND o.account_name = $2 AND owner_name = $3`, &newspaper, &author, &acc.Name).Scan(&newspaper)
	return err == nil, err
}

func CreateArticle(article *NewspaperArticle, special bool, newspaperName string) error {
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}
	defer rollback(tx)
	var id string
	err = tx.QueryRow(`SELECT id FROM newspaper_publication WHERE newspaper_name = $1 AND special = $2;`,
		&newspaperName, &special).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) && special {
		id, err = createSpecialPublication(tx, newspaperName)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	var name string
	err = tx.QueryRow(`SELECT name FROM account WHERE blocked = false AND name = $1`, &article.Author).Scan(&name)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO newspaper_article (id, title, subtitle, author, flair, html_body, raw_body, written, publication_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, helper.GetUniqueID(article.Author), &article.Title, &article.Subtitle,
		&article.Author, &article.Flair, &article.Body, &article.RawBody, time.Now().UTC(), &id)

	return tx.Commit()
}

func createSpecialPublication(tx *sql.Tx, name string) (string, error) {
	id := helper.GetUniqueID(name)
	_, err := tx.Exec(`INSERT INTO newspaper_publication (id, newspaper_name, special, published, publish_date) 
VALUES ($1, $2, true, false, $3)`, id, name, time.Now().UTC())
	if err != nil {
		return id, err
	}
	return id, nil
}

func PublishPublication(id string) error {
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}
	defer rollback(tx)
	var newspaperName string
	err = tx.QueryRow(`SELECT id FROM newspaper_article WHERE publication_id = $1 LIMIT 1;`, &id).Scan(&newspaperName)
	if err != nil {
		return err
	}
	var special bool
	err = tx.QueryRow(`UPDATE newspaper_publication 
SET published = true, publish_date = $2 WHERE id = $1 RETURNING special, newspaper_name`,
		&id, time.Now().UTC()).Scan(&special, &newspaperName)
	if err != nil {
		return err
	}
	if !special {
		_, err = tx.Exec(`INSERT INTO newspaper_publication (id, newspaper_name, special, published, publish_date) 
VALUES ($1, $2, false, false, $3)`, helper.GetUniqueID(newspaperName), newspaperName, time.Now().UTC())
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func GetPublicationForUser(id string, isAdmin bool) error {
	var pubID string
	err := postgresDB.QueryRow(`SELECT id FROM newspaper_publication 
          WHERE id = $1 AND (published = true OR $2 = true);`, &id, &isAdmin).Scan(&pubID)
	return err
}

func GetPublication(id string) (*Publication, []NewspaperArticle, error) {
	pub := &Publication{ID: id}
	err := postgresDB.QueryRow(`SELECT newspaper_name, special, published, publish_date 
FROM newspaper_publication WHERE id = $1`, &id).Scan(&pub.NewspaperName, &pub.Special, &pub.Published, &pub.PublishedDate)
	if err != nil {
		return nil, nil, err
	}
	result, err := postgresDB.Query(`SELECT id, title, subtitle, author, flair, html_body, written 
FROM newspaper_article WHERE publication_id = $1;`, &id)
	if err != nil {
		return nil, nil, err
	}
	defer closeRows(result)
	article := NewspaperArticle{}
	list := make([]NewspaperArticle, 0)
	for result.Next() {
		err = result.Scan(&article.ID, &article.Title, &article.Subtitle, &article.Author, &article.Flair,
			&article.Body, &article.Written)
		if err != nil {
			return nil, nil, err
		}
		list = append(list, article)
	}
	return pub, list, nil
}

func GetUnpublishedPublications() ([]Publication, error) {
	result, err := postgresDB.Query(`SELECT id, newspaper_name, special, publish_date 
FROM newspaper_publication WHERE published = false ORDER BY special DESC, publish_date;`)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	pub := Publication{Published: false}
	list := make([]Publication, 0)
	for result.Next() {
		err = result.Scan(&pub.ID, &pub.NewspaperName, &pub.Special, &pub.PublishedDate)
		if err != nil {
			return nil, err
		}
		list = append(list, pub)
	}
	return list, nil
}

func GetPublishedNewspaper(amount int, page int, newspaper string) ([]Publication, error) {
	result, err := postgresDB.Query(`SELECT id, newspaper_name, special, publish_date FROM newspaper_publication 
WHERE published = true AND newspaper_name LIKE '%' || $3 || '%' ORDER BY publish_date DESC OFFSET $1 LIMIT $2;`,
		(page-1)*amount, amount, newspaper)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	pub := Publication{Published: true}
	list := make([]Publication, 0)
	for result.Next() {
		err = result.Scan(&pub.ID, &pub.NewspaperName, &pub.Special, &pub.PublishedDate)
		if err != nil {
			return nil, err
		}
		list = append(list, pub)
	}
	return list, nil
}

type ArticleRejectionTransaction struct {
	tx            *sql.Tx
	NewspaperName string
	Article       NewspaperArticle
}

func RejectableArticle(id string) (*ArticleRejectionTransaction, error) {
	reject := &ArticleRejectionTransaction{}
	var err error
	reject.tx, err = postgresDB.Begin()
	if err != nil {
		return nil, err
	}

	err = reject.tx.QueryRow(`SELECT newspaper_name, newspaper_article.id, author, title, subtitle, flair, 
       html_body, raw_body, written  FROM newspaper_article 
    INNER JOIN newspaper_publication np on np.id = newspaper_article.publication_id 
         WHERE newspaper_article.id = $1 and published = false;`, &id).Scan(&reject.NewspaperName,
		&reject.Article.ID, &reject.Article.Author,
		&reject.Article.Title, &reject.Article.Subtitle, &reject.Article.Flair,
		&reject.Article.Body, &reject.Article.RawBody, &reject.Article.Written)
	if err != nil {
		_ = reject.tx.Rollback()
		return nil, err
	}
	return reject, nil
}

func (a *ArticleRejectionTransaction) DeleteArticle() error {
	var publicationID string
	err := a.tx.QueryRow(`DELETE FROM newspaper_article WHERE id = $1 RETURNING publication_id;`,
		&a.Article.ID).Scan(&publicationID)
	if err != nil {
		_ = a.tx.Rollback()
		return err
	}
	err = a.tx.QueryRow(`SELECT publication_id FROM newspaper_article WHERE publication_id = $1;`,
		&publicationID).Scan(&publicationID)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = a.tx.Exec(`DELETE FROM newspaper_publication WHERE id = $1 AND special = true;`, &publicationID)
		if err != nil {
			_ = a.tx.Rollback()
			return err
		}
	} else if err != nil {
		_ = a.tx.Rollback()
		return err
	}
	return nil
}

func (a *ArticleRejectionTransaction) CreateLetter(letter *Letter) error {
	err := createLetter(a.tx, letter)
	if err != nil {
		return err
	}
	return a.tx.Commit()
}
