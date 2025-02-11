package caseconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConvertPascalToSnake(t *testing.T) {
	tcs := map[string]string{
		"PascalCase": "pascal_case",
		"ID":         "id",
		"UserID":     "user_id",
		"I":          "i",
		"IDs":        "ids",
	}

	for input, expect := range tcs {
		t.Run(input, func(t *testing.T) {
			output := ConvertPascalToSnake(input)
			assert.Equal(t, expect, output)
		})
	}
}

func Test_ConvertSnakeToPascal(t *testing.T) {
	tcs := map[string]string{
		"pascal_case": "PascalCase",
		"id":          "ID",
		"user_id":     "UserID",
		"i":           "I",
		"ids":         "IDs",
		"idss":        "IDss", // 仕様 :(
	}

	for input, expect := range tcs {
		t.Run(input, func(t *testing.T) {
			output := ConvertSnakeToPascal(input)
			assert.Equal(t, expect, output)
		})
	}
}
