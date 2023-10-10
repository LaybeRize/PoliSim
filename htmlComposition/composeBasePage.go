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
func GetBasePage(pageTitle string, role database.RoleLevel, loadURL string, loadURLAddition string) Node {
	return HTML(LANG(language),
		HEAD(
			META(CHARSET("UTF-8")),
			TITLE(Text(pageTitle)),
			META(NAME("viewport"), CONTENT("width=device-width, initial-scale=1")),
			LINK(REL("icon"), HREF(Configuration["logo"])),
			LINK(REL("stylesheet"), HREF("/public/jsdelivr.css")),
			LINK(REL("stylesheet"), HREF("/public/design.css")),
			SCRIPT(SRC("/public/_hyperscript_0.9.11.js")),
			SCRIPT(SRC("/public/htmx_1.9.5.js")),
		),
		BODY(CLASS("bg-slate-800 min-h-screen text-slate-200"),
			DIV(ID(InformationID), HIDDEN(),
				INPUT(NAME("personalRoleLevel"), VALUE(strconv.Itoa(int(role))), TYPE("hidden")),
				INPUT(NAME("currentPageURL"), VALUE(loadURL), TYPE("hidden")),
			),
			DIV(CLASS("flex flex-row"),
				getSidebar(role, nil),
				DIV(ID(MainBodyID), HXGET("/"+APIPreRoute+loadURL+loadURLAddition), HXTRIGGER("load"), HXSWAP("outerHTML"), HXINCLUDE("#"+InformationID)),
			),
		))
}

// getBasePageWrapper wraps the children in the MainBodyID div (and now an addition div to fucking standardize the fade in effect).
func getBasePageWrapper(children ...Node) Node {
	return DIV(ID(MainBodyID), CLASS("flex items-center flex-col basePadding w-full minSizeBase"),
		DIV(CLASS("flex items-center flex-col h-full fadeMeIn"), Group(children...)),
	)
}

// getCustomPageHeader returns a <h1> element for the header text on the page
func getCustomPageHeader(text string) Node {
	return H1(CLASS("text-3xl font-bold mt-3"), Text(text))
}

// getPageHeader returns a <h1> element for the header text filled with the PageTitleMap value for the given HttpUrl
func getPageHeader(url HttpUrl) Node {
	return getCustomPageHeader(PageTitleMap[url])
}

// GetInfoDiv returns a hx-swap-oob <div> element containing the role and pageURL as inputs
func GetInfoDiv(role database.RoleLevel, pageURL HttpUrl) Node {
	return DIV(ID(InformationID), HXSWAPOOB("true"), HIDDEN(),
		INPUT(NAME("personalRoleLevel"), VALUE(strconv.Itoa(int(role))), TYPE("hidden")),
		INPUT(NAME("currentPageURL"), VALUE(string(pageURL)), TYPE("hidden")))
}

// GetMessage returns the message div, colored and filled correctly based on the parameter. (invisible when
// the string is empty, green when the message is positive)
func GetMessage(val dataValidation.ValidationMessage) Node {
	return DIV(ID(MessageID),
		P(If(val.Message == "", HIDDEN()),
			IfElse(val.Positive, CLASS("text-white p-2 mt-2 bg-emerald-800"),
				CLASS("text-white p-2 mt-2 bg-rose-800")),
			Text(val.Message),
		))
}

// GetTitleReplacement returns a new <title> element for htmx to swap with the correct PageTitle
func GetTitleReplacement(url HttpUrl) Node {
	return TITLE(HXSWAPOOB("true"), Text(PageTitleMap[url]))
}

// GetSidebarReplacement gets a new sidebar based on the database.RoleLevel that has the hx-swap-oob parameter
func GetSidebarReplacement(level database.RoleLevel) Node {
	return getSidebar(level, HXSWAPOOB("true"))
}

func GetErrorPage(errorText string) Node {
	return getBasePageWrapper(DIV(CLASS("h-full flex items-center justify-center"),
		DIV(STYLE("padding: 0.5em; line-height: 1; justify-content: center; align-items: center;--clr-border: rgb(159 18 57); background-size: 4px 100%, 100% 4px, 4px 100% , 100% 4px;"),
			CLASS("box box-e flex-col flex"),
			P(STYLE("font-size: 5em; margin-top: 3px; margin-left: 10px; margin-right: 10px"), CLASS("text-rose-600"),
				Text(Translation["errorPageTitle"])),
			P(STYLE("font-size: 2em; margin-top: 8px; margin-bottom: 21px;margin-left: 10px; margin-right: 10px"), CLASS("text-rose-600"),
				Text(errorText)),
		),
	))
}
