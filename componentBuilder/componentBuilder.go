package componentBuilder

import (
	"html/template"
	"io"
)

// RenderDoctype writes the html doctype to the io.Writer
// and returns the error, if one occures.
func RenderDoctype(w io.Writer) error {
	_, err := w.Write([]byte("<!DOCTYPE html>"))
	return err
}

type Node interface {
	Render(w io.Writer) error
}

type ElementFunc func(io.Writer) error
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
			err = renderAttributes(w, c)
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
			err = renderElements(w, c)
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
			_, err = w.Write([]byte("=\"" + str[0] + "\""))
		}
		return
	})
}

// Text creates a text DOM Node that Renders the escaped string t.
func Text(t string) Node {
	return ElementFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(template.HTMLEscapeString(t)))
		return err
	})
}

// Raw creates a text DOM Node that just Renders the unescaped string t.
func Raw(t string) Node {
	return ElementFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(t))
		return err
	})
}

func renderElements(w io.Writer, c Node) error {
	switch c.(type) {
	case ElementFunc:
		return c.Render(w)
	case AttributeFunc:
		return nil
	default:
		panic("Node type is not allowed!")
	}
	return nil
}

func renderAttributes(w io.Writer, c Node) error {
	switch c.(type) {
	case ElementFunc:
		return nil
	case AttributeFunc:
		return c.Render(w)
	default:
		panic("Node type is not allowed!")
	}
	return nil
}
