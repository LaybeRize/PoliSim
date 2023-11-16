package builder

type elementType string
type attributeType string
type HttpUrl string

var (
	DIV      = elementWrapper("div")
	P        = elementWrapper("p")
	I        = elementWrapper("i")
	A        = elementWrapper("a")
	SPAN     = elementWrapper("span")
	H1       = elementWrapper("h1")
	HTML     = elementWrapper("html")
	HEAD     = elementWrapper("head")
	BODY     = elementWrapper("body")
	TITLE    = elementWrapper("title")
	SCRIPT   = elementWrapper("script")
	OPTION   = elementWrapper("option")
	DATALIST = elementWrapper("datalist")
	LABEL    = elementWrapper("label")
	SELECT   = elementWrapper("select")
	TEXTAREA = elementWrapper("textarea")
	BUTTON   = elementWrapper("button")
	FORM     = elementWrapper("form")
	TABLE    = elementWrapper("table")
	TH       = elementWrapper("th")
	TR       = elementWrapper("tr")
	TD       = elementWrapper("td")

	AREA    = elementWrapper(areaTag)
	BASE    = elementWrapper(baseTag)
	BR      = elementWrapper(brTag)
	COL     = elementWrapper(colTag)
	COMMAND = elementWrapper(commandTag)
	EMBED   = elementWrapper(embedTag)
	HR      = elementWrapper(hrTag)
	IMG     = elementWrapper(imgTag)
	INPUT   = elementWrapper(inputTag)
	KEYGEN  = elementWrapper(keygenTag)
	LINK    = elementWrapper(linkTag)
	META    = elementWrapper(metaTag)
	PARAM   = elementWrapper(paramTag)
	SOURCE  = elementWrapper(sourceTag)
	TRACK   = elementWrapper(trackTag)
	WBR     = elementWrapper(wbrTag)

	HXDELETE    = attributeWrapper("hx-delete")
	HXPATCH     = attributeWrapper("hx-patch")
	HXPOST      = attributeWrapper("hx-post")
	HXGET       = attributeWrapper("hx-get")
	HXTRIGGER   = attributeWrapper("hx-trigger")
	HXEXTEND    = attributeWrapper("hx-ext")
	TEST        = attributeWrapper("data-testid")
	HXSWAPOOB   = attributeWrapper("hx-swap-oob")
	HXTARGET    = attributeWrapper("hx-target")
	HXSWAP      = attributeWrapper("hx-swap")
	HXPUSHURL   = attributeWrapper("hx-push-url")
	SRC         = attributeWrapper("src")
	ALT         = attributeWrapper("alt")
	CHARSET     = attributeWrapper("charset")
	NAME        = attributeWrapper("name")
	CONTENT     = attributeWrapper("content")
	REL         = attributeWrapper("rel")
	HREF        = attributeWrapper("href")
	CLASS       = attributeWrapper("class")
	ROWSPAN     = attributeWrapper("rowspan")
	STYLE       = attributeWrapper("style")
	ID          = attributeWrapper("id")
	HYPERSCRIPT = attributeWrapper("_")
	HIDDEN      = attributeWrapper("hidden")
	LANG        = attributeWrapper("lang")
	VALUE       = attributeWrapper("value")
	FOR         = attributeWrapper("for")
	TYPE        = attributeWrapper("type")
	LIST        = attributeWrapper("list")
	DISABLED    = attributeWrapper("disabled")
	SELECTED    = attributeWrapper("selected")
	CHECKED     = attributeWrapper("checked")
)

const (
	areaTag    elementType = "area"
	baseTag    elementType = "base"
	brTag      elementType = "br"
	colTag     elementType = "col"
	commandTag elementType = "command"
	embedTag   elementType = "embed"
	hrTag      elementType = "hr"
	imgTag     elementType = "img"
	inputTag   elementType = "input"
	keygenTag  elementType = "keygen"
	linkTag    elementType = "link"
	metaTag    elementType = "meta"
	paramTag   elementType = "param"
	sourceTag  elementType = "source"
	trackTag   elementType = "track"
	wbrTag     elementType = "wbr"

	hxValue attributeType = "hx-vals"
)

var voidElements = map[elementType]struct{}{
	areaTag:    {},
	baseTag:    {},
	brTag:      {},
	colTag:     {},
	commandTag: {},
	embedTag:   {},
	hrTag:      {},
	imgTag:     {},
	inputTag:   {},
	keygenTag:  {},
	linkTag:    {},
	metaTag:    {},
	paramTag:   {},
	sourceTag:  {},
	trackTag:   {},
	wbrTag:     {},
}
