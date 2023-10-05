package htmlComposition

import (
	. "PoliSim/componentHelper"
	"PoliSim/database"
)

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
		getSidebarButton(level, database.HeadAdmin, Test),
		If(database.User <= level, getSidebarBreaker()),
		If(database.MediaAdmin <= level, getSidebarBreaker()),
		If(database.Admin <= level, getSidebarBreaker()),
	)
}

func getSidebarBreaker() Node {
	return El(DIV, Attr(CLASS, "my-2 bg-gray-600 h-[1px]"))
}

func getSidebarButton(userLevel database.RoleLevel, minimumLevel database.RoleLevel, url HttpUrl) Node {
	return generalSidebarButton(minimumLevel > userLevel, url, false)
}

func GetSidebarButtonForSwap(hidden bool, url HttpUrl) Node {
	return generalSidebarButton(hidden, url, true)
}

func generalSidebarButton(hidden bool, url HttpUrl, addSwap bool) Node {
	if hidden {
		return El(A, Attr(ID, string(url)+SidebarID), If(addSwap, Attr(HXSWAPOOB, "true")), Attr(HIDDEN))
	}
	return El(A, Attr(HXGET, "/"+APIPreRoute+string(url)), Attr(ID, string(url)+SidebarID), Attr(HXTARGET, "#"+MainBodyID),
		Attr(HXINCLUDE, "#"+InformationID),
		Attr(HYPERSCRIPT, getClickAction(url)), Attr(CLASS, "p-2.5 mt-3 flex items-center px-4 duration-300 cursor-pointer text-white hover:bg-blue-600"),
		If(addSwap, Attr(HXSWAPOOB, "true")),
		El(SPAN, Attr(CLASS, "text-[15px] ml-4 text-gray-200 font-bold"), Text(SidebarTitleMap[url])),
	)
}

func getClickAction(link HttpUrl) string {
	return "on auxclick[button==1] call window.open('/" + string(link) + "', '_blank')"
}
