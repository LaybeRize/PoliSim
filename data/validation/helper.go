package validation

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
)

// Message provides a generic sturct for responding to a request with a message (either as an error or as successful)
type Message struct {
	Message  string
	Positive bool
}

func GetDisplayNameArray(accs *[]database.Account) []string {
	strs := make([]string, len(*accs))
	for i, acc := range *accs {
		strs[i] = acc.DisplayName
	}
	return strs
}

// isRoleValid checks if the role is not database.NotLoggedIn but still one of the existing roles
func isRoleValid(level int) bool {
	return level >= int(database.PressAccount) && level != int(database.NotLoggedIn) && level <= int(database.HeadAdmin)
}

// isValidString checks if the string is empty or has at maximum the specified length.
// it also returns true on any length that is unequal to 0 if the length is set to -1.
func isValidString(str string, length int) bool {
	return str != "" && (length == -1 || len([]rune(str)) <= length)
}

func isOrgStatusValid(str string) bool {
	_, ok := database.StatusTranslation[database.StatusString(str)]
	return ok
}

func IsAccountValidForUser(userID int64, accountDisplayName string) (*database.Account, bool, error) {
	acc, err := extraction.GetAccountByDisplayName(accountDisplayName)
	if err != nil || acc.Suspended {
		return acc, false, err
	}
	if (acc.Role == database.PressAccount && acc.Linked.Int64 == userID) || acc.ID == userID {
		return acc, true, err
	}
	return acc, false, err
}

// IsOrganisationValidForAccount returns the organisation and a boolean indicating if the user is an admin or not
func IsOrganisationValidForAccount(userID int64, organisationName string) (*database.Organisation, bool, error) {
	org, err := extraction.GetOrganisationForWithUserInIt(userID, organisationName)
	for _, account := range org.Admins {
		if account.ID == userID {
			return org, true, err
		}
	}
	return org, false, err
}
