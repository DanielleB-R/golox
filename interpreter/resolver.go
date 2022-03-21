package interpreter

import (
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
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          []map[string]bool{},
		currentFunction: NO_FUNCTION,
	}
}

func (r *Resolver) Resolve(statements []ast.Stmt) {
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
	r.Resolve(stmt.Statements)
	r.endScope()
}

func (r *Resolver) VisitClass(stmt *ast.Class) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
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
		panic("Can't return from top-level code")
	}

	if stmt.Value != nil {
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

func (r *Resolver) VisitAssign(expr *ast.Assign) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinary(expr *ast.Binary) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCall(expr *ast.Call) interface{} {
	r.resolveExpr(expr.Callee)

	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitGrouping(expr *ast.Grouping) interface{} {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteral(expr *ast.Literal) interface{} {
	return nil
}

func (r *Resolver) VisitLogical(expr *ast.Logical) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnary(expr *ast.Unary) interface{} {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariable(expr *ast.Variable) interface{} {
	if len(r.scopes) > 0 {
		if value, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; ok && value == false {
			// TODO: error handling is getting messy
			panic("Can't read local variable in its own initializer")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
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
		// TODO: this should be better
		panic("Already a variable with this name in this scope")
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

	r.beginScope()
	for _, param := range stmt.Params {
		r.declare(param)
		r.define(param)

	}
	r.Resolve(stmt.Body)
	r.endScope()

	r.currentFunction = previousFunctionType
}
