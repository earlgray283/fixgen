package caseconv

import (
	"unicode"
)

const (
	SepSnake = '_'
)

type (
	ConvertFunc            func(r rune) rune
	IsSeparatePositionFunc func(prev, current rune) bool
)

func ConvertPascalToSnake(s string) string {
	return ConvertToExpandCase(s, SepSnake, unicode.ToLower, func(prev, current rune) bool {
		return !unicode.IsUpper(prev) && unicode.IsUpper(current)
	})
}

func ConvertToExpandCase(s string, sep rune, convert ConvertFunc, isSeparatePosition IsSeparatePositionFunc) string {
	runeSlice := []rune(s)
	t := make([]rune, 0)
	for i, c := range runeSlice {
		if i > 0 && isSeparatePosition(runeSlice[i-1], c) {
			t = append(t, sep)
		}
		t = append(t, convert(c))
	}
	return string(t)
}
