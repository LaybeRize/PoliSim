package composition

import (
	"PoliSim/data/extraction"
	. "PoliSim/html/builder"
	"fmt"
)

// getCheckBoxWithHideScript is the same as getCheckBox but parses the name as the id and adds a hyperscript tag,
// that toggels the hidden flag on the element with the provided id in the idToHide parameter.
func getCheckBoxWithHideScript(checked bool, value string, name string, labelText string, idToHide string) Node {
	return getCheckBox(name, checked, value, name, labelText, HYPERSCRIPT("on click toggle .hidden on #"+idToHide))
}

// getStandardCheckBox is the same as getStandardCheckBoxID but uses the name as the id value
func getStandardCheckBox(checked bool, value string, name string, labelText string) Node {
	return getStandardCheckBoxID(name, checked, value, name, labelText)
}

const (
	divCheckBoxClass   = "form-check mt-2"
	inputCheckBoxClass = `form-check-input appearance-none h-4 w-4 border border-gray-300 rounded-sm bg-white 
 checked:bg-blue-600 checked:border-blue-600 focus:outline-none transition duration-200 mt-1 align-top bg-no-repeat 
 bg-center bg-contain float-left mr-2 cursor-pointer`
	lableCheckBoxClass = "form-check-label inline-block"
)

// getStandardCheckBoxID os the sa,e as getCheckBox but without a possible extra Node to add
func getStandardCheckBoxID(id string, checked bool, value string, name string, labelText string) Node {
	inputID := id + "Input"
	return DIV(CLASS(divCheckBoxClass), ID(id),
		INPUT(CLASS(inputCheckBoxClass), TYPE("checkbox"), VALUE(value), NAME(name),
			ID(inputID), TEST(inputID), CONVERTTO("bool"), If(checked, CHECKED())),
		LABEL(CLASS(lableCheckBoxClass), FOR(inputID), Text(labelText)),
	)
}

// getCheckBox returns a checkbox that automatically transforms its value into a number when using the json-enc
// extension.
func getCheckBox(id string, checked bool, value string, name string, labelText string, hyperscript Node) Node {
	inputID := id + "Input"
	return DIV(CLASS(divCheckBoxClass), ID(id),
		INPUT(CLASS(inputCheckBoxClass), TYPE("checkbox"), VALUE(value), NAME(name),
			ID(inputID), TEST(inputID), hyperscript, CONVERTTO("bool"), If(checked, CHECKED())),
		LABEL(CLASS(lableCheckBoxClass), FOR(inputID),
			Text(labelText)),
	)
}

// getRadioButton returns a radio button that automatically transforms its value into a number when using the
// json-enc extension.
func getRadioButton(id string, checked bool, value string, name string, labelText string) Node {
	inputID := id + "Input"
	return DIV(CLASS(divCheckBoxClass), ID(id),
		INPUT(CLASS(inputCheckBoxClass),
			TYPE("radio"), VALUE(value), NAME(name), ID(inputID), TEST(inputID), CONVERTTO("number"),
			If(checked, CHECKED())),
		LABEL(CLASS(lableCheckBoxClass), FOR(inputID),
			Text(labelText)),
	)
}

const inputFieldBackgroundColor = "bg-slate-700 "

// getDropDown only works correct if the type t used also has the fmt.Stringer interface implemented.
func getDropDown[t comparable](name string, id string, labelText string, disable bool, arr []t, m map[t]string, selectedItem t) Node {
	return DIV(CLASS("mt-2"),
		LABEL(FOR(id), Text(labelText)),
		SELECT(If(disable, DISABLED()), ID(id), TEST(id), NAME(name),
			CLASS(inputFieldBackgroundColor+"appearance-none w-full py-2 px-3"),
			getOptions(arr, m, selectedItem),
		),
	)
}

const userDropDownID = "authorAccount"

// getUserDropdown gets a dropdown with all valid accounts for the parsed user. it has "authorAccount" as the name.
// Parse a specific account by display name to selectedAccount to activate the selected Tag on the option
// with that name.
func getUserDropdown(user *extraction.AccountAuth, selectedAccount string, labelText string) Node {
	return DIV(CLASS("mt-2"),
		LABEL(FOR(userDropDownID), Text(labelText)),
		SELECT(If(user.ID == 0, DISABLED()), ID(userDropDownID), TEST(userDropDownID),
			NAME(userDropDownID), CLASS(inputFieldBackgroundColor+"appearance-none w-full py-2 px-3"),
			getUserOptions(user, selectedAccount),
		),
	)
}

func getUserOptions(user *extraction.AccountAuth, selectedAccount string) Node {
	children, _ := extraction.GetAllChildrenDisplayNames(user.ID)
	nodes := make([]Node, len(*children)+1)
	nodes[0] = OPTION(VALUE(user.DisplayName), If(user.DisplayName == selectedAccount, SELECTED()),
		Text(user.DisplayName))
	for i, item := range *children {
		nodes[i+1] = OPTION(VALUE(item.DisplayName), If(item.DisplayName == selectedAccount, SELECTED()),
			Text(item.DisplayName))
	}
	return Group(nodes...)
}

// getOptions is the reason why getDropDown only works correctly if the type t used also has the fmt.Stringer
// interface implemented. It takes the value formats it to a string and puts it as the input value.
func getOptions[t comparable](arr []t, m map[t]string, selectedItem t) Node {
	nodes := make([]Node, len(arr))
	for index, item := range arr {
		strItem := any(item).(fmt.Stringer).String()
		nodes[index] = OPTION(VALUE(strItem), If(item == selectedItem, SELECTED()),
			Text(m[item]))
	}
	return Group(nodes...)
}

// getDataList creates a <datalist> element containing all items in listItems as <option> tags
// and the listName as the id.
func getDataList(listName string, listItems []string) Node {
	options := make([]Node, len(listItems)+1)
	options[0] = ID(listName)
	for i, str := range listItems {
		options[i+1] = OPTION(VALUE(str))
	}
	return DATALIST(options...)
}

// getDataListFromMap takes in any map that has a string as it's key and writes every
// key as an option into the datalist that is returned.
func getDataListFromMap[t any](listName string, listMap map[string]t) Node {
	options := make([]Node, len(listMap)+1)
	options[0] = ID(listName)
	i := 1
	for str := range listMap {
		options[i] = OPTION(VALUE(str))
		i++
	}
	return DATALIST(options...)
}

// getTextArea returns a styled text area for a form. content is the text filled into the area.
func getTextArea(id string, name string, content string, labelText string, patchURL HttpUrl) Node {
	return DIV(CLASS("mt-2"),
		LABEL(FOR(id), Text(labelText)), BR(),
		TEXTAREA(NAME(name), ID(id), TEST(id),
			CLASS(inputFieldBackgroundColor+"appearance-none w-full h-[200px] py-2 px-3"),
			If(patchURL != "", Group(
				HXPATCH("/"+APIPreRoute+string(patchURL)),
				HXTARGET("#"+DisplayID),
				HXTRIGGER("keyup changed delay:1s"),
				HXSWAP("outerHTML"))),
			Text(content),
		),
	)
}

// getSimpleTextInput calls getInput with the typStr "text" and list and extraClass parameter empty.
func getSimpleTextInput(id string, name string, value string, labelText string) Node {
	return getInput(id, name, value, labelText, "text", "", "")
}

// getInput returns an <input> element filled with the id, name, value, type (here typeStr), the used list for
// suggestions and addition css parameter with extraClass.
func getInput(id string, name string, value string, labelText string, typeStr string, list string, extraClass string, others ...Node) Node {
	return DIV(CLASS("mt-2 "+extraClass), ID(id+"Div"),
		LABEL(FOR(id), Text(labelText)),
		INPUT(TYPE(typeStr), NAME(name), ID(id), TEST(id), VALUE(value), Group(others...),
			CLASS(inputFieldBackgroundColor+"appearance-none w-full py-2 px-3"),
			If(list != "", LIST(list))),
	)
}

const buttonClassAttribute = inputFieldBackgroundColor + "text-white p-2 disable-selection"

// getSubmitButton returns the standard form submit button
func getSubmitButton(id string, buttonText string) Node {
	return BUTTON(TYPE("submit"), CLASS(buttonClassAttribute+" mt-2 mr-2"),
		ID(id), TEST(id), Text(buttonText))
}

type formType int

const (
	GET formType = iota
	POST
	PATCH
	DELETE
)

// getSubmitButtonOverwriteURL returns a submit button that also overwrites the form
// hx-get/hx-post/hx-patch/hx-delete attribute with the desired new url and submission type.
func getSubmitButtonOverwriteURL(buttonText string, submit formType, url string) Node {
	hx := Node(nil)
	switch submit {
	case GET:
		hx = HXGET(url)
	case POST:
		hx = HXPOST(url)
	case PATCH:
		hx = HXPATCH(url)
	case DELETE:
		hx = HXDELETE(url)
	}
	return BUTTON(TYPE("submit"), CLASS(buttonClassAttribute+" mt-2 mr-2"), hx,
		Text(buttonText))
}

// getFormStandardForm wraps all children in a <form> element with the needed htmx parameter based on the
// submit type and the url.
func getFormStandardForm(id string, submit formType, url string, children ...Node) Node {
	hx := Node(nil)
	switch submit {
	case GET:
		hx = HXGET(url)
	case POST:
		hx = HXPOST(url)
	case PATCH:
		hx = HXPATCH(url)
	case DELETE:
		hx = HXDELETE(url)
	}
	return FORM(hx, ID(id), HXTARGET("#"+MainBodyID), HXSWAP("outerHTML"), Group(children...))
}

// getEditableList returns a button that adds on click new inputs with the same name. Can be used to let the user
// make a list of times. If items should be already present in the array, provide the content wit the array of
// items as strings, that should be displayed.
func getEditableList(content []string, nameSpace string, listName string, addButtonText string, basicDivStyling string) Node {
	nodes := make([]Node, len(content))
	for i, str := range content {
		nodes[i] = getEditDiv(listName, nameSpace, str, "")
	}
	return DIV(CLASS(basicDivStyling),
		BUTTON(CLASS("bg-gray-900 text-white p-2 mt-2 disable-selection"), TYPE("button"),
			HYPERSCRIPT("on click tell next <div/> from me put you as HTML after you"+
				" then toggle .hidden on next <div/> from you"), Text(addButtonText)),
		getEditDiv(listName, nameSpace, "", "hidden"),
		Group(nodes...),
	)
}

func getEditDiv(listName string, nameSpace string, value string, extraClass string) Node {
	return DIV(CLASS("flex flex-row "+extraClass),
		INPUT(LIST(listName), CLASS(inputFieldBackgroundColor+"appearance-none w-full py-2 px-3 mt-2"),
			NAME(nameSpace), VALUE(value)),
		BUTTON(CLASS(inputFieldBackgroundColor+"text-white p-4 ml-2 mt-2 hover:bg-rose-800 disable-selection"),
			TYPE("button"), HYPERSCRIPT("on click tell me.parentElement remove yourself"),
			Text(Translation["deleteEditableInput"]),
		),
	)
}
