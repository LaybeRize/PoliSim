package htmlComposition

import (
	. "PoliSim/componentHelper"
	"PoliSim/dataValidation"
	"PoliSim/database"
	"os"
	"strconv"
	"strings"
)

var language = strings.ToLower(os.Getenv("LANG"))

func GetBasePage(pageTitle string, role database.RoleLevel, loadURL string) Node {
	return El(HTML, Attr(LANG, language),
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
			El(DIV, Attr(ID, InformationID), Attr(HIDDEN)),
			El(DIV, Attr(CLASS, "flex flex-row"),
				getSidebar(role, nil),
				El(DIV, Attr(ID, MainBodyID), Attr(HXGET, "/htmx"+loadURL), Attr(HXTRIGGER, "load"), Attr(HXSWAP, "outerHTML")),
			)))
}

func getBasePageWrapper(children ...Node) Node {
	return El(DIV,
		append(children, Attr(ID, MainBodyID), Attr(CLASS, "flex items-center basePadding flex-col w-full"))...,
	)
}

func getCustomPageHeader(text string) Node {
	return El(H1, Attr(CLASS, "text-3xl font-bold mt-3"), Text(text))
}

func getPageHeader(url HttpUrl) Node {
	return getCustomPageHeader(PageTitleMap[url])
}

func GetInfoDiv(role database.RoleLevel, pageURL HttpUrl) Node {
	return El(DIV, Attr(ID, InformationID), Attr(HXSWAPOOB, "true"), Attr(HIDDEN),
		El(INPUT, Attr(NAME, "personalRoleLevel"), Attr(VALUE, strconv.Itoa(int(role))), Attr(TYPE, "hidden")),
		El(INPUT, Attr(NAME, "currentPageURL"), Attr(VALUE, string(pageURL)), Attr(TYPE, "hidden")))
}

func GetMessage(val dataValidation.ValidationMessage) Node {
	return El(DIV, Attr(ID, MessageID),
		El(P, If(val.Message == "", Attr(HIDDEN)),
			IfElse(val.Positive, Attr(CLASS, "text-white p-2 mt-2 bg-emerald-800"),
				Attr(CLASS, "text-white p-2 mt-2 bg-rose-800")),
			Text(val.Message),
		))
}

func GetTitleReplacement(url HttpUrl) Node {
	return El(TITLE, Attr(HXSWAPOOB, "true"), Text(PageTitleMap[url]))
}

func GetSidebarReplacement(level database.RoleLevel) Node {
	return getSidebar(level, Attr(HXSWAPOOB, "true"))
}
