package database

import (
	"PoliSim/helper"
	"html/template"
	"time"
)

type NewspaperArticle struct {
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
		`CREATE (:Newspaper {name: $name});`, map[string]any{
			"name": newspaper.Name})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Run(ctx,
		`MATCH (t:Newspaper) WHERE t.name = $newspaper 
CREATE (:Publication {id: $id, special: $special, 
published: $published, published_date: $publishedDate})-[:PUBLISHED]->(t);`,
		map[string]any{
			"newspaper":     newspaper.Name,
			"id":            helper.GetUniqueID(newspaper.Name),
			"special":       false,
			"published":     false,
			"publishedDate": time.Now().UTC()})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Run(ctx, `MATCH (a:Account), (t:Newspaper) WHERE a.name IN $names  
AND t.name = $newspaper CREATE (a)-[:AUTHOR]->(t);`, map[string]any{
		"newspaper": newspaper.Name,
		"names":     newspaper.Authors})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func CreateArticle(article *NewspaperArticle, newspaperName string) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `MATCH (p:Publication)-[:PUBLISHED]->(t:Newspaper) WHERE t.name = $newspaper 
AND p.special = false AND p.published = false RETURN p;`,
		map[string]any{"newspaper": newspaperName})
	if result.Record() == nil || result.Peek(ctx) || err != nil {
		_ = tx.Rollback(ctx)
		return notFoundError
	}

	_, err = tx.Run(ctx,
		`CREATE (:Article {title: $title , subtitle: $subtitle , author: $Author , flair: $Flair, 
written: $written , raw_body: $rawbody , body: $Body});`, map[string]any{
			"title":    article.Title,
			"subtitle": article.Subtitle,
			"Author":   article.Author,
			"Flair":    article.Flair,
			"written":  article.Written,
			"rawbody":  article.RawBody,
			"Body":     article.Body})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}
