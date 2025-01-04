package helper

import (
	"fmt"
	"strings"
	"time"
)

func GetUniqueID(author string) string {
	authorRunes := []rune(author)
	sum := 0
	for i, singleRune := range authorRunes {
		if i > 3 {
			break
		}
		sum += int(singleRune)
	}
	return fmt.Sprintf("%x-%x", sum, time.Now().UnixNano()/1000000)
}

func MakeCommaSeperatedStringToList(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return make([]string, 0)
	}
	arr := strings.Split(input, ",")
	result := make([]string, 0, len(arr))
	for _, element := range arr {
		element = strings.TrimSpace(element)
		if element != "" {
			result = append(result, element)
		}
	}
	return result
}
