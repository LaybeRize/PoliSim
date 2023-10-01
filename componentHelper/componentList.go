package componentHelper

type ElementType string
type AttributeType string

const (
	DIV    ElementType = "div"
	P      ElementType = "p"
	HTML   ElementType = "html"
	HEAD   ElementType = "head"
	BODY   ElementType = "body"
	TITLE  ElementType = "title"
	SCRIPT ElementType = "script"

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

	HXPOST    AttributeType = "hx-post"
	HXGET     AttributeType = "hx-get"
	HXTRIGGER AttributeType = "hx-trigger"
	HXINCLUDE AttributeType = "hx-include"
	HXSWAP    AttributeType = "hx-swap"
	SRC       AttributeType = "src"
	CHARSET   AttributeType = "charset"
	NAME      AttributeType = "name"
	CONTENT   AttributeType = "content"
	REL       AttributeType = "rel"
	HREF      AttributeType = "href"
	CLASS     AttributeType = "class"
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
