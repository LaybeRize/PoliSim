package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"strings"
)

func GetViewTitelPage() Node {
	listing := make([]Node, len(extraction.TitleGroupMap))
	counter := 0
	for _, outer := range extraction.TitleMainGroupList {
		innerCounter := 0
		innerListing := make([]Node, len(extraction.TitleGroupMap[outer]))
		for _, inner := range extraction.TitleGroupMap[outer] {
			innerListing[innerCounter] = Group(BUTTON(
				CLASS("text-2xl mt-2 w-full text-left"), Text(inner),
				HXGET("/"+APIPreRoute+string(getTitleSubGroup)+outer+"/"+inner),
				HXTARGET("#out-"+outer+"-in-"+inner), ID("out-"+outer+"-in-"+inner+"-button"),
				HXSWAP("outerHTML"),
			),
				DIV(ID("out-"+outer+"-in-"+inner)))
			innerCounter++
		}
		listing[counter] = Group(BUTTON(CLASS("text-3xl mt-2 w-full text-left"), Text(outer),
			HYPERSCRIPT("on click toggle .hidden on #outer-"+outer)),
			DIV(ID("outer-"+outer), CLASS("border-l-4 border-white pl-6 mt-2 collapse-all hidden"),
				Group(innerListing...)))
		counter++
	}

	return getBasePageWrapper(
		getPageHeader(ViewTitles),
		DIV(CLASS("flex flex-row w-[600px]"),
			BUTTON(TYPE("button"), CLASS("bg-slate-700 text-white p-2 m-2"),
				HYPERSCRIPT("on click add .hidden to .collapse-all"), Text(Translation["collapseAll"])),
		),
		DIV(CLASS("mt-4 w-[600px]"),
			Group(listing...)),
	)
}

func GetViewSubGroupOfTitles(mainGroup string, subGroup string) Node {
	list, err := extraction.GetAllInSubGroup(mainGroup, subGroup)
	var newDiv Node = nil
	if err != nil {
		newDiv = DIV(ID("out-"+mainGroup+"-in-"+subGroup), CLASS("border-l-4 border-slate-400 pl-6 mt-2 collapse-all"),
			P(STYLE("font-size: 2em;"), CLASS("text-rose-600"), Text(Translation["errorWhileQueryingTitles"])))
	} else {
		nodeList := make([]Node, len(*list))
		for i, title := range *list {
			holderText := strings.Join(validation.GetDisplayNameArray(&title.Holder), ", ")
			nodeList[i] = DIV(CLASS("mt-2"),
				P(CLASS("text-xl"), Text(title.Name)),
				If(title.Flair.Valid, P(CLASS("text-base mt-2"), Text(Translation["viewTitleFlairText"], title.Flair.String))),
				P(CLASS(""), IfElse(holderText == "", Text(Translation["viewTitleNoHolderText"]),
					Text(Translation["viewTitleHolderText"], holderText))))

		}
		newDiv = DIV(ID("out-"+mainGroup+"-in-"+subGroup), CLASS("border-l-4 border-slate-400 pl-6 mt-2 collapse-all"),
			Group(nodeList...))
	}
	return Group(BUTTON(
		CLASS("text-2xl mt-2 w-full text-left"), Text(subGroup),
		ID("out-"+mainGroup+"-in-"+subGroup+"-button"),
		HXSWAPOOB("true"),
		HYPERSCRIPT("on click toggle .hidden on #out-"+mainGroup+"-in-"+subGroup),
	), newDiv)
}
