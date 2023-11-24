package serving

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventGenerator(t *testing.T) {
	event, err := formatServerSentEvent[string]("", false, "", func(info string, isAdmin bool, uuidStr string) (*SendEventStruct, error) {
		return &SendEventStruct{
			HTML:       "<test attr=\"test\"></test>",
			HTMXupdate: "abc",
			Target:     "def",
		}, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "event: change\ndata: {\"data\":\"<test attr=\\\"test\\\"></test>\",\"updateID\":\"abc\",\"targetID\":\"def\"}\n\n", event)
}
