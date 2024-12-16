package database

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

type AccountRole int

type Account struct {
	Name     string
	Username string
	Password string
	Role     AccountRole
	Blocked  bool
}

const (
	HEAD_ADMIN AccountRole = iota
	ADMIN
	PRESS_ADMIN
	USER
	PRESS_USER

	CON_OWNER = "OWNER"

	DB_ACC_NAME     = "name"
	DB_ACC_USERNAME = "username"
	DB_ACC_PASSWORD = "password"
	DB_ACC_ROLE     = "role"
	DB_ACC_BLOCKED  = "blocked"
)

func CreateAccount(acc Account) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:Account {`+DB_ACC_NAME+`: $name , `+DB_ACC_USERNAME+`: $username , `+
			DB_ACC_PASSWORD+`: $password , `+DB_ACC_ROLE+`: $role , `+DB_ACC_BLOCKED+`: $blocked });`,
		map[string]any{"name": acc.Name,
			"username": acc.Username,
			"password": acc.Password,
			"role":     acc.Role,
			"blocked":  acc.Blocked}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func GetAccountByUsername(username string) (*Account, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) WHERE a.`+DB_ACC_USERNAME+` = $name RETURN a;`,
		map[string]any{"name": username}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return getSingleAccount("a", result.Records)
}

func GetAccountByName(name string) (*Account, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (a:Account) WHERE a.`+DB_ACC_NAME+` = $name RETURN a;`,
		map[string]any{"name": name}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return getSingleAccount("a", result.Records)
}

func UpdateAccount(acc Account) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (a:Account)  WHERE a.`+DB_ACC_NAME+` = $name SET a.`+DB_ACC_BLOCKED+` = $blocked , a.`+DB_ACC_PASSWORD+` = $password , a.`+DB_ACC_ROLE+` = $role RETURN a;`,
		map[string]any{"name": acc.Name,
			"password": acc.Password,
			"role":     acc.Role,
			"blocked":  acc.Blocked}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err == nil {
		updateAccount(&acc)
	}
	return err
}

func MakeOwner(ownerName string, targetName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, "MATCH (a:Account), (t:Account) WHERE a."+DB_ACC_NAME+" = $owner AND t."+
		DB_ACC_NAME+" = $target CREATE (a)-[:"+CON_OWNER+"]->(t);",
		map[string]any{"owner": ownerName,
			"target": targetName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func RemoveOwner(targetName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, "MATCH (t:Account) WHERE t."+DB_ACC_NAME+" = $target MATCH (a:Account)-[r:"+CON_OWNER+"]->(t) DELETE r;",
		map[string]any{"target": targetName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	return err
}

func GetAllowedAsUser(ownerName string, orgName string) ([]Account, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (o:Organisation) WHERE a.`+DB_ORG_NAME+` = $org MATCH (a:Account)-[:`+CON_ADMIN+`|`+CON_USER+`]->(o) 
WHERE a.`+DB_ACC_NAME+` = $owner
RETURN a UNION 
MATCH (o:Organisation), (acc:Account) WHERE o.`+DB_ORG_NAME+` = $org AND acc.`+DB_ACC_NAME+` = $owner MATCH (acc)-[:`+CON_OWNER+`]->(a:Account)-[:`+CON_ADMIN+`|`+CON_USER+`]->(o) 
RETURN a ORDER BY a.`+DB_ACC_NAME+`;`,
		map[string]any{"org": orgName, "owner": ownerName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return getArrayOfAccounts("a", result.Records), err
}

func GetAllowedAsAdmin(ownerName string, orgName string) ([]Account, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `MATCH (o:Organisation) WHERE a.`+DB_ORG_NAME+` = $org MATCH (a:Account)-[:`+CON_ADMIN+`]->(o) 
WHERE a.`+DB_ACC_NAME+` = $owner
RETURN a UNION 
MATCH (o:Organisation), (acc:Account) WHERE o.`+DB_ORG_NAME+` = $org AND acc.`+DB_ACC_NAME+` = $owner MATCH (acc)-[:`+CON_OWNER+`]->(a:Account)-[:`+CON_ADMIN+`]->(o) 
RETURN a ORDER BY a.`+DB_ACC_NAME+`;`,
		map[string]any{"org": orgName, "owner": ownerName}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return getArrayOfAccounts("a", result.Records), err
}

func getArrayOfAccounts(letter string, records []*neo4j.Record) []Account {
	arr := make([]Account, 0, len(records))
	for _, record := range records {
		result, exists := record.Get(letter)
		if !exists {
			continue
		}
		node := result.(neo4j.Node)
		arr = append(arr, Account{
			Name:     node.Props[DB_ACC_NAME].(string),
			Username: node.Props[DB_ACC_USERNAME].(string),
			Password: node.Props[DB_ACC_PASSWORD].(string),
			Role:     node.Props[DB_ACC_ROLE].(AccountRole),
			Blocked:  node.Props[DB_ACC_BLOCKED].(bool),
		})
	}
	return arr
}

func getSingleAccount(letter string, records []*neo4j.Record) (*Account, error) {
	if len(records) == 0 {
		return nil, notFoundError
	} else if len(records) > 1 {
		return nil, multipleItemsError
	}
	result, exists := records[0].Get(letter)
	if !exists {
		return nil, notFoundError
	}
	node := result.(neo4j.Node)
	return &Account{
		Name:     node.Props[DB_ACC_NAME].(string),
		Username: node.Props[DB_ACC_USERNAME].(string),
		Password: node.Props[DB_ACC_PASSWORD].(string),
		Role:     node.Props[DB_ACC_ROLE].(AccountRole),
		Blocked:  node.Props[DB_ACC_BLOCKED].(bool),
	}, nil
}
