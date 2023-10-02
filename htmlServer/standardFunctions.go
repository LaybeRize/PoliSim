package htmlServer

import (
	"PoliSim/componentHelper"
	"PoliSim/dataExtraction"
	"PoliSim/dataValidation"
	"PoliSim/database"
	"PoliSim/htmlComposition"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// GetFullPage returns a http.HandlerFunc that writes a full page to the response
// with a div that automatically requests the URL via htmx.
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

// extractValuesForFields reads in all the fields of a struct and that has the "input" tag
// it looks up that tag as a form field and writes the value depending on the type of the struct field back into it.
// Any fields without the "input" tag will be ignored.
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

// getText extracs the first string in http.Request PostForm field
// and trims it before returning.
func getText(r *http.Request, fieldName string) string {
	return strings.TrimSpace(r.PostFormValue(fieldName))
}

// getSliceAsValue reads in the string slice of the requested PostForm field
// trims every entry and removes any empty entries and doubled entries.
func getSliceAsValue(r *http.Request, fieldName string) reflect.Value {
	slice := dataValidation.DeleteMultiplesAndEmpty(r.PostForm[fieldName])
	return reflect.ValueOf(slice)
}

// getBool checks if the requested field contains the text "true" and returns that value
func getBool(r *http.Request, fieldName string) bool {
	return getText(r, fieldName) == "true"
}

// getInt reads in the field then tries to transform it to a number. On sucess it returns the
// read in number on error it returns the onError int provided.
func getInt(r *http.Request, fieldName string, onError int64) int64 {
	i, err := strconv.Atoi(getText(r, fieldName))
	if err != nil {
		return onError
	}
	return int64(i)
}

func CheckUserPrivilges(w http.ResponseWriter, r *http.Request, roleString ...database.RoleString) (*dataExtraction.AccountAuth, bool) {
	inCookie, err := r.Cookie("token")
	if err != nil {
		return &dataExtraction.AccountAuth{Role: database.NotLoggedIn}, false
	}

	acc, cookie := dataValidation.ValidateToken(inCookie.Value)
	if cookie != nil {
		http.SetCookie(w, cookie)
	}

	return acc, CheckIfHasRole(acc, roleString...)
}

func CheckIfHasRole(acc *dataExtraction.AccountAuth, roles ...database.RoleString) bool {
	for _, r := range roles {
		if r == acc.Role {
			return true
		}
	}
	return false
}
