package componentHelper

type ElementType string
type AttributeType string

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

	BASE    = elementWrapper(BaseTag)
	BR      = elementWrapper(BrTag)
	COL     = elementWrapper(ColTag)
	COMMAND = elementWrapper(CommandTag)
	EMBED   = elementWrapper(EmbedTag)
	HR      = elementWrapper(HrTag)
	IMG     = elementWrapper(ImgTag)
	INPUT   = elementWrapper(InputTag)
	KEYGEN  = elementWrapper(KeygenTag)
	LINK    = elementWrapper(LinkTag)
	META    = elementWrapper(MetaTag)
	PARAM   = elementWrapper(ParamTag)
	SOURCE  = elementWrapper(SourceTag)
	TRACK   = elementWrapper(TrackTag)
	WBR     = elementWrapper(WbrTag)

	HXDELETE    = attributeWrapper("hx-delete")
	HXPATCH     = attributeWrapper("hx-patch")
	HXPOST      = attributeWrapper("hx-post")
	HXGET       = attributeWrapper("hx-get")
	HXTRIGGER   = attributeWrapper("hx-trigger")
	HXINCLUDE   = attributeWrapper("hx-include")
	HXVALS      = attributeWrapper(HxValue)
	HXSWAPOOB   = attributeWrapper("hx-swap-oob")
	HXTARGET    = attributeWrapper("hx-target")
	HXSWAP      = attributeWrapper("hx-swap")
	SRC         = attributeWrapper("src")
	ALT         = attributeWrapper("alt")
	CHARSET     = attributeWrapper("charset")
	NAME        = attributeWrapper("name")
	CONTENT     = attributeWrapper("content")
	REL         = attributeWrapper("rel")
	HREF        = attributeWrapper("href")
	CLASS       = attributeWrapper("class")
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
)

const (
	AreaTag    ElementType = "area"
	BaseTag    ElementType = "base"
	BrTag      ElementType = "br"
	ColTag     ElementType = "col"
	CommandTag ElementType = "command"
	EmbedTag   ElementType = "embed"
	HrTag      ElementType = "hr"
	ImgTag     ElementType = "img"
	InputTag   ElementType = "input"
	KeygenTag  ElementType = "keygen"
	LinkTag    ElementType = "link"
	MetaTag    ElementType = "meta"
	ParamTag   ElementType = "param"
	SourceTag  ElementType = "source"
	TrackTag   ElementType = "track"
	WbrTag     ElementType = "wbr"

	HxValue AttributeType = "hx-vals"
)

var voidElements = map[ElementType]struct{}{
	AreaTag:    {},
	BaseTag:    {},
	BrTag:      {},
	ColTag:     {},
	CommandTag: {},
	EmbedTag:   {},
	HrTag:      {},
	ImgTag:     {},
	InputTag:   {},
	KeygenTag:  {},
	LinkTag:    {},
	MetaTag:    {},
	ParamTag:   {},
	SourceTag:  {},
	TrackTag:   {},
	WbrTag:     {},
}
