package database

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Organisation struct {
	Name       string
	Visibility OrganisationVisibility
	MainType   string
	SubType    string
	Flair      string
}

type OrganisationVisibility string

const (
	PUBLIC  OrganisationVisibility = "public"
	PRIVATE OrganisationVisibility = "private"
	SECRET  OrganisationVisibility = "secret"
	HIDDEN  OrganisationVisibility = "hidden"
)

func CreateOrganisation(org *Organisation) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:Organisation {name: $name , visibility: $visibility , main_type: $maintype , 
sub_type: $subtype , flair: $flair});`,
		map[string]any{"name": org.Name,
			"visibility": org.Visibility,
			"maintype":   org.MainType,
			"subtype":    org.SubType,
			"flair":      org.Flair}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func CreateUserConnection(orgName string, userName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account), (o:Organisation) WHERE a.name = $user_name 
AND o.name = $org_name CREATE (a)-[:USER]->(o);`,
		map[string]any{"org_name": orgName,
			"user_name": userName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func CreateAdminConnection(orgName string, userName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account), (o:Organisation) WHERE a.name = $user_name 
AND o.name = $org_name CREATE (a)-[:ADMIN]->(o);`,
		map[string]any{"org_name": orgName,
			"user_name": userName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func DeleteAllConnectionsToOrganisation(orgName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (o:Organisation) WHERE o.name = $org_name 
MATCH (a:Account)-[r]->(o) DELETE r;`,
		map[string]any{"org_name": orgName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func DeleteUserConnectionsToOrganisation(orgName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (o:Organisation) WHERE o.name = $org_name 
MATCH (a:Account)-[r:USER]->(o) DELETE r;`,
		map[string]any{"org_name": orgName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func GetOrganisationsForUserView(name string) ([]Organisation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CALL { MATCH (a:Account) WHERE a.name = $name MATCH (a)-[:USER|ADMIN|OWNER*1..2]->(org:Organisation) 
		RETURN o UNION MATCH (o:Organisation) WHERE o.visibility = $public OR o.visibility = $private RETURN o 
		} RETURN o ORDER BY o.main_type, o.sub_type, o.name;`,
		map[string]any{"name": name,
			"private": PRIVATE,
			"public":  PUBLIC}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	return getArrayOfOrganisations("o", result.Records), err
}

func GetAllVisibleOrganisations() ([]Organisation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (o:Organisation) WHERE o.visibility != $hidden RETURN o 
ORDER BY o.main_type, o.sub_type, o.name;`,
		map[string]any{
			"hidden": HIDDEN,
		}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	return getArrayOfOrganisations("o", result.Records), err
}

func GetAllInvisibleOrganisations() ([]Organisation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (o:Organisation) WHERE o.visibility = $hidden RETURN o 
ORDER BY o.main_type, o.sub_type, o.name;`,
		map[string]any{
			"hidden": HIDDEN,
		}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	return getArrayOfOrganisations("o", result.Records), err
}

func getArrayOfOrganisations(letter string, records []*neo4j.Record) []Organisation {
	arr := make([]Organisation, 0, len(records))
	for _, record := range records {
		result, exists := record.Get(letter)
		if !exists {
			continue
		}
		node := result.(neo4j.Node)
		arr = append(arr, Organisation{
			Name:       node.Props["name"].(string),
			Visibility: OrganisationVisibility(node.Props["visibility"].(string)),
			MainType:   node.Props["main_type"].(string),
			SubType:    node.Props["sub_type"].(string),
			Flair:      node.Props["flair"].(string),
		})
	}
	return arr
}
