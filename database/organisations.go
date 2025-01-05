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

func (o *Organisation) Exists() bool {
	return o != nil
}

func (o *Organisation) VisibilityIsValid() bool {
	return o.Visibility == PUBLIC || o.Visibility == PRIVATE ||
		o.Visibility == SECRET || o.Visibility == HIDDEN
}

func (o *Organisation) ClearInvalidFlair() {
	if o.Visibility == SECRET || o.Visibility == HIDDEN {
		o.Flair = ""
	}
}

func (o *Organisation) HasFlair() bool {
	return o.Flair != ""
}

func (o *Organisation) GetClassType() string {
	switch o.Visibility {
	case PUBLIC:
		return "bi-public"
	case PRIVATE:
		return "bi-private"
	case SECRET:
		return "bi-secret"
	}
	return ""
}

type OrganisationVisibility string

const (
	PUBLIC  OrganisationVisibility = "public"
	PRIVATE OrganisationVisibility = "private"
	SECRET  OrganisationVisibility = "secret"
	HIDDEN  OrganisationVisibility = "hidden"
)

func CreateOrganisation(org *Organisation, userNames []string, adminNames []string) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx,
		`CREATE (:Organisation {name: $name , visibility: $visibility , main_type: $maintype , 
sub_type: $subtype , flair: $flair});`,
		map[string]any{"name": org.Name,
			"visibility": org.Visibility,
			"maintype":   org.MainType,
			"subtype":    org.SubType,
			"flair":      org.Flair})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `
MATCH (o:Organisation) WHERE o.name = $org 
MATCH (a:Account) WHERE a.name IN $aNames
CREATE (a)-[:ADMIN]->(o);`, map[string]any{
		"org":    org.Name,
		"aNames": adminNames})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `
MATCH (o:Organisation) WHERE o.name = $org
MATCH (u:Account) WHERE u.name IN $uNames AND (NOT u.name IN $aNames)
CREATE (u)-[:USER]->(o);`, map[string]any{
		"org":    org.Name,
		"uNames": userNames,
		"aNames": adminNames})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func UpdateOrganisation(oldName string, org *Organisation) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx,
		`MATCH (o:Organisation) WHERE o.name = $oldName 
SET o.name = $name , o.main_type = $maintype , o.visibility = $visibility, 
o.sub_type = $subtype , o.flair = $flair;`,
		map[string]any{
			"oldName":    oldName,
			"name":       org.Name,
			"visibility": org.Visibility,
			"maintype":   org.MainType,
			"subtype":    org.SubType,
			"flair":      org.Flair})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `
MATCH (a:Account)-[r:ADMIN|USER]->(o:Organisation) WHERE o.name = $org 
DELETE r;`, map[string]any{"org": org.Name})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if org.Visibility == HIDDEN {
		_, err = tx.Run(ctx, `
MATCH (:Account)-[r:FAVORITE]->(o:Organisation) WHERE o.name = $org 
DELETE r;`, map[string]any{"org": org.Name})
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	return err
}

func AddOrganisationMember(org *Organisation, userNames []string, adminNames []string) error {
	if org.Visibility == HIDDEN {
		return nil
	}
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx, `
MATCH (o:Organisation) WHERE o.name = $org 
MATCH (a:Account) WHERE a.name IN $aNames
CREATE (a)-[:ADMIN]->(o);`, map[string]any{
		"org":    org.Name,
		"aNames": adminNames})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `
MATCH (o:Organisation) WHERE o.name = $org
MATCH (u:Account) WHERE u.name IN $uNames AND (NOT u.name IN $aNames)
CREATE (u)-[:USER]->(o);`, map[string]any{
		"org":    org.Name,
		"uNames": userNames,
		"aNames": adminNames})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func GetFullOrganisationInfo(name string) (*Organisation, []string, []string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Organisation) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, nil, nil, err
	}
	organisation, err := getSingleOrganisation("t", result.Records)
	if err != nil {
		return nil, nil, nil, err
	}
	result, err = neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account)-[:USER]->(t:Organisation) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, nil, nil, err
	}
	user := make([]string, len(result.Records))
	for i, record := range result.Records {
		user[i] = record.Values[0].(string)
	}
	result, err = neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account)-[:ADMIN]->(t:Organisation) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, nil, nil, err
	}
	admin := make([]string, len(result.Records))
	for i, record := range result.Records {
		admin[i] = record.Values[0].(string)
	}
	return organisation, user, admin, err
}

func GetOrganisationByName(name string) (*Organisation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Organisation) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return getSingleOrganisation("t", result.Records)
}

func GetOrganisationNameList() ([]string, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Organisation) RETURN t.name AS name;`,
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

func GetOrganisationsForUserView(account *Account) ([]Organisation, error) {
	name := ""
	if account.Exists() {
		name = account.Name
	}
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CALL { 
MATCH (a:Account)-[:USER|ADMIN|OWNER*1..2]->(o:Organisation) WHERE a.name = $name RETURN o 
UNION 
MATCH (o:Organisation) WHERE o.visibility = $public OR o.visibility = $private RETURN o 
		} RETURN o ORDER BY o.main_type, o.sub_type, o.name;`,
		map[string]any{"name": name,
			"private": PRIVATE,
			"public":  PUBLIC}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	return getArrayOfOrganisations("o", result.Records), err
}

func GetFullOrganisationInfoForUserView(account *Account, orgName string) (*Organisation, []string, []string, error) {
	name := ""
	if account.Exists() {
		name = account.Name
	}
	result, err := neo4j.ExecuteQuery(ctx, driver, `CALL { 
MATCH (a:Account)-[:USER|ADMIN|OWNER*1..2]->(o:Organisation) WHERE a.name = $name AND o.name = $orgName RETURN o 
UNION 
MATCH (o:Organisation) WHERE (o.visibility = $public OR o.visibility = $private) AND o.name = $orgName RETURN o 
		} RETURN o;`,
		map[string]any{"name": name,
			"orgName": orgName,
			"private": PRIVATE,
			"public":  PUBLIC}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, nil, nil, err
	}
	organisation, err := getSingleOrganisation("o", result.Records)
	if err != nil {
		return nil, nil, nil, err
	}
	result, err = neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account)-[:USER]->(t:Organisation) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": orgName}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, nil, nil, err
	}
	user := make([]string, len(result.Records))
	for i, record := range result.Records {
		user[i] = record.Values[0].(string)
	}
	result, err = neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account)-[:ADMIN]->(t:Organisation) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": orgName}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, nil, nil, err
	}
	admin := make([]string, len(result.Records))
	for i, record := range result.Records {
		admin[i] = record.Values[0].(string)
	}
	return organisation, user, admin, err
}

func GetOrganisationMapForUser(account *Account) (map[string]map[string][]Organisation, error) {
	var list []Organisation
	var err error
	if account.Exists() && account.Role <= Admin {
		list, err = GetAllVisibleOrganisations()
	} else {
		list, err = GetOrganisationsForUserView(account)
	}
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return make(map[string]map[string][]Organisation), nil
	}
	mapping := make(map[string]map[string][]Organisation)
	mapping[list[0].MainType] = make(map[string][]Organisation)
	mapping[list[0].MainType][list[0].SubType] = []Organisation{list[0]}
	for i := range len(list) - 1 {
		if list[i].MainType != list[i+1].MainType {
			mapping[list[i+1].MainType] = make(map[string][]Organisation)
			mapping[list[i+1].MainType][list[i+1].SubType] = make([]Organisation, 0)
		} else if list[i].SubType != list[i+1].SubType {
			mapping[list[i+1].MainType][list[i+1].SubType] = make([]Organisation, 0)
		}
		mapping[list[i+1].MainType][list[i+1].SubType] = append(mapping[list[i+1].MainType][list[i+1].SubType],
			list[i+1])
	}
	return mapping, nil
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

func getSingleOrganisation(letter string, records []*neo4j.Record) (*Organisation, error) {
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
	title := &Organisation{
		Name:       node.Props["name"].(string),
		Visibility: OrganisationVisibility(node.Props["visibility"].(string)),
		MainType:   node.Props["main_type"].(string),
		SubType:    node.Props["sub_type"].(string),
		Flair:      node.Props["flair"].(string),
	}

	return title, nil
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
