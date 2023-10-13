package builder

import (
	"fmt"
	"html/template"
	"io"
)

// RenderHTMLDoc returns a Node that writes the html doctype to the io.Writer
func RenderHTMLDoc() Node {
	return Raw("<!DOCTYPE html>")
}

type Node interface {
	Render(w io.Writer) error
}

// elementFunc is the Node function to return, if you want the text to
// be rendered inside the parent element
type elementFunc func(io.Writer) error

// attributeFunc is the Node function to return, if you want it rendered as an
// attribute on the parent element.
type attributeFunc func(io.Writer) error

type elementReturn func(nodes ...Node) Node

func elementWrapper(str elementType) elementReturn {
	return func(nodes ...Node) Node {
		return el(str, nodes...)
	}
}

type attributeReturn func(str ...string) Node

func attributeWrapper(name attributeType) attributeReturn {
	return func(str ...string) Node {
		return attr(name, str...)
	}
}

// groupFunc is the Node function to return, if you want it group a bunch of children
// into a single node. This kind of Node can be put into an elementFunc and gets rendered
// correctly anyway.
type groupFunc func(w io.Writer, f func(io.Writer, Node) error) error

func (n elementFunc) Render(w io.Writer) error {
	return n(w)
}

func (n attributeFunc) Render(w io.Writer) error {
	return n(w)
}

func (n groupFunc) Render(w io.Writer) error {
	return n(w, renderElementChild)
}

func (n groupFunc) renderAttr(w io.Writer) error {
	return n(w, renderAttributeChild)
}

// Group groups children into a group that has the same hierarchy level.
// the Render function for this Node only renders elements. But it can be used in
// elements itself and will render all attributes in the group on the parent element of the group, while
// only rendering the elements correctly in the element itself.
func Group(children ...Node) Node {
	return groupFunc(func(w io.Writer, f func(io.Writer, Node) error) (err error) {

		for _, c := range children {
			err = f(w, c)
			if err != nil {
				return
			}
		}
		return nil
	})
}

// el creates an element DOM Node with a name and child Nodes.
// See https://dev.w3.org/html5/spec-LC/syntax.html#elements-0 for how elements are rendered.
// No tags are ever omitted from normal tags, even though it's allowed for elements given at
// https://dev.w3.org/html5/spec-LC/syntax.html#optional-tags
// if an element is a void element, non-attribute children nodes are ignored.
func el(name elementType, children ...Node) Node {
	return elementFunc(func(w io.Writer) (err error) {

		_, err = w.Write([]byte("<" + name))
		if err != nil {
			return
		}

		for _, c := range children {
			err = renderAttributeChild(w, c)
			if err != nil {
				return
			}
		}

		_, err = w.Write([]byte(">"))
		if err != nil {
			return
		}

		if _, ok := voidElements[name]; ok {
			return
		}

		for _, c := range children {
			err = renderElementChild(w, c)
			if err != nil {
				return
			}
		}

		_, err = w.Write([]byte("</" + name + ">"))
		return
	})
}

// attr creates an attribute DOM Node with a name and optional value.
// If only a name is passed, it's a name-only (boolean) attribute (like "required").
// If a name and value are passed, it's a name-value attribute (like `class="header"`).
// attr ignores more than the first provided parameter.
func attr(name attributeType, str ...string) Node {
	return attributeFunc(func(w io.Writer) (err error) {
		_, err = w.Write([]byte(" " + name))
		if err != nil {
			return
		}
		if len(str) > 0 {
			// special case for hx-vals because it uses json as input text
			if name == hxValue {
				_, err = w.Write([]byte("='" + str[0] + "'"))
				return
			}
			_, err = w.Write([]byte("=\"" + str[0] + "\""))
		}
		return
	})
}

// Text creates a text DOM Node that Renders the escaped string format
// run through fmt.Sprintf beforehand with the given args.
func Text(format string, args ...any) Node {
	return elementFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(template.HTMLEscapeString(fmt.Sprintf(format, args...))))
		return err
	})
}

// Raw creates a text DOM Node that just Renders the unescaped string format
// run through fmt.Sprintf beforehand with the given args.
func Raw(format string, args ...any) Node {
	return elementFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(fmt.Sprintf(format, args...)))
		return err
	})
}

// If returns the Node if the statement is true, otherwise returns nil.
// for returning a different Node on a false statement see IfElse.
func If(statement bool, node Node) Node {
	if statement {
		return node
	}
	return nil
}

// IfElse returns the whenTrue Node when the statement evaluates to true
// otherwise it returns the whenFalse Node.
func IfElse(statement bool, whenTrue Node, whenFalse Node) Node {
	if statement {
		return whenTrue
	}
	return whenFalse
}

// renderElementChild renders the child only if it is a elementFunc or if it is a elementFunc in a groupFunc.
func renderElementChild(w io.Writer, c Node) error {
	switch c.(type) {
	case elementFunc:
		return c.Render(w)
	case groupFunc:
		return c.(groupFunc).Render(w)
	}
	return nil
}

// renderAttributeChild renders the child only if it is a attributeFunc or if it is a attributeFunc in a groupFunc.
func renderAttributeChild(w io.Writer, c Node) error {
	switch c.(type) {
	case attributeFunc:
		return c.Render(w)
	case groupFunc:
		return c.(groupFunc).renderAttr(w)
	}
	return nil
}
