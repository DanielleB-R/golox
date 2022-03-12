package token

import "fmt"

const (
	// Single-char tokens
	LEFT_PAREN = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two char tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type Token struct {
	TokenType int
	Lexeme    string
	Literal   interface{}
	Line      int
}

func NewToken(tokenType int, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%d %s %v", t.TokenType, t.Lexeme, t.Literal)
}
