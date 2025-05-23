package database

import (
	loc "PoliSim/localisation"
	"database/sql"
	"github.com/lib/pq"
	"html/template"
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
	FontSize int
	TimeZone *time.Location
}

func (a *Account) GetName() string {
	if a == nil {
		return ""
	}
	return a.Name
}

func (a *Account) GetFontSize() int {
	if a == nil {
		return 100
	}
	return a.FontSize
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

func (a *Account) QueryEscapeName() template.URL {
	return template.URL(template.URLQueryEscaper(a.Name))
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

var BlockedAccountChannel = make(chan string)
var OwnerChangeOnAccountChannel = make(chan []string)

func CreateAccount(acc *Account) error {
	if loc.IsAdministrationName(acc.Name) {
		return NotAllowedError
	}
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}
	defer rollback(tx)
	_, err = tx.Exec(`
INSERT INTO account(name, username, password, role, blocked, font_size, time_zone) VALUES ($1,$2,$3,$4,$5,100,'UTC');`,
		&acc.Name, &acc.Username, &acc.Password, &acc.Role, &acc.Blocked)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO ownership(account_name, owner_name) VALUES ($1, $1);`,
		acc.Name)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err == nil {
		createdAccountForNotifications <- acc.Name
	}
	return err
}

func GetAccountByUsername(username string) (*Account, error) {
	if username == loc.AdministrationAccountUsername {
		return nil, sql.ErrNoRows
	}
	acc := &Account{}
	timeZoneStr := ""
	err := postgresDB.QueryRow(`SELECT name,username,password,role,blocked,font_size,time_zone FROM account 
                                      WHERE username = $1;`,
		&username).Scan(&acc.Name, &acc.Username, &acc.Password, &acc.Role, &acc.Blocked, &acc.FontSize,
		&timeZoneStr)
	if err != nil {
		return nil, err
	}
	acc.TimeZone, err = time.LoadLocation(timeZoneStr)
	return acc, err
}

func GetAccountByName(name string) (*Account, error) {
	acc := &Account{}
	timeZoneStr := ""
	err := postgresDB.QueryRow(`SELECT name,username,password,role,blocked,font_size,time_zone FROM account
                                      WHERE name = $1;`,
		&name).Scan(&acc.Name, &acc.Username, &acc.Password, &acc.Role, &acc.Blocked, &acc.FontSize, &timeZoneStr)
	if err != nil {
		return nil, err
	}
	acc.TimeZone, err = time.LoadLocation(timeZoneStr)
	return acc, err
}

func GetAccountByRole(role AccountRole) (*Account, error) {
	acc := &Account{}
	timeZoneStr := ""
	err := postgresDB.QueryRow(`SELECT name,username,password,role,blocked,font_size,time_zone FROM account
                                      WHERE role = $1;`,
		role).Scan(&acc.Name, &acc.Username, &acc.Password, &acc.Role, &acc.Blocked, &acc.FontSize, &timeZoneStr)
	if err != nil {
		return nil, err
	}
	acc.TimeZone, err = time.LoadLocation(timeZoneStr)
	return acc, err
}

func UpdateAccount(acc *Account) error {
	var err error
	var tx *sql.Tx
	if acc.Blocked {
		tx, err = postgresDB.Begin()
		if err != nil {
			return err
		}
		defer rollback(tx)
		_, err = tx.Exec(`UPDATE account SET role = $2, blocked = true WHERE name = $1;`, &acc.Name, &acc.Role)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`UPDATE organisation SET users = array_remove(users, $1), admins = array_remove(admins, $1) 
                                FROM organisation_to_account ota 
                                WHERE  ota.organisation_name = organisation.name AND ota.account_name = $1;`, &acc.Name)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`DELETE FROM organisation_to_account WHERE account_name = $1;`, &acc.Name)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`DELETE FROM title_to_account WHERE account_name = $1;`, &acc.Name)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`DELETE FROM newspaper_to_account WHERE account_name = $1;`, &acc.Name)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`UPDATE ownership SET owner_name = $1 WHERE account_name = $1;`, &acc.Name)
		if err != nil {
			return err
		}
		err = tx.Commit()
		if err == nil {
			BlockedAccountChannel <- acc.Name
			blockedAccountChannelForNotifications <- acc.Name
		}
	} else {
		_, err = postgresDB.Exec(`UPDATE account SET role = $2, blocked = false WHERE name = $1;`,
			&acc.Name, &acc.Role)
	}
	if err == nil {
		updateAccount(acc)
	}
	return err
}

func UpdatePassword(acc *Account) error {
	_, err := postgresDB.Exec(`UPDATE account SET password = $2 
                                     WHERE name = $1;`, &acc.Name, &acc.Password)
	if err == nil {
		updateAccount(acc)
	}
	return err
}

func SetPersonalSettings(acc *Account) error {
	timeZoneStr := acc.TimeZone.String()
	_, err := postgresDB.Exec(`UPDATE account SET font_size = $2, time_zone = $3 
                                     WHERE name = $1;`, &acc.Name, &acc.FontSize, &timeZoneStr)
	if err == nil {
		updateAccount(acc)
	}
	return err
}

func GetAccountAndOwnerByAccountName(name string) (account *Account, owner *Account, err error) {
	owner = nil
	account, err = GetAccountByName(name)
	if err != nil {
		return
	}
	ownerName := ""
	err = postgresDB.QueryRow(`SELECT owner_name FROM ownership WHERE account_name = $1`, &name).Scan(&ownerName)
	if err != nil || ownerName == name {
		return
	}
	owner, err = GetAccountByName(ownerName)
	return
}

func IsAccountAllowedToPostWith(user *Account, poster string) (bool, error) {
	ownerName := ""
	err := postgresDB.QueryRow(`SELECT owner_name FROM ownership WHERE account_name = $1`, &poster).Scan(&ownerName)
	if err != nil {
		return false, err
	}
	return user.Name == ownerName, nil
}

// GetNames returns first an array of Names, then the array of Usernames and then the error, if one occurred
func GetNames() ([]string, []string, error) {
	result, err := postgresDB.Query(`SELECT name, username FROM account`)
	if err != nil {
		return nil, nil, err
	}
	defer closeRows(result)
	names := make([]string, 0)
	usernames := make([]string, 0)
	name := ""
	username := ""
	for result.Next() {
		err = result.Scan(&name, &username)
		if err != nil {
			return nil, nil, err
		}
		names = append(names, name)
		usernames = append(usernames, username)
	}

	return names, usernames, err
}

func GetNonBlockedNames() ([]string, error) {
	result, err := postgresDB.Query(`SELECT name FROM account WHERE blocked = false ORDER BY name;`)
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

func FilterNameListForNonBlocked(list []string, extraEmptyEntries int) ([]string, error) {
	result, err := postgresDB.Query(`
SELECT name FROM account WHERE blocked = false AND name = ANY($1) ORDER BY name;`, pq.Array(list))
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

	for range extraEmptyEntries {
		names = append(names, "")
	}
	return names, err
}

func GetNamesForActiveUsers() ([]string, error) {
	result, err := postgresDB.Query(`SELECT name FROM account WHERE blocked = false AND role <> $1 ORDER BY name;`, PressUser)
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

func GetOwnedAccountNames(owner *Account) ([]string, error) {
	result, err := postgresDB.Query(`SELECT account_name FROM ownership 
                    WHERE owner_name = $1 AND account_name <> $1 ORDER BY account_name;`, &owner.Name)
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

func GetMyAccountNames(owner *Account) ([]string, error) {
	result, err := postgresDB.Query(`SELECT account_name FROM ownership 
                    WHERE owner_name = $1 AND account_name <> $1 ORDER BY account_name;`, &owner.Name)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	names := make([]string, 1)
	names[0] = owner.Name
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

func UpdateOwner(ownerName string, targetName string) error {
	_, err := postgresDB.Exec(`UPDATE ownership SET owner_name = $1 WHERE account_name = $2`, &ownerName, &targetName)
	if err == nil {
		OwnerChangeOnAccountChannel <- []string{ownerName, targetName}
		ownerChangeOnAccountChannelForNotifications <- []string{ownerName, targetName}
	}
	return err
}

func GetAccountFlairs(acc *Account) (string, error) {
	flairArr := make([]string, 0)
	err := postgresDB.QueryRow(`SELECT ARRAY(
SELECT flair FROM organisation
    LEFT JOIN organisation_to_account ota on organisation.name = ota.organisation_name
    WHERE account_name = $1 AND flair <> ''
UNION ALL
SELECT flair FROM title
    LEFT JOIN title_to_account tta on title.name = tta.title_name
    WHERE account_name = $1 AND flair <> '' ORDER BY flair)`, acc.Name).Scan(pq.Array(&flairArr))
	if err != nil {
		return "", err
	}
	return strings.Join(flairArr, ", "), nil
}
