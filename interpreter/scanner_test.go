package interpreter

import (
	"testing"

	"github.com/DanielleB-R/golox/interpreter/token"
	"github.com/stretchr/testify/require"
)

func scan(source string) ([]*token.Token, error) {
	scanner := NewSourceScanner(source)
	return scanner.ScanTokens()
}

func TestBasicTokens(t *testing.T) {
	cases := map[string]*token.Token{
		"(":        token.NewToken(token.LEFT_PAREN, "(", nil, 1),
		")":        token.NewToken(token.RIGHT_PAREN, ")", nil, 1),
		"+":        token.NewToken(token.PLUS, "+", nil, 1),
		"!":        token.NewToken(token.BANG, "!", nil, 1),
		"!=":       token.NewToken(token.BANG_EQUAL, "!=", nil, 1),
		"=":        token.NewToken(token.EQUAL, "=", nil, 1),
		" -\t":     token.NewToken(token.MINUS, "-", nil, 1),
		"\"test\"": token.NewToken(token.STRING, "\"test\"", "test", 1),
		"1.2":      token.NewToken(token.NUMBER, "1.2", float64(1.2), 1),
		"test":     token.NewToken(token.IDENTIFIER, "test", nil, 1),
		"var":      token.NewToken(token.VAR, "var", nil, 1),
	}
	eofToken := token.NewToken(token.EOF, "", nil, 1)

	for input, token := range cases {
		t.Run(input, func(t *testing.T) {
			tokens, err := scan(input)
			require.NoError(t, err)
			require.Len(t, tokens, 2)
			require.Equal(t, token, tokens[0])
			require.Equal(t, eofToken, tokens[1])
		})
	}
}

func tokenTypes(tokens []*token.Token) []int {
	types := make([]int, len(tokens))
	for i, tok := range tokens {
		types[i] = tok.TokenType
	}
	return types
}

func TestEmptySource(t *testing.T) {
	tokens, err := scan("")
	require.NoError(t, err)
	require.Equal(t, []int{token.EOF}, tokenTypes(tokens))
}

func TestWhitespaceOnly(t *testing.T) {
	tokens, err := scan(" \t\r\n\n ")
	require.NoError(t, err)
	require.Equal(t, []int{token.EOF}, tokenTypes(tokens))
	require.Equal(t, 3, tokens[0].Line)
}

func TestCommentConsumesToEndOfLine(t *testing.T) {
	tokens, err := scan("1 // comment ( ) + -\n2")
	require.NoError(t, err)
	require.Equal(t, []int{token.NUMBER, token.NUMBER, token.EOF}, tokenTypes(tokens))
	require.Equal(t, 1, tokens[0].Line)
	require.Equal(t, 2, tokens[1].Line)
}

func TestCommentAtEOFWithNoTrailingNewline(t *testing.T) {
	tokens, err := scan("// only a comment")
	require.NoError(t, err)
	require.Equal(t, []int{token.EOF}, tokenTypes(tokens))
}

func TestMaximalMunchOnOperators(t *testing.T) {
	tokens, err := scan("!===<=>=")
	require.NoError(t, err)
	require.Equal(t, []int{
		token.BANG_EQUAL, token.EQUAL_EQUAL, token.LESS_EQUAL, token.GREATER_EQUAL, token.EOF,
	}, tokenTypes(tokens))
}

func TestMultilineString(t *testing.T) {
	tokens, err := scan("\"line1\nline2\" x")
	require.NoError(t, err)
	require.Equal(t, "line1\nline2", tokens[0].Literal)
	require.Equal(t, 2, tokens[0].Line)
	require.Equal(t, 2, tokens[1].Line, "line count should carry over after multiline string")
}

func TestEmptyString(t *testing.T) {
	tokens, err := scan("\"\"")
	require.NoError(t, err)
	require.Equal(t, token.STRING, tokens[0].TokenType)
	require.Equal(t, "", tokens[0].Literal)
}

func TestUnterminatedStringIsError(t *testing.T) {
	_, err := scan("\"unterminated")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Unterminated string")
}

func TestUnterminatedStringAtEOFAfterNewline(t *testing.T) {
	_, err := scan("\"abc\ndef")
	require.Error(t, err)
	require.Contains(t, err.Error(), "[line 2]")
}

func TestIntegerWithoutFraction(t *testing.T) {
	tokens, err := scan("123")
	require.NoError(t, err)
	require.Equal(t, float64(123), tokens[0].Literal)
}

func TestTrailingDotIsNotConsumedByNumber(t *testing.T) {
	// A digit must follow the '.' for it to be part of the number literal.
	tokens, err := scan("123.")
	require.NoError(t, err)
	require.Equal(t, []int{token.NUMBER, token.DOT, token.EOF}, tokenTypes(tokens))
	require.Equal(t, float64(123), tokens[0].Literal)
}

func TestLeadingDotIsNotPartOfNumber(t *testing.T) {
	tokens, err := scan(".5")
	require.NoError(t, err)
	require.Equal(t, []int{token.DOT, token.NUMBER, token.EOF}, tokenTypes(tokens))
}

func TestKeywordIsNotMatchedAsPrefix(t *testing.T) {
	tokens, err := scan("classify")
	require.NoError(t, err)
	require.Equal(t, token.IDENTIFIER, tokens[0].TokenType)
	require.Equal(t, "classify", tokens[0].Lexeme)
}

func TestIdentifierWithUnderscoreAndDigits(t *testing.T) {
	tokens, err := scan("_foo1_bar2")
	require.NoError(t, err)
	require.Equal(t, token.IDENTIFIER, tokens[0].TokenType)
	require.Equal(t, "_foo1_bar2", tokens[0].Lexeme)
}

func TestUnexpectedCharacterIsError(t *testing.T) {
	_, err := scan("@")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Unexpected character")
}

func TestMultipleErrorsAreAllCollected(t *testing.T) {
	_, err := scan("@\n#\n$")
	require.Error(t, err)
	errs, ok := err.(SourceErrors)
	require.True(t, ok)
	require.Len(t, errs, 3)
	require.Contains(t, err.Error(), "[line 1]")
	require.Contains(t, err.Error(), "[line 2]")
	require.Contains(t, err.Error(), "[line 3]")
}
