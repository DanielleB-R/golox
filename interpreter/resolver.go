package interpreter

import (
	"fmt"

	"github.com/DanielleB-R/golox/interpreter/ast"
	"github.com/DanielleB-R/golox/interpreter/token"
)

var (
	_ ast.ExprVisitor = (*Resolver)(nil)
	_ ast.StmtVisitor = (*Resolver)(nil)
)

type FunctionType = int

const (
	NO_FUNCTION FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

type ClassType = int

const (
	NO_CLASS ClassType = iota
	CLASS
	SUBCLASS
)

type Scope = map[string]bool

type Resolver struct {
	interpreter     *Interpreter
	scopes          []Scope
	currentFunction FunctionType
	currentClass    ClassType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          nil,
		currentFunction: NO_FUNCTION,
		currentClass:    NO_CLASS,
	}
}

func (r *Resolver) Resolve(statements []ast.Stmt) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		if readerErr, ok := err.(*ReaderError); ok {
			fmt.Println(readerErr.Error())
			return
		}
		panic(err)
	}()

	r.resolveStmts(statements)
}

func (r *Resolver) resolveStmts(statements []ast.Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) VisitBlock(stmt *ast.Block) {
	r.beginScope()
	r.resolveStmts(stmt.Statements)
	r.endScope()
}

func (r *Resolver) VisitClass(stmt *ast.Class) {
	enclosingClass := r.currentClass
	defer func() { r.currentClass = enclosingClass }()
	r.currentClass = CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil {
		if stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
			panic(&ReaderError{token: stmt.Name, message: "A class cannot inherit from itself"})
		}
		r.currentClass = SUBCLASS
		r.resolveExpr(stmt.Superclass)
		r.beginScope()
		defer r.endScope()
		r.scopes[len(r.scopes)-1]["super"] = true
	}

	r.beginScope()
	defer r.endScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range stmt.Methods {
		declaration := METHOD
		if method.Name.Lexeme == "init" {
			declaration = INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) {
	r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitFunction(stmt *ast.Function) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FUNCTION)
}

func (r *Resolver) VisitIf(stmt *ast.If) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
}

func (r *Resolver) VisitPrint(stmt *ast.Print) {
	r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitReturn(stmt *ast.Return) {
	if r.currentFunction == NO_FUNCTION {
		panic(&ReaderError{token: stmt.Keyword, message: "Can't return from top-level code"})
	}

	if stmt.Value != nil {
		if r.currentFunction == INITIALIZER {
			panic(&ReaderError{token: stmt.Keyword, message: "Can't return a value from an initializer"})
		}

		r.resolveExpr(stmt.Value)
	}
}

func (r *Resolver) VisitVar(stmt *ast.Var) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
}
func (r *Resolver) VisitWhile(stmt *ast.While) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitAssign(expr *ast.Assign) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinary(expr *ast.Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCall(expr *ast.Call) any {
	r.resolveExpr(expr.Callee)

	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitGet(expr *ast.Get) any {
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitGrouping(expr *ast.Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteral(expr *ast.Literal) any {
	return nil
}

func (r *Resolver) VisitLogical(expr *ast.Logical) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitSet(expr *ast.Set) any {
	r.resolveExpr(expr.Object)
	r.resolveExpr(expr.Value)
	return nil
}

func (r *Resolver) VisitSuper(expr *ast.Super) any {
	if r.currentClass == NO_CLASS {
		panic(&ReaderError{token: expr.Keyword, message: "Cannot use 'super' outside of a class"})
	}
	if r.currentClass != SUBCLASS {
		panic(&ReaderError{token: expr.Keyword, message: "Can't use 'super' in a class with no superclasses"})
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitThis(expr *ast.This) any {
	if r.currentClass == NO_CLASS {
		panic(&ReaderError{token: expr.Keyword, message: "Cannot use 'this' outside of a class"})
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitUnary(expr *ast.Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariable(expr *ast.Variable) any {
	if len(r.scopes) > 0 {
		if value, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; ok && value == false {
			panic(&ReaderError{token: expr.Name, message: "Can't read local variable in its own initializer"})
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, Scope{})
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[0:(len(r.scopes) - 1)]
}

func (r *Resolver) declare(name *token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.Lexeme]; ok {
		panic(&ReaderError{token: name, message: "Already a variable with this name in this scope"})
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name *token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(stmt *ast.Function, functionType FunctionType) {
	previousFunctionType := r.currentFunction
	r.currentFunction = functionType
	defer func() { r.currentFunction = previousFunctionType }()

	r.beginScope()
	defer r.endScope()
	for _, param := range stmt.Params {
		r.declare(param)
		r.define(param)

	}
	r.resolveStmts(stmt.Body)
}
