package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func CreateDocumentPage(acc *extraction.AccountAuth, document *validation.CreateDocument, val validation.Message) Node {
	node, err := UpdateUserOrganisations(acc, &extraction.AccountModification{ID: acc.ID,
		DisplayName: acc.DisplayName}, document.Organisation, "admin")
	if err != nil {
		val.Message = Translation["errorRetrievingOrganisationForAccount"] + "\n" + val.Message
	}
	return getBasePageWrapper(
		getPageHeader(CreateTextDocument),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateTextDocument), CLASS("w-[800px]"),
			node,
			getSimpleTextInput("title", "title", document.Title, Translation["titleTextDocument"]),
			getSimpleTextInput("subtitle", "subtitle", document.Subtitle, Translation["subtitleTextDocument"]),
			getTextArea("content", "content", document.Content, Translation["contentTextDocument"], true),

			getSimpleTextInput("tag", "tag", document.TagText, Translation["tagTextDocument"]),
			getInput("color", "color", document.TagColor, Translation["tagColorTextDocument"], "color",
				"", "", STYLE("min-height: 20px;")),
			getSubmitButton(Translation["createTextDocumentButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}

func ViewDocumentPage(acc *extraction.AccountAuth, uuidStr string) Node {
	return getBasePageWrapper(
		Text(uuidStr),
		Text(acc.DisplayName),
	)
}

func UpdateUserOrganisations(baseAccount *extraction.AccountAuth, account *extraction.AccountModification, organisationName string, isAdmin string) (Node, error) {
	searchForAdmin := isAdmin == "admin"
	orgList, err := extraction.GetOrganisationsForWithUserInIt(account.ID, searchForAdmin)
	nodes := make([]Node, len(*orgList))
	for i, item := range *orgList {
		nodes[i] = OPTION(VALUE(item.Name), If(item.Name == organisationName, SELECTED()),
			Text(item.Name))
	}
	return DIV(ID(UserSelectionID), DIV(CLASS("mt-2"),
		LABEL(FOR("authorAccount"), Text(Translation["accountDocument"])),
		SELECT(ID("authorAccount"), NAME("authorAccount"), CLASS("bg-slate-700 appearance-none w-full py-2 px-3"),
			HXPATCH("/"+APIPreRoute+string(updateUserSelectionLink)+isAdmin), HXTRIGGER("change"),
			HXTARGET("#"+UserSelectionID), HXSWAP("outerHTML"),
			getUserOptions(baseAccount, account.DisplayName),
		),
	), DIV(CLASS("mt-2"),
		LABEL(FOR("authorAccount"), Text(Translation["organisationDocument"])),
		SELECT(ID("organisation"), NAME("organisation"),
			CLASS("bg-slate-700 appearance-none w-full py-2 px-3"),
			Group(nodes...),
		),
	)), err
}
