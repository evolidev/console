package color

import (
	"fmt"
	"regexp"
)

func Text(code int, value interface{}) string {
	return fmt.Sprintf("\u001b[38;5;%dm%s\u001b[0m", code, value)
}

func Bg(code int, value interface{}) string {
	return fmt.Sprintf("\u001b[48;5;%dm%s\u001b[0m", code, value)
}

func RemoveAsciiColors(s string) string {
	re := regexp.MustCompile("\x1b\\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]")
	return re.ReplaceAllString(s, "")
}
