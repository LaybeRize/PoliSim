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
		acc, _ := CheckUserPrivilges(w, r)
		addition := "?" + r.URL.RawQuery
		if addition == "?" {
			addition = ""
		}
		html := htmlComposition.GetBasePage(pageTitle, acc.Role, r.URL.Path[1:], addition)
		renderRequest(w, true, html.Render)
	}
}

// renderRequest renders the request to the response writers, if the write fails a
// http.StatusInternalServerError is put in the provided http.ResponseWriter.
// addDoc adds the <!DOCTYPE html> to the start of the response.
func renderRequest(w http.ResponseWriter, addDoc bool, funcs ...func(io.Writer) error) {
	if addDoc {
		err := componentHelper.RenderHTMLDoc(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	for _, f := range funcs {
		if f == nil {
			continue
		}
		err := f(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// extractFormValuesForFields reads in all the fields of a struct and that has the "input" tag
// it looks up that tag as a form field and writes the value depending on the type of the struct field back into it.
// Any fields without the "input" tag will be ignored.
func extractFormValuesForFields(object any, r *http.Request, onError int64) error {
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

// extractUrlValuesForFields does the same for urls what extractFormValuesForFields does for forms
func extractUrlValuesForFields(object any, r *http.Request, onError int64) error {
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
			v.Elem().Field(i).SetString(getTextFromURL(r, input))
		case reflect.Bool:
			v.Elem().Field(i).SetBool(getBoolFromURL(r, input))
		case reflect.Int, reflect.Int64:
			v.Elem().Field(i).SetInt(getIntFromURL(r, input, onError))
		}
	}
	return nil
}

// getTextFromURL reads from the url the specified field and returns the text
func getTextFromURL(r *http.Request, urlField string) string {
	return r.URL.Query().Get(urlField)
}

// getBoolFromURL reads from the url the specified field and returns true if the text is "true"
func getBoolFromURL(r *http.Request, urlField string) bool {
	return r.URL.Query().Get(urlField) == "true"
}

// getIntFromURL reads from the url the specified field and returns it as an int64. if the text can't be converted return onError
func getIntFromURL(r *http.Request, urlField string, onError int64) int64 {
	i, err := strconv.Atoi(r.URL.Query().Get(urlField))
	if err != nil {
		return onError
	}
	return int64(i)
}

// CheckUserPrivilges gets the account from the cookie and returns the dataExtraction.AccountAuth if one is found for the cookie, and if
// the account has one of the specified roles returns true. Otherwise, returns false
func CheckUserPrivilges(w http.ResponseWriter, r *http.Request, roleString ...database.RoleLevel) (*dataExtraction.AccountAuth, bool) {
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

// CheckIfHasRole checks if the referenced account has one of the provided roles
func CheckIfHasRole(acc *dataExtraction.AccountAuth, roles ...database.RoleLevel) bool {
	for _, r := range roles {
		if r == acc.Role {
			return true
		}
	}
	return false
}

// onlySwapMessage retargets the request and only replaces the Message <div> on the page via htmx
func onlySwapMessage(w http.ResponseWriter, val dataValidation.ValidationMessage, f func(io.Writer) error) {
	w.Header().Set("HX-Retarget", "#"+htmlComposition.MessageID)
	html := htmlComposition.GetMessage(val)
	renderRequest(w, false, f, html.Render)
}

type UserInformation struct {
	RoleLevel int    `input:"personalRoleLevel"`
	Url       string `input:"currentPageURL"`
	PushURL   bool   `input:"pushURL"`
}

// updateInformation extracts the current roleLevel and pageURL via submitted form/url and if one of these are different from what
// is expected of the page it will replace the title/sidebar according to the new page/accountLevel. This is added as extra
func updateInformation(w http.ResponseWriter, r *http.Request, level database.RoleLevel, currentPage htmlComposition.HttpUrl) func(io.Writer) error {
	fields := &UserInformation{}
	var err error
	if r.Method == http.MethodGet {
		err = extractUrlValuesForFields(fields, r, 0)
	} else {
		err = extractFormValuesForFields(fields, r, 0)
	}
	if fields.PushURL {
		w.Header().Add("HX-Push-Url", "/"+string(currentPage))
	}
	if err != nil || (fields.RoleLevel == int(level) && fields.Url == string(currentPage)) {
		return func(w io.Writer) error {
			return nil
		}
	}
	return func(w io.Writer) error {
		var internalError error
		if fields.RoleLevel != int(level) {
			internalError = htmlComposition.GetSidebarReplacement(level).Render(w)
		}
		if internalError != nil {
			return internalError
		}
		if fields.Url != string(currentPage) {
			internalError = htmlComposition.GetTitleReplacement(currentPage).Render(w)
		}
		if internalError != nil {
			return internalError
		}
		internalError = htmlComposition.GetInfoDiv(level, currentPage).Render(w)
		return internalError
	}
}
