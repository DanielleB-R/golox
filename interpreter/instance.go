package interpreter

import (
	"fmt"

	"github.com/DanielleB-R/golox/interpreter/token"
)

var (
	_ fmt.Stringer = (*LoxInstance)(nil)
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: map[string]any{},
	}
}

func (l *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", l.class.name)
}

func (l *LoxInstance) Get(name *token.Token) (any, error) {
	if value, ok := l.fields[name.Lexeme]; ok {
		return value, nil
	}

	method := l.class.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(l), nil
	}

	return nil, &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined property '%s'", name.Lexeme),
	}
}

func (l *LoxInstance) Set(name *token.Token, value any) {
	l.fields[name.Lexeme] = value
}
