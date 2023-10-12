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
			getCheckBox("searchByUsername", acc.SearchByUsername, false, "true", "searchByUsername", Translation["searchByUsername"],
				HYPERSCRIPT("on click toggle .hidden on #usernameDiv then toggle .hidden on #displayNameDiv")),
			getInput("username", "username", acc.Username, Translation["username"], "text", "userNames", "hidden"),
			getInput("displayName", "displayName", acc.DisplayName, Translation["displayName"], "text", "displayNames", ""),
			getSubmitButtonOverwriteURL(Translation["searchAccountButton"], PATCH, "/"+APIPreRoute+string(SearchUser)),
			getCheckBox("changeFlair", acc.ChangeFlair, false, "true", "changeFlair", Translation["changeFlair"], nil),
			getSimpleTextInput("flair", "flair", acc.Flair, Translation["flair"]),
			getCheckBox("suspended", acc.Suspended, false, "true", "suspended", Translation["suspended"], nil),
			getDropDown("role", "role", Translation["role"], acc.Role == int(database.PressAccount),
				database.Roles, database.RoleTranslation, database.RoleLevel(acc.Role)),
			getInput("linked", "linked", strconv.Itoa(int(acc.Linked)), Translation["linked"], "number", "", hideLinked),
			getSubmitButton(Translation["changeAccountButton"])),
		GetMessage(val),
	)
}

func GetViewAccountList(id string) Node {
	i, err := strconv.Atoi(id)
	if err != nil {
		i = 0
	}
	arr, err := dataExtraction.ReturnAccountList(int64(i))
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
				getClickableLink("/"+APIPreRoute+link, "/"+link,
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
				getTableHeader(StartPos, 0, "ID"),
				getTableHeader(MiddlePos, 1, "Anzeigename"),
				getTableHeader(MiddlePos, 2, "Nutzername"),
				getTableHeader(MiddlePos, 3, "Flair"),
				getTableHeader(MiddlePos, 4, "Rolle"),
				getTableHeader(MiddlePos, 5, "Status"),
				getTableHeader(EndPos, 6, "Verlinkt mit"),
			),
			Group(nodes...),
		),
	)
}

type Position int

const (
	StartPos Position = iota
	MiddlePos
	EndPos
)

func getStandardTable(id string, children ...Node) Node {
	return TABLE(ID(id), CLASS("table-auto mt-4 w-[800px]"),
		Group(children...))
}

func getTableHeader(p Position, sortPos int, text string) Node {
	var hScript Node = nil
	if sortPos >= 0 {
		hScript = HYPERSCRIPT("on click call sortTable(" + strconv.Itoa(sortPos) + ")")
	}
	switch p {
	case StartPos:
		return TH(CLASS("p-2 border-r-2 border-gray-600"), If(hScript != nil, STYLE("cursor: pointer;")), hScript,
			Text(text))
	case MiddlePos:
		return TH(CLASS("p-2 border-r-2 border-gray-600"), If(hScript != nil, STYLE("cursor: pointer;")), hScript,
			Text(text))
	case EndPos:
		return TH(CLASS("p-2"), If(hScript != nil, STYLE("cursor: pointer;")), hScript,
			Text(text))
	}
	return nil
}
func getTableElement(p Position, rowSpan int, node Node) Node {
	switch p {
	case StartPos:
		return TD(ROWSPAN(strconv.Itoa(rowSpan)), CLASS("p-2 border-r-2 border-t-2 border-gray-600"), node)
	case MiddlePos:
		return TD(ROWSPAN(strconv.Itoa(rowSpan)), CLASS("p-2 border-r-2 border-t-2 border-gray-600"), node)
	case EndPos:
		return TD(ROWSPAN(strconv.Itoa(rowSpan)), CLASS("p-2 border-t-2 border-gray-600"), node)
	}
	return nil
}

var tableNode = Raw(`<script>
        function sortTable(n) {
            var table, rows, switching, i, x, y, shouldSwitch, dir, switchcount = 0;
            table = document.getElementById("sortTable");
            switching = true;
            //Set the sorting direction to ascending:
            dir = "asc";
            /*Make a loop that will continue until
            no switching has been done:*/
            while (switching) {
                //start by saying: no switching is done:
                switching = false;
                rows = table.rows;
                /*Loop through all table rows (except the
                first, which contains table headers):*/
                for (i = 1; i < (rows.length - 1); i++) {
                    //start by saying there should be no switching:
                    shouldSwitch = false;
                    /*Get the two elements you want to compare,
                    one from current row and one from the next:*/
                    x = rows[i].getElementsByTagName("TD")[n];
                    y = rows[i + 1].getElementsByTagName("TD")[n];
                    /*check if the two rows should switch place,
                    based on the direction, asc or desc:*/
                    if (dir == "asc") {
                        if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase()) {
                            //if so, mark as a switch and break the loop:
                            shouldSwitch= true;
                            break;
                        }
                    } else if (dir == "desc") {
                        if (x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase()) {
                            //if so, mark as a switch and break the loop:
                            shouldSwitch = true;
                            break;
                        }
                    }
                }
                if (shouldSwitch) {
                    /*If a switch has been marked, make the switch
                    and mark that a switch has been done:*/
                    rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
                    switching = true;
                    //Each time a switch is done, increase this count by 1:
                    switchcount ++;
                } else {
                    /*If no switching has been done AND the direction is "asc",
                    set the direction to "desc" and run the while loop again.*/
                    if (switchcount == 0 && dir == "asc") {
                        dir = "desc";
                        switching = true;
                    }
                }
            }
        }
    </script>`)
