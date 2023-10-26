package helper

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
	"log"
	"os"
	"regexp"
	"strings"
)

// ArrayContainsString returns true if the string is contained in the array.
func ArrayContainsString(s *[]string, str string) bool {
	for _, v := range *s {
		if v == str {
			return true
		}
	}

	return false
}

// GetPositionOfString returns the position of the first occurrence in the array or -1 if it is not contained within.
func GetPositionOfString(input *[]string, value string) int {
	for p, v := range *input {
		if v == value {
			return p
		}
	}
	return -1
}

// RemoveFromArray takes the position of the given the array, writes the value from position
// zero into it and cuts the first value of the array. It does not change the array when i == -1
// and empties the array if the length is 1 and i == 0.
// RemoveFromArray(["1","2","3","4","3"],2) -> ["2","1","4","3"]
func RemoveFromArray(s *[]string, i int) {
	if i == -1 {
		return
	}
	if i == 0 && len(*s) == 1 {
		*s = []string{}
		return
	}
	(*s)[i] = (*s)[0]
	*s = (*s)[1:]
	return
}

// RemoveFirstStringOccurrenceFromArray removes the first occurrence of the str parameter from the given array.
// if there is no such string in the array, the array will not be modified. For how the modification is processed view RemoveFromArray and GetPositionOfString.
func RemoveFirstStringOccurrenceFromArray(s *[]string, str string) {
	i := GetPositionOfString(s, str)
	RemoveFromArray(s, i)
}

// TrimSuffix measures the length of the suffix and takes that length away from the end of string.
// TrimSuffix("test","st") -> "te" but also TrimSuffix("test","s") -> "tes"
func TrimSuffix(s, suffix string) string {
	s = s[:len(s)-len(suffix)]
	return s
}

// TrimPrefix measures the length of the prefix and takes that length away from the start of the string.
// TrimPrefix("test","te") -> "st" but also TrimPrefix("test","s") -> "est"
func TrimPrefix(s, prefix string) string {
	s = s[len(prefix):]
	return s
}

// ClearStringArray takes in the pointer to a slice, resets
// the array and adds every unique string (after being trimmed) back
// into the array
// ["abc", "\n ", "abc ", "ab", "ab", ""] -> ["abc","ab"]
func ClearStringArray(array *[]string) {
	clone := make([]string, len(*array))
	copy(clone, *array)
	*array = make([]string, 0, len(clone))
	lookUp := make(map[string]struct{})
	lookUp[""] = struct{}{}
	for _, str := range clone {
		str = strings.TrimSpace(str)
		if _, ok := lookUp[str]; !ok {
			*array = append(*array, str)
			lookUp[str] = struct{}{}
		}
	}
}

// RemoveEntriesFromList removes the first occurrence in list of every element form toRemove
func RemoveEntriesFromList(list *[]string, toRemove []string) {
	clone := make([]string, len(*list))
	copy(clone, *list)
	*list = make([]string, 0, len(clone))
	// create lookup table for items not to take with
	lookUp := make(map[string]struct{})
	for _, str := range toRemove {
		lookUp[str] = struct{}{}
	}
	//look at every item in the original list and if it is not on the lookup table add it back to the list
	for _, str := range clone {
		if _, ok := lookUp[str]; !ok {
			*list = append(*list, str)
		}
	}
}

// CreateHTML creates correctly formated html from the markdown input
func CreateHTML(md string) string {
	intermediate := markdown.NormalizeNewlines([]byte(md))
	maybeUnsafeHTML := markdown.ToHTML(intermediate, nil, nil)
	htmlResult := bluemonday.UGCPolicy().Sanitize(string(maybeUnsafeHTML))
	return updateHtmlResult(htmlResult)
}

var replacerMap map[string]string

func updateHtmlResult(htmlResult string) string {
	htmlResult = strings.ReplaceAll(htmlResult, "<code>\n", "<code>")
	for key, val := range replacerMap {
		var withAttr = regexp.MustCompile(`(?m)(<` + regexp.QuoteMeta(key) + ` )`)
		var withoutAttr = regexp.MustCompile(`(?m)(<` + regexp.QuoteMeta(key) + `)>`)
		intermediate := fmt.Sprintf("$1 %s ", val)
		htmlResult = withAttr.ReplaceAllString(htmlResult, intermediate)
		intermediate = fmt.Sprintf("$1 %s>", val)
		htmlResult = withoutAttr.ReplaceAllString(htmlResult, intermediate)
	}
	return htmlResult
}

// UpdateAttributes updates the added attributes to the html tags for markdown formatting
func UpdateAttributes() {
	replacerMap = make(map[string]string)
	var re = regexp.MustCompile(`(?m)<(\w*?) (.*?)>`)
	var getTemplate = regexp.MustCompile(`(?s)<!-- Test start -->(.*)<!-- Test end -->`)
	b, err := os.ReadFile("resources/markdown.html")
	if err != nil {
		log.Fatalln(err)
	}
	b = getTemplate.FindAllSubmatch(b, -1)[0][1]
	for _, match := range re.FindAllSubmatch(b, -1) {
		replacerMap[string(match[1])] = string(match[2])
	}
}
