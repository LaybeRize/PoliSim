package database

import (
	"PoliSim/helper"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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

func CreateArticle(article *NewspaperArticle, newspaperName string) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `MATCH (t:Newspaper)-[:PUBLISHED]->(p:Publication) WHERE t.name = $newspaper 
AND p.special = false AND p.published = false RETURN p;`,
		map[string]any{"newspaper": newspaperName})
	if result.Record() == nil || result.Peek(ctx) || err != nil {
		_ = tx.Rollback(ctx)
		return notFoundError
	}
	id := result.Record().Values[0].(string)

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
			"written":   article.Written,
			"rawbody":   article.RawBody,
			"Body":      article.Body})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
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
