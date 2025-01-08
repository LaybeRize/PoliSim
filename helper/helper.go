package helper

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

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

// GetFormEntry strips the whitespace from the response and returns the result
func GetFormEntry(request *http.Request, field string) string {
	return strings.TrimSpace(request.Form.Get(field))
}

// GetFormList returns the list
func GetFormList(request *http.Request, field string) []string {
	userNames := request.Form[field]
	if userNames == nil {
		return []string{""}
	}
	for i, str := range userNames {
		userNames[i] = strings.TrimSpace(str)
	}
	return userNames
}
