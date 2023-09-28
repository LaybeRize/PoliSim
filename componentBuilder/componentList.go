package componentBuilder

type ElementType string
type AttributeType string

const (
	DIV   ElementType = "div"
	HTML  ElementType = "html"
	HEAD  ElementType = "head"
	BODY  ElementType = "body"
	TITLE ElementType = "title"

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

	HXPOST AttributeType = "hx-post"
	SRC    AttributeType = "src"
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
