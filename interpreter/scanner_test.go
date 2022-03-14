package interpreter

import (
	"testing"

	"github.com/DanielleB-R/golox/interpreter/token"
	"github.com/stretchr/testify/require"
)

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
			scanner := NewSourceScanner(input)
			tokens, err := scanner.ScanTokens()
			require.NoError(t, err)
			require.Len(t, tokens, 2)
			require.Equal(t, token, tokens[0])
			require.Equal(t, eofToken, tokens[1])
		})
	}
}
