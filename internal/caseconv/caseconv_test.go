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
