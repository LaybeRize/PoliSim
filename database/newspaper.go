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
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx,
		`CREATE (:Newspaper {name: $name})-[:PUBLISHED]->(:Publication {id: $id, special: $special, 
published: $published, published_date: $publishedDate});`, map[string]any{
			"name":          newspaper.Name,
			"id":            helper.GetUniqueID(newspaper.Name),
			"special":       false,
			"published":     false,
			"publishedDate": time.Now().UTC()})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	/*
			_, err = tx.Run(ctx, `MATCH (a:Account), (t:Newspaper) WHERE a.name IN $names
		AND t.name = $newspaper CREATE (a)-[:AUTHOR]->(t);`, map[string]any{
				"newspaper": newspaper.Name,
				"names":     newspaper.Authors})
			if err != nil {
				_ = tx.Rollback(ctx)
				return err
			}
	*/
	err = tx.Commit(ctx)
	return err
}

func GetFullNewspaperInfo(name string) (*Newspaper, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Newspaper) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil || len(result.Records) != 1 {
		return nil, notFoundError
	}

	newspaper := &Newspaper{Name: name}

	result, err = neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account)-[:AUTHOR]->(t:Newspaper) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	newspaper.Authors = make([]string, len(result.Records))
	for i, record := range result.Records {
		newspaper.Authors[i] = record.Values[0].(string)
	}

	return newspaper, err
}

func GetNewspaperNameList() ([]string, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Newspaper) RETURN t.name AS name;`,
		nil, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	names := make([]string, len(queryResult.Records))
	for i, record := range queryResult.Records {
		names[i] = record.Values[0].(string)
	}
	return names, err
}

func GetNewspaperNameListForAccount(name string) ([]string, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account)-[r:AUTHOR]->(t:Newspaper) 
WHERE a.name = $name RETURN t.name AS name;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	names := make([]string, len(queryResult.Records))
	for i, record := range queryResult.Records {
		names[i] = record.Values[0].(string)
	}
	return names, err
}

func RemoveAccountsFromNewspaper(newspaper *Newspaper) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account)-[r:AUTHOR]->(t:Newspaper) 
WHERE t.name = $newspaper DELETE r;`, map[string]any{
		"newspaper": newspaper.Name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func UpdateNewspaper(newspaper *Newspaper) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account), (t:Newspaper) WHERE a.name IN $names
		AND t.name = $newspaper CREATE (a)-[:AUTHOR]->(t);`, map[string]any{
		"newspaper": newspaper.Name,
		"names":     newspaper.Authors}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func CheckIfUserAllowedInNewspaper(acc *Account, author string, newspaper string) (bool, error) {
	var result *neo4j.EagerResult
	var err error
	if acc.Name == author {
		result, err = neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account)-[:AUTHOR]->(t:Newspaper) 
WHERE t.name = $newspaper AND a.name = $author RETURN a, t;`, map[string]any{
			"newspaper": newspaper,
			"author":    author}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	} else {
		result, err = neo4j.ExecuteQuery(ctx, driver, `
MATCH (b:Account)-[:OWNER]->(a:Account)-[:AUTHOR]->(t:Newspaper) 
WHERE t.name = $newspaper AND a.name = $author AND b.name = $owner 
RETURN b, a, t;`, map[string]any{
			"newspaper": newspaper,
			"author":    author,
			"owner":     acc.Name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	}
	return len(result.Records) == 1, err
}

func CreateArticle(article *NewspaperArticle, special bool, newspaperName string) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	var id string
	result, err := tx.Run(ctx, `MATCH (t:Newspaper)-[:PUBLISHED]->(p:Publication) WHERE t.name = $newspaper 
AND p.special = $special AND p.published = false RETURN p.id;`,
		map[string]any{"newspaper": newspaperName, "special": special})

	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	} else if result.Next(ctx); result.Record() == nil && special {
		id, err = createSpecialPublication(tx, newspaperName)
		if err != nil {
			return err
		}
	} else if result.Record() == nil {
		return notFoundError
	} else {
		id = result.Record().Values[0].(string)
	}

	_, err = tx.Run(ctx,
		`MATCH (p:Publication) WHERE p.id = $id
CREATE (a:Article {id: $articleID, title: $title , subtitle: $subtitle , author: $Author , flair: $Flair, 
written: $written , raw_body: $rawbody , body: $Body})
MERGE (a)-[:IN]->(p);`, map[string]any{
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
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func createSpecialPublication(tx neo4j.ExplicitTransaction, name string) (string, error) {
	result, err := tx.Run(ctx, `MATCH (t:Newspaper) WHERE t.name = $newspaper 
RETURN t;`,
		map[string]any{"newspaper": name})
	if result.Next(ctx); result.Record() == nil || err != nil {
		_ = tx.Rollback(ctx)
		return "", notFoundError
	}
	id := helper.GetUniqueID(name)
	_, err = tx.Run(ctx,
		`MATCH (n:Newspaper) WHERE n.name = $name
CREATE (p:Publication {id: $id, special: $special, 
published: $published, published_date: $publishedDate}) 
MERGE (n)-[:PUBLISHED]->(p);`, map[string]any{
			"name":          name,
			"id":            id,
			"special":       true,
			"published":     false,
			"publishedDate": time.Now().UTC()})
	if err != nil {
		_ = tx.Rollback(ctx)
		return "", err
	}
	return id, err
}

func PublishPublication(id string) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `MATCH (n:Newspaper)-[:PUBLISHED]->(p:Publication) 
WHERE p.id = $id SET p.published = true, 
 p.published_date = $publishedDate RETURN p.special, n.name;`,
		map[string]any{"id": id, "publishedDate": time.Now().UTC()})
	if result.Next(ctx); result.Record() == nil || err != nil {
		_ = tx.Rollback(ctx)
		return notFoundError
	}

	if list := result.Record().Values; !list[0].(bool) {
		name := list[1].(string)
		_, err = tx.Run(ctx,
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
			_ = tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	return err
}

func GetPublicationForUser(id string, isAdmin bool) (bool, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (p:Publication) 
WHERE p.id = $id AND (p.published = true OR $admin = true) RETURN p;`, map[string]any{
		"id":    id,
		"admin": isAdmin}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return len(result.Records) == 1, err
}

func GetPublication(id string) (*Publication, []NewspaperArticle, error) {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return nil, nil, err
	}
	result, err := tx.Run(ctx,
		`MATCH (t:Newspaper)-[:PUBLISHED]->(p:Publication) 
WHERE p.id = $id 
RETURN t, p;`, map[string]any{
			"id": id})
	if result.Next(ctx); err != nil || result.Record() == nil {
		slog.Debug("", "Error", err, "ID", id)
		_ = tx.Rollback(ctx)
		return nil, nil, notFoundError
	}
	pub := getArrayOfPublications("p", "t", []*neo4j.Record{result.Record()})[0]
	result, err = tx.Run(ctx,
		`MATCH (a:Article)-[:IN]->(p:Publication) 
WHERE p.id = $id 
RETURN a;`, map[string]any{
			"id": id})
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, nil, err
	}

	results := make([]*neo4j.Record, 0)
	for result.Next(ctx) {
		results = append(results, result.Record())
	}

	err = tx.Commit(ctx)
	return &pub, getArrayOfArticles("a", results), err
}

func GetUnpublishedPublications() ([]Publication, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Newspaper)-[:PUBLISHED]->(p:Publication) 
WHERE p.published = false RETURN p, t  ORDER BY p.special, p.published_date;`,
		nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return getArrayOfPublications("p", "t", result.Records), err
}

func GetPublishedNewspaper(amount int, page int, newspaper string) ([]Publication, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Newspaper)-[:PUBLISHED]->(p:Publication) 
WHERE t.name CONTAINS $newspaper 
RETURN t, p ORDER BY p.published_date DESC SKIP $skip LIMIT $amount;`,
		map[string]any{
			"amount":    amount,
			"skip":      (page - 1) * amount,
			"newspaper": newspaper,
		}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return getArrayOfPublications("p", "t", result.Records), err
}

func getArrayOfPublications(pubLetter string, newsLetter string, records []*neo4j.Record) []Publication {
	arr := make([]Publication, 0, len(records))
	for _, record := range records {
		result, exists := record.Get(pubLetter)
		if !exists {
			continue
		}
		news, exists := record.Get(newsLetter)
		if !exists {
			continue
		}
		node := result.(neo4j.Node)
		arr = append(arr, Publication{
			NewspaperName: news.(neo4j.Node).Props["name"].(string),
			ID:            node.Props["id"].(string),
			Special:       node.Props["special"].(bool),
			Published:     node.Props["published"].(bool),
			PublishedDate: node.Props["published_date"].(time.Time),
		})
	}
	return arr
}

func getArrayOfArticles(letter string, records []*neo4j.Record) []NewspaperArticle {
	arr := make([]NewspaperArticle, 0, len(records))
	for _, record := range records {
		if record == nil {
			continue
		}
		result, exists := record.Get(letter)
		if !exists {
			continue
		}
		node := result.(neo4j.Node)
		arr = append(arr, NewspaperArticle{
			ID:       node.Props["id"].(string),
			Title:    node.Props["title"].(string),
			Subtitle: node.Props["subtitle"].(string),
			Author:   node.Props["author"].(string),
			Flair:    node.Props["flair"].(string),
			Body:     template.HTML(node.Props["body"].(string)),
			Written:  node.Props["written"].(time.Time),
		})
	}
	return arr
}
