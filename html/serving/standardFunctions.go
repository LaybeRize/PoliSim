package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"encoding/json"
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
		acc, _ := CheckUserPrivileges(w, r)
		addition := "?" + r.URL.RawQuery
		if addition == "?" {
			addition = ""
		}
		html := composition.GetBasePage(pageTitle, acc.Role, r.URL.Path[1:], addition)
		renderRequest(w, true, html.Render)
	}
}

// renderRequest renders the request to the response writer, if the write fails a
// http.StatusInternalServerError is put in the provided http.ResponseWriter.
// addDoc adds the <!DOCTYPE html> to the start of the response.
func renderRequest(w http.ResponseWriter, addDoc bool, f func(io.Writer) error) {
	if addDoc {
		err := builder.RenderHTMLDoc().Render(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err := f(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
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

// getText extracts the first string in http.Request PostForm field
// and trims it before returning.
func getText(r *http.Request, fieldName string) string {
	return strings.TrimSpace(r.PostFormValue(fieldName))
}

// getSliceAsValue reads in the string slice of the requested PostForm field
// trims every entry and removes any empty entries and doubled entries.
func getSliceAsValue(r *http.Request, fieldName string) reflect.Value {
	slice := r.PostForm[fieldName]
	helper.ClearStringArray(&slice)
	return reflect.ValueOf(slice)
}

// getBool checks if the requested field contains the text "true" and returns that value
func getBool(r *http.Request, fieldName string) bool {
	return getText(r, fieldName) == "true"
}

// getInt reads in the field then tries to transform it to a number. On success, it returns the
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

// CheckUserPrivileges gets the account from the cookie and returns the extraction.AccountAuth if one is found for the cookie, and if
// the account has one of the specified roles returns true. Otherwise, returns false
func CheckUserPrivileges(w http.ResponseWriter, r *http.Request, roleString ...database.RoleLevel) (*extraction.AccountAuth, bool) {
	acc := validation.ValidateToken(w, r)

	return acc, CheckIfHasRole(acc, roleString...)
}

// CheckIfHasRole checks if the referenced account has one of the provided roles
func CheckIfHasRole(acc *extraction.AccountAuth, roles ...database.RoleLevel) bool {
	for _, r := range roles {
		if r == acc.Role {
			return true
		}
	}
	return false
}

// onlySwapMessage retargets the request and only replaces the Message <div> on the page via htmx
func onlySwapMessage(w http.ResponseWriter, val validation.Message, node builder.Node) {
	w.Header().Set("HX-Retarget", "#"+composition.MessageID)
	html := composition.GetMessage(val)
	renderRequest(w, false, groupNodes(node, html))
}

type UserInformation struct {
	RoleLevel int    `input:"personalRoleLevel" json:"personalRoleLevel"`
	Url       string `input:"currentPageURL" json:"currentPageURL"`
}

// updateInformation extracts the current roleLevel and pageURL via submitted form/url and if one of these are different from what
// is expected of the page it will replace the title/sidebar according to the new page/accountLevel. It covers GET request, form requests and json requests
// It gets all it's information from the UserInformation struct.
func updateInformation(r *http.Request, level database.RoleLevel, currentPage composition.HttpUrl) builder.Node {
	fields := &UserInformation{}
	var err error
	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		err = extractUrlValuesForFields(fields, r, 0)
	} else if r.Header.Get("Content-Type") == "application/json" {
		err = extractAsJson(r, fields)
	} else {
		err = extractFormValuesForFields(fields, r, 0)
	}
	switch true {
	case err != nil || (fields.RoleLevel == int(level) && fields.Url == string(currentPage)):
		return nil
	case fields.RoleLevel != int(level) && fields.Url != string(currentPage):
		return builder.Group(composition.GetSidebarReplacement(level),
			composition.GetTitleReplacement(currentPage),
			composition.GetInfoDiv(level, currentPage))
	case fields.RoleLevel != int(level):
		return builder.Group(composition.GetSidebarReplacement(level),
			composition.GetInfoDiv(level, currentPage))
	case fields.Url != string(currentPage):
		return builder.Group(composition.GetTitleReplacement(currentPage),
			composition.GetInfoDiv(level, currentPage))
	}
	return nil
}

// extractAsJson extracts the needed information from the json request
// and transforms them to the parsed UserInformation struct.
func extractAsJson(r *http.Request, fields *UserInformation) error {
	temp := &struct {
		RoleLevel string `json:"personalRoleLevel"`
		Url       string `json:"currentPageURL"`
	}{}
	err := json.NewDecoder(r.Body).Decode(temp)
	if err != nil {
		return err
	}

	var i int
	i, err = strconv.Atoi(temp.RoleLevel)
	if err != nil {
		i = 0
	}
	fields.RoleLevel = i
	fields.Url = temp.Url
	return nil
}

// groupNodes returns the render function ofa group of nodes
func groupNodes(children ...builder.Node) func(io.Writer) error {
	return builder.Group(children...).Render
}

// genericRenderer returns a generics render function for a typical urls by parsing the htmlComposition.HttpUrl
func genericRenderer(currentPage composition.HttpUrl) func(w http.ResponseWriter,
	r *http.Request, level database.RoleLevel, node builder.Node) {
	return func(w http.ResponseWriter, r *http.Request, level database.RoleLevel, node builder.Node) {
		renderRequest(w, false, groupNodes(updateInformation(r, level, currentPage),
			node))
	}
}

// genericMessageSwapper returns a generics render function that only swaps the message
// element for a typical urls by parsing the htmlComposition.HttpUrl
func genericMessageSwapper(currentPage composition.HttpUrl) func(w http.ResponseWriter,
	r *http.Request, val validation.Message, level database.RoleLevel) {
	return func(w http.ResponseWriter, r *http.Request, val validation.Message, level database.RoleLevel) {
		onlySwapMessage(w, val, updateInformation(r, level, currentPage))
	}
}
