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

// GetBasePage returns the typical <html> frame for the page.
// it includes the correct sidebar for the parsed role and an element which automatically
// calls the correct partial for the page content.
func GetBasePage(pageTitle string, role database.RoleLevel, loadURL string, loadURLAddtion string) Node {
	return El(HTML, Attr(LANG, language),
		El(HEAD,
			El(META, Attr(CHARSET, "UTF-8")),
			El(TITLE, Text(pageTitle)),
			El(META, Attr(NAME, "viewport"), Attr(CONTENT, "width=device-width, initial-scale=1")),
			El(LINK, Attr(REL, "icon"), Attr(HREF, Configuration["logo"])),
			El(LINK, Attr(REL, "stylesheet"), Attr(HREF, "/public/jsdelivr.css")),
			El(LINK, Attr(REL, "stylesheet"), Attr(HREF, "/public/design.css")),
			El(SCRIPT, Attr(SRC, "/public/_hyperscript_0.9.11.js")),
			El(SCRIPT, Attr(SRC, "/public/htmx_1.9.5.js")),
		),
		El(BODY, Attr(CLASS, "bg-slate-800 min-h-screen text-slate-200"),
			El(DIV, Attr(ID, InformationID), Attr(HIDDEN),
				El(INPUT, Attr(NAME, "personalRoleLevel"), Attr(VALUE, strconv.Itoa(int(role))), Attr(TYPE, "hidden")),
				El(INPUT, Attr(NAME, "currentPageURL"), Attr(VALUE, loadURL), Attr(TYPE, "hidden")),
			),
			El(DIV, Attr(CLASS, "flex flex-row"),
				getSidebar(role, nil),
				El(DIV, Attr(ID, MainBodyID), Attr(HXGET, "/"+APIPreRoute+loadURL+loadURLAddtion), Attr(HXTRIGGER, "load"), Attr(HXSWAP, "outerHTML"), Attr(HXINCLUDE, "#"+InformationID)),
			),
		))
}

// getBasePageWrapper wraps the children in the MainBodyID div (and now an addition div to fucking standardize the fade in affect).
func getBasePageWrapper(children ...Node) Node {
	return El(DIV, Attr(ID, MainBodyID), Attr(CLASS, "flex items-center flex-col basePadding w-full"),
		El(DIV, Attr(CLASS, "flex items-center flex-col h-full fadeMeIn"), Group(children...)),
	)
}

// getCustomPageHeader returns a <h1> element for the header text on the page
func getCustomPageHeader(text string) Node {
	return El(H1, Attr(CLASS, "text-3xl font-bold mt-3"), Text(text))
}

// getPageHeader returns a <h1> element for the header text filled with the PageTitleMap value for the given HttpUrl
func getPageHeader(url HttpUrl) Node {
	return getCustomPageHeader(PageTitleMap[url])
}

// GetInfoDiv returns a hx-swap-oob <div> element containing the role and pageURL as inputs
func GetInfoDiv(role database.RoleLevel, pageURL HttpUrl) Node {
	return El(DIV, Attr(ID, InformationID), Attr(HXSWAPOOB, "true"), Attr(HIDDEN),
		El(INPUT, Attr(NAME, "personalRoleLevel"), Attr(VALUE, strconv.Itoa(int(role))), Attr(TYPE, "hidden")),
		El(INPUT, Attr(NAME, "currentPageURL"), Attr(VALUE, string(pageURL)), Attr(TYPE, "hidden")))
}

// GetMessage returns the message div, colored and filled correctly based on the parameter. (invisiable when
// the string is empty, green when the message is positive)
func GetMessage(val dataValidation.ValidationMessage) Node {
	return El(DIV, Attr(ID, MessageID),
		El(P, If(val.Message == "", Attr(HIDDEN)),
			IfElse(val.Positive, Attr(CLASS, "text-white p-2 mt-2 bg-emerald-800"),
				Attr(CLASS, "text-white p-2 mt-2 bg-rose-800")),
			Text(val.Message),
		))
}

// GetTitleReplacement returns a new <title> element for htmx to swap with the correct PageTitle
func GetTitleReplacement(url HttpUrl) Node {
	return El(TITLE, Attr(HXSWAPOOB, "true"), Text(PageTitleMap[url]))
}

// GetSidebarReplacement gets a new sidebar based on the database.RoleLevel that has the hx-swap-oob parameter
func GetSidebarReplacement(level database.RoleLevel) Node {
	return getSidebar(level, Attr(HXSWAPOOB, "true"))
}
