package htmlComposition

import (
	. "PoliSim/componentHelper"
)

// getDataList creates a <datalist> element containing all items in listItems as <option> tags
// and the listName as the id.
func getDataList(listName string, listItems []string) Node {
	options := make([]Node, len(listItems)+1)
	options[0] = Attr(ID, listName)
	for i, str := range listItems {
		options[i+1] = El(OPTION, Attr(VALUE, str))
	}
	return El(DATALIST, options...)
}

// getTextArea returns a styled text area for a form. content is the text filled into the area.
func getTextArea(id string, name string, content string, labelText string) Node {
	return El(DIV, Attr(CLASS, "mt-2"),
		El(LABEL, Attr(FOR, id), Text(labelText)), El(BR),
		El(TEXTAREA, Attr(NAME, name), Attr(ID, id), Attr(CLASS, "bg-slate-700 appearance-none w-full h-[200px] py-2 px-3"),
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
func getInput(id string, name string, value string, labelText string, typeStr string, list string, extraClass string) Node {
	if extraClass != "" {
		extraClass = " " + extraClass
	}
	return El(DIV, Attr(CLASS, "mt-2"),
		El(LABEL, Attr(FOR, id), Text(labelText)),
		El(INPUT, Attr(TYPE, typeStr), Attr(NAME, name), Attr(ID, id), Attr(VALUE, value),
			Attr(CLASS, "bg-slate-700 appearance-none w-full py-2 px-3"+extraClass),
			If(list != "", Attr(LIST, list))),
	)
}

var buttonClassAttribute = "bg-slate-700 text-white p-2"

// getSubmitButton returns the standard form submit button
func getSubmitButton(buttonText string) Node {
	return El(BUTTON, Attr(TYPE, "submit"), Attr(CLASS, buttonClassAttribute+" mt-2 mr-2"),
		Text(buttonText))
}

type FormType int

const (
	GET FormType = iota
	POST
	PATCH
	DELETE
)

// getSubmitButtonOverwriteURL returns a submit button that also overwrites the form hx-get/hx-post/hx-patch/hx-delete attribute
// with the desired new url and submission type.
func getSubmitButtonOverwriteURL(buttonText string, submit FormType, url string) Node {
	hx := Node(nil)
	switch submit {
	case GET:
		hx = Attr(HXGET, url)
	case POST:
		hx = Attr(HXPOST, url)
	case PATCH:
		hx = Attr(HXPATCH, url)
	case DELETE:
		hx = Attr(HXDELETE, url)
	}
	return El(BUTTON, Attr(TYPE, "submit"), Attr(CLASS, buttonClassAttribute+" mt-2 mr-2"), hx,
		Text(buttonText))
}

// getFormStandardForm wraps all children in a <form> element with the needed htmx parameter based on the submit type and the url.
func getFormStandardForm(id string, submit FormType, url string, children ...Node) Node {
	hx := Node(nil)
	switch submit {
	case GET:
		hx = Attr(HXGET, url)
	case POST:
		hx = Attr(HXPOST, url)
	case PATCH:
		hx = Attr(HXPATCH, url)
	case DELETE:
		hx = Attr(HXDELETE, url)
	}
	return El(FORM, append(children, hx, Attr(ID, id), Attr(HXTARGET, "#"+MainBodyID), Attr(HXSWAP, "outerHTML"), Attr(HXINCLUDE, "#"+InformationID))...)
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

templ AddCheckbox(id string, checked bool, hidden bool, value string, name string, hyperscript string) {
    <div class={SafeCSS("form-check mt-2"), EnableCSS("hidden", hidden)} id={ id }>
        <input class="form-check-input appearance-none h-4 w-4 border border-gray-300 rounded-sm bg-white checked:bg-blue-600 checked:border-blue-600 focus:outline-none
            transition duration-200 mt-1 align-top bg-no-repeat bg-center bg-contain float-left mr-2 cursor-pointer"
            type="checkbox" value={ value } name={ name } id={ id+"Input" }
            _={ hyperscript }
            checked?={ checked && !hidden } />
        <label class="form-check-label inline-block" for={ id+"Input" }>
            { children... }
        </label>
    </div>
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


templ DropDownSelectedAccount(name string, list database.AccountList, selectedAccount string, disable bool) {
    <div class="mt-2">
        <label for="selectedAccount">{ children... }</label>
        <select name={ name } id="selectedAccount" class="bg-slate-700 appearance-none w-full py-2 px-3"
            disabled?={ disable }>
            @addOptionsForAccounts(list, selectedAccount)
        </select>
    </div>
}

templ addOptionsForAccounts(list database.AccountList, selectedAccount string) {
    for _, acc := range list {
        <option value={ acc.DisplayName } selected?={ selectedAccount == acc.DisplayName }>{ acc.DisplayName }</option>
    }
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
