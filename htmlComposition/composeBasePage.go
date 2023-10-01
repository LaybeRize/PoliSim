package htmlComposition

import (
	. "PoliSim/componentHelper"
)

func GetBasePage(pageTitle string, loadURL string) Node {
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
			getSidebar(),
			El(DIV, Attr(HXGET, "/htmx"+loadURL), Attr(HXTRIGGER, "load"), Attr(HXSWAP, "outerHTML"))))
}
