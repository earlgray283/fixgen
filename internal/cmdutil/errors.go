package cmdutil

import "fmt"

type ExitCodeError struct {
	ExitCode int
	Message  string
}

func (e *ExitCodeError) Error() string {
	return fmt.Sprintf("exit code %d: %s", e.ExitCode, e.Message)
}

func NewExitCodeError(exitCode int, message string) *ExitCodeError {
	return &ExitCodeError{exitCode, message}
}
