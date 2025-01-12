package database

import (
	loc "PoliSim/localisation"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"strings"
	"time"
)

type AccountRole int

type Account struct {
	Name     string
	Username string
	Password string
	Role     AccountRole
	Blocked  bool
	FontSize int64
	TimeZone *time.Location
}

func (a *Account) Exists() bool {
	return a != nil
}

func (a *Account) IsAtLeastPressAdmin() bool {
	return a != nil && a.Role <= PressAdmin
}

func (a *Account) IsAtLeastAdmin() bool {
	return a != nil && a.Role <= Admin
}

func (a *Account) IsAtLeastHeadAdmin() bool {
	return a != nil && a.Role <= HeadAdmin
}

func (a *Account) IsPressUser() bool {
	return a != nil && a.Role == PressUser
}

func (a *Account) IsUser() bool {
	return a != nil && a.Role == User
}

func (a *Account) IsPressAdmin() bool {
	return a != nil && a.Role == PressAdmin
}

func (a *Account) IsAdmin() bool {
	return a != nil && a.Role == Admin
}

func (a *Account) IsHeadAdmin() bool {
	return a.IsAtLeastHeadAdmin()
}

const (
	Special AccountRole = iota - 1
	RootAdmin
	HeadAdmin
	Admin
	PressAdmin
	User
	PressUser
)

func CreateAccount(acc *Account) error {
	if acc.Name == loc.AdminstrationName {
		return notAllowedError
	}
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:Account {name: $name , username: $username ,
                password: $password , role: $role , blocked: $blocked, fontSize: 100, timezone: 'UTC' });`,
		map[string]any{"name": acc.Name,
			"username": acc.Username,
			"password": acc.Password,
			"role":     acc.Role,
			"blocked":  acc.Blocked}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func GetAccountByUsername(username string) (*Account, error) {
	if username == loc.AdminstrationAccountUsername {
		return nil, notFoundError
	}
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) WHERE a.username = $name RETURN a;`,
		map[string]any{"name": username}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return getSingleAccount("a", result.Records)
}

func GetAccountByName(name string) (*Account, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) WHERE a.name = $name RETURN a;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return getSingleAccount("a", result.Records)
}

func UpdateAccount(acc *Account) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (a:Account)  WHERE a.name = $name 
MATCH (s:Account)-[r]->() WHERE s.name = $name AND type(r) <> 'WRITTEN' AND $blocked 
DELETE r 
SET a.blocked = $blocked , a.role = $role 
RETURN a;`,
		map[string]any{"name": acc.Name,
			"role":    acc.Role,
			"blocked": acc.Blocked}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err == nil {
		updateAccount(acc)
	}
	return err
}

func UpdatePassword(acc *Account) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (a:Account)  WHERE a.name = $name SET a.password = $password RETURN a;`,
		map[string]any{"name": acc.Name,
			"password": acc.Password}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err == nil {
		updateAccount(acc)
	}
	return err
}

func SetPersonalSettings(acc *Account) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (a:Account)  WHERE a.name = $name SET a.fontSize = $fontSize, a.timezone = $timezone RETURN a;`,
		map[string]any{"name": acc.Name,
			"fontSize": acc.FontSize,
			"timezone": acc.TimeZone.String()}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err == nil {
		updateAccount(acc)
	}
	return err
}

func GetAccountAndOwnerByAccountName(name string) (account *Account, owner *Account, err error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) WHERE a.name = $name 
OPTIONAL MATCH (t:Account)-[:OWNER]->(a) RETURN a, t;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return
	}
	account, err = getSingleAccount("a", result.Records)
	if err != nil {
		return
	}
	owner, err = getSingleAccount("t", result.Records)
	if errors.Is(err, notFoundError) {
		err = nil
	}
	return
}

func IsAccountAllowedToPostWith(user *Account, poster string) (bool, error) {
	if user.Name == poster {
		return true, nil
	}
	result, err := makeRequest(`MATCH (t:Account)-[:OWNER]->(a:Account) 
WHERE a.name = $name AND t.name = $owner RETURN a;`, map[string]any{"name": poster, "owner": user.Name})
	if err != nil {
		return false, err
	}
	return len(result.Records) >= 1, err
}

// GetNames returns first an array of Names, then the array of Usernames and then the error, if one occurred
func GetNames() ([]string, []string, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) 
RETURN a.name AS name, a.username AS username;`,
		nil, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, nil, err
	}
	names := make([]string, len(queryResult.Records))
	usernames := make([]string, len(queryResult.Records))
	for i, record := range queryResult.Records {
		names[i] = record.Values[0].(string)
		usernames[i] = record.Values[1].(string)
	}

	return names, usernames, err
}

func GetNonBlockedNames() ([]string, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) 
WHERE a.blocked = false 
RETURN a.name AS name ORDER BY name;`,
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

func FilterNameListForNonBlocked(list []string) ([]string, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) 
WHERE a.blocked = false AND a.name IN $list
RETURN a.name AS name ORDER BY name;`,
		map[string]any{"list": list}, neo4j.EagerResultTransformer,
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

func GetNamesForActiveUsers() ([]string, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) WHERE a.role <> $role 
AND a.blocked = false RETURN a.name AS name;`,
		map[string]any{"role": PressUser}, neo4j.EagerResultTransformer,
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

func GetOwnerName(acc *Account) (string, error) {
	queryResult, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) WHERE a.name = $name 
MATCH (t:Account)-[:OWNER]->(a) RETURN t.name AS name;`,
		map[string]any{"name": acc.Name, "role": PressUser}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return "", err
	}
	if len(queryResult.Records) == 0 {
		return "", err
	}
	return queryResult.Records[0].Values[0].(string), err
}

func GetOwnedAccountNames(owner *Account) ([]string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (o:Account) WHERE o.name = $name 
MATCH (o)-[:OWNER]->(a:Account) 
RETURN a.name AS name ORDER BY name;`,
		map[string]any{"name": owner.Name}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	names := make([]string, len(result.Records))
	for i, record := range result.Records {
		names[i] = record.Values[0].(string)
	}
	return names, err
}

func GetMyAccountNames(owner *Account) ([]string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (o:Account) WHERE o.name = $name 
MATCH (o)-[:OWNER]->(a:Account) 
RETURN a.name AS name ORDER BY name;`,
		map[string]any{"name": owner.Name}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	names := make([]string, len(result.Records)+1)
	names[0] = owner.Name
	for i, record := range result.Records {
		names[i+1] = record.Values[0].(string)
	}
	return names, err
}

func MakeOwner(ownerName string, targetName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account), (t:Account) WHERE a.name = $owner 
AND t.name = $target MERGE (a)-[:OWNER]->(t);`,
		map[string]any{"owner": ownerName,
			"target": targetName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func RemoveOwner(targetName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (t:Account) WHERE t.name = $target 
MATCH (a:Account)-[r:OWNER]->(t) DELETE r;`,
		map[string]any{"target": targetName}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func getSingleAccount(letter string, records []*neo4j.Record) (*Account, error) {
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
	acc := &Account{
		Name:     node.Props["name"].(string),
		Username: node.Props["username"].(string),
		Role:     AccountRole(node.Props["role"].(int64)),
		Blocked:  node.Props["blocked"].(bool),
		Password: node.Props["password"].(string),
		FontSize: node.Props["fontSize"].(int64),
	}
	acc.TimeZone, _ = time.LoadLocation(node.Props["timezone"].(string))

	return acc, nil
}

func GetAccountFlairs(acc *Account) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `CALL {
MATCH (a:Account)-[:HAS]->(n:Title)
WHERE a.name = $name
RETURN n
UNION
MATCH (a:Account)-[:USER|ADMIN]->(n:Organisation)
WHERE a.name = $name
RETURN n }
RETURN n.flair AS flair ORDER BY flair;`,
		map[string]any{"name": acc.Name}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(""))
	if err != nil || len(result.Records) == 0 {
		return "", err
	}
	flairs := make([]string, 0, len(result.Records))
	for _, record := range result.Records {
		if flair := strings.TrimSpace(record.Values[0].(string)); flair != "" {
			flairs = append(flairs, flair)
		}
	}
	return strings.Join(flairs, ", "), err
}
