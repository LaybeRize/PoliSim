package helper

import (
	"bytes"
	"fmt"
	"log/slog"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var generator = rand.New(rand.NewSource(time.Now().UnixNano()))
var matchColor = regexp.MustCompile(`(?m)^#[A-Fa-f0-9]{6}$`)

func init() {
	setupDiscord()
	setupMarkdown()
}

func Shutdown() {
	err := discord.Close()
	if err != nil {
		slog.Error(err.Error())
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
