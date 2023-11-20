package composition

import (
	"PoliSim/data/database"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
	"time"
)

func getDocumentHead(doc *database.Document, isAdmin bool, extra ...Node) Node {
	addition := ""
	if doc.Blocked {
		addition = " text-rose-600"
	}
	return DIV(CLASS("w-[800px]"), ID("document-header-div"),
		If(isAdmin, A(HXPATCH("/"+APIPreRoute+string(BlockDocumentLink)+url.PathEscape(doc.UUID)), HXTARGET("#"+MessageID),
			ID("block-button-document"), TEST("block-button-document"), HXSWAP("outerHTML"),
			CLASS("bg-slate-700 text-white p-2"),
			IfElse(doc.Blocked, Text(Translation["unblockDocument"]), Text(Translation["blockDocument"])),
		)),
		H1(CLASS("text-3xl underline decoration-2 underline-offset-2 mt-2"+addition), If(doc.Blocked, Text(Translation["documentBlockedText"])), Text(doc.Title)),
		If(doc.Subtitle.Valid, H1(CLASS("text-2xl"), Text(doc.Subtitle.String))),
		P(CLASS("my-2"), I(Text(fmt.Sprintf(doc.Written.Format(Translation["authorSummaryDocument"]), doc.Organisation, doc.Author))),
			If(doc.Flair != "", Group(I(Text("; ")), Text(doc.Flair)))),
		Group(extra...),
	)
}

func getDocumentBody(doc *database.Document) Node {
	return DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
		Raw(doc.HTMLContent),
	)
}

func scriptForUpdateOnEnd(doc *database.Document, httpUrl HttpUrl) Node {
	return Group(DIV(ID("trigger-me-on-document-finish"),
		HXTRIGGER("pageReloaded"), HXGET("/"+APIPreRoute+string(httpUrl)+url.PathEscape(doc.UUID)),
		HXTARGET("#"+MessageID), HXSWAP("outerHTML")),
		SCRIPT(Raw(`
		function timeForRefresh() {
			htmx.trigger("#trigger-me-on-document-finish", "pageReloaded");
   		}
		var timeEnd = new Date("`+doc.Info.Finishing.Format(time.RFC3339)+`").getTime();
		var currentTime = new Date().getTime();
		var subtractMilliSecondsValue = timeEnd - currentTime;
		if (subtractMilliSecondsValue < 0) {
			timeForRefresh();
		} else {
			setTimeout(timeForRefresh, subtractMilliSecondsValue);
		}
`)))
}

func GetNewDocumentHeader(doc *database.Document) Node {
	extra := Node(nil)
	if doc.Type == database.LegislativeText && len(doc.Info.Post) != 0 {
		extra = DIV(CLASS("mt-2 w-[800px]"),
			renderTag(doc.Info.Post[0], ""))
	}
	return getDocumentHead(doc, true, HXSWAPOOB("true"), extra)
}
