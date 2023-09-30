package htmlServer

import (
	"PoliSim/componentHelper"
	"io"
	"net/http"
)

// renderRequest renders the request to the response writers, if the write fails a
// http.StatusInternalServerError is put in the provided http.ResponseWriter.
// addDoc adds the <!DOCTYPE html> to the start of the response.
func renderRequest(w http.ResponseWriter, addDoc bool, f func(io.Writer) error) {
	if addDoc {
		err := componentHelper.RenderHTMLDoc(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	err := f(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
