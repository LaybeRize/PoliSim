package dataValidation

import (
	"PoliSim/database"
	"strings"
)

func ArrayContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetPositionOfString(input []string, value string) int {
	for p, v := range input {
		if v == value {
			return p
		}
	}
	return -1
}

func RemoveFromArray(s []string, i int) []string {
	if i == -1 {
		return s
	}
	if i == 0 && len(s) == 1 {
		return []string{}
	}
	s[i] = s[0]
	return s[1:]
}

func RemoveFirstStringOccurrenceFromArray(s []string, str string) []string {
	i := GetPositionOfString(s, str)
	return RemoveFromArray(s, i)
}

func TrimSuffix(s, suffix string) string {
	s = s[:len(s)-len(suffix)]
	return s
}

func TrimPrefix(s, prefix string) string {
	s = s[len(prefix):]
	return s
}

func RemoveDuplicates(array []string) []string {
	var result []string
	result = []string{}
	for _, val := range array {
		if GetPositionOfString(result, val) == -1 {
			result = append(result, val)
		}
	}
	return result
}

func ClearStringArray(array *[]string) {
	clone := make([]string, len(*array))
	copy(clone, *array)
	*array = []string{}
	for _, str := range clone {
		str = strings.TrimSpace(str)
		if str != "" && GetPositionOfString(*array, str) == -1 {
			*array = append(*array, str)
		}
	}
}

func DeleteMultiplesAndEmpty(a []string) []string {
	ClearStringArray(&a)
	return a
}

// RemoveEntriesFromList removes the first occurrence in list of every element form toRemove
func RemoveEntriesFromList(list []string, toRemove []string) []string {
	for _, str := range toRemove {
		list = RemoveFirstStringOccurrenceFromArray(list, str)
	}
	return list
}

type ValidationMessage struct {
	Message  string
	Positive bool
}

// isRoleValid checks if the role not database.NotLoggedIn
func isRoleValid(level int) bool {
	return level >= int(database.PressAccount) && level != int(database.NotLoggedIn) && level <= int(database.HeadAdmin)
}

func isEmptyOrNotInRange(str string, length int) bool {
	return str == "" || len([]rune(str)) > length
}
