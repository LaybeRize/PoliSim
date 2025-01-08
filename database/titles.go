package database

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"slices"
)

type Title struct {
	Name     string
	MainType string
	SubType  string
	Flair    string
}

func (t *Title) Exists() bool {
	return t != nil
}

var TitleMap = make(map[string]map[string][]string)

func init() {
	titles, err := GetAllTitles()
	if err != nil {
		log.Fatalf("Titel query error: %v", err)
	}
	for _, title := range titles {
		addTitleToMap(&title)
	}
}

func addTitleToMap(title *Title) {
	if _, exists := TitleMap[title.MainType]; !exists {
		TitleMap[title.MainType] = make(map[string][]string)
	}
	if _, exists := TitleMap[title.MainType][title.SubType]; !exists {
		TitleMap[title.MainType][title.SubType] = make([]string, 0)
	}
	TitleMap[title.MainType][title.SubType] = append(TitleMap[title.MainType][title.SubType], title.Name)
	slices.Sort(TitleMap[title.MainType][title.SubType])
}

func removeOldTitleFromMap(titleName string) {
	for mainKey, mainMap := range TitleMap {
		for subKey, subMap := range mainMap {
			for index, str := range subMap {
				if titleName == str {
					TitleMap[mainKey][subKey] = append(TitleMap[mainKey][subKey][:index],
						TitleMap[mainKey][subKey][index+1:]...)
					return
				}
			}
		}
	}
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
	if err == nil {
		addTitleToMap(title)
	}
	return err
}

func UpdateTitle(oldtitle string, title *Title) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx,
		`MATCH (t:Title) WHERE t.name = $oldName 
SET t.name = $name , t.main_type = $maintype , 
t.sub_type = $subtype , t.flair = $flair;`,
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
	_, err = tx.Run(ctx, `MATCH (a:Account)-[r:HAS]->(t:Title) WHERE t.name = $title 
DELETE r;`, map[string]any{"title": title.Name})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	if err == nil {
		removeOldTitleFromMap(oldtitle)
		addTitleToMap(title)
	}
	return err
}

func AddTitleHolder(title *Title, holderNames []string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account), (t:Title) WHERE a.name IN $names  
AND t.name = $title MERGE (a)-[r:HAS]->(t);`,
		map[string]any{"title": title.Name, "names": holderNames}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
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
WHERE t.name = $name RETURN a.name AS name;`,
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

func GetAllTitles() ([]Title, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Title) RETURN t;`,
		nil, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	titles := make([]Title, len(queryResult.Records))
	for i, record := range queryResult.Records {
		node := record.Values[0].(neo4j.Node)
		titles[i] = Title{
			Name:     node.Props["name"].(string),
			SubType:  node.Props["sub_type"].(string),
			MainType: node.Props["main_type"].(string),
		}
	}
	return titles, err
}
