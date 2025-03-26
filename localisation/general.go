package loc

import (
	"strings"
)

func SetHomePage(text []byte) {
	replaceMap["_home"]["$$home-page$$"] = string(text)
}

func LocaliseTemplateString(input []byte, name string) string {
	result := string(input)
	for key, value := range replaceMap {
		if name == key {
			for left, right := range value {
				result = strings.ReplaceAll(result, left, right)
			}
		}
	}
	return result
}

func CleanUpMap() {
	replaceMap = nil
}

func IsAdministrationName(name string) bool {
	res := true
	switch name {
	case "Administration": //DE, EN
	default:
		res = false
	}
	return res
}
