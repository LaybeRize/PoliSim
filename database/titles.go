package database

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

type Title struct {
	Name     string
	MainType string
	SubType  string
	Flair    string
}

func (t *Title) Exists() bool {
	return t != nil
}

func CreateTitle(title *Title, holderNames []string) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx,
		`CREATE (:Title {name: $name , main_type: $maintype , 
sub_type: $subtype , flair: $flair});`, map[string]any{
			"name":     title.Name,
			"maintype": title.MainType,
			"subtype":  title.SubType,
			"flair":    title.Flair})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `MATCH (a:Account), (t:Title) WHERE a.name IN $names  
AND t.name = $title CREATE (a)-[:HAS]->(t);`, map[string]any{
		"title": title.Name,
		"names": holderNames})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func UpdateTitle(oldtitle string, title *Title, holderNames []string) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx,
		`MATCH (t:Title) WHERE name = $oldName 
SET name = $name , main_type = $maintype , 
sub_type = $subtype , flair = $flair;`,
		map[string]any{
			"oldName":  oldtitle,
			"name":     title.Name,
			"maintype": title.MainType,
			"subtype":  title.SubType,
			"flair":    title.Flair})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `MATCH (:Account)-[r]->(t:Title) WHERE t.name = $title 
DELETE r;`, map[string]any{"title": title.Name})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `MATCH (a:Account), (t:Title) WHERE a.name IN $names  
AND t.name = $title CREATE (a)-[:HAS]->(t);`, map[string]any{
		"title": title.Name,
		"names": holderNames})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func GetTitleNameList() ([]string, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Title) RETURN t.name AS name;`,
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

func GetTitleByName(name string) (*Title, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Title) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return getSingleTitle("t", result.Records)
}

func GetTitleAndHolder(name string) (*Title, []string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Title) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, nil, err
	}
	title, err := getSingleTitle("t", result.Records)
	if err != nil {
		return title, nil, err
	}
	result, err = neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account)-[:HAS]->(t:Title) 
WHERE t.name != $name RETURN a.name AS name;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return title, nil, err
	}
	names := make([]string, len(result.Records))
	for i, record := range result.Records {
		names[i] = record.Values[0].(string)
	}
	return title, names, err
}

func getSingleTitle(letter string, records []*neo4j.Record) (*Title, error) {
	if len(records) == 0 {
		return nil, notFoundError
	} else if len(records) > 1 {
		return nil, multipleItemsError
	}
	result, exists := records[0].Get(letter)
	if !exists || result == nil {
		return nil, notFoundError
	}
	node := result.(neo4j.Node)
	title := &Title{
		Name:     node.Props["name"].(string),
		MainType: node.Props["main_type"].(string),
		SubType:  node.Props["sub_type"].(string),
		Flair:    node.Props["flair"].(string),
	}

	return title, nil
}
