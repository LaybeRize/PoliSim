//go:build EN

package loc

import "strings"

const (
	AdminstrationName            = "Administration"
	AdminstrationAccountName     = "John Administrator"
	AdminstrationAccountUsername = ""
	AdminstrationAccountPassword = ""
)

var replaceMap = map[string]string{}

func LocaliseTemplateString(input []byte) string {
	result := string(input)
	for key, value := range replaceMap {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}
