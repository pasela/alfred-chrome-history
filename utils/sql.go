package utils

import "strings"

func EscapeLike(value string, escChar rune) string {
	esc := string(escChar)
	repr := strings.NewReplacer(
		"%", esc+"%",
		"_", esc+"_",
		esc, esc+esc,
	)
	return repr.Replace(value)
}
