package htmlComposition

import (
	. "PoliSim/componentHelper"
	"PoliSim/dataValidation"
	"PoliSim/database"
	"strconv"
)

func GetCreateAccountPage(acc *dataValidation.AccountModification, val dataValidation.ValidationMessage) Node {
	return getBasePageWrapper(
		getPageHeader(CreateUser),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateUser), CLASS("w-[800px]"),
			getSimpleTextInput("username", "username", acc.Username, Translation["username"]),
			getSimpleTextInput("displayName", "displayName", acc.DisplayName, Translation["displayName"]),
			getSimpleTextInput("password", "password", acc.Password, Translation["password"]),
			getSimpleTextInput("flair", "flair", acc.Flair, Translation["flair"]),
			getDropDown("role", "role", Translation["role"], false,
				database.Roles, database.RoleTranslation, database.RoleLevel(acc.Role)),
			getInput("linked", "linked", strconv.Itoa(int(acc.Linked)), Translation["linked"], "number", "", ""),
			getSubmitButton(Translation["createButton"])),
		GetMessage(val),
	)
}
