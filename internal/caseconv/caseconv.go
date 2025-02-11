package caseconv

import (
	"strings"
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

func ConvertSnakeToPascal(s string) string {
	return ConvertToPascal(s, SepSnake)
}

func ConvertToPascal(s string, sep rune) string {
	runeSlice := []rune(s)
	t := make([]rune, 0, len(s))
	shouldConvert := true
	for i, c := range runeSlice {
		if i > 0 && c == sep {
			shouldConvert = true
			continue
		}
		if shouldConvert {
			c = unicode.ToUpper(c)
		}
		t = append(t, c)
		shouldConvert = false
	}

	for i := 0; i+1 < len(t); {
		sw := isSpecialWord(t[i:])
		if sw != "" {
			for j := i; j < i+len(sw); j++ {
				t[j] = unicode.ToUpper(t[j])
			}
			i += len(sw)
		} else {
			i++
		}
	}

	return string(t)
}

var specialWords = []string{
	"Id", "Http", "Https",
	"Url", "Uri", "Url",
}

func isSpecialWord(s []rune) string {
	sstr := string(s)
	runeSlice := []rune(s)
	for _, sw := range specialWords {
		swLen := len([]rune(sw))
		if strings.HasPrefix(sstr, sw) && (len(runeSlice) == swLen || unicode.IsUpper(runeSlice[swLen]) || runeSlice[swLen] == 's') {
			return sw
		}
	}
	return ""
}
