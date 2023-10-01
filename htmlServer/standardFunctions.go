package htmlServer

import (
	"PoliSim/componentHelper"
	"PoliSim/dataValidation"
	"PoliSim/htmlComposition"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func GetFullPage(pageTitle string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addition := "?" + r.URL.RawQuery
		if addition == "?" {
			addition = ""
		}
		html := htmlComposition.GetBasePage(pageTitle, r.URL.Path+addition)
		renderRequest(w, true, html.Render)
	}
}

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

func extractValuesForFields(object any, r *http.Request, onError int64) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	v := reflect.ValueOf(object)
	iterate := reflect.Indirect(v).NumField()
	for i := 0; i < iterate; i++ {
		input, ok := v.Type().Elem().Field(i).Tag.Lookup("input")
		if !ok {
			continue
		}
		kind := reflect.Indirect(v).Field(i).Kind()
		switch kind {
		case reflect.String:
			v.Elem().Field(i).SetString(getText(r, input))
		case reflect.Slice:
			v.Elem().Field(i).Set(getSliceAsValue(r, input))
		case reflect.Bool:
			v.Elem().Field(i).SetBool(getBool(r, input))
		case reflect.Int, reflect.Int64:
			v.Elem().Field(i).SetInt(getInt(r, input, onError))
		}
	}
	return nil
}

func getText(r *http.Request, fieldName string) string {
	return strings.TrimSpace(r.PostFormValue(fieldName))
}

func getSliceAsValue(r *http.Request, fieldName string) reflect.Value {
	slice := dataValidation.DeleteMultiplesAndEmpty(r.PostForm[fieldName])
	return reflect.ValueOf(slice)
}

func getBool(r *http.Request, fieldName string) bool {
	return getText(r, fieldName) == "true"
}

func getInt(r *http.Request, fieldName string, onError int64) int64 {
	i, err := strconv.Atoi(getText(r, fieldName))
	if err != nil {
		return onError
	}
	return int64(i)
}
