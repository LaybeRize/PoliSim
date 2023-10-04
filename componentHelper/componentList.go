package componentHelper

type ElementType string
type AttributeType string

const (
	DIV      ElementType = "div"
	P        ElementType = "p"
	A        ElementType = "a"
	SPAN     ElementType = "SPAN"
	H1       ElementType = "h1"
	HTML     ElementType = "html"
	HEAD     ElementType = "head"
	BODY     ElementType = "body"
	TITLE    ElementType = "title"
	SCRIPT   ElementType = "script"
	OPTION   ElementType = "option"
	DATALIST ElementType = "datalist"
	LABEL    ElementType = "label"
	TEXTAREA ElementType = "textarea"
	BUTTON   ElementType = "button"
	FORM     ElementType = "form"

	AREA    ElementType = "area"
	BASE    ElementType = "base"
	BR      ElementType = "br"
	COL     ElementType = "col"
	COMMAND ElementType = "command"
	EMBED   ElementType = "embed"
	HR      ElementType = "hr"
	IMG     ElementType = "img"
	INPUT   ElementType = "input"
	KEYGEN  ElementType = "keygen"
	LINK    ElementType = "link"
	META    ElementType = "meta"
	PARAM   ElementType = "param"
	SOURCE  ElementType = "source"
	TRACK   ElementType = "track"
	WBR     ElementType = "wbr"

	HXDELETE    AttributeType = "hx-delete"
	HXPATCH     AttributeType = "hx-patch"
	HXPOST      AttributeType = "hx-post"
	HXGET       AttributeType = "hx-get"
	HXTRIGGER   AttributeType = "hx-trigger"
	HXINCLUDE   AttributeType = "hx-include"
	HXSWAPOOB   AttributeType = "hx-swap-oob"
	HXTARGET    AttributeType = "hx-target"
	HXSWAP      AttributeType = "hx-swap"
	SRC         AttributeType = "src"
	ALT         AttributeType = "alt"
	CHARSET     AttributeType = "charset"
	NAME        AttributeType = "name"
	CONTENT     AttributeType = "content"
	REL         AttributeType = "rel"
	HREF        AttributeType = "href"
	CLASS       AttributeType = "class"
	ID          AttributeType = "id"
	HYPERSCRIPT AttributeType = "_"
	HIDDEN      AttributeType = "hidden"
	LANG        AttributeType = "lang"
	VALUE       AttributeType = "value"
	FOR         AttributeType = "for"
	TYPE        AttributeType = "type"
	LIST        AttributeType = "list"
)

var voidElements = map[ElementType]struct{}{
	AREA:    {},
	BASE:    {},
	BR:      {},
	COL:     {},
	COMMAND: {},
	EMBED:   {},
	HR:      {},
	IMG:     {},
	INPUT:   {},
	KEYGEN:  {},
	LINK:    {},
	META:    {},
	PARAM:   {},
	SOURCE:  {},
	TRACK:   {},
	WBR:     {},
}
