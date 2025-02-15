package caseconv

import (
	"strings"
	"unicode"
)

func ConvertPascalToSnake(s string) string {
	t := &strings.Builder{}
	t.Grow(len(s))

	for i, c := range s {
		if i > 0 && unicode.IsLower(rune(s[i-1])) && unicode.IsUpper(rune(c)) {
			t.WriteByte('_')
		}
		t.WriteRune(unicode.ToLower(c))
	}

	return t.String()
}
