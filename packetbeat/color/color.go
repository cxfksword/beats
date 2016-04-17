package color

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/shiena/ansicolor"
)

var w = ansicolor.NewAnsiColorWriter(os.Stdout)
var regSpace = regexp.MustCompile(`\s`)

const (
	Black = uint8(iota + 30)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Default      = uint8(39)
	Gray         = uint8(90)
	LightRed     = uint8(91)
	LightGreen   = uint8(92)
	LightYellow  = uint8(93)
	LightBlue    = uint8(94)
	LightMagenta = uint8(95)
	LightCyan    = uint8(96)
	LightWhite   = uint8(97)

	EndColor         = "\033[0m"
	printV           = "A"
	contentJsonRegex = `application/json`
)

func Color(i interface{}, color uint8) string {
	str := fmt.Sprintf("%v", i)
	return fmt.Sprintf("%s%s%s", ColorStart(color), str, EndColor)
}

func ColorStart(color uint8) string {
	return fmt.Sprintf("\033[%dm", color)
}

func StatusCodeColor(statusCode int) string {
	switch statusCode {
	case 301:
		fallthrough
	case 302:
		fallthrough
	case 304:
		return Color(statusCode, Cyan)
	case 403:
		fallthrough
	case 404:
		return Color(statusCode, Yellow)
	case 500:
		fallthrough
	case 503:
		return Color(statusCode, Red)
	default:
		return strconv.Itoa(statusCode)
	}
}

func MethodColor(method string) string {
	switch method {
	case "GET":
		return Color(method, White)
	case "POST":
		return Color(method, LightGreen)
	case "DELETE":
		return Color(method, Red)
	case "PUT":
		return Color(method, LightBlue)
	default:
		return Color(method, Magenta)
	}
}

func ContentColor(content string) string {
	content = regSpace.ReplaceAllString(content, "")
	if len(content) > 50 {
		content = content[:50] + "..."
	}
	return Color(content, Gray)
}

func Println(str string, color uint8) {
	fmt.Fprintln(w, Color(str, color))
}

func Print(str string, color uint8) {
	fmt.Fprint(w, Color(str, color))
}

func Printf(format string, color uint8, params ...interface{}) {
	fmt.Fprint(w, Color(fmt.Sprintf(format, params...), color))
}

func ColorfulJson(str string) string {
	var rsli []rune
	var key, val, startcolor, endcolor, startsemicolon bool
	var prev rune
	for _, char := range []rune(str) {
		switch char {
		case ' ':
			rsli = append(rsli, char)
		case '{':
			startcolor = true
			key = true
			val = false
			rsli = append(rsli, char)
		case '}':
			startcolor = false
			endcolor = false
			key = false
			val = false
			rsli = append(rsli, char)
		case '"':
			if startsemicolon && prev == '\\' {
				rsli = append(rsli, char)
			} else {
				if startcolor {
					rsli = append(rsli, char)
					if key {
						rsli = append(rsli, []rune(ColorStart(Magenta))...)
					} else if val {
						rsli = append(rsli, []rune(ColorStart(Cyan))...)
					}
					startsemicolon = true
					key = false
					val = false
					startcolor = false
				} else {
					rsli = append(rsli, []rune(EndColor)...)
					rsli = append(rsli, char)
					endcolor = true
					startsemicolon = false
				}
			}
		case ',':
			if !startsemicolon {
				startcolor = true
				key = true
				val = false
				if !endcolor {
					rsli = append(rsli, []rune(EndColor)...)
					endcolor = true
				}
			}
			rsli = append(rsli, char)
		case ':':
			if !startsemicolon {
				key = false
				val = true
				startcolor = true
				if !endcolor {
					rsli = append(rsli, []rune(EndColor)...)
					endcolor = true
				}
			}
			rsli = append(rsli, char)
		case '\n', '\r', '[', ']':
			rsli = append(rsli, char)
		default:
			if !startsemicolon {
				if key && startcolor {
					rsli = append(rsli, []rune(ColorStart(Magenta))...)
					key = false
					startcolor = false
					endcolor = false
				}
				if val && startcolor {
					rsli = append(rsli, []rune(ColorStart(Cyan))...)
					val = false
					startcolor = false
					endcolor = false
				}
			}
			rsli = append(rsli, char)
		}
		prev = char
	}
	return string(rsli)
}

func ColorfulHTML(str string) string {
	return Color(str, White)
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}
