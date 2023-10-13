package composition

import (
	. "PoliSim/html/builder"
)

func GetNotFoundPage() Node {
	return getBasePageWrapper(
		DIV(CLASS("h-full flex items-center"),
			DIV(CLASS("box box-e flex-col flex"),
				STYLE("padding: 0.5em; line-height: 1; justify-content: center; align-items: center;"),
				P(STYLE("font-size: 10em; margin-top: -10px"), Text("404")),
				P(STYLE("font-size: 2em; margin-bottom: 5px; margin-left: 10px; margin-right: 10px"),
					Text(Translation["pageNotFoundText"]),
				),
			),
		),
	)
}
