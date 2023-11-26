package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/logic"
	. "PoliSim/html/builder"
	"fmt"
	"net/url"
)

func GetZwitschers(extra *extraction.ExtraZwitscherInfo) Node {
	view, err := logic.GetZwitschers(extra)
	if err != nil {
		return GetErrorPage(Translation["errorLoadingLetters"])
	}
	nodes := make([]Node, len(*view.ZwitscherList))
	for i, item := range *view.ZwitscherList {
		link := string(ViewSingleZwitscherLink) + url.PathEscape(item.UUID)
		nodes[i] = getClickableLink("/"+HTMXPreRouter+link, "/"+link, Group(getStandardBoxClass,
			IfElse(item.Blocked, STYLE("--clr-border: rgb(40 51 69);"), STYLE("--clr-border: rgb(140 140 140);")),
			H1(CLASS("text-2xl"), Text(item.HTMLContent)),
			P(Text(Translation["authorShortFormLetter"], item.Author))))
	}
	before := fmt.Sprintf("%s?uuid=%s&amount=%d&before=true", string(ViewZwitscher),
		url.QueryEscape(view.BeforeUUID), extra.Amount)
	next := fmt.Sprintf("%s?uuid=%s&amount=%d", string(ViewZwitscher),
		url.QueryEscape(view.NextUUID), extra.Amount)
	extraStr := GetExtraStringForZwitscher(extra)
	return getBasePageWrapper(
		getPageHeader(ViewLetter),
		Group(nodes...),
		pagerFooter(view.BeforeUUID, view.NextUUID,
			before+extraStr, next+extraStr),
	)
}

func GetAuthorQueryString(author string) string {
	return "author=" + url.QueryEscape(author)
}

func GetExtraStringForZwitscher(extra *extraction.ExtraZwitscherInfo) string {
	result := ""
	if extra.HideBlock {
		result += "&hideblock=true"
	}
	if extra.ShowOnlyReplies {
		result += "&onlyreplies=true"
	}
	if extra.ShowOnlyZwitscher {
		result += "&onlyzwitscher=true"
	}
	if extra.Author != "" {
		result += "&" + GetAuthorQueryString(extra.Author)
	}
	return result
}
