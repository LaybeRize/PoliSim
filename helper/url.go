package helper

import (
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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

func GetAdvancedFormValuesWithoutDebugLogger(request *http.Request) (AdvancedValues, error) {
	err := request.ParseForm()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
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

const ISOTimeFormat = "2006-01-02T15:04:05.999999"

func (a AdvancedValues) GetUTCTime(field string, onExceptionNow bool) (time.Time, bool) {
	vs := a[field]
	if len(vs) == 0 {
		if onExceptionNow {
			return time.Now().UTC(), false
		}
		return time.Time{}, false
	}
	val, err := time.ParseInLocation(ISOTimeFormat, strings.TrimSpace(vs[0]), time.UTC)
	if err != nil {
		if onExceptionNow {
			return time.Now().UTC(), true
		}
		return time.Time{}, true
	}
	return val, true
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

func (a AdvancedValues) Encode() string {
	return url.Values(a).Encode()
}

func (a AdvancedValues) Has(field string) bool {
	return len(a[field]) != 0
}

func (a AdvancedValues) Exists(field string) bool {
	_, exists := a[field]
	return exists
}

func (a AdvancedValues) DeleteEmptyFields(fields []string) {
	for _, field := range fields {
		vs := a[field]
		if len(vs) == 0 {
			continue
		}
		if strings.TrimSpace(vs[0]) == "" {
			delete(a, field)
		}
	}
}

func (a AdvancedValues) DeleteFields(fields []string) {
	for _, field := range fields {
		delete(a, field)
	}
}
