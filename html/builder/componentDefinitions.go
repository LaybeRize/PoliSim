package builder

import (
	"bytes"
	"io"
)

type Node interface {
	//Render renders all Elements in the group into one string which is then written to the
	//parsed io.Writer. It returns an error if it encounters one
	Render(w io.Writer) error
	//PureNode returns a Node which has prerended itself into a byte buffer to improve performance.
	PureNode() Node
}

type GroupNode interface {
	//Render renders all Elements in the group into one string which is then written to the
	//parsed io.Writer. It returns an error if it encounters one
	Render(w io.Writer) error
	//PureNode returns a Node which has prerended itself into a byte buffer to improve performance.
	PureNode() Node
	//PureAttributes is only available for groups to prerender attributes to be used in elements.
	PureAttributes() Node
}

// elementFunc is the Node function to return, if you want the text to
// be rendered inside the parent element
type elementFunc func(io.Writer) error

// attributeFunc is the Node function to return, if you want it rendered as an
// attribute on the parent element.
type attributeFunc func(io.Writer) error

// groupFunc is the Node function to return, if you want it group a bunch of children
// into a single node. This kind of Node can be put into an elementFunc and gets rendered
// correctly anyway.
type groupFunc func(w io.Writer, f func(io.Writer, Node) error) error

func (n elementFunc) Render(w io.Writer) error {
	return n(w)
}

func (n elementFunc) PureNode() Node {
	buff := &bytes.Buffer{}
	err := n(buff)
	if err != nil {
		return nil
	}
	return elementFunc(func(w io.Writer) error {
		_, elErr := w.Write(buff.Bytes())
		return elErr
	})
}

func (n attributeFunc) Render(w io.Writer) error {
	return n(w)
}

func (n attributeFunc) PureNode() Node {
	buff := &bytes.Buffer{}
	err := n(buff)
	if err != nil {
		return nil
	}
	return attributeFunc(func(w io.Writer) error {
		_, elErr := w.Write(buff.Bytes())
		return elErr
	})
}

func (n groupFunc) Render(w io.Writer) error {
	return n(w, renderElementChild)
}

func (n groupFunc) PureNode() Node {
	buff := &bytes.Buffer{}
	err := n(buff, renderElementChild)
	if err != nil {
		return nil
	}
	return elementFunc(func(w io.Writer) error {
		_, elErr := w.Write(buff.Bytes())
		return elErr
	})
}

func (n groupFunc) PureAttributes() Node {
	buff := &bytes.Buffer{}
	err := n(buff, renderAttributeChild)
	if err != nil {
		return nil
	}
	return attributeFunc(func(w io.Writer) error {
		_, elErr := w.Write(buff.Bytes())
		return elErr
	})
}

func (n groupFunc) renderAttr(w io.Writer) error {
	return n(w, renderAttributeChild)
}

func elementWrapper(str string) func(nodes ...Node) Node {
	return func(nodes ...Node) Node {
		return el(str, nodes...)
	}
}

func attributeWrapper(name string) func(str ...string) Node {
	return func(str ...string) Node {
		return attr(name, str...)
	}
}
