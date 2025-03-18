package database

import (
	"github.com/lib/pq"
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

// Todo: transfer to migration
const titleTableDefinition = `CREATE TABLE title(
    name TEXT PRIMARY KEY,
    main_group TEXT NOT NULL,
    sub_group TEXT NOT NULL,
    flair TEXT NOT NULL
);
CREATE TABLE title_to_account(
    title_name TEXT NOT NULL,
    account_name TEXT NOT NULL,
    CONSTRAINT fk_organisation_name
        FOREIGN KEY(title_name) REFERENCES title(name) ON UPDATE CASCADE,
    CONSTRAINT fk_account_name
        FOREIGN KEY(account_name) REFERENCES account(name)
);`

func CreateTitle(title *Title, holderNames []string) error {
	_, err := postgresDB.Exec(`
INSERT INTO title (name, main_group, sub_group, flair) 
	VALUES ($1, $2, $3, $4);
-- Insert into connection table
INSERT INTO title_to_account (title_name, account_name) 
SELECT $1 AS title_name, name FROM account
WHERE name = ANY($5) AND blocked = false;`,
		&title.Name, &title.MainType, &title.SubType, &title.Flair,
		pq.Array(holderNames))
	if err == nil {
		addTitleToMap(title)
	}
	return err
}

func UpdateTitle(oldTitle string, title *Title) error {
	_, err := postgresDB.Exec(`
DELETE FROM title_to_account WHERE title_name = $1;
UPDATE title SET name = $2, main_group = $3, sub_group = $4, flair = $5 WHERE name = $1;`,
		&oldTitle, &title.Name, &title.MainType, &title.SubType, &title.Flair)
	if err == nil {
		removeOldTitleFromMap(oldTitle)
		addTitleToMap(title)
	}
	return err
}

func AddTitleHolder(title *Title, holderNames []string) error {
	_, err := postgresDB.Exec(`
INSERT INTO title_to_account (title_name, account_name) 
SELECT $1 AS title_name, name FROM account
WHERE name = ANY($2) AND blocked = false;`,
		&title.Name, pq.Array(holderNames))
	return err
}

func GetTitleNameList() ([]string, error) {
	result, err := postgresDB.Query(`SELECT name FROM title ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
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

func GetTitleByName(name string) (*Title, error) {
	title := &Title{}
	err := postgresDB.QueryRow(`SELECT name, main_group, sub_group, flair FROM title WHERE name = $1`,
		&name).Scan(&title.Name, &title.MainType, &title.SubType, &title.Flair)
	if err != nil {
		return nil, err
	}
	return title, err
}

func GetTitleAndHolder(name string) (*Title, []string, error) {
	title := &Title{}
	names := make([]string, 0)
	err := postgresDB.QueryRow(`SELECT name, main_group, sub_group, flair, 
       ARRAY(SELECT account_name FROM title_to_account WHERE title_name = $1 ORDER BY account_name) AS account 
       FROM title WHERE name = $1`,
		&name).Scan(&title.Name, &title.MainType, &title.SubType, &title.Flair, pq.Array(names))
	if err != nil {
		return nil, nil, err
	}
	return title, names, err
}

func GetAllTitles() ([]Title, error) {
	result, err := postgresDB.Query(`SELECT name, main_group, sub_group FROM title`)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	titles := make([]Title, 0)
	title := Title{}
	for result.Next() {
		err = result.Scan(&title.Name, &title.MainType, &title.SubType)
		if err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}
	return titles, err
}
