package helper

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var generator = rand.New(rand.NewSource(time.Now().UnixNano()))
var matchColor = regexp.MustCompile(`(?m)^#[A-Fa-f0-9]{6}$`)

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

func FilterList(list []string) []string {
	result := make([]string, 0, len(list))
	for _, element := range list {
		if strings.TrimSpace(element) != "" {
			result = append(result, element)
		}
	}
	return result
}

func StringIsAColor(input string) bool {
	return matchColor.FindString(input) != ""
}

type AdvancedValues map[string][]string

func GetAdvancedFormValues(request *http.Request) (AdvancedValues, error) {
	err := request.ParseForm()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	slog.Debug("Reading Form: ", "URL", request.URL.EscapedPath(), "Mapping", request.Form)
	return AdvancedValues(request.Form), nil
}

func GetAdvancedURLValues(request *http.Request) AdvancedValues {
	slog.Debug("Reading Form: ", "URL", request.URL.EscapedPath(), "Mapping", request.URL.Query())
	return AdvancedValues(request.URL.Query())
}

func (a AdvancedValues) MergeIntoMe(otherValues AdvancedValues) AdvancedValues {
	for key, value := range otherValues {
		a[key] = append(a[key], value...)
	}
	return a
}

func (a AdvancedValues) GetString(field string) string {
	vs := a[field]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (a AdvancedValues) GetTrimmedString(field string) string {
	vs := a[field]
	if len(vs) == 0 {
		return ""
	}
	return strings.TrimSpace(vs[0])
}

func (a AdvancedValues) GetArray(field string) []string {
	vs := a[field]
	if len(vs) == 0 {
		return []string{}
	}
	return vs
}

func (a AdvancedValues) GetTrimmedArray(field string) []string {
	vs := a[field]
	if len(vs) == 0 {
		return []string{}
	}
	for i, str := range vs {
		vs[i] = strings.TrimSpace(str)
	}
	return vs
}

func (a AdvancedValues) GetCommaSeperatedArray(field string) []string {
	vs := a[field]
	if len(vs) == 0 {
		return []string{}
	}
	return MakeCommaSeperatedStringToList(vs[0])
}

func (a AdvancedValues) GetFilteredArray(field string) []string {
	result := make([]string, 0, len(a[field]))
	for _, element := range a[field] {
		if str := strings.TrimSpace(element); str != "" {
			result = append(result, str)
		}
	}
	return result
}

func (a AdvancedValues) GetBool(field string) bool {
	vs := a[field]
	if len(vs) == 0 {
		return false
	}
	return strings.TrimSpace(vs[0]) == "true"
}

func (a AdvancedValues) GetInt(field string) int {
	vs := a[field]
	if len(vs) == 0 {
		return -1
	}
	res, err := strconv.Atoi(vs[0])
	if err != nil {
		return -1
	}
	return res
}

func (a AdvancedValues) GetTime(field string, format string, location *time.Location) time.Time {
	vs := a[field]
	if len(vs) == 0 {
		return time.Time{}
	}
	val, err := time.ParseInLocation(format, vs[0], location)
	if err != nil {
		return time.Time{}
	}
	return val
}

func (a AdvancedValues) Has(field string) bool {
	return len(a[field]) != 0
}

func (a AdvancedValues) Exists(field string) bool {
	_, exists := a[field]
	return exists
}
