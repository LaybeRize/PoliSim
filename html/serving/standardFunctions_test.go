package serving

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

type Test struct {
	Test1 string   `input:"test"`
	Test2 bool     `input:"bazinga"`
	Test3 []string `input:"sliceStrings"`
	Test4 int      `input:"goingWrong"`
	Test5 int64    `input:"gettingCorrect"`
	Test7 string
	Test8 bool
}

func TestFillStruct(t *testing.T) {
	req := httptest.NewRequest("GET", "https://google.com", strings.NewReader(""))
	err := req.ParseForm()
	assert.Nil(t, err)
	req.PostForm = map[string][]string{
		"test":           {"    testString \n", "not necessary"},
		"bazinga":        {"\ttrue  ", "false"},
		"sliceStrings":   {"test1", "test2", "    test1\t"},
		"goingWrong":     {"test"},
		"gettingCorrect": {"512"},
	}
	val := &Test{
		Test1: "asdsad",
		Test2: false,
		Test3: []string{"rttesfd", "basd"},
		Test4: 1242,
		Test5: 514213,
		Test7: "notChanged",
		Test8: false,
	}
	err = extractFormValuesForFields(val, req, 12)
	assert.Nil(t, err)
	assert.Equal(t, Test{
		Test1: "testString",
		Test2: true,
		Test3: []string{"test1", "test2"},
		Test4: 12,
		Test5: 512,
		Test7: "notChanged",
		Test8: false,
	}, *val)
}
