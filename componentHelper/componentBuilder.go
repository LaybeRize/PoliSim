package componentHelper

import (
	"fmt"
	"html/template"
	"io"
)

// RenderHTMLDoc writes the html doctype to the io.Writer
func RenderHTMLDoc(w io.Writer) error {
	_, err := w.Write([]byte("<!DOCTYPE html>"))
	return err
}

type Node interface {
	Render(w io.Writer) error
}

// ElementFunc is the Node function to return, if you want the text to
// be rendered inside the parent element
type ElementFunc func(io.Writer) error

// AttributeFunc is the Node function to return, if you want it rendered as an
// on the parent element.
type AttributeFunc func(io.Writer) error

func (n ElementFunc) Render(w io.Writer) error {
	return n(w)
}

func (n AttributeFunc) Render(w io.Writer) error {
	return n(w)
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
			err = renderChild[AttributeFunc](w, c)
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
			err = renderChild[ElementFunc](w, c)
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

// renderChild renders the child, if it has the given function type.
// For only rendering Elements use renderChild[ElementFunc](writer, node)
func renderChild[t AttributeFunc | ElementFunc](w io.Writer, c Node) error {
	switch c.(type) {
	case t:
		return c.Render(w)
	}
	return nil
}
