package htmlComposition

import (
	. "PoliSim/componentHelper"
	"PoliSim/database"
)

// getSidebar returns a full <div> with every button needed for navigation
func getSidebar(level database.RoleLevel, specialNode Node) Node {
	return El(DIV, specialNode, Attr(ID, SidebarID), Attr(CLASS, "lg:left-0 p-2 sidebarSize min-h-screen text-center bg-gray-900"),
		El(DIV, Attr(CLASS, "text-gray-100 text-xl"),
			El(DIV, Attr(CLASS, "p-2.5 mt-1 flex items-center"),
				El(IMG, Attr(SRC, Configuration["logo"]), Attr(ALT, "Logo"),
					Attr(CLASS, "h-[25px] ml-3")),
				El(H1, Attr(CLASS, "font-bold text-gray-200 text-[15px] ml-3"), Text(Configuration["projectName"])),
			),
			getSidebarBreaker(),
		),
		//here come the inputs
		getSidebarButton(level, database.NotLoggedIn, Start),
		getSidebarButton(level, database.HeadAdmin, "test"),
		// TODO: remove me ^
		If(database.User <= level, getSidebarBreaker()),
		If(database.MediaAdmin <= level, getSidebarBreaker()),
		If(database.Admin <= level, getSidebarBreaker()),
		getSideBarSubMenu(level, database.HeadAdmin, Translation["accountSubMenu"],
			getSidebarSubMenuButton(level, database.HeadAdmin, "test"),
		),
	)
}

// getSidebarBreaker returns a <div> that functions as a space break
func getSidebarBreaker() Node {
	return El(DIV, Attr(CLASS, "my-2 bg-gray-600 h-[1px]"))
}

// getSidebarButton returns a button if the userLevel is as high or higher than minimumLevel for the given url
func getSidebarButton(userLevel database.RoleLevel, minimumLevel database.RoleLevel, url HttpUrl) Node {
	if minimumLevel > userLevel {
		return El(A, Attr(ID, string(url)+SidebarID), Attr(HIDDEN))
	}
	return El(A, Attr(HXGET, "/"+APIPreRoute+string(url)), Attr(ID, string(url)+SidebarID), Attr(HXTARGET, "#"+MainBodyID),
		Attr(HXINCLUDE, "#"+InformationID), Attr(HXVALS, `{"pushURL": "true"}`), Attr(HXSWAP, "outerHTML"),
		Attr(HYPERSCRIPT, getClickAction(url)), Attr(CLASS, "p-2.5 mt-3 flex items-center px-4 duration-300 cursor-pointer text-white hover:bg-blue-600"),
		El(SPAN, Attr(CLASS, "text-[15px] ml-4 text-gray-200 font-bold"), Text(SidebarTitleMap[url])),
	)
}

// getSideBarSubMenu returns a wrapper for submenu buttons. It can hide and show the children buttons via a click
func getSideBarSubMenu(userLevel database.RoleLevel, minimumLevel database.RoleLevel, subMenuName string, children ...Node) Node {
	if minimumLevel > userLevel {
		return El(DIV, Attr(ID, subMenuName+SidebarID), Attr(HIDDEN))
	}
	return El(DIV, Attr(ID, subMenuName+SidebarID),
		El(DIV, Attr(CLASS, "p-2.5 mt-3 flex items-center px-4 duration-300 cursor-pointer text-white hover:bg-blue-600"),
			Attr(HYPERSCRIPT, "on click toggle .hidden on next <div/> from me then toggle .rotate-180 on last <span/> in first <div/> in me"),
			El(DIV, Attr(CLASS, "flex justify-between w-full items-center"),
				El(SPAN, Attr(CLASS, "text-[15px] ml-4 text-gray-200 font-bold"), Text(subMenuName)),
				El(SPAN, Attr(CLASS, "text-sm"),
					El(I, Attr(CLASS, "bi bi-chevron-down")),
				),
			),
		),
		El(DIV, Group(children...), Attr(CLASS, "text-left text-sm mt-2 w-4/5 mx-auto text-gray-200 font-bold")),
	)
}

// getSidebarSubMenuButton returns a button specially made for the getSideBarSubMenu wrapper.
func getSidebarSubMenuButton(userLevel database.RoleLevel, minimumLevel database.RoleLevel, url HttpUrl) Node {
	if minimumLevel > userLevel {
		return El(A, Attr(ID, string(url)+SidebarID), Attr(HIDDEN))
	}
	return El(A, Attr(HXGET, "/"+APIPreRoute+string(url)), Attr(ID, string(url)+SidebarID), Attr(HXTARGET, "#"+MainBodyID),
		Attr(HXINCLUDE, "#"+InformationID), Attr(HXVALS, `{"pushURL": "true"}`), Attr(HXSWAP, "outerHTML"),
		Attr(HYPERSCRIPT, getClickAction(url)),
		El(H1, Attr(CLASS, "cursor-pointer p-2 mt-1 w-full hover:bg-blue-600"), Text(SidebarTitleMap[url])),
	)
}

// getClickAction returns the hyperscript for the middle click for the button to open a second tab
func getClickAction(link HttpUrl) string {
	return "on auxclick[button==1] call window.open('/" + string(link) + "', '_blank')"
}
