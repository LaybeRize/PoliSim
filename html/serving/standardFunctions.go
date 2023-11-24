package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	"PoliSim/data/validation"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// GetFullPage returns a http.HandlerFunc that writes a full page to the response
// with a div that automatically requests the URL via htmx.
func GetFullPage(pageTitle string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acc, _ := CheckUserPrivileges(r)
		addition := "?" + r.URL.RawQuery
		if addition == "?" {
			addition = ""
		}
		validation.AddCookie(w, r, acc.Session)
		html := composition.GetBasePage(pageTitle, acc, r.URL.Path[1:], addition)
		renderRequest(w, builder.RenderHTMLDoc(), html)
	}
}

// renderRequest renders the request to the response writer, if the write fails a
// http.StatusInternalServerError is put in the provided http.ResponseWriter.
// addDoc adds the <!DOCTYPE html> to the start of the response.
func renderRequest(w http.ResponseWriter, nodes ...builder.Node) {
	err := builder.Group(nodes...).Render(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func extractURLFieldValues(object any, r *http.Request, min int64, standard int64, max int64) {
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
			v.Elem().Field(i).SetString(r.URL.Query().Get(input))
		case reflect.Bool:
			value := r.URL.Query().Get(input) == "true"
			v.Elem().Field(i).SetBool(value)
		case reflect.Int, reflect.Int64:
			num, err := strconv.ParseInt(r.URL.Query().Get(input), 10, 64)
			if err != nil {
				num = standard
			} else if num < min {
				num = min
			} else if num > max {
				num = max
			}
			v.Elem().Field(i).SetInt(num)
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
	recursiveCall(r, v, onError)
	return nil
}

func recursiveCall(r *http.Request, v reflect.Value, onError int64) {
	iterate := reflect.Indirect(v).NumField()
	for i := 0; i < iterate; i++ {
		input, ok := v.Type().Elem().Field(i).Tag.Lookup("input")
		if !ok {
			continue
		} else if input == "_struct_" {
			recursiveCall(r, v.Elem().Field(i).Addr(), onError)
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

// CheckUserPrivileges gets the account from the cookie and returns the extraction.AccountAuth if one is found for the cookie, and if
// the account has one of the specified roles returns true. Otherwise, returns false
func CheckUserPrivileges(r *http.Request, roleString ...database.RoleLevel) (*extraction.AccountAuth, bool) {
	acc := validation.ValidateToken(r)

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

// updateInformation extracts the current roleLevel and pageURL via submitted form/url and if one of these are different from what
// is expected of the page it will replace the title/sidebar according to the new page/accountLevel. It covers GET request, form requests and json requests
// It gets all it's information from the Cookie.
func updateInformation(w http.ResponseWriter, r *http.Request, acc *extraction.AccountAuth, currentPage builder.HttpUrl) builder.Node {
	go logic.UpdateLetterNotification(acc.ID)
	role, ok := acc.Session.Values["role"].(int)
	if !ok {
		role = -100
	}
	acc.Session.Values["role"] = int(acc.Role)
	validation.AddCookie(w, r, acc.Session)

	arr := []builder.Node{nil, composition.GetTitleReplacement(currentPage)}
	if role != int(acc.Role) {
		arr[0] = composition.GetSidebarReplacement(acc)
	} else if acc.Role != database.NotLoggedIn {
		arr[0] = composition.GetLetterSidebarButton(acc, true)
	}

	return builder.Group(arr...)
}

// genericRenderer returns a generics render function for a typical urls by parsing the htmlComposition.HttpUrl
func genericRenderer(currentPage builder.HttpUrl) func(w http.ResponseWriter,
	r *http.Request, acc *extraction.AccountAuth, node builder.Node) {
	return func(w http.ResponseWriter, r *http.Request, acc *extraction.AccountAuth, node builder.Node) {
		renderRequest(w, updateInformation(w, r, acc, currentPage), node)
	}
}

func pushURL(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Push-Url", url)
}

func retargetToMessage(w http.ResponseWriter) {
	w.Header().Set("HX-Retarget", "#"+composition.MessageID)
}

// genericMessageSwapper returns a generics render function that only swaps the message
// element for a typical urls by parsing the htmlComposition.HttpUrl. It retargets the request and only
// replaces the Message <div> on the page via htmx.
func genericMessageSwapper(currentPage builder.HttpUrl) func(w http.ResponseWriter,
	r *http.Request, val validation.Message, acc *extraction.AccountAuth) {
	return func(w http.ResponseWriter, r *http.Request, val validation.Message, acc *extraction.AccountAuth) {
		retargetToMessage(w)
		html := composition.GetMessage(val)
		renderRequest(w, updateInformation(w, r, acc, currentPage), html)
	}
}
