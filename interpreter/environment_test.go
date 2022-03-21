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
	environment := NewEnvironment(nil)

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
	require.Equal(t, "print", result)

	err = environment.Assign(tokenNamed("a"), 200)
	require.NoError(t, err)

	result, err = environment.Get(tokenNamed("a"))
	require.NoError(t, err)
	require.Equal(t, 200, result)

	err = environment.Assign(tokenNamed("c"), nil)
	require.Error(t, err)
}

func TestEnclosingEnvironments(t *testing.T) {
	outer := NewEnvironment(nil)

	outer.Define("a", 20)
	outer.Define("b", "print")

	inner := NewEnvironment(outer)

	result, err := inner.Get(tokenNamed("a"))
	require.NoError(t, err)
	require.Equal(t, 20, result)

	err = inner.Assign(tokenNamed("a"), 200)
	require.NoError(t, err)

	result, err = inner.Get(tokenNamed("a"))
	require.NoError(t, err)
	require.Equal(t, 200, result)

	result, err = outer.Get(tokenNamed("a"))
	require.NoError(t, err)
	require.Equal(t, 200, result)

	inner.Define("b", "var")

	result, err = inner.Get(tokenNamed("b"))
	require.NoError(t, err)
	require.Equal(t, "var", result)

	result, err = outer.Get(tokenNamed("b"))
	require.NoError(t, err)
	require.Equal(t, "print", result)

	result, err = inner.Get(tokenNamed("c"))
	require.Error(t, err)

	err = inner.Assign(tokenNamed("c"), nil)
	require.Error(t, err)
}

func TestAncestor(t *testing.T) {
	one := NewEnvironment(nil)
	two := NewEnvironment(one)
	three := NewEnvironment(two)

	require.Equal(t, three, three.ancestor(0))
	require.Equal(t, two, three.ancestor(1))
	require.Equal(t, one, three.ancestor(2))
}
