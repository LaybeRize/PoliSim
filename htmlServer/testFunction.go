package htmlServer

import (
	. "PoliSim/componentHelper"
	"net/http"
)

func ServeTestGet(w http.ResponseWriter, r *http.Request) {
	err := RenderHTMLDoc(w,
		El(HEAD, El(TITLE, Text("Test"))),
		El(BODY, Raw("<p>test</p>"), El(DIV, Attr(HXPOST, "/test"), Text(Translation["test"]))))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
