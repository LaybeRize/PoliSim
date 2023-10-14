package composition

import (
	. "PoliSim/html/builder"
	"fmt"
)

func getCheckBox(id string, checked bool, hidden bool, value string, name string, labelText string, hyperscript Node) Node {
	return DIV(IfElse(hidden, CLASS("form-check mt-2 hidden"), CLASS("form-check mt-2")), ID(id),
		INPUT(CLASS(`form-check-input appearance-none h-4 w-4 border border-gray-300 rounded-sm bg-white checked:bg-blue-600 checked:border-blue-600 focus:outline-none
            transition duration-200 mt-1 align-top bg-no-repeat bg-center bg-contain float-left mr-2 cursor-pointer`),
			TYPE("checkbox"), VALUE(value), NAME(name), ID(id+"Input"), hyperscript,
			If(checked && !hidden, CHECKED())),
		LABEL(CLASS("form-check-label inline-block"), FOR(id+"Input"),
			Text(labelText)),
	)
}

// getDropDown only works correct if the type t used also has the fmt.Stringer interface implemented.
func getDropDown[t comparable](name string, id string, labelText string, disable bool, arr []t, m map[t]string, selectedItem t) Node {
	return DIV(CLASS("mt-2"),
		LABEL(FOR(id), Text(labelText)),
		SELECT(If(disable, DISABLED()), ID(id), NAME(name), CLASS("bg-slate-700 appearance-none w-full py-2 px-3"),
			getOptions(arr, m, selectedItem),
		),
	)
}

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
func getTextArea(id string, name string, content string, labelText string) Node {
	return DIV(CLASS("mt-2"),
		LABEL(FOR(id), Text(labelText)), BR(),
		TEXTAREA(NAME(name), ID(id), CLASS("bg-slate-700 appearance-none w-full h-[200px] py-2 px-3"),
			Text(content),
		),
	)
}

// getSimpleTextInput calls getInput with the typStr "text" and list and extraClass parameter empty.
func getSimpleTextInput(id string, name string, value string, labelText string) Node {
	return getInput(id, name, value, labelText, "text", "", "")
}

// getInput returns an <input> element filled with the id, name, value, type (here typeStr), the used list for suggestions and
// addition css parameter with extraClass.
func getInput(id string, name string, value string, labelText string, typeStr string, list string, extraClass string, others ...Node) Node {
	return DIV(CLASS("mt-2 "+extraClass), ID(id+"Div"),
		LABEL(FOR(id), Text(labelText)),
		INPUT(TYPE(typeStr), NAME(name), ID(id), VALUE(value), Group(others...),
			CLASS("bg-slate-700 appearance-none w-full py-2 px-3"),
			If(list != "", LIST(list))),
	)
}

var buttonClassAttribute = "bg-slate-700 text-white p-2"

// getSubmitButton returns the standard form submit button
func getSubmitButton(buttonText string) Node {
	return BUTTON(TYPE("submit"), CLASS(buttonClassAttribute+" mt-2 mr-2"),
		Text(buttonText))
}

type formType int

const (
	GET formType = iota
	POST
	PATCH
	DELETE
)

// getSubmitButtonOverwriteURL returns a submit button that also overwrites the form hx-get/hx-post/hx-patch/hx-delete attribute
// with the desired new url and submission type.
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

// getFormStandardForm wraps all children in a <form> element with the needed htmx parameter based on the submit type and the url.
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
	return FORM(hx, ID(id), HXTARGET("#"+MainBodyID), HXSWAP("outerHTML"), HXINCLUDE("#"+InformationID), Group(children...))
}

func getEditableList(content []string, nameSpace string, listName string, addButtonText string, basicDivStyling string) Node {
	nodes := make([]Node, len(content))
	for i, str := range content {
		nodes[i] = getEditDiv(listName, nameSpace, str, "")
	}
	return DIV(CLASS(basicDivStyling),
		BUTTON(CLASS("bg-gray-900 text-white p-2 mt-2"), TYPE("button"),
			HYPERSCRIPT("on click tell next <div/> from me put you as HTML after you then toggle .hidden on next <div/> from you"), Text(addButtonText)),
		getEditDiv(listName, nameSpace, "", "hidden"),
		Group(nodes...),
	)
}

func getEditDiv(listName string, nameSpace string, value string, extraClass string) Node {
	return DIV(CLASS("flex flex-row "+extraClass),
		INPUT(LIST(listName), CLASS("bg-slate-700 appearance-none w-full py-2 px-3 mt-2"), NAME(nameSpace),
			VALUE(value)),
		BUTTON(CLASS("bg-slate-700 text-white p-4 ml-2 mt-2 hover:bg-rose-800"),
			HYPERSCRIPT("on click tell me.parentElement remove yourself"), Text(Translation["deleteEditableInput"]),
		),
	)
}

/*
package partHTML

import "API_MBundestag/database"

var editDiv = SafeCSS("flex flex-row")
var editInput = SafeCSS("bg-slate-700 appearance-none w-full py-2 px-3 mt-2")
var editButton = SafeCSS("bg-slate-700 text-white p-4 ml-2 mt-2 hover:bg-rose-800")
var deleteText = "Löschen"

templ AddEditableList(content []string, nameSpace string, listName string, divClass templ.ConstantCSSClass) {
    <div class={ divClass }>
        <button class="bg-gray-900 text-white p-2 mt-2" type="button" _="on click tell next <div/> from me put you as HTML after you then toggle .hidden on next <div/> from you">
            { children... }
        </button>
        <div class={ editDiv, SafeCSS("hidden") }>
            <input list={ listName } class={ editInput } name={ nameSpace } />
            <button class={ editButton } _="on click tell me.parentElement remove yourself">
                { deleteText }
            </button>
        </div>
        @editableButtonAndInput(content, nameSpace, listName)
    </div>
}

templ editableButtonAndInput(content []string, nameSpace string, listName string) {
    for _, str := range content {
        <div class={ editDiv }>
            <input list={ listName } class={ editInput } name={ nameSpace } value={ str }/>
            <button class={ editButton } _="on click tell me.parentElement remove yourself">
                { deleteText }
            </button>
        </div>
    }
}

templ AddSpecialSubmitButton(url string, formID string) {
    <button type="button" hx-post={ url } hx-include={ "#"+formID } class="bg-slate-700 text-white p-2 mt-2 mr-2">
        { children... }
    </button>
}


templ AddPreviewButton(include string) {
    <button type="button" class="bg-slate-700 text-white p-2 mt-2 mr-2"
        hx-post="/markdown-html"
        hx-push-url="false"
        hx-target="#previewDiv"
        hx-include={ include }>
        { children... }
    </button>
}

templ AddPreviewField() {
    <div class="w-[800px] mt-2">
        @Breaker("h-[3px] w-[800px]")
        <h1 class="text-2xl text-white mb-2">Vorschau</h1>
        @Breaker("w-[300px]")
        <div class="w-[800px]" id="previewDiv">

        </div>
        @Breaker("h-[3px] w-[800px]")
    </div>
}



templ AddDropDownSelection[K database.StatusString | database.RoleString | database.VoteType](name string, id string, disable bool, list []K, mapping map[K]string, selectedItem K) {
    <div class="mt-2">
        <label for={ id }>{ children... }</label>
        <select name={ name } id={ id } class="bg-slate-700 appearance-none w-full py-2 px-3"
            disabled?={ disable }>
            @addOptionsWithMap(list, mapping, selectedItem)
        </select>
    </div>
}

templ addOptionsWithMap[K database.StatusString | database.RoleString | database.VoteType](list []K, mapping map[K]string, selectedItem K) {
    for _, item := range list {
        <option value={ string(item) } selected?={ item == selectedItem }>{ mapping[item] }</option>
    }
}
*/
