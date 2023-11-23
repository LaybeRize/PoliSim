package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"net/url"
	"strings"
)

const (
	joinSeperator       = ", "
	collapseButtonClass = "bg-slate-700 text-white p-2 m-2 disable-selection"
)

func GetViewTitelPage() Node {
	listing := make([]Node, len(extraction.TitleGroupMap))
	counter := 0
	for _, outer := range extraction.TitleMainGroupList {
		innerCounter := 0
		innerListing := make([]Node, len(extraction.TitleGroupMap[outer]))
		for _, inner := range extraction.TitleGroupMap[outer] {
			innerListing[innerCounter] = Group(BUTTON(
				CLASS("text-2xl mt-2 w-full text-left"), Text(inner),
				HXGET("/"+APIPreRoute+string(getTitleSubGroup)+url.PathEscape(outer)+"/"+url.PathEscape(inner)),
				HXTARGET("#out-"+outer+"-in-"+inner), ID("out-"+outer+"-in-"+inner+"-button"),
				HXSWAP("outerHTML"),
			),
				DIV(ID("out-"+outer+"-in-"+inner)))
			innerCounter++
		}
		listing[counter] = Group(BUTTON(CLASS("text-3xl mt-2 w-full text-left"), Text(outer),
			HYPERSCRIPT("on click toggle .hidden on #outer-"+outer)),
			DIV(ID("outer-"+outer), CLASS("border-l-4 border-white pl-6 mt-2 collapse-all hidden"),
				Group(innerListing...)))
		counter++
	}

	return getBasePageWrapper(
		getPageHeader(ViewTitles),
		DIV(CLASS("flex flex-row w-[600px]"),
			BUTTON(TYPE("button"), CLASS(collapseButtonClass),
				HYPERSCRIPT("on click add .hidden to .collapse-all"), Text(Translation["collapseAll"])),
		),
		DIV(CLASS("mt-4 w-[600px]"),
			Group(listing...)),
	)
}

func GetViewSubGroupOfTitles(mainGroup string, subGroup string) Node {
	list, err := extraction.GetAllTitlesInSubGroup(mainGroup, subGroup)
	var newDiv Node = nil
	if err != nil {
		newDiv = DIV(ID("out-"+mainGroup+"-in-"+subGroup), CLASS("border-l-4 border-slate-400 pl-6 mt-2 collapse-all"),
			P(STYLE("font-size: 2em;"), CLASS("text-rose-600"), Text(Translation["errorWhileQueryingTitles"])))
	} else {
		nodeList := make([]Node, len(*list))
		for i, title := range *list {
			holderText := strings.Join(validation.GetDisplayNameArray(&title.Holder), joinSeperator)
			nodeList[i] = DIV(CLASS("mt-2"),
				P(CLASS("text-xl"), Text(title.Name)),
				If(title.Flair.Valid, P(CLASS("text-base mt-2"), Text(Translation["viewTitleFlairText"], title.Flair.String))),
				P(If(!title.Flair.Valid, CLASS("mt-2")), I(CLASS("bi bi-people-fill")), IfElse(holderText == "", Text(" "+Translation["viewTitleNoHolderText"]),
					Text(" "+Translation["viewTitleHolderText"], holderText))))
		}
		newDiv = DIV(ID("out-"+mainGroup+"-in-"+subGroup), CLASS("border-l-4 border-slate-400 pl-6 mt-2 collapse-all"),
			Group(nodeList...))
	}
	return Group(BUTTON(
		CLASS("text-2xl mt-2 w-full text-left"), Text(subGroup),
		ID("out-"+mainGroup+"-in-"+subGroup+"-button"),
		HXSWAPOOB("true"),
		HYPERSCRIPT("on click toggle .hidden on #out-"+mainGroup+"-in-"+subGroup),
	), newDiv)
}

func GetViewOrganisationPage(accountID int64, isAdmin bool) Node {
	var orgGroupings *[][]string
	var err error
	if isAdmin {
		orgGroupings, err = extraction.GetOrganisationGroupingsForAdmins()
	} else {
		orgGroupings, err = extraction.GetOrganisationGroupings(accountID)
	}

	if err != nil {
		return GetErrorPage(Translation["errorWhileLoadingOrganisations"])
	}
	listing := make([]Node, len((*orgGroupings)[0]))
	for i, outer := range (*orgGroupings)[0] {
		innerListing := make([]Node, len((*orgGroupings)[i+1]))
		for pos, inner := range (*orgGroupings)[i+1] {
			innerListing[pos] = Group(BUTTON(
				CLASS("text-2xl mt-2 w-full text-left"), Text(inner),
				HXGET("/"+APIPreRoute+string(getOrganisationSubGroup)+url.PathEscape(outer)+"/"+url.PathEscape(inner)),
				HXTARGET("#out-"+outer+"-in-"+inner), ID("out-"+outer+"-in-"+inner+"-button"),
				HXSWAP("outerHTML"),
			),
				DIV(ID("out-"+outer+"-in-"+inner)))
		}
		listing[i] = Group(BUTTON(CLASS("text-3xl mt-2 w-full text-left"), Text(outer),
			HYPERSCRIPT("on click toggle .hidden on #outer-"+outer)),
			DIV(ID("outer-"+outer), CLASS("border-l-4 border-white pl-6 mt-2 collapse-all hidden"),
				Group(innerListing...)))
	}
	return getBasePageWrapper(
		getPageHeader(ViewOrganisations),
		DIV(CLASS("flex flex-row w-[600px]"),
			BUTTON(TYPE("button"), CLASS(collapseButtonClass),
				HYPERSCRIPT("on click add .hidden to .collapse-all"), Text(Translation["collapseAll"])),
		),
		DIV(CLASS("mt-4 w-[600px]"),
			Group(listing...)),
	)
}

func GetViewSubGroupOfOrganisations(accountID int64, isAdmin bool, mainGroup string, subGroup string) Node {
	var orgs *database.OrganisationList
	var err error
	if isAdmin {
		orgs, err = extraction.GetAllOrganisationsInSubGroupForAdmins(mainGroup, subGroup)
	} else {
		orgs, err = extraction.GetAllOrganisationsInSubGroup(accountID, mainGroup, subGroup)
	}

	var newDiv Node = nil
	if len(*orgs) == 0 || err != nil {
		newDiv = DIV(ID("out-"+mainGroup+"-in-"+subGroup), CLASS("border-l-4 border-slate-400 pl-6 mt-2 collapse-all"),
			P(STYLE("font-size: 2em;"), CLASS("text-rose-600"), Text(Translation["errorWhileQueryingOrganisations"])))
	} else {
		nodeList := make([]Node, len(*orgs))
		for i, organisation := range *orgs {
			memberText := strings.Join(validation.GetDisplayNameArray(&organisation.Members), joinSeperator)
			adminText := strings.Join(validation.GetDisplayNameArray(&organisation.Admins), joinSeperator)
			//<div class="flex items-center">
			//            <p class="text-xl">{ org.Name }</p>
			//    if org.Status == database.Secret {
			//            <i class="text-xl bi bi-eye-slash px-2"></i>
			//	} else {
			//            <i class="text-xl bi bi-eye px-2"></i>
			//    }
			//    if org.Status == database.Private {
			//            <i class="text-xl bi bi-file-lock"></i>
			//    }
			//            </div>
			nodeList[i] = DIV(CLASS("mt-2"),
				DIV(CLASS("flex items-center"),
					P(CLASS("text-xl"), Text(organisation.Name)),
					IfElse(organisation.Status == database.Secret, I(CLASS("text-xl bi bi-eye-slash px-2")),
						I(CLASS("text-xl bi bi-eye px-2"))),
					If(organisation.Status == database.Private, I(CLASS("text-xl bi bi-file-lock")))),
				If(organisation.Flair.Valid, P(CLASS("text-base mt-2"), Text(Translation["viewOrganisationFlairText"], organisation.Flair.String))),
				P(If(!organisation.Flair.Valid, CLASS("mt-2")), I(CLASS("bi bi-people-fill")), IfElse(memberText == "", Text(" "+Translation["viewOrganisationNoMemberText"]),
					Text(" "+Translation["viewOrganisationMemberText"], memberText))),
				P(I(CLASS("bi bi-person-fill-gear")), IfElse(adminText == "", Text(" "+Translation["viewOrganisationNoAdminText"]),
					Text(" "+Translation["viewOrganisationAdminText"], adminText))))

		}
		newDiv = DIV(ID("out-"+mainGroup+"-in-"+subGroup), CLASS("border-l-4 border-slate-400 pl-6 mt-2 collapse-all"),
			Group(nodeList...))
	}
	return Group(BUTTON(
		CLASS("text-2xl mt-2 w-full text-left"), Text(subGroup),
		ID("out-"+mainGroup+"-in-"+subGroup+"-button"),
		HXSWAPOOB("true"),
		HYPERSCRIPT("on click toggle .hidden on #out-"+mainGroup+"-in-"+subGroup),
	), newDiv)
}

func GetViewHiddenOrganisationPage() Node {
	orgs, err := extraction.GetHiddenOrganistaions()
	if err != nil {
		return GetErrorPage(Translation["errorWhileLoadingHiddenOrganisations"])
	}
	nodes := make([]Node, len(*orgs))
	counterMainGroups := 0
	counterSubGroups := 0
	for i := len(nodes) - 1; i > 0; i-- {
		if (*orgs)[i].MainGroup != (*orgs)[i-1].MainGroup {
			nodes[i] = TR(
				getTableElement(StartPos, counterMainGroups+1, Text((*orgs)[i].MainGroup)),
				getTableElement(MiddlePos, counterSubGroups+1, Text((*orgs)[i].SubGroup)),
				getTableElement(EndPos, 1, Text((*orgs)[i].Name)),
			)
			counterMainGroups = 0
			counterSubGroups = 0
			continue
		}
		if (*orgs)[i].SubGroup != (*orgs)[i-1].SubGroup {
			nodes[i] = TR(
				getTableElement(MiddlePos, counterSubGroups+1, Text((*orgs)[i].SubGroup)),
				getTableElement(EndPos, 1, Text((*orgs)[i].Name)),
			)
			counterMainGroups++
			counterSubGroups = 0
			continue
		}
		counterMainGroups++
		counterSubGroups++
		nodes[i] = TR(
			getTableElement(EndPos, 1, Text((*orgs)[i].Name)),
		)
	}
	if len(*orgs) != 0 {
		nodes[0] = TR(
			getTableElement(StartPos, counterMainGroups+1, Text((*orgs)[0].MainGroup)),
			getTableElement(MiddlePos, counterSubGroups+1, Text((*orgs)[0].SubGroup)),
			getTableElement(EndPos, 1, Text((*orgs)[0].Name)),
		)
	}

	return getBasePageWrapper(
		tableNode,
		getPageHeader(ViewHiddenOrganisations),
		getStandardTable("sortTable",
			TR(
				getTableHeader(StartPos, -1, Translation["organisationTableMainGroup"]),
				getTableHeader(MiddlePos, -1, Translation["organisationTableSubGroup"]),
				getTableHeader(EndPos, -1, Translation["organisationTableName"]),
			),
			Group(nodes...),
		),
	)
}
