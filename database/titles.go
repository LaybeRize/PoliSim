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

func CreateTitle(title *Title) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:Title {name: $name , main_type: $maintype , 
sub_type: $subtype , flair: $flair});`,
		map[string]any{"name": title.Name,
			"maintype": title.MainType,
			"subtype":  title.SubType,
			"flair":    title.Flair}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}
