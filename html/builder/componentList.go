package builder

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
	CODE     = elementWrapper("code")
	PRE      = elementWrapper("pre")

	BR    = elementWrapper(brTag)
	IMG   = elementWrapper(imgTag)
	INPUT = elementWrapper(inputTag)
	LINK  = elementWrapper(linkTag)
	META  = elementWrapper(metaTag)
	//AREA    = elementWrapper(areaTag)
	//BASE    = elementWrapper(baseTag)
	//COL     = elementWrapper(colTag)
	//COMMAND = elementWrapper(commandTag)
	//EMBED   = elementWrapper(embedTag)
	//HR      = elementWrapper(hrTag)
	//KEYGEN  = elementWrapper(keygenTag)
	//PARAM   = elementWrapper(paramTag)
	//SOURCE  = elementWrapper(sourceTag)
	//TRACK   = elementWrapper(trackTag)
	//WBR     = elementWrapper(wbrTag)

	HXDELETE    = attributeWrapper("hx-delete")
	HXPATCH     = attributeWrapper("hx-patch")
	HXPOST      = attributeWrapper("hx-post")
	HXGET       = attributeWrapper("hx-get")
	HXTRIGGER   = attributeWrapper("hx-trigger")
	HXINCLUDE   = attributeWrapper("hx-include")
	HXEXTEND    = attributeWrapper("hx-ext")
	HXSWAPOOB   = attributeWrapper("hx-swap-oob")
	HXTARGET    = attributeWrapper("hx-target")
	HXSWAP      = attributeWrapper("hx-swap")
	HXPUSHURL   = attributeWrapper("hx-push-url")
	SSECONNECT  = attributeWrapper("sse-connect")
	SSESWAP     = attributeWrapper("sse-swap")
	TEST        = attributeWrapper("data-testid")
	CONVERTTO   = attributeWrapper("data-convert")
	MIN         = attributeWrapper("min")
	MAX         = attributeWrapper("max")
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
	areaTag    = "area"
	baseTag    = "base"
	brTag      = "br"
	colTag     = "col"
	commandTag = "command"
	embedTag   = "embed"
	hrTag      = "hr"
	imgTag     = "img"
	inputTag   = "input"
	keygenTag  = "keygen"
	linkTag    = "link"
	metaTag    = "meta"
	paramTag   = "param"
	sourceTag  = "source"
	trackTag   = "track"
	wbrTag     = "wbr"
)

var voidElements = map[string]struct{}{
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
