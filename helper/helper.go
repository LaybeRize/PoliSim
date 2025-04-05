package helper

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
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

	prefix := make([]byte, 4)
	suffix := make([]byte, 8)
	generator.Read(prefix)
	for pos, singleRune := range []byte(author) {
		prefix[pos%4] += singleRune
	}

	timeNano := time.Now().UnixNano()
	suffix[0] += byte(timeNano)
	suffix[1] += byte(timeNano >> 8)
	suffix[2] += byte(timeNano >> 16)
	suffix[3] += byte(timeNano >> 24)
	suffix[4] += byte(timeNano >> 32)
	suffix[5] += byte(timeNano >> 40)
	suffix[6] += byte(timeNano >> 48)
	suffix[7] += byte(timeNano >> 56)

	return fmt.Sprintf("ID-%X-%X", prefix, suffix)
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

func EscapeStringForJSON(src string) string {
	var buf bytes.Buffer
	start := 0
	for i := 0; i < len(src); {
		if b := src[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			buf.WriteString(src[start:i])
			switch b {
			case '\\', '"':

				buf.Write([]byte{'\\', b})
			case '\b':
				buf.Write([]byte{'\\', 'b'})
			case '\f':
				buf.Write([]byte{'\\', 'f'})
			case '\n':
				buf.Write([]byte{'\\', 'n'})
			case '\r':
				buf.Write([]byte{'\\', 'r'})
			case '\t':
				buf.Write([]byte{'\\', 't'})
			default:
				// This encodes bytes < 0x20 except for \b, \f, \n, \r and \t.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				buf.Write([]byte{'\\', 'u', '0', '0', hex[b>>4], hex[b&0xF]})
			}
			i++
			start = i
			continue
		}
		n := len(src) - i
		if n > utf8.UTFMax {
			n = utf8.UTFMax
		}
		c, size := utf8.DecodeRuneInString(src[i : i+n])
		if c == utf8.RuneError && size == 1 {
			buf.WriteString(src[start:i] + `\ufffd`)
			i += size
			start = i
			continue
		}
		if c == '\u2028' || c == '\u2029' {
			buf.WriteString(src[start:i])
			buf.Write([]byte{'\\', 'u', '2', '0', '2', hex[c&0xF]})
			i += size
			start = i
			continue
		}
		i += size
	}
	buf.WriteString(src[start:])
	return buf.String()
}

const hex = "0123456789abcdef"

var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}
