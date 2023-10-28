package composition

import (
	"PoliSim/data/database"
	. "PoliSim/html/builder"
)

// getSidebar returns a full <div> with every button needed for navigation
func getSidebar(level database.RoleLevel, specialNode Node) Node {
	return DIV(specialNode, ID(SidebarID), CLASS("lg:left-0 p-2 sidebarSize min-h-screen text-center bg-gray-900"),
		DIV(CLASS("text-gray-100 text-xl"),
			DIV(CLASS("p-2.5 mt-1 flex items-center"),
				IMG(SRC(Configuration["logo"]), ALT("Logo"),
					CLASS("h-[25px] ml-3")),
				H1(CLASS("font-bold text-gray-200 text-[15px] ml-3"), Text(Configuration["projectName"])),
			),
			getSidebarBreaker(),
		),
		//here come the inputs
		getSidebarButton(level, database.NotLoggedIn, Start),
		getSidebarButton(level, database.NotLoggedIn, ViewTitles),
		getSidebarButton(level, database.NotLoggedIn, ViewOrganisations),
		If(database.User <= level, getSidebarBreaker()),
		getSidebarButton(level, database.User, CreatePressRelease),
		If(database.MediaAdmin <= level, getSidebarBreaker()),
		getSidebarButton(level, database.MediaAdmin, ViewHiddenNewspaperList),
		If(database.Admin <= level, getSidebarBreaker()),
		getSideBarSubMenu(level, database.Admin, Translation["organisationSubMenu"],
			getSidebarSubMenuButton(level, database.Admin, CreateOrganisation),
			getSidebarSubMenuButton(level, database.Admin, EditOrganisation),
			getSidebarSubMenuButton(level, database.Admin, ViewHiddenOrganisations),
		),
		getSideBarSubMenu(level, database.Admin, Translation["titleSubMenu"],
			getSidebarSubMenuButton(level, database.Admin, CreateTitle),
			getSidebarSubMenuButton(level, database.Admin, EditTitle),
		),
		getSideBarSubMenu(level, database.HeadAdmin, Translation["accountSubMenu"],
			getSidebarSubMenuButton(level, database.HeadAdmin, CreateUser),
			getSidebarSubMenuButton(level, database.HeadAdmin, EditUser),
			getSidebarSubMenuButton(level, database.HeadAdmin, ViewUser),
		),
	)
}

// getSidebarBreaker returns a <div> that functions as a space break
func getSidebarBreaker() Node {
	return DIV(CLASS("my-2 bg-gray-600 h-[1px]"))
}

// getSidebarButton returns a button if the userLevel is as high or higher than minimumLevel for the given url
func getSidebarButton(userLevel database.RoleLevel, minimumLevel database.RoleLevel, url HttpUrl) Node {
	if minimumLevel > userLevel {
		return A(ID(string(url)+SidebarID), HIDDEN())
	}
	return A(HXGET("/"+APIPreRoute+string(url)), ID(string(url)+SidebarID), HXTARGET("#"+MainBodyID),
		HXPUSHURL("/"+string(url)), HXSWAP("outerHTML"), HYPERSCRIPT(getClickAction(url)),
		CLASS("p-2.5 mt-3 flex items-center px-4 duration-300 cursor-pointer text-white hover:bg-blue-600"),
		SPAN(CLASS("text-[15px] ml-4 text-gray-200 font-bold"), Text(SidebarTitleMap[url])),
	)
}

// getSideBarSubMenu returns a wrapper for submenu buttons. It can hide and show the children buttons via a click
func getSideBarSubMenu(userLevel database.RoleLevel, minimumLevel database.RoleLevel, subMenuName string, children ...Node) Node {
	if minimumLevel > userLevel {
		return DIV(ID(subMenuName+SidebarID), HIDDEN())
	}
	return DIV(ID(subMenuName+SidebarID),
		DIV(CLASS("p-2.5 mt-3 flex items-center px-4 duration-300 cursor-pointer text-white hover:bg-blue-600"),
			HYPERSCRIPT("on click toggle .hidden on next <div/> from me then toggle .rotate-180 on last <span/> in first <div/> in me"),
			DIV(CLASS("flex justify-between w-full items-center"),
				SPAN(CLASS("text-[15px] ml-4 text-gray-200 font-bold"), Text(subMenuName)),
				SPAN(CLASS("text-sm rotate-180"),
					I(CLASS("bi bi-chevron-down")),
				),
			),
		),
		DIV(Group(children...), CLASS("text-left text-sm mt-2 w-4/5 mx-auto text-gray-200 font-bold hidden")),
	)
}

// getSidebarSubMenuButton returns a button specially made for the getSideBarSubMenu wrapper.
func getSidebarSubMenuButton(userLevel database.RoleLevel, minimumLevel database.RoleLevel, url HttpUrl) Node {
	if minimumLevel > userLevel {
		return A(ID(string(url)+SidebarID), HIDDEN())
	}
	return A(HXGET("/"+APIPreRoute+string(url)), ID(string(url)+SidebarID), HXTARGET("#"+MainBodyID),
		HXPUSHURL("/"+string(url)), HXSWAP("outerHTML"), HYPERSCRIPT(getClickAction(url)),
		H1(CLASS("cursor-pointer p-2 mt-1 w-full hover:bg-blue-600"), Text(SidebarTitleMap[url])),
	)
}

// getClickAction returns the hyperscript for the middle click for the button to open a second tab
func getClickAction(link HttpUrl) string {
	return "on auxclick[button==1] call window.open('/" + string(link) + "', '_blank')"
}
