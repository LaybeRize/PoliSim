package htmlComposition

import (
	. "PoliSim/componentHelper"
)

func GetNotFoundPage() Node {
	return getBasePageWrapper(
		El(DIV, Attr(CLASS, "h-full flex items-center fadeMeIn"),
			El(DIV, Attr(CLASS, "box box-e flex-col flex"),
				Attr(STYLE, "padding: 0.5em; line-height: 1; justify-content: center; align-items: center;"),
				El(P, Attr(STYLE, "font-size: 10em; margin-top: -10px"), Text("404")),
				El(P, Attr(STYLE, "font-size: 2em; margin-bottom: 5px; margin-left: 10px; margin-right: 10px"),
					Text(Translation["pageNotFoundText"]),
				),
			),
		),
	)
}
