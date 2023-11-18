package composition

import (
	"PoliSim/data/database"
	. "PoliSim/html/builder"
	"fmt"
)

func getDocumentHead(doc *database.Document, extra ...Node) Node {
	addition := ""
	if doc.Blocked {
		addition = " text-rose-600"
	}
	return DIV(CLASS("w-[800px]"),
		H1(CLASS("text-3xl underline decoration-2 underline-offset-2"+addition), If(doc.Blocked, Text(Translation["documentBlockedText"])), Text(doc.Title)),
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
