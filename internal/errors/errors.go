package errors

import "fmt"

type UnsupportedTypeError struct {
	fieldName string
	typeName  string
}

func (s *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("unsupported type error(field: %s, type: %s)", s.fieldName, s.typeName)
}

func NewUnsupportedTypeError(fieldName, typeName string) *UnsupportedTypeError {
	return &UnsupportedTypeError{fieldName, typeName}
}
