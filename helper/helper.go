package helper

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var generator = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

func GetUniqueID(author string) string {
	sum := time.Now().UnixNano()
	for _, singleRune := range []rune(author) {
		sum += int64(singleRune)
	}

	suffix := make([]byte, 4)
	prefix := make([]byte, 4)
	generator.Read(suffix)
	generator.Read(prefix)

	suffix[0] += byte(sum)
	suffix[1] += byte(sum >> 8)
	suffix[2] += byte(sum >> 16)
	suffix[3] += byte(sum >> 24)
	prefix[0] += byte(sum >> 32)
	prefix[1] += byte(sum >> 40)
	prefix[2] += byte(sum >> 48)
	prefix[3] += byte(sum >> 56)

	return fmt.Sprintf("%X-%X", suffix, prefix)
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

// GetPureFormEntry returns the unchanged string
func GetPureFormEntry(request *http.Request, field string) string {
	return request.Form.Get(field)
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
