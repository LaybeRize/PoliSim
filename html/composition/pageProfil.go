package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"sort"
	"strings"
)

func GetPersonalProfil(acc *database.AccountAuth) Node {
	val := validation.Message{}
	var list extraction.AccountList
	var err error
	nodes := make([]Node, 0, 20)
	list = extraction.AccountList{}
	if acc.ID != 0 {
		list, err = extraction.ReturnAccountList(acc.ID)
		if err != nil {
			errorAccountData(&val)
		}
	}
	var orgs *database.OrganisationList
	var titles *database.TitleList
	sort.Slice(list, func(i, j int) bool {
		if list[i].DisplayName == acc.DisplayName {
			return true
		}
		if list[j].DisplayName == acc.DisplayName {
			return false
		}
		return list[i].DisplayName < list[j].DisplayName
	})
	for _, item := range list {
		orgs, err = extraction.GetOrganisationsForUser(item.ID)
		if err != nil {
			errorAccountData(&val)
		}
		titles, err = extraction.GetTitlesForUser(item.ID)
		if err != nil {
			errorAccountData(&val)
		}
		nodes = append(nodes, getRowForUser(item.DisplayName, item.Flair, orgs, titles))
	}
	return getBasePageWrapper(
		getPageHeader(ViewSelf),
		getStandardTable("sortTable",
			TR(
				getTableHeader(StartPos, -1, Translation["selfTableDisplayName"]),
				getTableHeader(MiddlePos, -1, Translation["selfTableFlair"]),
				getTableHeader(MiddlePos, -1, Translation["selfTableOrganisations"]),
				getTableHeader(EndPos, -1, Translation["selfTableTitles"]),
			),
			Group(nodes...),
			TR(STYLE("height: 10px;")),
		),
		GetLoginThing(false),
		GetMessage(val),
	)
}

func getRowForUser(accName string, accFlair string, orgs *database.OrganisationList, titles *database.TitleList) Node {
	orgArray := make([]string, len(*orgs))
	titleArray := make([]string, len(*titles))
	for i, org := range *orgs {
		orgArray[i] = org.Name
	}
	for i, title := range *titles {
		titleArray[i] = title.Name
	}
	return TR(
		getTableElement(StartPos, 1, Text(accName)),
		getTableElement(MiddlePos, 1, Text(accFlair)),
		getTableElement(MiddlePos, 1, Text(strings.Join(orgArray, ", "))),
		getTableElement(EndPos, 1, Text(strings.Join(titleArray, ", "))),
	)
}

func errorAccountData(val *validation.Message) {
	if val.Message != "" {
		return
	}
	val.Message = Translation["errorAccountData"]
}

func GetLoginThing(swap bool) Node {
	return DIV(ID("password-div-id"), If(swap, HXSWAPOOB("true")),
		getFormStandardForm("form", PATCH, "/"+HTMXPreRouter+string(ChangePassword),
			getInput("ordPassword", "ordPassword", "", Translation["ordPassword"],
				"password", "", ""),
			getInput("newPassword", "newPassword", "", Translation["newPassword"],
				"password", "", ""),
			getInput("newPasswordAgain", "newPasswordAgain", "", Translation["newPasswordAgain"],
				"password", "", ""),
			getSubmitButton("loginButton", Translation["changePasswordButton"])),
	)
}
