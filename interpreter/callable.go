package interpreter

import (
	"fmt"
	"time"

	"github.com/DanielleB-R/golox/interpreter/ast"
)

var (
	_ Callable     = (*NativeFunction)(nil)
	_ fmt.Stringer = (*NativeFunction)(nil)
	_ Callable     = (*LoxFunction)(nil)
	_ fmt.Stringer = (*LoxFunction)(nil)
)

type Callable interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	Arity() int
}

type NativeFunction struct {
	arity     int
	behaviour func(*Interpreter, []interface{}) interface{}
}

func (*NativeFunction) String() string {
	return "<native fn>"
}

func (n *NativeFunction) Arity() int {
	return n.arity
}

func (n *NativeFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return n.behaviour(interpreter, arguments)
}

var Clock *NativeFunction = &NativeFunction{
	arity: 0,
	behaviour: func(interpreter *Interpreter, arguments []interface{}) interface{} {
		return time.Now().Unix()
	},
}

type LoxFunction struct {
	declaration *ast.Function
	closure     *Environment
}

func NewLoxFunction(declaration *ast.Function, closure *Environment) *LoxFunction {
	return &LoxFunction{
		declaration: declaration,
		closure:     closure,
	}
}

func (l *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	environment := NewEnvironment(l.closure)
	for i, param := range l.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(l.declaration.Body, environment)
	returnValue := interpreter.activeReturnValue
	interpreter.resetReturnValue()
	return returnValue
}

func (l *LoxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", l.declaration.Name.Lexeme)
}

func (l *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(l.closure)
	environment.Define("this", instance)
	return NewLoxFunction(l.declaration, environment)
}
