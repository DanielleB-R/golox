package interpreter

import (
	"fmt"
	"strings"

	"github.com/DanielleB-R/golox/interpreter/token"
)

var _ error = (*SourceError)(nil)
var _ error = (*ParseError)(nil)

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

type ParseError struct {
	token   *token.Token
	message string
}

func (p *ParseError) Error() string {
	if p.token.TokenType == token.EOF {
		return fmt.Sprintf("[line %d at end] Error: %s", p.token.Line, p.message)
	}
	return fmt.Sprintf("[line %d at '%s'] Error: %s", p.token.Line, p.token.Lexeme, p.message)
}

type ParseErrors []*ParseError

func (s ParseErrors) Error() string {
	errorStrings := []string{}
	for _, e := range s {
		errorStrings = append(errorStrings, e.Error())
	}

	return strings.Join(errorStrings, "\n")
}

type RuntimeError struct {
	token   *token.Token
	message string
}

func (r *RuntimeError) Error() string {
	return fmt.Sprintf("Runtime error line %d: %s", r.token.Line, r.message)
}
