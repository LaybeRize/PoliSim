package builder

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestEl(t *testing.T) {
	var buf bytes.Buffer
	err := el("test").Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, "<test></test>", buf.String())
}

func TestAttr(t *testing.T) {
	var buf bytes.Buffer
	err := attr("test").Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, " test", buf.String())

	buf = bytes.Buffer{}
	err = attr("test", "val", "ahsdk").Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, " test=\"val\"", buf.String())

	buf = bytes.Buffer{}
	err = attr("hxValue", "val\"", "test2").Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, " hxValue='val\"'", buf.String())
}

func TestGroup(t *testing.T) {
	var buf bytes.Buffer
	err := Group(el("test"),
		attr("hxValue", "val", "test2"),
		el("other")).Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, "<test></test><other></other>", buf.String())
}

func TestMixing(t *testing.T) {
	var buf bytes.Buffer
	err := el("abc", Group(el("test"),
		attr("hxValue", "val\"", "test2"),
		el("other"))).Render(&buf)
	assert.Nil(t, err)
	assert.Equal(t, "<abc hxValue='val\"'><test></test><other></other></abc>", buf.String())
}

func TestPathEscape(t *testing.T) {
	testStr := "\""
	assert.Equal(t, "%22", url.PathEscape(testStr))
}
