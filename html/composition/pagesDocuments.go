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
			getTextArea("content", "content", document.Content, Translation["contentTextDocument"],
				MarkdownFormPage),

			getSimpleTextInput("tag", "tag", document.TagText, Translation["tagTextDocument"]),
			getInput("color", "color", document.TagColor, Translation["tagColorTextDocument"], "color",
				"", "", STYLE("min-height: 20px;")),
			getSubmitButton(Translation["createTextDocumentButton"])),
		GetMessage(val),
		getPreviewElement(),
	)
}

const (
	documentTagDiv = "document-tag-div"
	currentTagDiv  = "current-document-tag-div"
)

func ViewDocumentPage(uuidStr string) Node {
	doc, err := extraction.GetDocumentIfNotPrivate(database.LegislativeText, uuidStr)
	if err != nil {
		return GetErrorPage(Translation["documentDoesNotExists"])
	}
	nodes := make([]Node, 0, len(doc.Info.Post))
	for _, post := range doc.Info.Post {
		if post.Hidden {
			continue
		}
		nodes = append(nodes, DIV(CLASS("mt-2 w-[800px]"),
			renderTag(post, "")))
	}
	return getBasePageWrapper(
		getPageHeader(ViewTextDocument),
		DIV(CLASS("w-[800px]"),
			H1(CLASS("text-3xl underline decoration-2 underline-offset-2"), Text(doc.Title)),
			If(doc.Subtitle.Valid, H1(CLASS("text-2xl"), Text(doc.Subtitle.String))),
			P(CLASS("my-2"), I(Text(fmt.Sprintf(doc.Written.Format(Translation["authorTextDocument"]), doc.Organisation, doc.Author))),
				If(doc.Flair != "", Group(I(Text("; ")), Text(doc.Flair)))),
			DIV(ID(currentTagDiv),
				If(len(nodes) != 0, nodes[0]),
			),
		),
		DIV(CLASS("w-[800px] box box-e p-2 mt-2"), STYLE("--clr-border: rgb(40 51 69);"),
			Raw(doc.HTMLContent),
		),
		DIV(ID(DocumentAdminPanel), HXGET("/"+APIPreRoute+string(AddTagDocumentLink)+url.PathEscape(doc.UUID)+"?org="+
			url.QueryEscape(doc.Organisation)), HXTRIGGER("load"), HXSWAP("outerHTML"), HXTARGET("#"+DocumentAdminPanel)),
		GetMessage(validation.Message{}),
		DIV(CLASS("w-[800px]"), ID(documentTagDiv), If(len(nodes) != 0, Group(nodes[1:]...))),
	)
}

func GetTagAdminPanel(uuid string, isAdmin bool) Node {
	doc, _ := extraction.GetDocumentIfNotPrivate(database.LegislativeText, uuid)
	if len(doc.Info.Post) == 0 {
		doc.Info.Post = append(doc.Info.Post, database.Posts{})
	}
	nodes := make([]Node, 0, len(doc.Info.Post))
	hiddenNodes := make([]Node, 0, len(doc.Info.Post))
	for _, post := range doc.Info.Post {
		if post.Hidden {
			hiddenNodes = append(hiddenNodes, DIV(CLASS("mt-2 w-[800px] grid grid-flow-col grid-cols-3 justify-stretch"),
				renderTag(post, "col-span-2"), getCustomRequestClickable(HXPATCH, "/"+APIPreRoute+string(ChangeTagDocumentLink)+
					url.PathEscape(uuid)+"/"+url.PathEscape(post.UUID), "", P(CLASS("bg-slate-700 text-white p-2 mt-2 ml-2"),
					STYLE("text-align: center;"), Text(Translation["showTagButtonText"])),
				)))
			continue
		}
		nodes = append(nodes, DIV(CLASS("mt-2 w-[800px] grid grid-flow-col grid-cols-3 justify-stretch"),
			renderTag(post, "col-span-2"), getCustomRequestClickable(HXPATCH, "/"+APIPreRoute+string(ChangeTagDocumentLink)+
				url.PathEscape(uuid)+"/"+url.PathEscape(post.UUID), "", P(CLASS("bg-slate-700 text-white p-2 mt-2 ml-2"),
				STYLE("text-align: center;"), Text(Translation["hideTagButtonText"])),
			)))
	}
	return Group(DIV(
		getFormStandardForm("form", PATCH, "/"+APIPreRoute+string(AddTagDocumentLink)+url.PathEscape(uuid), CLASS("w-[800px]"),
			getSimpleTextInput("tag", "tag", "", Translation["tagTextDocument"]),
			getInput("color", "color", "", Translation["tagColorTextDocument"], "color",
				"", "", STYLE("min-height: 20px;")),
			getSubmitButton(Translation["addTagButton"])),
	),
		If(isAdmin, DIV(HXSWAPOOB("true"), ID(currentTagDiv), CLASS("w-[800px]"),
			If(len(nodes) != 0, nodes[0]))),
		If(isAdmin, DIV(HXSWAPOOB("true"), ID(documentTagDiv), CLASS("w-[800px]"),
			If(len(nodes) != 0, Group(nodes[1:]...)),
			If(len(hiddenNodes) != 0, P(CLASS("my-2 text-xl"), Text(Translation["hiddenTags"]))),
			Group(hiddenNodes...))))
}

func renderTag(posts database.Posts, extraClass string) Node {
	return P(CLASS("p-2 "+extraClass), STYLE("background-color: "+posts.Color+";"),
		Text(posts.Info), BR(),
		I(STYLE("font-size: 0.875rem;"), Text(posts.Submitted.Format(Translation["postsTimeFormat"]))))
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
