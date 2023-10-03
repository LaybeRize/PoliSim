package htmlComposition

import (
	. "PoliSim/componentHelper"
	"PoliSim/database"
)

func GetBasePage(pageTitle string, role database.RoleLevel, loadURL string) Node {
	return El(HTML,
		El(HEAD,
			El(META, Attr(CHARSET, "UTF-8")),
			El(TITLE, Text(pageTitle)),
			El(META, Attr(NAME, "viewport"), Attr(CONTENT, "width=device-width, initial-scale=1")),
			El(LINK, Attr(REL, "icon"), Attr(HREF, Configuration["logo"])),
			El(LINK, Attr(REL, "stylesheet"), Attr(HREF, "/public/jsdelivr.css")),
			El(LINK, Attr(REL, "stylesheet"), Attr(HREF, "/public/design.css")),
			El(SCRIPT, Attr(SRC, "/public/_hyperscript_0.9.11.min.js")),
			El(SCRIPT, Attr(SRC, "/public/htmx_1.9.5.min.js")),
		),
		El(BODY, Attr(CLASS, "bg-slate-800 min-h-screen text-slate-200"),
			El(DIV, Attr(CLASS, "flex flex-row"),
				getSidebar(role),
				El(DIV, Attr(ID, MainBodyID), Attr(HXGET, "/htmx"+loadURL), Attr(HXTRIGGER, "load"), Attr(HXSWAP, "outerHTML")),
			)))
}

func getBasePageWrapper(children ...Node) Node {
	children = append(children, Attr(ID, MainBodyID), Attr(CLASS, "flex items-center pl-2 flex-col w-full"))
	return El(DIV, children...)
}

func getCustomPageHeader(text string) Node {
	return El(H1, Attr(CLASS, "text-3xl font-bold mt-3"), Text(text))
}

func getPageHeader(url HttpUrl) Node {
	return getCustomPageHeader(PageTitleMap[url])
}
