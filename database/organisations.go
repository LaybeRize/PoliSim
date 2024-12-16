package database

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Organisation struct {
	Name       string
	Visibility OrganisationVisibility
	MainType   string
	SubType    string
}

type OrganisationVisibility string

const (
	PUBLIC  OrganisationVisibility = "public"
	PRIVATE OrganisationVisibility = "private"
	SECRET  OrganisationVisibility = "secret"
	HIDDEN  OrganisationVisibility = "hidden"

	CON_USER  = "USER"
	CON_ADMIN = "ADMIN"

	DB_ORG_NAME       = "name"
	DB_ORG_VISIBILITY = "visibility"
	DB_ORG_MAIN_TYPE  = "main_type"
	DB_ORG_SUB_TYPE   = "sub_type"
)

func CreateOrganisation(org Organisation) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:Organisation {`+DB_ORG_NAME+`: $name , `+DB_ORG_VISIBILITY+`: $visibility , `+
			DB_ORG_MAIN_TYPE+`: $maintype , `+DB_ORG_SUB_TYPE+`: $subtype });`,
		map[string]any{"name": org.Name,
			"visibility": org.Visibility,
			"maintype":   org.MainType,
			"subtype":    org.SubType}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func CreateUserConnection(orgName string, userName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, "MATCH (a:Account), (o:Organisation) WHERE a."+DB_ACC_NAME+" = $user_name AND o."+
		DB_ORG_NAME+" = $org_name CREATE (a)-[:"+CON_USER+"]->(o);",
		map[string]any{"org_name": orgName,
			"user_name": userName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func CreateAdminConnection(orgName string, userName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, "MATCH (a:Account), (o:Organisation) WHERE a."+DB_ACC_NAME+" = $user_name AND o."+
		DB_ORG_NAME+" = $org_name CREATE (a)-[:"+CON_ADMIN+"]->(o);",
		map[string]any{"org_name": orgName,
			"user_name": userName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func DeleteAllConnectionsToOrganisation(orgName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (o:Organisation) WHERE o.`+DB_ORG_NAME+` = $org_name MATCH (a:Account)-[r]->(o) DELETE r;`,
		map[string]any{"org_name": orgName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func DeleteUserConnectionsToOrganisation(orgName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (o:Organisation) WHERE o.`+DB_ORG_NAME+` = $org_name MATCH (a:Account)-[r:`+CON_USER+`]->(o) DELETE r;`,
		map[string]any{"org_name": orgName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func GetOrganisationsForUserView(name string) ([]Organisation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CALL { MATCH (a:Account) WHERE a.`+DB_ACC_NAME+` = $name MATCH (a)-[:`+CON_USER+`|`+CON_ADMIN+`|`+CON_OWNER+`*1..2]->(o:Organisation) 
RETURN o UNION MATCH (o:Organisation) WHERE o.`+DB_ORG_VISIBILITY+` = "`+string(PUBLIC)+`" OR o.`+DB_ORG_VISIBILITY+` = "`+string(PRIVATE)+`" RETURN o 
} RETURN o ORDER BY o.`+DB_ORG_MAIN_TYPE+`, o.`+DB_ORG_SUB_TYPE+`, o.`+DB_ORG_NAME+`;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	return getArrayOfOrganisations("o", result.Records), err
}

func GetAllVisibleOrganisations() ([]Organisation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (o:Organisation) WHERE o.`+DB_ORG_VISIBILITY+` != "`+string(HIDDEN)+`" RETURN o 
ORDER BY o.`+DB_ORG_MAIN_TYPE+`, o.`+DB_ORG_SUB_TYPE+`, o.`+DB_ORG_NAME+`;`,
		nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	return getArrayOfOrganisations("o", result.Records), err
}

func GetAllInvisibleOrganisations() ([]Organisation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (o:Organisation) WHERE o.`+DB_ORG_VISIBILITY+` = "`+string(HIDDEN)+`" RETURN o 
ORDER BY o.`+DB_ORG_MAIN_TYPE+`, o.`+DB_ORG_SUB_TYPE+`, o.`+DB_ORG_NAME+`;`,
		nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
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
			Name:       node.Props[DB_ORG_NAME].(string),
			Visibility: node.Props[DB_ORG_VISIBILITY].(OrganisationVisibility),
			MainType:   node.Props[DB_ORG_MAIN_TYPE].(string),
			SubType:    node.Props[DB_ORG_SUB_TYPE].(string),
		})
	}
	return arr
}
