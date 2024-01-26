package database

func FindFullAccountByDisplayName(name string) (Account, error) {
	result := Account{}
	row, err := DB.NamedQuery(`SELECT account.id, account.display_name, account.flair, account.username, account.password, account.suspended, 
                                     account.login_tries, account.next_login_time, account.role, account.linked, account.has_letters, account.parent FROM 
                                     account WHERE 
                                     account.display_name=:name`, map[string]interface{}{"name": name})
	if err != nil {
		return result, err
	}
	row.Next()
	err = row.StructScan(&result)
	return result, err
}
func FindFullAccountByUsername(name string) ([]Account, error) {
	result := make([]Account, 0)
	rows, err := DB.NamedQuery(`SELECT account.id, account.display_name, account.flair, account.username, account.password, account.suspended, 
                                     account.login_tries, account.next_login_time, account.role, account.linked, account.has_letters, account.parent FROM 
                                     account WHERE 
                                     account.username=:name`, map[string]interface{}{"name": name})
	if err != nil {
		return result, err
	}
	pos := 0
	for rows.Next() {
		result = append(result, Account{})
		err = rows.StructScan(&result[pos])
		if err != nil {
			return result, err
		}
		pos++
	}
	return result, err
}
