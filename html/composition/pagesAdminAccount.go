package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"strconv"
)

func GetCreateAccountPage(acc *validation.AccountModification, val validation.Message) Node {
	return getBasePageWrapper(
		getPageHeader(CreateUser),
		getFormStandardForm("form", POST, "/"+HTMXPreRouter+string(CreateUser), CLASS("w-[800px]"),
			getSimpleTextInput("username", "username", acc.Username, Translation["username"]),
			getSimpleTextInput("displayName", "displayName", acc.DisplayName, Translation["displayName"]),
			getSimpleTextInput("password", "password", acc.Password, Translation["password"]),
			getSimpleTextInput("flair", "flair", acc.Flair, Translation["flair"]),
			getDropDown("role", "role", Translation["role"], false,
				database.Roles, database.RoleTranslation, database.RoleLevel(acc.Role)),
			getInput("linked", "linked", strconv.Itoa(int(acc.Linked)),
				Translation["linked"], "number", "", ""),
			getSubmitButton("createAccountButton", Translation["createButton"])),
		GetMessage(val),
	)
}

func GetModifyAccount(acc *validation.AccountModification, val validation.Message) Node {
	hideLinked := ""
	if acc.Role != int(database.PressAccount) {
		hideLinked = "hidden"
	}
	display, user, err := extraction.ReturnNames()
	if err != nil {
		val.Message += "\n" + Translation["errorWhileRetrievingNames"]
	}
	return getBasePageWrapper(
		getDataList("userNames", user),
		getDataList("displayNames", display),
		getPageHeader(EditUser),
		getFormStandardForm("form", POST, "/"+HTMXPreRouter+string(EditUser), CLASS("w-[800px]"),
			getCheckBox("searchByUsername", acc.SearchByUsername, "true", "searchByUsername",
				Translation["searchByUsername"],
				HYPERSCRIPT("on click toggle .hidden on #usernameDiv then toggle .hidden on #displayNameDiv")),
			getInput("username", "username", acc.Username, Translation["username"], "text",
				"userNames", "hidden"),
			getInput("displayName", "displayName", acc.DisplayName, Translation["displayName"], "text",
				"displayNames", ""),
			getSubmitButtonOverwriteURL(Translation["searchAccountButton"], PATCH, "/"+HTMXPreRouter+string(SearchUser)),
			getStandardCheckBox(acc.ChangeFlair, "true", "changeFlair", Translation["changeFlair"]),
			getSimpleTextInput("flair", "flair", acc.Flair, Translation["flair"]),
			getStandardCheckBox(acc.Suspended, "true", "suspended", Translation["suspended"]),
			getDropDown("role", "role", Translation["role"], acc.Role == int(database.PressAccount),
				database.Roles, database.RoleTranslation, database.RoleLevel(acc.Role)),
			getInput("linked", "linked", strconv.Itoa(int(acc.Linked)), Translation["linked"],
				"number", "", hideLinked),
			getSubmitButton("modifyAccountButton", Translation["changeAccountButton"])),
		GetMessage(val),
	)
}

func GetViewAccountList(id string) Node {
	i, err := strconv.Atoi(id)
	if err != nil {
		i = 0
	}
	arr, err := extraction.ReturnAccountList(int64(i))
	if err != nil {
		return GetErrorPage(Translation["errorWithDatabaseRequest"])
	}
	nodes := make([]Node, len(arr))
	for i, item := range arr {
		susSpan := SPAN(CLASS("text-sm"), I(CLASS("bi bi-check-lg")))
		if item.Suspended {
			susSpan = SPAN(CLASS("text-sm"), I(CLASS("bi bi-x-lg")))
		}
		link := string(ViewUser) + "?id=" + strconv.FormatInt(item.ID, 10)
		nodes[i] = TR(
			getTableElement(StartPos, 1, Text(strconv.FormatInt(item.ID, 10))),
			getTableElement(MiddlePos, 1, IfElse(item.Role == database.PressAccount,
				Text(item.DisplayName),
				getClickableLink("/"+HTMXPreRouter+link, "/"+link,
					Text(item.DisplayName)))),
			getTableElement(MiddlePos, 1, Text(item.Username)),
			getTableElement(MiddlePos, 1, Text(item.Flair)),
			getTableElement(MiddlePos, 1, Text(database.RoleTranslation[item.Role])),
			getTableElement(MiddlePos, 1, susSpan),
			getTableElement(EndPos, 1, IfElse(item.Linked.Valid,
				Text(strconv.FormatInt(item.Linked.Int64, 10)),
				Text(Translation["notLinked"]))),
		)
	}
	return getBasePageWrapper(
		tableNode,
		getPageHeader(ViewUser),
		getStandardTable("sortTable",
			TR(
				getTableHeader(StartPos, 0, Translation["accountTableID"]),
				getTableHeader(MiddlePos, 1, Translation["accountTableDisplayName"]),
				getTableHeader(MiddlePos, 2, Translation["accountTableUsername"]),
				getTableHeader(MiddlePos, 3, Translation["accountTableFlair"]),
				getTableHeader(MiddlePos, 4, Translation["accountTableRole"]),
				getTableHeader(MiddlePos, 5, Translation["accountTableStatus"]),
				getTableHeader(EndPos, 6, Translation["accountTableLinked"]),
			),
			Group(nodes...),
		),
	)
}
