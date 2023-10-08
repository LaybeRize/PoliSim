package componentHelper

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

// ElementFunc is the Node function to return, if you want the text to
// be rendered inside the parent element
type ElementFunc func(io.Writer) error

// AttributeFunc is the Node function to return, if you want it rendered as an
// attribute on the parent element.
type AttributeFunc func(io.Writer) error

// GroupFunc is the Node function to return, if you want it group a bunch of children
// into a single node. This kind of Node can be put into an ElementFunc and gets rendered
// correctly anyway.
type GroupFunc func(w io.Writer, f func(io.Writer, Node) error) error

func (n ElementFunc) Render(w io.Writer) error {
	return n(w)
}

func (n AttributeFunc) Render(w io.Writer) error {
	return n(w)
}

func (n GroupFunc) Render(w io.Writer) error {
	return n(w, renderElementChild)
}

func (n GroupFunc) renderAttr(w io.Writer) error {
	return n(w, renderAttributeChild)
}

// Group groups children into a group that has the same hierarchy level.
// the Render function for this Node only renders elements. But it can be used in
// elements itself and will render all attributes in the group on the parent element of the group, while
// only rendering the elements correctly in the element itself.
func Group(children ...Node) Node {
	return GroupFunc(func(w io.Writer, f func(io.Writer, Node) error) (err error) {

		for _, c := range children {
			err = f(w, c)
			if err != nil {
				return
			}
		}
		return nil
	})
}

// El creates an element DOM Node with a name and child Nodes.
// See https://dev.w3.org/html5/spec-LC/syntax.html#elements-0 for how elements are rendered.
// No tags are ever omitted from normal tags, even though it's allowed for elements given at
// https://dev.w3.org/html5/spec-LC/syntax.html#optional-tags
// If an element is a void element, non-attribute children nodes are ignored.
func El(name ElementType, children ...Node) Node {
	return ElementFunc(func(w io.Writer) (err error) {

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

// Attr creates an attribute DOM Node with a name and optional value.
// If only a name is passed, it's a name-only (boolean) attribute (like "required").
// If a name and value are passed, it's a name-value attribute (like `class="header"`).
// Attr ignores more than the first provided parameter.
func Attr(name AttributeType, str ...string) Node {
	return AttributeFunc(func(w io.Writer) (err error) {
		_, err = w.Write([]byte(" " + name))
		if err != nil {
			return
		}
		if len(str) > 0 {
			if name == HXVALS {
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
	return ElementFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(template.HTMLEscapeString(fmt.Sprintf(format, args...))))
		return err
	})
}

// Raw creates a text DOM Node that just Renders the unescaped string format
// run through fmt.Sprintf beforehand with the given args.
func Raw(format string, args ...any) Node {
	return ElementFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(fmt.Sprintf(format, args...)))
		return err
	})
}

// If returns the Node if the statment is true, otherwise returns nil.
// for returning a different Node on a false statment see IfElse.
func If(statement bool, node Node) Node {
	if statement {
		return node
	}
	return nil
}

// IfElse returns the whenTrue Node when the statment evaluates to true
// otherwhise it returns the whenFalse Node.
func IfElse(statment bool, whenTrue Node, whenFalse Node) Node {
	if statment {
		return whenTrue
	}
	return whenFalse
}

// renderElementChild renders the child only if it is a ElementFunc or if it is a ElementFunc in a GroupFunc.
func renderElementChild(w io.Writer, c Node) error {
	switch c.(type) {
	case ElementFunc:
		return c.Render(w)
	case GroupFunc:
		return c.(GroupFunc).Render(w)
	}
	return nil
}

// renderAttributeChild renders the child only if it is a AttributeFunc or if it is a AttributeFunc in a GroupFunc.
func renderAttributeChild(w io.Writer, c Node) error {
	switch c.(type) {
	case AttributeFunc:
		return c.Render(w)
	case GroupFunc:
		return c.(GroupFunc).renderAttr(w)
	}
	return nil
}
