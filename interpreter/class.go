package interpreter

import (
	"fmt"
)

var (
	_ Callable     = (*LoxClass)(nil)
	_ fmt.Stringer = (*LoxClass)(nil)
)

type LoxClass struct {
	name       string
	methods    map[string]*LoxFunction
	superclass *LoxClass
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		name:       name,
		methods:    methods,
		superclass: superclass,
	}
}

func (l *LoxClass) String() string {
	return l.name
}

func (l *LoxClass) Call(interpreter *Interpreter, arguments []any) any {
	instance := NewLoxInstance(l)
	initializer := l.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (l *LoxClass) Arity() int {
	initializer := l.FindMethod("init")
	if initializer != nil {
		return initializer.Arity()
	}

	return 0
}

func (l *LoxClass) FindMethod(name string) *LoxFunction {
	if method, ok := l.methods[name]; ok {
		return method
	}

	if l.superclass != nil {
		return l.superclass.FindMethod(name)
	}

	return nil
}
