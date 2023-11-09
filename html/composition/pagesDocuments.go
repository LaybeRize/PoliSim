package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
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

func ViewDocumentPage(uuidStr string) Node {
	doc, err := extraction.GetDocument(database.LegislativeText, uuidStr)
	if err != nil {
		return GetErrorPage(Translation["documentDoesNotExists"])
	}
	otherNodes := make([]Node, len(doc.Info.Post)-1)
	for i, post := range doc.Info.Post[1:] {
		otherNodes[i] = DIV(CLASS("mt-2 w-[800px]"),
			renderTag(post))
	}
	return getBasePageWrapper(
		getPageHeader(ViewTextDocument),
		DIV(CLASS("w-[800px]"),
			H1(CLASS("text-3xl underline decoration-2 underline-offset-2"), Text(doc.Title)),
			If(doc.Subtitle.Valid, H1(CLASS("text-2xl"), Text(doc.Subtitle.String))),
			P(CLASS("my-2"), I(Text(fmt.Sprintf(doc.Written.Format(Translation["authorTextDocument"]), doc.Organisation, doc.Author))),
				If(doc.Flair != "", Group(I(Text("; ")), Text(doc.Flair)))),
			renderTag(doc.Info.Post[0]),
		),
		DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
			Raw(doc.HTMLContent),
		),
		DIV(ID(DocumentAdminPanel), HXGET("/"+APIPreRoute+string(AddTagDocumentLink)+url.PathEscape(doc.UUID)+"?org="+
			url.QueryEscape(doc.Organisation)), HXTRIGGER("load"), HXSWAP("outerHTML"), HXTARGET("#"+DocumentAdminPanel)),
		GetMessage(validation.Message{}),
		DIV(CLASS("w-[800px]"), ID(DocumentTagDiv), Group(otherNodes...)),
	)
}

func renderTag(posts database.Posts) Node {
	return P(CLASS("p-2"), STYLE("background-color: "+posts.Color+";"),
		Text(posts.Info), BR(),
		I(STYLE("font-size: 0.875rem;"), Text(posts.Submitted.Format(Translation["postsTimeFormat"]))))
}

func GetTagAdminPanel(uuid string) Node {
	doc, _ := extraction.GetDocument(database.LegislativeText, uuid)
	if len(doc.Info.Post) == 0 {
		doc.Info.Post = append(doc.Info.Post, database.Posts{})
	}
	otherNodes := make([]Node, len(doc.Info.Post)-1)
	for i, post := range doc.Info.Post[1:] {
		otherNodes[i] = DIV(CLASS("mt-2 w-[800px]"),
			renderTag(post))
	}
	return Group(DIV(
		getFormStandardForm("form", PATCH, "/"+APIPreRoute+string(AddTagDocumentLink)+url.PathEscape(uuid), CLASS("w-[800px]"),
			getSimpleTextInput("tag", "tag", "", Translation["tagTextDocument"]),
			getInput("color", "color", "", Translation["tagColorTextDocument"], "color",
				"", "", STYLE("min-height: 20px;")),
			getSubmitButton(Translation["addTagButton"])),
	),
		DIV(HXSWAPOOB("true"), ID(DocumentTagDiv), CLASS("w-[800px]"), Group(otherNodes...)))
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
