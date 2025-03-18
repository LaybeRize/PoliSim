package database

import (
	"github.com/lib/pq"
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

	// Todo transfer to migration
	tableDefinition = `CREATE TABLE organisation(
    name TEXT PRIMARY KEY,
    main_group TEXT NOT NULL,
    sub_group TEXT NOT NULL,
    visibility INT NOT NULL,
    flair TEXT NOT NULL,
    users TEXT[] NOT NULL,
    admins TEXT[] NOT NULL
);
CREATE TABLE organisation_to_account(
    organisation_name TEXT NOT NULL,
    account_name TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL,
    CONSTRAINT fk_organisation_name
        FOREIGN KEY(organisation_name) REFERENCES organisation(name) ON UPDATE CASCADE,
    CONSTRAINT fk_account_name
        FOREIGN KEY(account_name) REFERENCES account(name)
);
CREATE VIEW organisation_linked AS
    SELECT organisation.*, ota.account_name, ota.is_admin, ownership.owner_name FROM organisation
    LEFT JOIN organisation_to_account ota ON organisation.name = ota.organisation_name
    LEFT JOIN ownership ON ota.account_name = ownership.account_name;`
)

func CreateOrganisation(org *Organisation, userNames []string, adminNames []string) error {
	var err error
	if org.Visibility == HIDDEN {
		_, err = postgresDB.Exec(`
INSERT INTO organisation (name, main_group, sub_group, visibility, flair, users, admins) 
	VALUES ($1, $2, $3, $4, $5, '{}', '{}')`, &org.Name, &org.MainType, &org.SubType, &org.Visibility, &org.Flair)
	} else {
		_, err = postgresDB.Exec(`
INSERT INTO organisation (name, main_group, sub_group, visibility, flair, users, admins) 
	VALUES ($1, $2, $3, $4, $5, 
	        ARRAY(SELECT name FROM account WHERE name = ANY($5) AND blocked = false), 
	        ARRAY(SELECT name FROM account WHERE name = ANY($6) AND (NOT (name = ANY($5))) AND blocked = false));
-- Insert into connection table
INSERT INTO organisation_to_account (organisation_name, account_name, is_admin) 
SELECT $1 AS organisation_name, name, false AS is_admin FROM account
WHERE name = ANY($5) AND blocked = false;
INSERT INTO organisation_to_account (organisation_name, account_name, is_admin) 
SELECT $1 AS organisation_name, name, true AS is_admin FROM account
WHERE name = ANY($6) AND (NOT (name = ANY($5))) AND blocked = false;`,
			&org.Name, &org.MainType, &org.SubType, &org.Visibility, &org.Flair,
			pq.Array(userNames), pq.Array(adminNames))
	}
	return err
}

func UpdateOrganisation(oldName string, org *Organisation) error {
	_, err := postgresDB.Exec(`
DELETE FROM organisation_to_account WHERE organisation_name = $1;
UPDATE organisation SET name = $2, main_group = $3, sub_group = $4, visibility = $5, flair = $6, 
                        users = '{}', admins = '{}' WHERE name = $1;`,
		&oldName, &org.Name, &org.MainType, &org.SubType, &org.Visibility, &org.Flair)
	return err
}

func AddOrganisationMember(org *Organisation, userNames []string, adminNames []string) error {
	if org.Visibility == HIDDEN {
		return nil
	}
	_, err := postgresDB.Exec(`
UPDATE organisation SET users = ARRAY(SELECT name FROM account WHERE name = ANY($2) AND blocked = false), 
                        admins = ARRAY(SELECT name FROM account WHERE name = ANY($3) AND (NOT (name = ANY($2))) AND blocked = false) 
                    WHERE name = $1;
INSERT INTO organisation_to_account (organisation_name, account_name, is_admin) 
SELECT $1 AS organisation_name, name, false AS is_admin FROM account
WHERE name = ANY($2) AND blocked = false;
INSERT INTO organisation_to_account (organisation_name, account_name, is_admin) 
SELECT $1 AS organisation_name, name, true AS is_admin FROM account
WHERE name = ANY($3) AND (NOT (name = ANY($2))) AND blocked = false;`, &org.Name, pq.Array(userNames), pq.Array(adminNames))
	return err
}

func GetFullOrganisationInfo(name string) (*Organisation, []string, []string, error) {
	organisation := &Organisation{}
	user := make([]string, 0)
	admin := make([]string, 0)
	err := postgresDB.QueryRow(`SELECT name, main_group, sub_group, visibility, flair, users, admins
    FROM organisation WHERE name = $1;`, &name).Scan(&organisation.Name, &organisation.MainType, &organisation.SubType,
		&organisation.Visibility, &organisation.Flair, pq.Array(user), pq.Array(admin))
	return organisation, user, admin, err
}

func GetOrganisationByName(name string) (*Organisation, error) {
	organisation := &Organisation{}
	err := postgresDB.QueryRow(`SELECT name, main_group, sub_group, visibility, flair
    FROM organisation WHERE name = $1;`, &name).Scan(&organisation.Name, &organisation.MainType, &organisation.SubType,
		&organisation.Visibility, &organisation.Flair)
	return organisation, err
}

func GetOrganisationNameList() ([]string, error) {
	result, err := postgresDB.Query(`SELECT name from organisation ORDER BY name;`)
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

func GetOrganisationsForUserView(account *Account) ([]Organisation, error) {
	name := ""
	if account.Exists() {
		name = account.Name
	}
	result, err := postgresDB.Query(`SELECT DISTINCT ON (name) name, main_group, sub_group, visibility, flair 
FROM organisation_linked WHERE visibility = $1 OR visibility = $2 OR owner_name = $3
ORDER BY main_group, sub_group, name;`, PUBLIC, PRIVATE, &name)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]Organisation, 0)
	org := Organisation{}
	for result.Next() {
		err = result.Scan(&org.Name, &org.MainType, &org.SubType, &org.Visibility, &org.Flair)
		if err != nil {
			return nil, err
		}
		arr = append(arr, org)
	}
	return arr, err
}

func GetOrganisationNamesAdminIn(name string) ([]string, error) {
	result, err := postgresDB.Query(`SELECT name, main_group, sub_group, visibility, flair 
FROM organisation_linked WHERE is_admin = true AND account_name = $1
ORDER BY name;`, &name)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	names := make([]string, 0)
	orgName := ""
	for result.Next() {
		err = result.Scan(&orgName)
		if err != nil {
			return nil, err
		}
		names = append(names, orgName)
	}
	return names, err
}

func GetFullOrganisationInfoForUserView(account *Account, orgName string) (*Organisation, []string, []string, error) {
	name := ""
	if account.Exists() {
		name = account.Name
	}
	organisation := &Organisation{}
	user := make([]string, 0)
	admin := make([]string, 0)
	err := postgresDB.QueryRow(`SELECT name, main_group, sub_group, visibility, flair, users, admins
    FROM organisation_linked WHERE (visibility = $1 OR visibility = $2 OR owner_name = $3) AND name = $4
    LIMIT 1;`, PUBLIC, PRIVATE, &name, &orgName).Scan(&organisation.Name, &organisation.MainType, &organisation.SubType,
		&organisation.Visibility, &organisation.Flair, pq.Array(user), pq.Array(admin))
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
	result, err := postgresDB.Query(`SELECT name, main_group, sub_group, visibility, flair 
FROM organisation WHERE visibility <> $1
ORDER BY main_group, sub_group, name;`, HIDDEN)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]Organisation, 0)
	org := Organisation{}
	for result.Next() {
		err = result.Scan(&org.Name, &org.MainType, &org.SubType, &org.Visibility, &org.Flair)
		if err != nil {
			return nil, err
		}
		arr = append(arr, org)
	}
	return arr, err
}
