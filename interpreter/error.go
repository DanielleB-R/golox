package interpreter

import (
	"fmt"
	"strings"
)

var _ error = (*SourceError)(nil)

type SourceError struct {
	line    int
	where   string
	message string
}

func NewSourceError(line int, where string, message string) *SourceError {
	return &SourceError{line, where, message}
}

func (s *SourceError) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s", s.line, s.where, s.message)
}

type SourceErrors []*SourceError

func (s SourceErrors) Error() string {
	errorStrings := []string{}
	for _, e := range s {
		errorStrings = append(errorStrings, e.Error())
	}

	return strings.Join(errorStrings, "\n")
}
