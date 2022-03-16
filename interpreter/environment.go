package interpreter

import (
	"fmt"

	"github.com/DanielleB-R/golox/interpreter/token"
)

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Assign(name *token.Token, value interface{}) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	return &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}

func (e *Environment) Get(name *token.Token) (interface{}, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}
	return nil, &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}
