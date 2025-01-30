package database

import (
	"PoliSim/helper"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"html/template"
	"log/slog"
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
		return n.Written.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return n.Written.Format("2006-01-02 15:04:05 MST")
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
		return n.PublishedDate.In(a.TimeZone).Format("2006-01-02 15:04:05 MST")
	}
	return n.PublishedDate.Format("2006-01-02 15:04:05 MST")
}

func CreateNewspaper(newspaper *Newspaper) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(
		`CREATE (:Newspaper {name: $name})-[:PUBLISHED]->(:Publication {id: $id, special: $special, 
published: $published, published_date: $publishedDate});`, map[string]any{
			"name":          newspaper.Name,
			"id":            helper.GetUniqueID(newspaper.Name),
			"special":       false,
			"published":     false,
			"publishedDate": time.Now().UTC()})
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func GetFullNewspaperInfo(name string) (*Newspaper, error) {
	result, err := makeRequest(`MATCH (t:Newspaper) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name})
	if err != nil || len(result) != 1 {
		return nil, notFoundError
	}

	newspaper := &Newspaper{Name: name}

	result, err = makeRequest(`MATCH (a:Account)-[:AUTHOR]->(t:Newspaper) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": name})
	if err != nil {
		return nil, err
	}

	newspaper.Authors = make([]string, len(result))
	for i, record := range result {
		newspaper.Authors[i] = record.Values[0].(string)
	}

	return newspaper, err
}

func GetNewspaperNameList() ([]string, error) {
	result, err := makeRequest(`MATCH (t:Newspaper) RETURN t.name AS name;`, nil)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(result))
	for i, record := range result {
		names[i] = record.Values[0].(string)
	}
	return names, err
}

func GetNewspaperNameListForAccount(name string) ([]string, error) {
	result, err := makeRequest(`MATCH (a:Account)-[r:AUTHOR]->(t:Newspaper) 
WHERE a.name = $name RETURN t.name AS name;`,
		map[string]any{"name": name})
	if err != nil {
		return nil, err
	}

	names := make([]string, len(result))
	for i, record := range result {
		names[i] = record.Values[0].(string)
	}
	return names, err
}

func RemoveAccountsFromNewspaper(newspaper *Newspaper) error {
	_, err := makeRequest(`MATCH (a:Account)-[r:AUTHOR]->(t:Newspaper) 
WHERE t.name = $newspaper DELETE r;`, map[string]any{
		"newspaper": newspaper.Name})
	return err
}

func UpdateNewspaper(newspaper *Newspaper) error {
	_, err := makeRequest(`MATCH (a:Account), (t:Newspaper) 
WHERE a.name IN $names AND a.blocked = false AND t.name = $newspaper 
MERGE (a)-[:AUTHOR]->(t);`, map[string]any{
		"newspaper": newspaper.Name,
		"names":     newspaper.Authors})
	return err
}

func CheckIfUserAllowedInNewspaper(acc *Account, author string, newspaper string) (bool, error) {
	var result []*neo4j.Record
	var err error
	if acc.Name == author {
		result, err = makeRequest(`MATCH (a:Account)-[:AUTHOR]->(t:Newspaper) 
WHERE t.name = $newspaper AND a.name = $author RETURN a, t;`, map[string]any{
			"newspaper": newspaper,
			"author":    author})
	} else {
		result, err = makeRequest(`
MATCH (b:Account)-[:OWNER]->(a:Account)-[:AUTHOR]->(t:Newspaper) 
WHERE t.name = $newspaper AND a.name = $author AND b.name = $owner 
RETURN b, a, t;`, map[string]any{
			"newspaper": newspaper,
			"author":    author,
			"owner":     acc.Name})
	}
	return len(result) == 1, err
}

func CreateArticle(article *NewspaperArticle, special bool, newspaperName string) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}

	var id string
	result, err := tx.Run(`MATCH (t:Newspaper)-[:PUBLISHED]->(p:Publication) WHERE t.name = $newspaper 
AND p.special = $special AND p.published = false RETURN p.id;`,
		map[string]any{"newspaper": newspaperName, "special": special})

	if err != nil {
		return err
	} else if !result.Next() && special {
		id, err = createSpecialPublication(tx, newspaperName)
		if err != nil {
			return err
		}
	} else if result.Record() == nil {
		return notFoundError
	} else {
		id = result.Record().Values[0].(string)
	}

	result, err = tx.Run(`MATCH (acc:Account) WHERE acc.name = $Author AND acc.blocked = false 
RETURN acc;`,
		map[string]any{"Author": article.Author})
	if err != nil {
		return err
	} else if !result.Next() {
		return notAllowedError
	}

	err = tx.RunWithoutResult(
		`MATCH (p:Publication) WHERE p.id = $id
MATCH (acc:Account) WHERE acc.name = $Author
CREATE (a:Article {id: $articleID, title: $title , subtitle: $subtitle , author: $Author , flair: $Flair, 
written: $written , raw_body: $rawbody , body: $Body})
MERGE (a)-[:IN]->(p) 
MERGE (acc)-[:WRITTEN]->(a);`, map[string]any{
			"id":        id,
			"articleID": helper.GetUniqueID(article.Author),
			"title":     article.Title,
			"subtitle":  article.Subtitle,
			"Author":    article.Author,
			"Flair":     article.Flair,
			"written":   time.Now().UTC(),
			"rawbody":   article.RawBody,
			"Body":      article.Body})
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func createSpecialPublication(tx *dbTransaction, name string) (string, error) {
	result, err := tx.Run(`MATCH (t:Newspaper) WHERE t.name = $newspaper 
RETURN t;`,
		map[string]any{"newspaper": name})
	if !result.Next() || err != nil {
		return "", notFoundError
	}
	id := helper.GetUniqueID(name)
	err = tx.RunWithoutResult(`MATCH (n:Newspaper) WHERE n.name = $name
CREATE (p:Publication {id: $id, special: $special, 
published: $published, published_date: $publishedDate}) 
MERGE (n)-[:PUBLISHED]->(p);`, map[string]any{
		"name":          name,
		"id":            id,
		"special":       true,
		"published":     false,
		"publishedDate": time.Now().UTC()})
	if err != nil {
		return "", err
	}
	return id, err
}

func PublishPublication(id string) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}

	result, err := tx.Run(`MATCH (n:Newspaper)-[:PUBLISHED]->(p:Publication)<-[:IN]-(:Article)
WHERE p.id = $id AND p.published = false SET p.published = true, 
 p.published_date = $publishedDate RETURN p.special, n.name;`,
		map[string]any{"id": id, "publishedDate": time.Now().UTC()})
	if !result.Next() || err != nil {
		return notFoundError
	}

	if list := result.Record().Values; !list[0].(bool) {
		name := list[1].(string)
		err = tx.RunWithoutResult(
			`MATCH (n:Newspaper) WHERE n.name = $name
CREATE (p:Publication {id: $id, special: $special, 
published: $published, published_date: $publishedDate}) 
MERGE (n)-[:PUBLISHED]->(p);`, map[string]any{
				"name":          name,
				"id":            helper.GetUniqueID(name),
				"special":       false,
				"published":     false,
				"publishedDate": time.Now().UTC()})
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func GetPublicationForUser(id string, isAdmin bool) (bool, error) {
	result, err := makeRequest(`MATCH (p:Publication) 
WHERE p.id = $id AND (p.published = true OR $admin = true) RETURN p;`, map[string]any{
		"id":    id,
		"admin": isAdmin})
	return len(result) == 1, err
}

func GetPublication(id string) (*Publication, []NewspaperArticle, error) {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return nil, nil, err
	}
	result, err := tx.Run(`MATCH (t:Newspaper)-[:PUBLISHED]->(p:Publication) 
WHERE p.id = $id 
RETURN t, p;`, map[string]any{
		"id": id})
	if !result.Next() || err != nil {
		slog.Debug("", "Error", err, "ID", id)
		return nil, nil, notFoundError
	}
	pub := getArrayOfPublications(1, 0, []*neo4j.Record{result.Record()})[0]

	result, err = tx.Run(
		`MATCH (a:Article)-[:IN]->(p:Publication) 
WHERE p.id = $id 
RETURN a;`, map[string]any{
			"id": id})
	if err != nil {
		return nil, nil, err
	}

	results := make([]*neo4j.Record, 0)
	for result.Next() {
		results = append(results, result.Record())
	}

	err = tx.Commit()
	return &pub, getArrayOfArticles(0, results), err
}

func GetUnpublishedPublications() ([]Publication, error) {
	result, err := makeRequest(`MATCH (t:Newspaper)-[:PUBLISHED]->(p:Publication) 
WHERE p.published = false RETURN p, t  ORDER BY p.special DESC, p.published_date;`, nil)
	if err != nil {
		return nil, err
	}
	return getArrayOfPublications(0, 1, result), err
}

func GetPublishedNewspaper(amount int, page int, newspaper string) ([]Publication, error) {
	result, err := makeRequest(`MATCH (t:Newspaper)-[:PUBLISHED]->(p:Publication) 
WHERE t.name CONTAINS $newspaper AND p.published = true 
RETURN t, p ORDER BY p.published_date DESC SKIP $skip LIMIT $amount;`,
		map[string]any{
			"amount":    amount,
			"skip":      (page - 1) * amount,
			"newspaper": newspaper})
	return getArrayOfPublications(1, 0, result), err
}

func getArrayOfPublications(pubPos int, newsPos int, records []*neo4j.Record) []Publication {
	arr := make([]Publication, 0, len(records))
	for _, record := range records {
		pubProps := GetPropsMapForRecordPosition(record, pubPos)
		if pubProps == nil {
			continue
		}
		newsProps := GetPropsMapForRecordPosition(record, newsPos)
		if newsProps == nil {
			continue
		}
		arr = append(arr, Publication{
			NewspaperName: newsProps.GetString("name"),
			ID:            pubProps.GetString("id"),
			Special:       pubProps.GetBool("special"),
			Published:     pubProps.GetBool("published"),
			PublishedDate: pubProps.GetTime("published_date"),
		})
	}
	return arr
}

func getArrayOfArticles(pos int, records []*neo4j.Record) []NewspaperArticle {
	arr := make([]NewspaperArticle, 0, len(records))
	for _, record := range records {
		props := GetPropsMapForRecordPosition(record, pos)
		if props == nil {
			continue
		}
		arr = append(arr, NewspaperArticle{
			ID:       props.GetString("id"),
			Title:    props.GetString("title"),
			Subtitle: props.GetString("subtitle"),
			Author:   props.GetString("author"),
			Flair:    props.GetString("flair"),
			Body:     template.HTML(props.GetString("body")),
			Written:  props.GetTime("written"),
		})
	}
	return arr
}

type ArticleRejectionTransaction struct {
	tx            *dbTransaction
	NewspaperName string
	Article       NewspaperArticle
}

func RejectableArticle(id string) (*ArticleRejectionTransaction, error) {
	reject := &ArticleRejectionTransaction{}
	var err error
	reject.tx, err = openTransaction()
	if err != nil {
		return nil, err
	}

	result, err := reject.tx.Run(`
MATCH (a:Article)-[:IN]->(p:Publication)<-[:PUBLISHED]-(n:Newspaper) 
WHERE a.id = $id AND p.published = false RETURN a, n.name;`, map[string]any{

		"id": id})
	if err != nil {
		reject.tx.Close()
		return nil, err
	} else if !result.Next() {
		reject.tx.Close()
		return nil, notFoundError
	}

	reject.NewspaperName = result.Record().Values[1].(string)
	reject.Article = getArrayOfArticles(0, []*neo4j.Record{result.Record()})[0]

	return reject, nil
}

func (a *ArticleRejectionTransaction) DeleteArticle() error {
	result, err := a.tx.Run(`MATCH (a:Article)-[:IN]->(:Publication)<-[:IN]-(r:Article) 
WHERE a.id = $id 
RETURN r;`, map[string]any{"id": a.Article.ID})
	if err != nil {
		a.tx.Close()
		return err
	} else if result.Next() {
		err = a.tx.RunWithoutResult(`MATCH (a:Article) 
WHERE a.id = $id 
DETACH DELETE a;`, map[string]any{"id": a.Article.ID})
		if err != nil {
			a.tx.Close()
			return err
		}
	} else {
		err = a.tx.RunWithoutResult(`MATCH (a:Article) WHERE a.id = $id 
OPTIONAL MATCH (a)-[:IN]->(p:Publication) WHERE p.special = true
DETACH DELETE a 
DETACH DELETE p;`, map[string]any{"id": a.Article.ID})
		if err != nil {
			a.tx.Close()
			return err
		}
	}
	return nil
}

func (a *ArticleRejectionTransaction) CreateLetter(letter *Letter) error {
	defer a.tx.Close()

	err := a.tx.RunWithoutResult(letterCreation, letter.GetCreationMap())
	if err != nil {
		return err
	}

	err = a.tx.RunWithoutResult(letterLinkage, letter.GetCreationMap())
	if err != nil {
		return err
	}

	err = a.tx.Commit()
	return err
}
