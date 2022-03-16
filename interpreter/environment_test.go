package interpreter

import (
	"testing"

	"github.com/DanielleB-R/golox/interpreter/token"
	"github.com/stretchr/testify/require"
)

func tokenNamed(name string) *token.Token {
	return &token.Token{
		Lexeme:    name,
		Line:      0,
		Literal:   nil,
		TokenType: token.IDENTIFIER,
	}
}

func TestBasicEnvironment(t *testing.T) {
	environment := NewEnvironment()

	require.Len(t, environment.values, 0)

	environment.Define("a", 20)

	result, err := environment.Get(tokenNamed("a"))
	require.NoError(t, err)
	require.Equal(t, 20, result)

	result, err = environment.Get(tokenNamed("b"))
	require.Error(t, err)

	environment.Define("b", "print")

	result, err = environment.Get(tokenNamed("b"))
	require.NoError(t, err)
	require.Equal(t, result, "print")

	err = environment.Assign(tokenNamed("a"), 200)
	require.NoError(t, err)

	result, err = environment.Get(tokenNamed("a"))
	require.NoError(t, err)
	require.Equal(t, result, 200)

	err = environment.Assign(tokenNamed("c"), nil)
	require.Error(t, err)
}
