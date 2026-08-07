package interpreter

import (
	"fmt"

	"github.com/DanielleB-R/golox/interpreter/token"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    nil,
		enclosing: enclosing,
	}
}

func (e *Environment) Define(name string, value any) {
	if e.values == nil {
		e.values = map[string]any{name: value}
	} else {
		e.values[name] = value
	}
}

func (e *Environment) Assign(name *token.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	return &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}

func (e *Environment) AssignAt(distance int, name *token.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e *Environment) Get(name *token.Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}

// This panics if the name doesn't exist, but we assume the resolver got this right
func (e *Environment) GetAt(distance int, name string) (any, error) {
	return e.ancestor(distance).values[name], nil
}

// Panics if the distance is beyond the number of environments in the chain
func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for range distance {
		environment = environment.enclosing
	}
	return environment
}
