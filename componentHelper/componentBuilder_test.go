package componentHelper

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEl(t *testing.T) {
	var buf bytes.Buffer
	err := El("test").Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, "<test></test>", buf.String())
}

func TestAttr(t *testing.T) {
	var buf bytes.Buffer
	err := Attr("test").Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, " test", buf.String())

	buf = bytes.Buffer{}
	err = Attr("test", "val", "ahsdk").Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, " test=\"val\"", buf.String())

	buf = bytes.Buffer{}
	err = Attr(HXVALS, "val", "test2").Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, " "+string(HXVALS)+"='val'", buf.String())
}

func TestGroup(t *testing.T) {
	var buf bytes.Buffer
	err := Group(El("test"),
		Attr(HXVALS, "val", "test2"),
		El("other")).Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, "<test></test><other></other>", buf.String())
}

func TestMixing(t *testing.T) {
	var buf bytes.Buffer
	err := El("abc", Group(El("test"),
		Attr(HXVALS, "val", "test2"),
		El("other"))).Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, "<abc "+string(HXVALS)+"='val'><test></test><other></other></abc>", buf.String())
}
