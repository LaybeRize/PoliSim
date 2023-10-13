package validation

import (
	"PoliSim/data/database"
)

func GetDisplayNameArray(accs *[]database.Account) []string {
	strs := make([]string, len(*accs))
	for i, acc := range *accs {
		strs[i] = acc.DisplayName
	}
	return strs
}

type Message struct {
	Message  string
	Positive bool
}

// isRoleValid checks if the role not database.NotLoggedIn
func isRoleValid(level int) bool {
	return level >= int(database.PressAccount) && level != int(database.NotLoggedIn) && level <= int(database.HeadAdmin)
}

func isEmptyOrNotInRange(str string, length int) bool {
	return str == "" || len([]rune(str)) > length
}
