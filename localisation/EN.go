//go:build EN

package loc

import "strings"

const (
	AdministrationName            = "Administration"
	AdministrationAccountName     = "John Administrator"
	AdministrationAccountUsername = ""
	AdministrationAccountPassword = ""
	StandardColorName             = "Standard Color"
)

var replaceMap = map[string]string{}

func LocaliseTemplateString(input []byte) string {
	result := string(input)
	for key, value := range replaceMap {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}
