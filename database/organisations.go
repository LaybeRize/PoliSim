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
	return o.Visibility >= PUBLIC && o.Visibility <= HIDDEN
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
	default:
		return ""
	}
}

func (o *Organisation) IsPublic() bool {
	return o.Visibility == PUBLIC
}
func (o *Organisation) IsPrivate() bool {
	return o.Visibility == PRIVATE
}
func (o *Organisation) IsSecret() bool {
	return o.Visibility == SECRET
}
func (o *Organisation) IsHidden() bool {
	return o.Visibility == HIDDEN
}

type OrganisationVisibility int

const (
	PUBLIC OrganisationVisibility = iota
	PRIVATE
	SECRET
	HIDDEN
)

func CreateOrganisation(org *Organisation, userNames []string, adminNames []string) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(
		`CREATE (:Organisation {name: $name , visibility: $visibility , main_type: $maintype , 
sub_type: $subtype , flair: $flair});`,
		map[string]any{"name": org.Name,
			"visibility": org.Visibility,
			"maintype":   org.MainType,
			"subtype":    org.SubType,
			"flair":      org.Flair})
	if err != nil {
		return err
	}
	if org.Visibility != HIDDEN {
		err = tx.RunWithoutResult(`
MATCH (o:Organisation) WHERE o.name = $org 
MATCH (a:Account) WHERE a.name IN $aNames 
MERGE (a)-[:ADMIN]->(o);`, map[string]any{
			"org":    org.Name,
			"aNames": adminNames})
		if err != nil {
			return err
		}
		err = tx.RunWithoutResult(`
MATCH (o:Organisation) WHERE o.name = $org 
MATCH (u:Account) WHERE u.name IN $uNames AND (NOT u.name IN $aNames) 
MERGE (u)-[:USER]->(o);`, map[string]any{
			"org":    org.Name,
			"uNames": userNames,
			"aNames": adminNames})
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

func UpdateOrganisation(oldName string, org *Organisation) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(`MATCH (o:Organisation) WHERE o.name = $oldName 
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
		return err
	}
	err = tx.RunWithoutResult(`
MATCH (a:Account)-[r:ADMIN|USER]->(o:Organisation) WHERE o.name = $org 
DELETE r;`, map[string]any{"org": org.Name})
	if err != nil {
		return err
	}

	if org.Visibility == HIDDEN {
		err = tx.RunWithoutResult(`
MATCH (:Account)-[r:FAVORITE]->(o:Organisation) WHERE o.name = $org 
DELETE r;`, map[string]any{"org": org.Name})
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func AddOrganisationMember(org *Organisation, userNames []string, adminNames []string) error {
	if org.Visibility == HIDDEN {
		return nil
	}
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(`
MATCH (o:Organisation) WHERE o.name = $org 
MATCH (a:Account) WHERE a.name IN $aNames AND a.blocked = false
MERGE (a)-[:ADMIN]->(o);`, map[string]any{
		"org":    org.Name,
		"aNames": adminNames})
	if err != nil {
		return err
	}
	err = tx.RunWithoutResult(`
MATCH (o:Organisation) WHERE o.name = $org
MATCH (u:Account) WHERE u.name IN $uNames AND (NOT u.name IN $aNames) AND u.blocked = false
MERGE (u)-[:USER]->(o);`, map[string]any{
		"org":    org.Name,
		"uNames": userNames,
		"aNames": adminNames})
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func GetFullOrganisationInfo(name string) (*Organisation, []string, []string, error) {
	result, err := makeRequest(`MATCH (t:Organisation) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name})
	if err != nil {
		return nil, nil, nil, err
	}
	organisation, err := getSingleOrganisation(0, result)
	if err != nil {
		return nil, nil, nil, err
	}
	result, err = makeRequest(`MATCH (a:Account)-[:USER]->(t:Organisation) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": name})
	if err != nil {
		return nil, nil, nil, err
	}
	user := make([]string, len(result))
	for i, record := range result {
		user[i] = record.Values[0].(string)
	}
	result, err = makeRequest(`MATCH (a:Account)-[:ADMIN]->(t:Organisation) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": name})
	if err != nil {
		return nil, nil, nil, err
	}
	admin := make([]string, len(result))
	for i, record := range result {
		admin[i] = record.Values[0].(string)
	}
	return organisation, user, admin, err
}

func GetOrganisationByName(name string) (*Organisation, error) {
	result, err := makeRequest(`MATCH (t:Organisation) WHERE t.name = $name RETURN t;`,
		map[string]any{"name": name})
	if err != nil {
		return nil, err
	}
	return getSingleOrganisation(0, result)
}

func GetOrganisationNameList() ([]string, error) {
	result, err := makeRequest(`MATCH (t:Organisation) RETURN t.name AS name;`,
		nil)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(result))
	for i, record := range result {
		names[i] = record.Values[0].(string)
	}
	return names, err
}

func GetOrganisationsForUserView(account *Account) ([]Organisation, error) {
	name := ""
	if account.Exists() {
		name = account.Name
	}
	result, err := makeRequest(
		`CALL { 
MATCH (a:Account)-[:USER|ADMIN|OWNER*1..2]->(o:Organisation) WHERE a.name = $name RETURN o 
UNION 
MATCH (o:Organisation) WHERE o.visibility = $public OR o.visibility = $private RETURN o 
		} RETURN o ORDER BY o.main_type, o.sub_type, o.name;`,
		map[string]any{"name": name,
			"private": PRIVATE,
			"public":  PUBLIC})
	if err != nil {
		return nil, err
	}

	return getArrayOfOrganisations(0, result), err
}

func GetOrganisationNamesAdminIn(name string) ([]string, error) {
	result, err := makeRequest(`MATCH (a:Account)-[:ADMIN]->(o:Organisation) 
WHERE a.name = $name RETURN o.name;`, map[string]any{"name": name})
	if err != nil {
		return nil, err
	}
	names := make([]string, len(result))
	for i, record := range result {
		names[i] = record.Values[0].(string)
	}
	return names, err
}

func GetFullOrganisationInfoForUserView(account *Account, orgName string) (*Organisation, []string, []string, error) {
	name := ""
	if account.Exists() {
		name = account.Name
	}
	result, err := makeRequest(`CALL { 
MATCH (a:Account)-[:USER|ADMIN|OWNER*1..2]->(o:Organisation) WHERE a.name = $name AND o.name = $orgName RETURN o 
UNION 
MATCH (o:Organisation) WHERE (o.visibility = $public OR o.visibility = $private) AND o.name = $orgName RETURN o 
		} RETURN o;`,
		map[string]any{"name": name,
			"orgName": orgName,
			"private": PRIVATE,
			"public":  PUBLIC})
	if err != nil {
		return nil, nil, nil, err
	}
	organisation, err := getSingleOrganisation(0, result)
	if err != nil {
		return nil, nil, nil, err
	}
	result, err = makeRequest(`MATCH (a:Account)-[:USER]->(t:Organisation) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": orgName})
	if err != nil {
		return nil, nil, nil, err
	}
	user := make([]string, len(result))
	for i, record := range result {
		user[i] = record.Values[0].(string)
	}
	result, err = makeRequest(`MATCH (a:Account)-[:ADMIN]->(t:Organisation) 
WHERE t.name = $name RETURN a.name AS name;`,
		map[string]any{"name": orgName})
	if err != nil {
		return nil, nil, nil, err
	}
	admin := make([]string, len(result))
	for i, record := range result {
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
	result, err := makeRequest(
		`MATCH (o:Organisation) WHERE o.visibility <> $hidden RETURN o 
ORDER BY o.main_type, o.sub_type, o.name;`,
		map[string]any{
			"hidden": HIDDEN,
		})
	if err != nil {
		return nil, err
	}

	return getArrayOfOrganisations(0, result), err
}

func getSingleOrganisation(pos int, records []*neo4j.Record) (*Organisation, error) {
	if len(records) == 0 {
		return nil, NotFoundError
	} else if len(records) > 1 {
		return nil, MultipleItemsError
	}
	props := GetPropsMapForRecordPosition(records[0], pos)
	if props == nil {
		return nil, NotFoundError
	}
	title := &Organisation{
		Name:       props.GetString("name"),
		Visibility: OrganisationVisibility(props.GetInt("visibility")),
		MainType:   props.GetString("main_type"),
		SubType:    props.GetString("sub_type"),
		Flair:      props.GetString("flair"),
	}

	return title, nil
}

func getArrayOfOrganisations(pos int, records []*neo4j.Record) []Organisation {
	arr := make([]Organisation, 0, len(records))
	for _, record := range records {
		props := GetPropsMapForRecordPosition(record, pos)
		if props == nil {
			continue
		}
		arr = append(arr, Organisation{
			Name:       props.GetString("name"),
			Visibility: OrganisationVisibility(props.GetInt("visibility")),
			MainType:   props.GetString("main_type"),
			SubType:    props.GetString("sub_type"),
			Flair:      props.GetString("flair"),
		})
	}
	return arr
}
