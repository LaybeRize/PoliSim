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
	defer tx.Close()
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(
		`CREATE (:Title {name: $name , main_type: $maintype , 
sub_type: $subtype , flair: $flair});`, map[string]any{
			"name":     title.Name,
			"maintype": title.MainType,
			"subtype":  title.SubType,
			"flair":    title.Flair})
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(`MATCH (a:Account), (t:Title) WHERE a.name IN $names  
AND t.name = $title CREATE (a)-[:HAS]->(t);`, map[string]any{
		"title": title.Name,
		"names": holderNames})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err == nil {
		addTitleToMap(title)
	}
	return err
}

func UpdateTitle(oldtitle string, title *Title) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(
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
		return err
	}
	err = tx.RunWithoutResult(`MATCH (a:Account)-[r:HAS]->(t:Title) WHERE t.name = $title 
DELETE r;`, map[string]any{"title": title.Name})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err == nil {
		removeOldTitleFromMap(oldtitle)
		addTitleToMap(title)
	}
	return err
}

func AddTitleHolder(title *Title, holderNames []string) error {
	_, err := makeRequest(`MATCH (a:Account), (t:Title) 
WHERE a.name IN $names AND a.blocked = false AND t.name = $title 
MERGE (a)-[r:HAS]->(t);`,
		map[string]any{"title": title.Name, "names": holderNames})
	return err
}

func GetTitleNameList() ([]string, error) {
	result, err := makeRequest(`MATCH (t:Title) RETURN t.name AS name;`, nil)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(result))
	for i, record := range result {
		names[i] = record.Values[0].(string)
	}
	return names, err
}

func GetTitleByName(name string) (*Title, error) {
	result, err := makeRequest(`MATCH (t:Title) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name})
	if err != nil {
		return nil, err
	}
	return getSingleTitle(0, result)
}

func GetTitleAndHolder(name string) (*Title, []string, error) {
	result, err := makeRequest(`MATCH (t:Title) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name})
	if err != nil {
		return nil, nil, err
	}
	title, err := getSingleTitle(0, result)
	if err != nil {
		return title, nil, err
	}
	result, err = makeRequest(`MATCH (a:Account)-[:HAS]->(t:Title) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": name})
	if err != nil {
		return title, nil, err
	}
	names := make([]string, len(result))
	for i, record := range result {
		names[i] = record.Values[0].(string)
	}
	return title, names, err
}

func getSingleTitle(pos int, records []*neo4j.Record) (*Title, error) {
	if len(records) == 0 {
		return nil, NotFoundError
	} else if len(records) > 1 {
		return nil, MultipleItemsError
	}
	props := GetPropsMapForRecordPosition(records[0], pos)
	if props == nil {
		return nil, NotFoundError
	}
	title := &Title{
		Name:     props.GetString("name"),
		MainType: props.GetString("main_type"),
		SubType:  props.GetString("sub_type"),
		Flair:    props.GetString("flair"),
	}

	return title, nil
}

func GetAllTitles() ([]Title, error) {
	result, err := makeRequest(`MATCH (t:Title) RETURN t;`, nil)
	if err != nil {
		return nil, err
	}
	titles := make([]Title, len(result))
	for i, record := range result {
		props := GetPropsMapForRecordPosition(record, 0)
		titles[i] = Title{
			Name:     props.GetString("name"),
			MainType: props.GetString("main_type"),
			SubType:  props.GetString("sub_type"),
		}
	}
	return titles, err
}
