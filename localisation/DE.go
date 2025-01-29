//go:build DE

package loc

import "strings"

const (
	AdministrationName            = "Administration"
	AdministrationAccountName     = "Max Musteradminstrator"
	AdministrationAccountUsername = ""
	AdministrationAccountPassword = ""
	StandardColorName             = "Standard Farbe"
)

var replaceMap = map[string]string{}

func LocaliseTemplateString(input []byte) string {
	result := string(input)
	for key, value := range replaceMap {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}
