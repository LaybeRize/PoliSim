package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClearStringArray(t *testing.T) {
	test := []string{"abc", "\n ", "abc ", "ab", "ab", ""}
	ClearStringArray(&test)
	assert.Equal(t, []string{"abc", "ab"}, test)
}

func TestRemoveEntriesFromList(t *testing.T) {
	list := []string{"a", "b", "c", "d", "e", "f"}
	toRemove := []string{"b", "c", "e", "f", "q", "l", "v"}
	RemoveEntriesFromList(&list, toRemove)
	assert.Equal(t, []string{"a", "d"}, list)
}

func TestTrimPrefix(t *testing.T) {
	assert.Equal(t, "st", TrimPrefix("test", "te"))
	assert.Equal(t, "est", TrimPrefix("test", "s"))
}

func TestTrimSuffix(t *testing.T) {
	assert.Equal(t, "te", TrimSuffix("test", "st"))
	assert.Equal(t, "tes", TrimSuffix("test", "s"))
}

func TestRemoveFromArray(t *testing.T) {
	arr := []string{"1", "2", "3", "4", "3"}
	RemoveFromArray(&arr, 2)
	assert.Equal(t, []string{"2", "1", "4", "3"}, arr)
}

func TestRemoveFirstStringOccurrenceFromArray(t *testing.T) {
	arr := []string{"1", "2", "3", "4", "3"}
	RemoveFirstStringOccurrenceFromArray(&arr, "3")
	assert.Equal(t, []string{"2", "1", "4", "3"}, arr)
}
