package htmlComposition

import (
	. "PoliSim/componentHelper"
	"PoliSim/dataExtraction"
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

func GetModifyAccount(acc *dataValidation.AccountModification, val dataValidation.ValidationMessage) Node {
	hideLinked := ""
	if acc.Role != int(database.PressAccount) {
		hideLinked = "hidden"
	}
	display, user, err := dataExtraction.ReturnNames()
	if err != nil {
		val.Message += "\n" + Translation["errorWhileRetrievingNames"]
	}
	return getBasePageWrapper(
		getDataList("userNames", user),
		getDataList("displayNames", display),
		getPageHeader(EditUser),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(EditUser), CLASS("w-[800px]"),
			getCheckBox("searchByUsername", false, false, "true", "searchByUsername", Translation["searchByUsername"],
				HYPERSCRIPT("on click toggle .hidden on #usernameDiv then toggle .hidden on #displayNameDiv")),
			getInput("username", "username", acc.Username, Translation["username"], "text", "userNames", "hidden"),
			getInput("displayName", "displayName", acc.DisplayName, Translation["displayName"], "text", "displayNames", ""),
			getSubmitButtonOverwriteURL(Translation["searchAccountButton"], PATCH, "/"+APIPreRoute+string(SearchUser)),
			getSimpleTextInput("flair", "flair", acc.Flair, Translation["flair"]),
			getDropDown("role", "role", Translation["role"], acc.Role == int(database.PressAccount),
				//exclude Press Accounts because they can only be created
				database.Roles[1:], database.RoleTranslation, database.RoleLevel(acc.Role)),
			getInput("linked", "linked", strconv.Itoa(int(acc.Linked)), Translation["linked"], "number", "", hideLinked),
			getSubmitButton(Translation["changeAccountButton"])),
		GetMessage(val),
	)
}
